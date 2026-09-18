package configuration_handler_v1

import (
	"io"
	"neomatica/neosync/infra/constants"
	"neomatica/neosync/infra/logger"
	"neomatica/neosync/infra/store/redis"
	"neomatica/neosync/middleware"
	"neomatica/neosync/pkg/adm"
	"neomatica/neosync/pkg/adm/admparser"
	"neomatica/neosync/pkg/httpx"
	"neomatica/neosync/pkg/httpx/httperr"
	"neomatica/neosync/pkg/permissions"
	"neomatica/neosync/types"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

func groupBySection(fields []adm.FieldConfigurationParsed) map[string]map[string]any {
	result := make(map[string]map[string]any)

	for _, f := range fields {
		schema, ok := constants.CfgSchema[f.RawUID]
		if !ok {
			continue
		}

		sections := []string{schema.Section}
		if strings.Contains(schema.Section, ";") {
			sections = strings.Split(schema.Section, ";")
		}

		for _, section := range sections {
			if section == "" {
				continue
			}

			if _, ok := result[section]; !ok {
				result[section] = make(map[string]any)
			}

			result[section][f.UID] = f.Value
		}
	}

	return result
}

/* Neosync HTTPx V1 */
/* Handler: импортирование конфигурации из .bin файла */

func (h *Handler) ConfigurationImportByImeiHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*types.AuthToken)

	/* Check access */
	if !authToken.User.Has(constants.AccessConfigurationApply) {
		return httperr.Forbidden(tr.TErr("config-access-restricted"))
	}

	imei := mux.Vars(r)["imei"]

	if imei == "" {
		return httperr.BadRequest(tr.TErr("device-imei-not-found"))
	}

	/* Device ownership verification */
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

	/* Get data from file */
	data, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Error("ConfigurationImportByImeiHandler_V1 req={%s} imei={%s}: %s", ctx.Value("XREQID").(string), imei, err.Error())
		return httperr.InternalServerError(tr.TErr("error-file-body"))
	}
	defer func() {
		_ = r.Body.Close()
	}()

	cfgParsed, errCfg := getConfiguration(tr, data)
	if errCfg != nil {
		logger.Error("ConfigurationImportByImeiHandler_V1 req={%s} imei={%s}: Failed get configuration: %s", ctx.Value("XREQID").(string), imei, errCfg.Error())
		return httperr.BadRequest(errCfg.Error())
	}

	grouped := groupBySection(cfgParsed)

	if err := redis.DeleteAllDrafts(authToken.User.UserContact.UserUUID, imei); err != nil {
		logger.Error("ConfigurationImportByImeiHandler_V1 req={%s} imei={%s}: Failed delete all drafts: %s", ctx.Value("XREQID").(string), imei, err.Error())
		return httperr.InternalServerError(tr.TErr("failed-to-save-draft"))
	}

	for section, changes := range grouped {
		draft := redis.ConfigurationDraftInsert{
			DeviceImei: imei,
			CfgHash:    admparser.GetCfgHash(data),
			Section:    section,
			Changes:    changes,
			Timestamp:  time.Now().Unix(),
		}

		if err := redis.SaveOrUpdateDraft(authToken.User.UserContact.UserUUID, imei, section, draft); err != nil {
			logger.Error("ConfigurationImportByImeiHandler_V1 req={%s} imei={%s} section={%s}: Failed save/update draft: %s", ctx.Value("XREQID").(string), imei, section, err.Error())
			return httperr.InternalServerError(tr.TErr("failed-to-save-draft"))
		}
	}

	httpx.HttpResponse(w, r, http.StatusOK, tr.T("success-import-cfg"))
	return nil
}
