package configuration_handler_v1

import (
	"neomatica/neosync/infra/constants"
	"neomatica/neosync/infra/logger"
	"neomatica/neosync/infra/store/redis"
	"neomatica/neosync/middleware"
	"neomatica/neosync/pkg/httpx"
	"neomatica/neosync/pkg/httpx/httperr"
	"neomatica/neosync/pkg/permissions"
	gtypes "neomatica/neosync/types"
	"net/http"

	"github.com/gorilla/mux"
)

func (h *Handler) GetConfigurationTelemetryHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*gtypes.AuthToken)

	/* Check access */
	if !authToken.User.Has(constants.AccessConfigurationRead) {
		return nil
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

	if !permissions.CanAccessDevice(ctx, h.Store.GroupMembers, device, device.ID, &authToken.User, authToken.RoleCode, permissions.ShareReadConfig) {
		return httperr.Forbidden(tr.TErr("device-not-owned"))
	}

	if !*device.Activated {
		return httperr.New(tr.TErr("device-not-activated"), http.StatusUnprocessableEntity)
	}

	tl, err := redis.GetTelemetry(imei)
	if err != nil {
		return httperr.Redis(ctx, err)
	}

	httpx.HttpResponseWithETag(w, r, http.StatusOK, tl)
	return nil
}

func (h *Handler) RebootConfigurationTelemetryHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*gtypes.AuthToken)

	/* Check access */
	if !authToken.User.Has(constants.AccessConfigurationRead) {
		return nil
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

	if !permissions.CanAccessDevice(ctx, h.Store.GroupMembers, device, device.ID, &authToken.User, authToken.RoleCode, permissions.ShareReadConfig) {
		return httperr.Forbidden(tr.TErr("device-not-owned"))
	}

	if !*device.Activated {
		return httperr.New(tr.TErr("device-not-activated"), http.StatusUnprocessableEntity)
	}

	if !device.Status {
		return httperr.New(tr.TErr("device-not-connected-to-server"), http.StatusUnprocessableEntity)
	}

	if redis.TryAcquireTelemetryRefresh(imei) {
		go func() {
			commands := []string{
				constants.COM0_COMMAND,
				constants.ADM20INFO_COMMAND,
				constants.BLESENSORINFO_COMMAND,
				constants.FUELINFO_COMMAND,
			}

			for _, cmd := range commands {
				rabbitData := gtypes.RabbitMQ_TransitBinary{
					Type:  constants.ADM_RC_TYPE_STRING,
					Imei:  imei,
					NResp: true,
					Data:  []byte(cmd),
				}

				if err := h.RMQ.SendToRabbitAsync(ctx, rabbitData); err != nil {
					logger.Error("RebootConfigurationTelemetryHandler_V1 req={%s}: failed to send to RabbitMQ: %s", ctx.Value("XREQID").(string), err.Error())
					break
				}
			}
		}()
	}

	httpx.HttpResponseWithETag(w, r, http.StatusOK, tr.T("sensor-data-updated"))
	return nil
}
