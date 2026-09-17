package device_handler_v1

import (
	"neomatica/neosync/infra/constants"
	"neomatica/neosync/middleware"
	"neomatica/neosync/pkg/httpx"
	"neomatica/neosync/pkg/httpx/httperr"
	"neomatica/neosync/pkg/permissions"
	"neomatica/neosync/types"
	"net/http"

	"github.com/go-playground/validator"
)

/* Neosync HTTPx V1 */
/* Handler: удаление массива устройств по imei */

func (h *Handler) DeviceMassiveDeleteHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*types.AuthToken)

	/* Check access */
	if !authToken.User.Has(constants.AccessTrekerDelete) {
		return httperr.Forbidden(tr.TErr("access-treker-delete"))
	}

	var payload *DeleteDevicesPayload

	if err := httpx.HttpParse(r, &payload); err != nil {
		return httperr.BadRequest(err.Error())
	}

	if err := httpx.Validate.Struct(payload); err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			return httperr.BadRequest(httpx.ValidateMsg(tr, err))
		}

		return httperr.BadRequest(tr.TErr("fields-not-filled"))
	}

	if len(payload.IMEIs) == 0 {
		return httperr.BadRequest(tr.TErr("device-list-empty"))
	}

	var imeisToDelete []string
	var imeisToUnlink []string
	var imeisToDisassociate []string

	for _, imei := range payload.IMEIs {
		/* Device ownership verification */
		device, err := h.Store.Devices.Get_DeviceByImei(ctx, imei)
		if err != nil {
			return httperr.Db(ctx, err)
		}

		if device == nil {
			return httperr.NotFound(tr.TErr("device-not-found"))
		}

		if !permissions.IsMainRole(authToken.RoleCode) && !permissions.DeviceBelongsToUser(device, &authToken.User) {
			return httperr.Forbidden(tr.TErr("device-not-owned"))
		}

		switch {
		case permissions.IsMainRole(authToken.RoleCode):
			imeisToDelete = append(imeisToDelete, imei)
		case authToken.RoleCode == constants.Role_User:
			imeisToDisassociate = append(imeisToDisassociate, imei)
		default:
			imeisToUnlink = append(imeisToUnlink, imei)
		}
	}

	if len(imeisToDisassociate) > 0 {
		if err := h.Store.Devices.Update_DevicesUuidByImeiToNull(ctx, imeisToDisassociate); err != nil {
			return httperr.Db(ctx, err)
		}
	}

	if len(imeisToDelete) > 0 {
		if err := h.Store.Devices.Delete_DevicesByImei(ctx, imeisToDelete); err != nil {
			return httperr.Db(ctx, err)
		}
	}

	var userUuid string = authToken.User.UserContact.UserUUID
	if authToken.RoleCode == constants.Role_AdminL2Support {
		if authToken.User.ParentUUID != nil {
			userUuid = *authToken.User.ParentUUID
		}
	}

	if len(imeisToUnlink) > 0 {
		err := h.Store.Devices.Update_DevicesUnlinkFromAccountByImei(ctx, imeisToUnlink, userUuid)
		if err != nil {
			return httperr.Db(ctx, err)
		}
	}

	httpx.HttpResponse(w, r, http.StatusOK, tr.T("devices-deleted"))
	return nil
}
