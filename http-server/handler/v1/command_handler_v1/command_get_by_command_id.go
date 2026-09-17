package command_handler_v1

import (
	"neomatica/neosync/infra/logger"
	"neomatica/neosync/middleware"
	"neomatica/neosync/pkg/httpx"
	"neomatica/neosync/pkg/httpx/httperr"
	gtypes "neomatica/neosync/types"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

/* Neosync HTTPx V1 */
/* Handler: получении подробной информации команды */

func (h *Handler) GetCommandByCommandIdHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*gtypes.AuthToken)

	commandIdStr := mux.Vars(r)["id"]
	if commandIdStr == "" {
		return httperr.NotFound(tr.TErr("invalid-command-id"))
	}

	commandId, err := strconv.Atoi(commandIdStr)
	if err != nil {
		logger.Error("GetCommandByCommandIdHandler_V1 req={%s}: %s", ctx.Value("XREQID").(string), err.Error())
		return httperr.BadRequest(tr.TErr("invalid-command-id"))
	}

	history, err := h.Store.DeviceCommands.Get_DeviceCommandWithExecutionsById(ctx, authToken.User.UserContact.UserUUID, uint64(commandId))
	if err != nil {
		return httperr.Db(ctx, err)
	}

	httpx.HttpResponseWithETag(w, r, http.StatusOK, history)
	return nil
}
