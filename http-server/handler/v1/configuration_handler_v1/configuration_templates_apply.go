package configuration_handler_v1

import (
	"encoding/json"
	"fmt"
	"neomatica/neosync/infra/constants"
	"neomatica/neosync/infra/logger"
	"neomatica/neosync/infra/store/postgres/models"
	"neomatica/neosync/infra/store/redis"
	"neomatica/neosync/middleware"
	"neomatica/neosync/pkg/adm"
	"neomatica/neosync/pkg/adm/admparser"
	"neomatica/neosync/pkg/httpx"
	"neomatica/neosync/pkg/httpx/httperr"
	"neomatica/neosync/pkg/permissions"
	gtypes "neomatica/neosync/types"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

/* Neosync HTTPx V1 */
/* Handler: применение шаблона конфигурации к устройству */

func (h *Handler) ApplyConfigurationTemplateHandler_V1(w http.ResponseWriter, r *http.Request) error { //nolint
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*gtypes.AuthToken)

	if !authToken.User.Has(constants.AccessConfigurationApply) {
		return httperr.Forbidden(tr.TErr("config-edit-access-restricted"))
	}

	imei := mux.Vars(r)["imei"]
	templateIDStr := mux.Vars(r)["id"]

	if imei == "" {
		return httperr.BadRequest(tr.TErr("device-imei-not-found"))
	}

	templateID, err := strconv.ParseUint(templateIDStr, 10, 64)
	if err != nil || templateID == 0 {
		return httperr.BadRequest(tr.TErr("fields-not-filled"))
	}

	device, err := h.Store.Devices.Get_DeviceByImei(ctx, imei)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	if device == nil {
		return httperr.NotFound(tr.TErr("device-not-found"))
	}

	if !permissions.CanAccessDevice(ctx, h.Store.GroupMembers, device, device.ID, &authToken.User, authToken.RoleCode, permissions.ShareEditConfig) {
		return httperr.Forbidden(tr.TErr("device-not-owned"))
	}

	if !*device.Activated {
		return httperr.New(tr.TErr("device-not-activated"), http.StatusUnprocessableEntity)
	}

	if !permissions.HasDeviceAccount(device.UserUUID) {
		return httperr.New(tr.TErr("device-has-no-account"), http.StatusUnprocessableEntity)
	}

	template, err := h.Store.ConfigurationTemplates.Get_ConfigurationTemplateByIdAndUserUuid(ctx, templateID, authToken.User.UserContact.UserUUID)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	if template == nil {
		return httperr.NotFound(tr.TErr("configuration-template-not-found"))
	}

	cfgData, err := h.Store.ConfigurationTemplates.Get_ConfigurationTemplateCfgDataByIdAndUserUuid(ctx, templateID, authToken.User.UserContact.UserUUID)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	if cfgData == nil {
		return httperr.NotFound(tr.TErr("configuration-template-not-found"))
	}

	createAnalyticRecord := func(cfgHash uint32, rmqPublished bool, errReason *string) {
		analytic := &models.AnalyticRecord{
			UserUUID:     authToken.User.UserContact.UserUUID,
			DeviceID:     device.ID,
			CfgHash:      cfgHash,
			RMQPublished: rmqPublished,
			DeviceOnline: device.Status,
			ErrorReason:  errReason,
		}

		if err := h.Store.Analytics.Create_AnalyticsRecord(ctx, analytic); err != nil {
			logger.Error("ApplyConfigurationTemplateHandler_V1 req={%s}: Failed create analytics record: %s", ctx.Value("XREQID").(string), err.Error())
		}
	}

	writeAnalyticError := func(cfgHash uint32, message string) {
		reason := message
		createAnalyticRecord(cfgHash, false, &reason)
	}

	var newCfg []byte
	var cfgConverted []adm.PacketFieldConfiguration
	var hash uint32

	switch template.TypeOfSaving {
	case "full":
		newCfg = cfgData
		cfgHashApplied := admparser.GetCfgHash(newCfg)

		cfgParsed, errCfg := getConfiguration(tr, newCfg)
		if errCfg == nil && cfgParsed != nil {
			cfgConverted = adm.ParsedToPacketFields(cfgParsed)
		}

		hash = cfgHashApplied
	case "modified_ones":
		var changesData ConfigurationTemplateChangesData
		if err := json.Unmarshal(cfgData, &changesData); err != nil {
			return httperr.BadRequest(tr.TErr("cfg-parse-error-device"))
		}

		if len(changesData.Changes) == 0 {
			return httperr.BadRequest(tr.TErr("no-changes-provided"))
		}

		var cfgParsed []adm.FieldConfigurationParsed
		var errCfg error

		cfgRedisBytes, err := redis.GetCacheConfiguration(imei)
		if err == nil && cfgRedisBytes != nil {
			cfgParsed, errCfg = getConfiguration(tr, cfgRedisBytes)
			hash = admparser.GetCfgHash(cfgRedisBytes)
		}

		rabbitData := gtypes.RabbitMQ_TransitBinary{
			Type: constants.ADM_RC_TYPE_GET_CFG,
			Imei: imei,
			Data: nil,
		}

		if cfgParsed == nil || errCfg != nil {
			response, err := h.RMQ.SendToRabbitAndWait(ctx, rabbitData, h.Session)
			if err == nil && response.Data != nil {
				cfgParsed, errCfg = getConfiguration(tr, response.Data)
				hash = admparser.GetCfgHash(response.Data)
			}
		}

		if cfgParsed == nil || errCfg != nil {
			configuration, err := h.Store.Configurations.Get_ConfigurationByDeviceId(ctx, device.ID)
			if err != nil {
				message := fmt.Sprintf("ApplyConfigurationTemplateHandler_V1 req={%s}: Failed get configuration from db: %s", ctx.Value("XREQID").(string), err.Error())
				logger.Error(message)
				writeAnalyticError(hash, message)
				return httperr.Conflict(tr.TErr("config-fetch-error"))
			}

			if configuration.CfgData == nil {
				message := fmt.Sprintf("ApplyConfigurationTemplateHandler_V1 req={%s}: Configuration NULL", ctx.Value("XREQID").(string))
				logger.Error(message)
				writeAnalyticError(hash, message)
				return httperr.Conflict(tr.TErr("config-fetch-error"))
			}

			cfgParsed, errCfg = getConfiguration(tr, configuration.CfgData)
			if errCfg != nil {
				message := fmt.Sprintf("ApplyConfigurationTemplateHandler_V1 req={%s}: Failed get configuration parsed: %s", ctx.Value("XREQID").(string), errCfg.Error())
				logger.Error(message)
				writeAnalyticError(hash, message)
				return httperr.BadRequest(tr.TErr("cfg-parse-error-device"))
			}

			hash = configuration.CfgHash
		}

		draft := &redis.ConfigurationDraftInsert{
			DeviceImei: imei,
			CfgHash:    changesData.CfgHash,
			Changes:    changesData.Changes,
		}

		if hash != draft.CfgHash {
			draft.CfgHash = hash
		}

		cfgConverted = convertDraftToParsedWithSize(draft, cfgParsed)
		admInst := adm.Adm{}

		newCfg, err = admInst.ParseJsonToBinary(tr, cfgConverted)
		if err != nil {
			message := fmt.Sprintf("ApplyConfigurationTemplateHandler_V1 req={%s}: Failed build configuration: %s", ctx.Value("XREQID").(string), err.Error())
			logger.Error(message)
			writeAnalyticError(hash, message)
			return httperr.BadRequest(err.Error())
		}
	default:
		return httperr.BadRequest(tr.TErr("fields-not-filled"))
	}

	if newCfg == nil {
		message := fmt.Sprintf("ApplyConfigurationTemplateHandler_V1 req={%s}: New builded configuration is NULL", ctx.Value("XREQID").(string))
		logger.Error(message)
		writeAnalyticError(hash, message)
		return httperr.BadRequest(tr.TErr("cfg-parse-error-device"))
	}

	if err := h.Store.Configurations.Update_ConfigurationByDeviceId(ctx, device.ID, newCfg); err != nil {
		message := fmt.Sprintf("ApplyConfigurationTemplateHandler_V1 req={%s}: Failed update configuration in db: %s", ctx.Value("XREQID").(string), err.Error())
		logger.Error(message)
		writeAnalyticError(admparser.GetCfgHash(newCfg), message)
		return httperr.Db(ctx, err)
	}

	if err := h.Store.Configurations.Update_ConfigurationSyncStatusByDeviceId(ctx, device.ID, constants.CFG_SYNC_STATUS_PENDING, nil); err != nil {
		message := fmt.Sprintf("ApplyConfigurationTemplateHandler_V1 req={%s}: Failed update configuration sync status: %s", ctx.Value("XREQID").(string), err.Error())
		logger.Error(message)
		writeAnalyticError(admparser.GetCfgHash(newCfg), message)
		return httperr.Db(ctx, err)
	}

	cfgHashApplied := admparser.GetCfgHash(newCfg)
	if len(cfgConverted) > 0 {
		admInst := adm.Adm{}
		if profValues, err := admInst.ParseProfValues(cfgConverted); err != nil {
			logger.Error("ApplyConfigurationTemplateHandler_V1 req={%s}: Failed build parsed configuration fields: %s", ctx.Value("XREQID").(string), err.Error())
		} else if err := h.Store.Configurations.Replace_ConfigurationProf(ctx, device.ID, cfgHashApplied, profValues); err != nil {
			logger.Error("ApplyConfigurationTemplateHandler_V1 req={%s}: Failed save parsed configuration: %s", ctx.Value("XREQID").(string), err.Error())
		}
	}

	_ = redis.SetCacheConfiguration(imei, constants.Redis_ConfigurationTTL, newCfg)

	rabbitData2 := gtypes.RabbitMQ_TransitBinary{
		Type: constants.ADM_RC_TYPE_SET_CFG,
		Imei: device.IMEI,
		Data: newCfg,
	}

	if err := h.RMQ.SendToRabbitAsync(ctx, rabbitData2); err != nil {
		message := fmt.Sprintf("ApplyConfigurationTemplateHandler_V1 req={%s}: Failed publish configuration in rabbitmq: %s", ctx.Value("XREQID").(string), err.Error())
		logger.Error(message)
		reason := message
		createAnalyticRecord(cfgHashApplied, false, &reason)
		return httperr.New(tr.TErr("rabbitmq-publish-failed"), http.StatusConflict)
	}

	createAnalyticRecord(cfgHashApplied, true, nil)

	var resp string
	if device.Status {
		resp = tr.T("config-sent-successfully")
	} else {
		resp = tr.T("config-sent-status-false-successfully")
	}

	httpx.HttpResponse(w, r, http.StatusAccepted, resp)
	return nil
}
