package configuration_handler_v1

import (
	"database/sql"
	"encoding/json"
	"errors"
	"neomatica/neosync/infra/constants"
	"neomatica/neosync/infra/logger"
	"neomatica/neosync/middleware"
	"neomatica/neosync/pkg/httpx"
	"neomatica/neosync/pkg/httpx/httperr"
	gtypes "neomatica/neosync/types"
	"net/http"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
)

/* Neosync HTTPx V1 */
/* Handler: обновление шаблона конфигурации */

func (h *Handler) UpdateConfigurationTemplateHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*gtypes.AuthToken)

	if !authToken.User.Has(constants.AccessConfigurationApply) {
		return httperr.Forbidden(tr.TErr("config-edit-access-restricted"))
	}

	templateID, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil || templateID == 0 {
		return httperr.BadRequest(tr.TErr("fields-not-filled"))
	}

	existingTemplate, err := h.Store.ConfigurationTemplates.Get_ConfigurationTemplateByIdAndUserUuid(ctx, templateID, authToken.User.UserContact.UserUUID)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	if existingTemplate == nil {
		return httperr.NotFound(tr.TErr("configuration-template-not-found"))
	}

	var payload *ConfigurationTemplateUpdate

	if err := httpx.HttpParse(r, &payload); err != nil {
		return httperr.BadRequest(err.Error())
	}

	if err := httpx.Validate.Struct(payload); err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			return httperr.BadRequest(httpx.ValidateMsg(tr, err))
		}

		return httperr.BadRequest(tr.TErr("fields-not-filled"))
	}

	if len(payload.Changes) == 0 {
		return httperr.BadRequest(tr.TErr("no-changes-provided"))
	}

	typeOfSaving := existingTemplate.TypeOfSaving
	if payload.TypeOfSaving != "" {
		typeOfSaving = payload.TypeOfSaving
	}

	existingCfgData, err := h.Store.ConfigurationTemplates.Get_ConfigurationTemplateCfgDataByIdAndUserUuid(ctx, templateID, authToken.User.UserContact.UserUUID)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	var cfgData []byte

	switch typeOfSaving {
	case "modified_ones":
		changesData := ConfigurationTemplateChangesData{
			CfgHash: payload.CfgHash,
			Changes: map[string]any{},
		}

		if len(existingCfgData) > 0 {
			if err := json.Unmarshal(existingCfgData, &changesData); err != nil {
				return httperr.BadRequest(tr.TErr("cfg-parse-error-device"))
			}
		}

		if changesData.Changes == nil {
			changesData.Changes = map[string]any{}
		}

		for key, value := range payload.Changes {
			changesData.Changes[key] = value
		}

		if payload.CfgHash != 0 {
			changesData.CfgHash = payload.CfgHash
		}

		cfgData, err = json.Marshal(changesData)
		if err != nil {
			logger.Error("UpdateConfigurationTemplateHandler_V1 req={%s}: Failed marshal changes: %s", ctx.Value("XREQID").(string), err.Error())
			return httperr.InternalServerError(tr.TErr("failed-to-save-data"))
		}
	case "full":
		cfgParsed := getDefaultConfigurationAll()

		if len(existingCfgData) > 0 {
			parsedExisting, errCfg := getConfiguration(tr, existingCfgData)
			if errCfg == nil && parsedExisting != nil {
				cfgParsed = parsedExisting
			}
		}

		cfgData, err = buildConfigurationFromChanges(tr, payload.Changes, cfgParsed)
		if err != nil {
			return httperr.BadRequest(err.Error())
		}
	default:
		return httperr.BadRequest(tr.TErr("fields-not-filled"))
	}

	name := existingTemplate.Name
	if payload.Name != "" {
		name = payload.Name
	}

	if err := h.Store.ConfigurationTemplates.Update_ConfigurationTemplateByIdAndUserUuid(ctx, templateID, authToken.User.UserContact.UserUUID, name, typeOfSaving, cfgData); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return httperr.NotFound(tr.TErr("configuration-template-not-found"))
		}

		return httperr.Db(ctx, err)
	}

	httpx.HttpResponse(w, r, http.StatusOK, tr.T("configuration-template-updated"))
	return nil
}
