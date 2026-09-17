package device_handler_v1

import (
	"neomatica/neosync/pkg/httpx"
	"neomatica/neosync/pkg/httpx/httperr"
	"net/http"
)

/* Neosync HTTPx V1 */
/* Handler: получение ВСЕХ устройств */

func (h *Handler) GetDevicesFullHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	devices, err := h.Store.Devices.Get_Devices(ctx)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	httpx.HttpResponseWithETag(w, r, http.StatusOK, devices)
	return nil
}
