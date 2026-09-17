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

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
)

/* Neosync HTTPx V1 */
/* Handler: сохранение черновика конфигурации */

func (h *Handler) SetConfigurationDraftHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*types.AuthToken)

	/* Check access */
	if !authToken.User.Has(constants.AccessConfigurationApply) {
		return httperr.Forbidden(tr.TErr("config-edit-access-restricted"))
	}

	imei := mux.Vars(r)["imei"]
	section := mux.Vars(r)["section"]

	if imei == "" {
		return httperr.BadRequest(tr.TErr("device-imei-not-found"))
	}

	if section == "" {
		return httperr.BadRequest(tr.TErr("device-section-not-found"))
	}

	var payload *ConfigurationDraftInsert

	if err := httpx.HttpParse(r, &payload); err != nil {
		return httperr.BadRequest(err.Error())
	}

	if err := httpx.Validate.Struct(payload); err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			return httperr.BadRequest(httpx.ValidateMsg(tr, err))
		}

		return httperr.BadRequest(tr.TErr("fields-not-filled"))
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

	if len(payload.Changes) == 0 {
		return httperr.BadRequest(tr.TErr("no-changes-provided"))
	}

	changes := redis.ConfigurationDraftInsert{
		DeviceImei: payload.DeviceImei,
		CfgHash:    payload.CfgHash,
		Section:    section,
		Changes:    payload.Changes,
		Timestamp:  payload.Timestamp,
	}

	if err := redis.SaveOrUpdateDraft(authToken.User.UserContact.UserUUID, imei, section, changes); err != nil {
		logger.Error("SetConfigurationDraftHandler_V1 req={%s} imei={%s} section={%s}: Failed save/update draft: %s", ctx.Value("XREQID").(string), imei, section, err.Error())
		return httperr.InternalServerError(tr.TErr("failed-to-save-draft"))
	}

	httpx.HttpResponse(w, r, http.StatusNoContent, nil)
	return nil
}
