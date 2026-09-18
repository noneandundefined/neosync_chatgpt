package command_handler_v1

import (
	"neomatica/neosync/infra/constants"
	"neomatica/neosync/infra/logger"
	"neomatica/neosync/middleware"
	"neomatica/neosync/pkg/httpx"
	"neomatica/neosync/pkg/httpx/httperr"
	gtypes "neomatica/neosync/types"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (h *Handler) CommandCancelHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*gtypes.AuthToken)

	/* Check access */
	if !authToken.User.Has(constants.AccessCommandSend) {
		return httperr.Forbidden(tr.TErr("commands-access-restricted"))
	}

	commandIdStr := mux.Vars(r)["id"]
	if commandIdStr == "" {
		return httperr.NotFound(tr.TErr("invalid-command-id"))
	}

	commandId, err := strconv.Atoi(commandIdStr)
	if err != nil {
		logger.Error("CommandCancelHandler_V1 req={%s}: %s", ctx.Value("XREQID").(string), err.Error())
		return httperr.BadRequest(tr.TErr("invalid-command-id"))
	}

	if err := h.Store.DeviceCommands.Update_DeviceCommandCancelByCommandId(ctx, authToken.User.UserContact.UserUUID, uint64(commandId)); err != nil {
		return httperr.Db(ctx, err)
	}

	httpx.HttpResponse(w, r, http.StatusOK, tr.T("command-cancel-successfully"))
	return nil
}
