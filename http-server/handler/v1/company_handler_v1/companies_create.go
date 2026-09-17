package company_handler_v1

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"neomatica/neosync/infra/constants"
	"neomatica/neosync/infra/locale"
	"neomatica/neosync/infra/store/postgres/models"
	"neomatica/neosync/infra/worker"
	"neomatica/neosync/middleware"
	"neomatica/neosync/pkg/httpx"
	"neomatica/neosync/pkg/httpx/httperr"
	"neomatica/neosync/pkg/permissions"
	gtypes "neomatica/neosync/types"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator"
)

const maxCompanyDevices = 250
const maxCompanyConfigurationFileSize = 5 * 1024

/* Neosync HTTPx V1 */
/* Handler: создание массовой настройки и задач для устройств */

func normalizeConfigurationSource(source string) string {
	switch source {
	case "device":
		return "device_sources"
	case "template":
		return "template_sources"
	case "file":
		return "file_sources"
	default:
		return source
	}
}

func (h *Handler) CompanyCreateHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*gtypes.AuthToken)

	if !authToken.User.Has(constants.AccessConfigurationApply) {
		return httperr.Forbidden(tr.TErr("config-edit-access-restricted"))
	}

	var payload *CreateCompanyPayload

	if err := httpx.HttpParse(r, &payload); err != nil {
		return httperr.BadRequest(err.Error())
	}

	if err := httpx.Validate.Struct(payload); err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			return httperr.BadRequest(httpx.ValidateMsg(tr, err))
		}

		return httperr.BadRequest(tr.TErr("fields-not-filled"))
	}

	payload.IncompatibleAction = "exclude"

	configurationSource := normalizeConfigurationSource(payload.ConfigurationSource)
	if configurationSource == "" {
		configurationSource = normalizeConfigurationSource(payload.ConfigurationCource)
	}

	if configurationSource == "" {
		return httperr.BadRequest(tr.TErr("fields-not-filled"))
	}

	if len(payload.IMEIs) > maxCompanyDevices {
		return httperr.BadRequest(tr.TErr("provisioning-company-max-devices-exceeded"))
	}

	configurationSourceData, err := h.resolveConfigurationSourceData(ctx, tr, authToken, configurationSource, payload)
	if err != nil {
		return err
	}

	existingTaskAction := payload.ExistingTaskAction
	if existingTaskAction == "replace" {
		existingTaskAction = "overwrite"
	}

	ttl, err := strconv.Atoi(payload.TTL)
	if err != nil {
		return httperr.BadRequest(tr.TErr("fields-not-filled"))
	}

	devices, _, priorityModel, err := h.resolveEligibleDevices(ctx, authToken, payload.IMEIs, payload.PriorityModel)
	if err != nil {
		return err
	}

	if err := h.validateConfigurationSource(ctx, tr, authToken, configurationSource, payload, devices, priorityModel); err != nil {
		return err
	}

	if configurationSource == "device_sources" {
		sourceImei := strings.TrimSpace(*payload.ConfigurationDeviceCource)
		sourceFound := false

		for _, device := range devices {
			if device.IMEI == sourceImei {
				sourceFound = true
				break
			}
		}

		if !sourceFound {
			return httperr.BadRequest(tr.TErr("provisioning-configuration-device-source-not-found"))
		}
	}

	if len(devices) == 0 {
		return httperr.BadRequest(tr.TErr("provisioning-no-compatible-devices"))
	}

	deviceIDs := make([]uint64, 0, len(devices))
	for _, device := range devices {
		deviceIDs = append(deviceIDs, device.ID)
	}

	if existingTaskAction == "skip" {
		pendingDeviceIDs, err := h.Store.Companies.Get_PendingCompanyTaskDeviceIds(ctx, deviceIDs)
		if err != nil {
			return httperr.Db(ctx, err)
		}

		if len(pendingDeviceIDs) > 0 {
			pendingSet := make(map[uint64]struct{}, len(pendingDeviceIDs))
			for _, deviceID := range pendingDeviceIDs {
				pendingSet[deviceID] = struct{}{}
			}

			filteredDevices := make([]*models.Device, 0, len(devices))
			for _, device := range devices {
				if _, exists := pendingSet[device.ID]; exists {
					continue
				}

				filteredDevices = append(filteredDevices, device)
			}

			devices = filteredDevices
			deviceIDs = deviceIDs[:0]

			for _, device := range devices {
				deviceIDs = append(deviceIDs, device.ID)
			}
		}
	}

	if len(deviceIDs) == 0 {
		if existingTaskAction == "skip" {
			return httperr.BadRequest(tr.TErr("provisioning-all-devices-skipped-active-tasks"))
		}

		return httperr.BadRequest(tr.TErr("provisioning-no-devices-to-create"))
	}

	tx, err := h.Db.BeginTx(ctx, nil)
	if err != nil {
		return httperr.Db(ctx, httperr.Err_DbNetwork)
	}

	defer func() {
		_ = tx.Rollback()
	}()

	if existingTaskAction == "overwrite" {
		if err := h.Store.Companies.Cancel_PendingCompanyTasksByDeviceIds(ctx, tx, deviceIDs); err != nil {
			return httperr.Db(ctx, err)
		}
	}

	metadata := map[string]interface{}{
		"version":        1,
		"priority_model": priorityModel,
		"selected_imeis": payload.IMEIs,
	}

	switch configurationSource {
	case "device_sources":
		metadata["configuration_device_cource"] = payload.ConfigurationDeviceCource

	case "template_sources":
		metadata["configuration_template_id"] = payload.ConfigurationTemplateID

	case "file_sources":
		metadata["configuration_file_name"] = payload.ConfigurationFileName

	}

	sourceMetadata, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	company := &models.Company{
		SourceMetadata:          sourceMetadata,
		UserUUID:                authToken.User.UserContact.UserUUID,
		Name:                    payload.Name,
		Status:                  "pending",
		TTL:                     ttl,
		ExistingTaskAction:      existingTaskAction,
		LaunchMode:              payload.LaunchMode,
		IncompatibleAction:      payload.IncompatibleAction,
		ConfigurationSource:     configurationSource,
		ConfigurationSourceData: configurationSourceData,
	}

	companyID, err := h.Store.Companies.Create_Company(ctx, tx, company)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	if err := tx.Commit(); err != nil {
		return httperr.Conflict(tr.TErr("failed-to-save-data"))
	}

	if h.CompanyWorker != nil {
		h.CompanyWorker.Enqueue(worker.CompanyTasksJob{
			CompanyID: companyID,
			DeviceIDs: deviceIDs,
			UserUUID:  authToken.User.UserContact.UserUUID,
		})
	}

	response := CreateCompanyResponse{
		ID:           companyID,
		TasksCreated: len(deviceIDs),
		Message:      tr.T("company-created-successfully"),
	}

	httpx.HttpResponse(w, r, http.StatusCreated, response)
	return nil
}

