package configuration_handler_v1

import (
	"encoding/json"
	"neomatica/neosync/infra/constants"
	"neomatica/neosync/infra/logger"
	"neomatica/neosync/infra/store/postgres/models"
	"neomatica/neosync/middleware"
	"neomatica/neosync/pkg/httpx"
	"neomatica/neosync/pkg/httpx/httperr"
	gtypes "neomatica/neosync/types"
	"net/http"

	"github.com/go-playground/validator"
)

/* Neosync HTTPx V1 */
/* Handler: создание шаблона конфигурации */

func (h *Handler) CreateConfigurationTemplateHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*gtypes.AuthToken)

	if !authToken.User.Has(constants.AccessConfigurationApply) {
		return httperr.Forbidden(tr.TErr("config-edit-access-restricted"))
	}

	var payload *ConfigurationTemplateCreate

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

	var cfgData []byte
	var err error

	switch payload.TypeOfSaving {
	case "modified_ones":
		changesData := ConfigurationTemplateChangesData{
			CfgHash: payload.CfgHash,
			Changes: payload.Changes,
		}

		cfgData, err = json.Marshal(changesData)
		if err != nil {
			logger.Error("CreateConfigurationTemplateHandler_V1 req={%s}: Failed marshal changes: %s", ctx.Value("XREQID").(string), err.Error())
			return httperr.InternalServerError(tr.TErr("failed-to-save-data"))
		}
	case "full":
		cfgParsed := getDefaultConfigurationAll()

		cfgData, err = buildConfigurationFromChanges(tr, payload.Changes, cfgParsed)
		if err != nil {
			return httperr.BadRequest(err.Error())
		}
	default:
		return httperr.BadRequest(tr.TErr("fields-not-filled"))
	}

	template := &models.ConfigurationTemplate{
		UserUUID:     authToken.User.UserContact.UserUUID,
		Name:         payload.Name,
		Model:        nil,
		TypeOfSaving: payload.TypeOfSaving,
	}

	id, err := h.Store.ConfigurationTemplates.Create_ConfigurationTemplate(ctx, template, cfgData)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	response := ConfigurationTemplateCreateResponse{
		ID:      id,
		Message: tr.T("configuration-template-created"),
	}

	httpx.HttpResponse(w, r, http.StatusCreated, response)
	return nil
}
