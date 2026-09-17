package command_handler_v1

import (
	"context"
	"errors"
	"neomatica/neosync/infra/constants"
	"neomatica/neosync/infra/logger"
	"neomatica/neosync/middleware"
	"neomatica/neosync/pkg/httpx"
	"neomatica/neosync/pkg/httpx/httperr"
	"neomatica/neosync/pkg/permissions"
	gtypes "neomatica/neosync/types"
	"net/http"
	"strings"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
)

/* Neosync HTTPx V1 */
/* Handler: отправка команды на терминал */

func (h *Handler) CommandSendAndWaitHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*gtypes.AuthToken)

	imei := mux.Vars(r)["imei"]
	if imei == "" {
		return httperr.NotFound(tr.TErr("device-imei-not-found"))
	}

	/* Check access */
	if !authToken.User.Has(constants.AccessCommandSend) {
		return httperr.Forbidden(tr.TErr("commands-access-restricted"))
	}

	var payload *CommandWaitPayload

	if err := httpx.HttpParse(r, &payload); err != nil {
		return httperr.BadRequest(err.Error())
	}

	if err := httpx.Validate.Struct(payload); err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			return httperr.BadRequest(httpx.ValidateMsg(tr, err))
		}

		return httperr.BadRequest(tr.TErr("fields-not-filled"))
	}

	/* Device ownership verification */
	device, err := h.Store.Devices.Get_DeviceByImei(ctx, imei)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	if device == nil {
		return httperr.NotFound(tr.TErr("device-not-found"))
	}

	if !permissions.CanAccessDevice(ctx, h.Store.GroupMembers, device, device.ID, &authToken.User, authToken.RoleCode, permissions.ShareSendCommand) {
		return httperr.Forbidden(tr.TErr("device-not-owned"))
	}

	if !*device.Activated {
		return httperr.New(tr.TErr("device-not-activated"), http.StatusUnprocessableEntity)
	}

	payload.Command = strings.TrimSpace(payload.Command)
	if payload.Command == "" {
		return httperr.BadRequest(tr.TErr("fields-not-filled"))
	}

	if constants.IsFirmwareUpdateCommand(payload.Command) {
		if err := h.Store.Syncs.Update_FirmwareUpdateStatusByDeviceId(ctx, device.ID, constants.FW_UPDATE_STATUS_PENDING); err != nil {
			logger.Error("CommandSendAndWaitHandler_V1 req={%s}: Failed to mark firmware update pending: %s", ctx.Value("XREQID").(string), err.Error())
		}
	}

	rabbitData := gtypes.RabbitMQ_TransitBinary{
		Type:  constants.ADM_RC_TYPE_STRING,
		Imei:  imei,
		NResp: true,
		Data:  []byte(payload.Command),
	}

	rabbitmqResp, err := h.RMQ.SendToRabbitAndWait(ctx, rabbitData, h.Session)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			logger.Warning("CommandCreateAndSendHandler_V1 timeout req={%s}: device did not respond in %f sec", ctx.Value("XREQID").(string), constants.RabbitMQ_TimeoutRead.Seconds())
			return httperr.New(tr.TErr("response-timeout-13s"), http.StatusGatewayTimeout)
		}

		logger.Error("CommandCreateAndSendHandler_V1 req={%s}: failed to send message to RabbitMQ: %s", ctx.Value("XREQID").(string), err.Error())
		return httperr.Conflict(tr.TErr("rabbitmq-publish-failed"))
	}

	httpx.HttpResponseWithETag(w, r, http.StatusOK, string(rabbitmqResp.Data))
	return nil
}