func (h *Handler) resolveConfigurationSourceData(ctx context.Context, tr locale.Translator, authToken *gtypes.AuthToken, configurationSource string, payload *CreateCompanyPayload) ([]byte, error) {
	switch configurationSource {
	case "device_sources":
		if payload.ConfigurationDeviceCource == nil || strings.TrimSpace(*payload.ConfigurationDeviceCource) == "" {
			return nil, httperr.BadRequest(tr.TErr("provisioning-configuration-device-source-required"))
		}

		sourceDevice, err := h.Store.Devices.Get_DeviceByImei(ctx, strings.TrimSpace(*payload.ConfigurationDeviceCource))
		if err != nil {
			return nil, httperr.Db(ctx, err)
		}

		if sourceDevice == nil {
			return nil, httperr.BadRequest(tr.TErr("provisioning-configuration-device-source-not-found"))
		}

		configuration, err := h.Store.Configurations.Get_ConfigurationByDeviceId(ctx, sourceDevice.ID)
		if err != nil {
			return nil, httperr.Db(ctx, err)
		}

		if configuration == nil || configuration.CfgData == nil || len(configuration.CfgData) == 0 {
			return nil, httperr.BadRequest(tr.TErr("provisioning-configuration-device-source-no-config"))
		}

		if configuration.CfgSyncStatus != constants.CFG_SYNC_STATUS_CONFIRMED {
			return nil, httperr.BadRequest(tr.TErr("provisioning-configuration-device-source-not-confirmed"))
		}

		return configuration.CfgData, nil
	case "template_sources":
		if payload.ConfigurationTemplateID == nil || *payload.ConfigurationTemplateID == 0 {
			return nil, httperr.BadRequest(tr.TErr("provisioning-configuration-template-source-required"))
		}

		template, err := h.Store.ConfigurationTemplates.Get_ConfigurationTemplateByIdAndUserUuid(ctx, *payload.ConfigurationTemplateID, authToken.User.UserContact.UserUUID)
		if err != nil {
			return nil, httperr.Db(ctx, err)
		}

		if template == nil {
			return nil, httperr.BadRequest(tr.TErr("provisioning-configuration-template-source-not-found"))
		}

		cfgData, err := h.Store.ConfigurationTemplates.Get_ConfigurationTemplateCfgDataByIdAndUserUuid(ctx, *payload.ConfigurationTemplateID, authToken.User.UserContact.UserUUID)
		if err != nil {
			return nil, httperr.Db(ctx, err)
		}

		if cfgData == nil {
			return nil, httperr.BadRequest(tr.TErr("provisioning-configuration-template-source-not-found"))
		}

		return cfgData, nil
	case "file_sources":
		if payload.ConfigurationFileBase64 == nil || *payload.ConfigurationFileBase64 == "" {
			return nil, httperr.BadRequest(tr.TErr("provisioning-configuration-file-source-required"))
		}

		fileData, err := base64.StdEncoding.DecodeString(*payload.ConfigurationFileBase64)
		if err != nil {
			return nil, httperr.BadRequest(tr.TErr("provisioning-configuration-file-source-invalid"))
		}

		if len(fileData) == 0 || len(fileData) > maxCompanyConfigurationFileSize {
			return nil, httperr.BadRequest(tr.TErr("provisioning-configuration-file-source-invalid"))
		}

		if payload.ConfigurationFileName == nil || strings.TrimSpace(*payload.ConfigurationFileName) == "" {
			return nil, httperr.BadRequest(tr.TErr("provisioning-configuration-file-source-required"))
		}

		return fileData, nil
	default:
		return nil, httperr.BadRequest(tr.TErr("fields-not-filled"))
	}
}

