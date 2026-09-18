package firmware_handler_v1

import (
	"neomatica/neosync/pkg/httpx"
	"neomatica/neosync/pkg/httpx/httperr"
	"net/http"
)

/* Neosync HTTPx V1 */
/* Handler: получение всех моделей на сервисе */

func (h *Handler) GetFirmwareModelsHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	models, err := h.Store.DeviceModelSources.Get_Models(ctx)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	httpx.HttpCache(w, 43200) // 12h.
	httpx.HttpResponse(w, r, http.StatusOK, models)
	return nil
}
