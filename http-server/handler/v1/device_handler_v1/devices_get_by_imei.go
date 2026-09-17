package device_handler_v1

import (
	"neomatica/neosync/middleware"
	"neomatica/neosync/pkg/httpx"
	"neomatica/neosync/pkg/httpx/httperr"
	"neomatica/neosync/pkg/permissions"
	"neomatica/neosync/types"
	"net/http"

	"github.com/gorilla/mux"
)

/* Neosync HTTPx V1 */
/* Handler: получение устройства у пользователя */

func (h *Handler) GetDeviceByImeiHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*types.AuthToken)

	imei := mux.Vars(r)["imei"]
	if imei == "" {
		return httperr.NotFound(tr.TErr("device-imei-not-found"))
	}

	if len(imei) != 15 {
		return httperr.BadRequest(tr.TErr("invalid-or-nonexistent-imei"))
	}

	device, err := h.Store.Devices.Get_DeviceAndConfByImei(ctx, imei)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	if device == nil {
		return httperr.NotFound(tr.TErr("device-not-found"))
	}

	if !permissions.CanAccessDevice(ctx, h.Store.GroupMembers, device, device.ID, &authToken.User, authToken.RoleCode, permissions.ShareViewDevice) {
		return httperr.Forbidden(tr.TErr("device-not-owned"))
	}

	w.Header().Set("Cache-Control", "no-store")
	httpx.HttpResponseWithETag(w, r, http.StatusOK, device)
	return nil
}