func (h *Handler) resolveEligibleDevices(ctx context.Context, authToken *gtypes.AuthToken, imeis []string, requestedModel string) ([]*models.Device, int, string, error) {
	tr := middleware.TranslatorFromContext(ctx)

	devices := make([]*models.Device, 0, len(imeis))
	modelCount := make(map[string]int)
	incompatibleCount := 0

	for _, imei := range imeis {
		device, err := h.Store.Devices.Get_DeviceByImei(ctx, imei)
		if err != nil {
			return nil, 0, "", httperr.Db(ctx, err)
		}

		if device == nil {
			return nil, 0, "", httperr.NotFound(fmt.Sprintf("%s: %s", tr.TErr("device-not-found"), imei))
		}

		if !permissions.HasDeviceAccount(device.UserUUID) {
			return nil, 0, "", httperr.New(fmt.Sprintf("%s: %s", tr.TErr("device-has-no-account"), imei), http.StatusUnprocessableEntity)
		}

		if !permissions.IsMainRole(authToken.RoleCode) && !authToken.User.CanViewChildGroups && !permissions.CanAccessDevice(ctx, h.Store.GroupMembers, device, device.ID, &authToken.User, authToken.RoleCode, permissions.ShareEditConfig) {
			return nil, 0, "", httperr.Forbidden(fmt.Sprintf("%s: %s", tr.TErr("device-not-owned"), imei))
		}

		devices = append(devices, device)

		if device.DeviceModel != nil && *device.DeviceModel != "" {
			modelCount[*device.DeviceModel]++
		}
	}

	priorityModel, valid := selectPriorityModel(modelCount, strings.TrimSpace(requestedModel))
	if !valid {
		return nil, 0, "", httperr.BadRequest(tr.TErr("provisioning-primary-model-not-selected"))
	}

	eligibleDevices := make([]*models.Device, 0, len(devices))

	for _, device := range devices {
		if priorityModel != "" && (device.DeviceModel == nil || *device.DeviceModel != priorityModel) {
			incompatibleCount++
			continue
		}

		eligibleDevices = append(eligibleDevices, device)
	}

	return eligibleDevices, incompatibleCount, priorityModel, nil
}

