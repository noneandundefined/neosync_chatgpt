package firmware_handler_v1

import (
	"context"
	"errors"
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

/* Neosync HTTPx V1 */
/* Handler: обновление устройства command:UPDATE */

func (h *Handler) FirmwareUpdateFirmwareVersionByImeiHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*gtypes.AuthToken)

	imei := mux.Vars(r)["imei"]
	if imei == "" {
		return httperr.NotFound(tr.TErr("device-not-found"))
	}

	/* Device ownership verification */
	device, err := h.Store.Devices.Get_DeviceByImei(ctx, imei)
	if err != nil || device == nil {
		return httperr.NotFound(tr.TErr("device-not-found"))
	}

	if !permissions.CanAccessDevice(ctx, h.Store.GroupMembers, device, device.ID, &authToken.User, authToken.RoleCode, permissions.ShareSendCommand) {
		return httperr.Forbidden(tr.TErr("device-not-owned"))
	}

	if !*device.Activated {
		return httperr.New(tr.TErr("device-not-activated"), http.StatusUnprocessableEntity)
	}

	/* Redis cache UPDATE */
	exists, err := redis.UpdDeviceExists(imei)
	if err != nil {
		logger.Error("FirmwareUpdateFirmwareVersionByImeiHandler_V1 req={%s}: %s", ctx.Value("XREQID").(string), err.Error())
		return httperr.InternalServerError(tr.TErr("command-cache-check-error"))
	}

	if !exists {
		if err := redis.SetUpdDevice(imei); err != nil {
			logger.Error("FirmwareUpdateFirmwareVersionByImeiHandler_V1 req={%s}: %s", ctx.Value("XREQID").(string), err.Error())
			return httperr.InternalServerError(tr.TErr("command-cache-fetch-error"))
		}
	}

	if err := h.Store.Syncs.Update_FirmwareUpdateStatusByDeviceId(ctx, device.ID, constants.FW_UPDATE_STATUS_PENDING); err != nil {
		logger.Error("FirmwareUpdateFirmwareVersionByImeiHandler_V1 req={%s}: Failed to mark firmware update pending: %s", ctx.Value("XREQID").(string), err.Error())
		return httperr.Db(ctx, err)
	}

	/* RabbitMQ */
	rabbitData := gtypes.RabbitMQ_TransitBinary{
		Type: constants.ADM_RC_TYPE_STRING,
		Imei: imei,
		Data: []byte(constants.UPDATE_COMMAND),
	}

	response, err := h.RMQ.SendToRabbitAndWait(ctx, rabbitData, h.Session)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return httperr.New(tr.TErr("response-timeout-13s"), http.StatusBadGateway)
		}

		return httperr.New(tr.TErr("rabbitmq-publish-failed"), http.StatusConflict)
	}

	if response == nil {
		return httperr.Conflict(tr.TErr("response-empty"))
	}

	if response.StatCode == http.StatusBadGateway {
		logger.Error("FirmwareUpdateFirmwareVersionByImeiHandler_V1 req={%s}: %s", ctx.Value("XREQID").(string), err.Error())
		return httperr.New(string(response.Data), http.StatusBadGateway)
	}

	httpx.HttpResponse(w, r, http.StatusAccepted, tr.T("firmware-uploaded-waiting"))
	return nil
}
