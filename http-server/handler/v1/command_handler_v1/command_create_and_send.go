package command_handler_v1

import (
	"context"
	"fmt"
	"neomatica/neosync/infra/constants"
	"neomatica/neosync/infra/logger"
	"neomatica/neosync/infra/store/postgres/models"
	"neomatica/neosync/middleware"
	"neomatica/neosync/pkg/httpx"
	"neomatica/neosync/pkg/httpx/httperr"
	"neomatica/neosync/pkg/permissions"
	gtypes "neomatica/neosync/types"
	"net/http"
	"strings"

	"github.com/go-playground/validator"
)

/* Neosync HTTPx V1 */
/* Handler: отправка команды на терминал и в БД */

func (h *Handler) CommandCreateAndSendHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*gtypes.AuthToken)

	/* Check access */
	if !authToken.User.Has(constants.AccessCommandSend) {
		return httperr.Forbidden(tr.TErr("commands-access-restricted"))
	}

	var payload *CommandPayload

	if err := httpx.HttpParse(r, &payload); err != nil {
		return httperr.BadRequest(err.Error())
	}

	if err := httpx.Validate.Struct(payload); err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			return httperr.BadRequest(httpx.ValidateMsg(tr, err))
		}

		return httperr.BadRequest(tr.TErr("fields-not-filled"))
	}

	var devices []*models.Device

	for _, imei := range payload.Imeis {
		/* Device ownership verification */
		device, err := h.Store.Devices.Get_DeviceByImei(ctx, imei)
		if err != nil {
			return httperr.Db(ctx, err)
		}

		if device == nil {
			return httperr.NotFound(fmt.Sprintf("%s: %s", tr.TErr("device-not-found"), imei))
		}

		if !permissions.IsMainRole(authToken.RoleCode) && !authToken.User.CanViewChildGroups && !permissions.CanAccessDevice(ctx, h.Store.GroupMembers, device, device.ID, &authToken.User, authToken.RoleCode, permissions.ShareSendCommand) {
			return httperr.Forbidden(fmt.Sprintf("%s: %s", tr.TErr("device-not-owned"), imei))
		}

		if !*device.Activated {
			return httperr.New(fmt.Sprintf("%s: %s", tr.TErr("device-not-activated"), imei), http.StatusUnprocessableEntity)
		}

		if payload.SendMode == constants.CMD_SEND_MODE_INSTANT && !device.Status {
			return httperr.New(tr.TErr("device-not-connected-to-server"), http.StatusUnprocessableEntity)
		}

		devices = append(devices, device)
	}

	payload.Command = strings.TrimSpace(payload.Command)
	if payload.Command == "" {
		return httperr.BadRequest(tr.TErr("fields-not-filled"))
	}

	tx, err := h.Db.BeginTx(ctx, nil)
	if err != nil {
		return httperr.Db(ctx, httperr.Err_DbNetwork)
	}

	defer func() {
		_ = tx.Rollback()
	}()

	command := &models.DeviceCommand{
		SessionID: payload.SessId,
		UserUUID:  authToken.User.UserContact.UserUUID,
		SendMode:  payload.SendMode,
		Command:   payload.Command,
	}

	commandId, err := h.Store.DeviceCommands.Create_DeviceCommand(ctx, tx, command)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	if err := h.Store.DeviceCommands.Create_Batch_DeviceCommandExecution(ctx, tx, commandId, payload.Imeis); err != nil {
		return httperr.Db(ctx, err)
	}

	if err := tx.Commit(); err != nil {
		return httperr.Conflict(tr.TErr("failed-to-save-data"))
	}

	if constants.IsFirmwareUpdateCommand(payload.Command) {
		for _, device := range devices {
			if err := h.Store.Syncs.Update_FirmwareUpdateStatusByDeviceId(ctx, device.ID, constants.FW_UPDATE_STATUS_PENDING); err != nil {
				logger.Error("CommandCreateAndSendHandler_V1 req={%s}: Failed to mark firmware update pending imei={%s}: %s", ctx.Value("XREQID").(string), device.IMEI, err.Error())
			}
		}
	}

	/* If mode ON_CONNECT return success response */
	if payload.SendMode == constants.CMD_SEND_MODE_ON_CONNECT {
		httpx.HttpResponse(w, r, http.StatusNoContent, nil)
		return nil
	}

	go func(imei []string) {
		for _, imei := range imei {
			rabbitData := gtypes.RabbitMQ_TransitBinary{
				Type: constants.ADM_RC_TYPE_STRING,
				Imei: imei,
				Data: []byte(payload.Command),
			}

			if err := h.RMQ.SendToRabbitAsync(context.Background(), rabbitData); err != nil {
				logger.Error("CommandCreateAndSendHandler_V1 req={%s}: failed to send message to RabbitMQ: %s", ctx.Value("XREQID").(string), err.Error())
			}
		}
	}(payload.Imeis)

	httpx.HttpResponse(w, r, http.StatusNoContent, nil)
	return nil
}
