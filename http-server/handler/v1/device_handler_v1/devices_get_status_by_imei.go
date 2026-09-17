package device_handler_v1

import (
	"neomatica/neosync/middleware"
	"neomatica/neosync/pkg/httpx"
	"neomatica/neosync/pkg/httpx/httperr"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
)

/* Neosync HTTPx V1 */
/* Handler: получение статуса устройства у пользователя */

func (h *Handler) GetDeviceStatusByImeiHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)

	imei := mux.Vars(r)["imei"]
	if imei == "" {
		return httperr.NotFound(tr.TErr("device-imei-not-found"))
	}

	if len(imei) != 15 {
		return httperr.BadRequest(tr.TErr("invalid-or-nonexistent-imei"))
	}

	deviceStatus, err := h.Store.Devices.Get_DeviceStatusByImei(ctx, imei)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	if deviceStatus == nil {
		return httperr.NotFound(tr.TErr("device-not-found"))
	}

	w.Header().Set("Cache-Control", "no-store")
	httpx.HttpResponseWithETag(w, r, http.StatusOK, deviceStatus)
	return nil
}

/* Neosync HTTPx V1 */
/* Handler: получение статуса устройств по IDs у пользователя */

func (h *Handler) GetDeviceDetailedStatusByImeisHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)

	var payload *DeviceDetailedStatusPayload

	if err := httpx.HttpParse(r, &payload); err != nil {
		return httperr.BadRequest(err.Error())
	}

	if err := httpx.Validate.Struct(payload); err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			return httperr.BadRequest(httpx.ValidateMsg(tr, err))
		}

		return httperr.BadRequest(tr.TErr("fields-not-filled"))
	}

	deviceStatus, err := h.Store.Devices.Get_DeviceDetailedStatusByImeis(ctx, payload.IMEIs)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	httpx.HttpResponseWithETag(w, r, http.StatusOK, deviceStatus)
	return nil
}
