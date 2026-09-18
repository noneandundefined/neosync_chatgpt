package configuration_handler_v1

import (
	"fmt"
	"neomatica/neosync/infra/constants"
	"neomatica/neosync/infra/store/redis"
	"neomatica/neosync/middleware"
	"neomatica/neosync/pkg/httpx/httperr"
	"neomatica/neosync/pkg/permissions"
	"neomatica/neosync/types"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

/* Neosync HTTPx V1 */
/* Handler: получение статуса изменной конфигурации устройства */
/* Используется SSE подключение */

func (h *Handler) GetConfigurationDraftCheckHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*types.AuthToken)

	/* Check access */
	if !authToken.User.Has(constants.AccessConfigurationRead) {
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

	if !permissions.CanAccessDevice(ctx, h.Store.GroupMembers, device, device.ID, &authToken.User, authToken.RoleCode, permissions.ShareReadConfig) {
		return httperr.Forbidden(tr.TErr("device-not-owned"))
	}

	if !*device.Activated {
		return httperr.New(tr.TErr("device-not-activated"), http.StatusUnprocessableEntity)
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	flusher, ok := w.(http.Flusher)
	if !ok {
		return httperr.Conflict("streaming unsupported")
	}

	/* SSE event. connected */
	fmt.Fprintf(w, "event: open\ndata: ok\n\n")
	flusher.Flush()

	var lastExists bool

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		exists, _ := redis.HasDraft(authToken.User.UserContact.UserUUID, imei)

		if exists != lastExists {
			if exists {
				fmt.Fprintf(w, "data: true\n\n")
			} else {
				fmt.Fprintf(w, "data: false\n\n")
			}
			flusher.Flush()
			lastExists = exists
		}

		select {
		case <-ctx.Done():
			return nil
		case <-time.After(1 * time.Second):
		}
	}
}