// selectPriorityModel keeps the automatic fallback for clients without a model selection.
func selectPriorityModel(counts map[string]int, requestedModel string) (string, bool) {
	if requestedModel != "" {
		return requestedModel, counts[requestedModel] > 0
	}
	model := ""
	maxCount := 0
	for candidate, count := range counts {
		if count > maxCount || (count == maxCount && candidate < model) {
			model = candidate
			maxCount = count
		}
	}
	return model, true
}

func (h *Handler) validateConfigurationSource(ctx context.Context, tr locale.Translator, authToken *gtypes.AuthToken, configurationSource string, payload *CreateCompanyPayload, eligibleDevices []*models.Device, priorityModel string) error {
	switch configurationSource {
	case "device_sources":
		if payload.ConfigurationDeviceCource == nil || strings.TrimSpace(*payload.ConfigurationDeviceCource) == "" {
			return httperr.BadRequest(tr.TErr("provisioning-configuration-device-source-required"))
		}

		sourceImei := strings.TrimSpace(*payload.ConfigurationDeviceCource)
		var sourceDevice *models.Device

		for _, device := range eligibleDevices {
			if device.IMEI == sourceImei {
				sourceDevice = device
				break
			}
		}

		if sourceDevice == nil {
			return httperr.BadRequest(tr.TErr("provisioning-configuration-device-source-incompatible"))
		}

		if priorityModel != "" && sourceDevice.DeviceModel != nil && *sourceDevice.DeviceModel != priorityModel {
			return httperr.BadRequest(tr.TErr("provisioning-configuration-device-source-model-mismatch"))
		}

		dominantFirmware, err := h.dominantFirmwareVersion(ctx, eligibleDevices)
		if err != nil {
			return httperr.Db(ctx, err)
		}

		if dominantFirmware > 0 {
			sourceSync, err := h.Store.Syncs.Get_SyncByDeviceId(ctx, sourceDevice.ID)
			if err != nil {
				return httperr.Db(ctx, err)
			}

			if sourceSync == nil || sourceSync.FirmwareVersion == 0 {
				return httperr.BadRequest(tr.TErr("provisioning-configuration-device-source-firmware-unknown"))
			}

			if sourceSync.FirmwareVersion != dominantFirmware {
				return httperr.BadRequest(tr.TErr("provisioning-configuration-device-source-firmware-mismatch"))
			}
		}

		return nil
	case "template_sources":
		if payload.ConfigurationTemplateID == nil || *payload.ConfigurationTemplateID == 0 {
			return httperr.BadRequest(tr.TErr("provisioning-configuration-template-source-required"))
		}

		template, err := h.Store.ConfigurationTemplates.Get_ConfigurationTemplateByIdAndUserUuid(ctx, *payload.ConfigurationTemplateID, authToken.User.UserContact.UserUUID)
		if err != nil {
			return httperr.Db(ctx, err)
		}

		if template == nil {
			return httperr.BadRequest(tr.TErr("provisioning-configuration-template-source-not-found"))
		}

		if priorityModel != "" && template.Model != nil && strings.TrimSpace(*template.Model) != "" && *template.Model != priorityModel {
			return httperr.BadRequest(tr.TErr("provisioning-configuration-template-model-mismatch"))
		}

		return nil
	default:
		return nil
	}
}

func (h *Handler) dominantFirmwareVersion(ctx context.Context, devices []*models.Device) (uint16, error) {
	counts := make(map[uint16]int)

	for _, device := range devices {
		sync, err := h.Store.Syncs.Get_SyncByDeviceId(ctx, device.ID)
		if err != nil {
			return 0, err
		}

		if sync == nil || sync.FirmwareVersion == 0 {
			continue
		}

		counts[sync.FirmwareVersion]++
	}

	dominant := uint16(0)
	maxCount := 0

	for firmware, count := range counts {
		if count > maxCount || (count == maxCount && firmware < dominant) {
			maxCount = count
			dominant = firmware
		}
	}

	return dominant, nil
}
