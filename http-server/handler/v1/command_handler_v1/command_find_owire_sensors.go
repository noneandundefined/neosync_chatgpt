package command_handler_v1

import (
	"neomatica/neosync/infra/constants"
	"neomatica/neosync/infra/logger"
	"neomatica/neosync/middleware"
	"neomatica/neosync/pkg/httpx"
	"neomatica/neosync/pkg/httpx/httperr"
	"neomatica/neosync/pkg/permissions"
	gtypes "neomatica/neosync/types"
	"net/http"

	"github.com/gorilla/mux"
)

func (h *Handler) CommandOWireHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*gtypes.AuthToken)

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

	if !permissions.CanAccessDevice(ctx, h.Store.GroupMembers, device, device.ID, &authToken.User, authToken.RoleCode, permissions.ShareSendCommand) {
		return httperr.Forbidden(tr.TErr("device-not-owned"))
	}

	if !*device.Activated {
		return httperr.New(tr.TErr("device-not-activated"), http.StatusUnprocessableEntity)
	}

	rabbitData := gtypes.RabbitMQ_TransitBinary{
		Type:  constants.ADM_RC_TYPE_STRING,
		Imei:  imei,
		NResp: true,
		Data:  []byte(constants.ONEWIRE_FIND_COMMAND),
	}

	if err := h.RMQ.SendToRabbitAsync(ctx, rabbitData); err != nil {
		logger.Error("CommandOWireHandler_V1 req={%s}: failed to send to RabbitMQ: %s", ctx.Value("XREQID").(string), err.Error())
		return httperr.Conflict(tr.TErr("rabbitmq-publish-failed"))
	}

	httpx.HttpResponse(w, r, http.StatusOK, tr.T("device-onewire-wait"))
	return nil
}
