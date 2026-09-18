package configuration_handler_v1

import (
	"neomatica/neosync/infra/constants"
	"neomatica/neosync/infra/logger"
	"neomatica/neosync/infra/store/redis"
	"neomatica/neosync/middleware"
	"neomatica/neosync/pkg/httpx"
	"neomatica/neosync/pkg/httpx/httperr"
	"neomatica/neosync/pkg/permissions"
	"neomatica/neosync/types"
	"net/http"

	"github.com/gorilla/mux"
)

func getSchemaForChanges(changes map[string]any) []constants.FieldSchema {
	var schemaList []constants.FieldSchema

	for _, schema := range constants.CfgSchema {
		if _, ok := changes[schema.Name]; ok {
			schemaList = append(schemaList, schema)
		}
	}

	return schemaList
}

/* Neosync HTTPx V1 */
/* Handler: получение черновика конфигурации устройства */

func (h *Handler) GetConfigurationDraftHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*types.AuthToken)

	/* Check access */
	if !authToken.User.Has(constants.AccessConfigurationRead) {
		return httperr.Forbidden(tr.TErr("config-access-restricted"))
	}

	imei := mux.Vars(r)["imei"]
	section := mux.Vars(r)["section"]

	if imei == "" {
		return httperr.BadRequest(tr.TErr("device-imei-not-found"))
	}

	if section == "" {
		return httperr.BadRequest(tr.TErr("device-section-not-found"))
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

	draft, err := redis.GetDraft(authToken.User.UserContact.UserUUID, imei, section)
	if err != nil {
		logger.Error("GetConfigurationDraftHandler_V1 req={%s} imei={%s} section={%s}: Failed get cache Redis draft: %s", ctx.Value("XREQID").(string), imei, section, err.Error())
		return httperr.InternalServerError(tr.TErr("failed-to-get-draft"))
	}

	if draft == nil {
		httpx.HttpResponse(w, r, http.StatusOK, nil)
		return nil
	}

	draftResp := ConfigurationDraftGet{
		DeviceImei: draft.DeviceImei,
		CfgHash:    draft.CfgHash,
		Section:    draft.Section,
		Changes:    draft.Changes,
		Schema:     getSchemaForChanges(draft.Changes),
		Timestamp:  draft.Timestamp,
	}

	httpx.HttpResponseWithETag(w, r, http.StatusOK, draftResp)
	return nil
}
