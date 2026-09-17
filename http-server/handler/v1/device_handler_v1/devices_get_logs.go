package device_handler_v1

import (
	"neomatica/neosync/infra/constants"
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
/* Handler: получение логов устройства у пользователя */

func (h *Handler) GetDeviceLogsHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*types.AuthToken)

	/* Check access */
	if !authToken.User.Has(constants.AccessLogRead) {
		return httperr.Forbidden(tr.TErr("log-access-restricted"))
	}

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

	logs, err := redis.ReadAdmLog(imei)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	httpx.HttpResponseWithETag(w, r, http.StatusOK, logs)
	return nil
}
