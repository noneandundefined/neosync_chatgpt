package device_handler_v1

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

/* Neosync HTTPx V1 */
/* Handler: удаление устройства по imei */

func (h *Handler) DeleteDeviceHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*types.AuthToken)

	/* Check access */
	if !authToken.User.Has(constants.AccessTrekerDelete) {
		return httperr.Forbidden(tr.TErr("access-treker-delete"))
	}

	imei := mux.Vars(r)["imei"]
	if imei == "" {
		return httperr.NotFound(tr.TErr("device-imei-not-found"))
	}

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

	/* If ROLE is SUPERADMIN/SUPPORT — hard delete from the system */
	if permissions.IsMainRole(authToken.RoleCode) {
		if err := redis.ClearCacheConfiguration(imei); err != nil {
			logger.Error("DeleteDeviceHandler_V1 req={%s}: Failed clear cache cfg: %s", ctx.Value("XREQID").(string), err.Error())
		}

		if err := h.Store.Devices.Delete_DeviceByImei(ctx, imei); err != nil {
			return httperr.Db(ctx, err)
		}

		httpx.HttpResponse(w, r, http.StatusOK, tr.T("device-deleted-successfully"))
		return nil
	}

	/* If ROLE is USER — unlink from own account, stay on dealer */
	if authToken.RoleCode == constants.Role_User {
		if err := h.Store.Devices.Update_DeviceOwnerUuidByImei(ctx, imei, nil); err != nil {
			return httperr.Db(ctx, err)
		}

		httpx.HttpResponse(w, r, http.StatusOK, tr.T("device-deleted-successfully"))
		return nil
	}

	var userUuid string = authToken.User.UserContact.UserUUID
	if authToken.RoleCode == constants.Role_AdminL2Support {
		if authToken.User.ParentUUID != nil {
			userUuid = *authToken.User.ParentUUID
		}
	}

	/* If ROLE is DEALER/DEALER_SUPPORT — unlink from account, keep in the system */
	if err := h.Store.Devices.Update_DeviceUnlinkFromAccountByImei(ctx, imei, userUuid); err != nil {
		return httperr.Db(ctx, err)
	}

	httpx.HttpResponse(w, r, http.StatusOK, tr.T("device-deleted-successfully"))
	return nil
}
