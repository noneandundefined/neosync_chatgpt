package command_handler_v1

import (
	"neomatica/neosync/pkg/httpx"
	"neomatica/neosync/pkg/httpx/httperr"
	gtypes "neomatica/neosync/types"
	"net/http"
)

/* Neosync HTTPx V1 */
/* Handler: получение истории команд у пользователя, также получать по сессии ?session_id */

func (h *Handler) GetCommandHistoryHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	authToken := ctx.Value("identity").(*gtypes.AuthToken)

	history, err := h.Store.DeviceCommands.Get_DeviceCommandWithExecutionsByUuid(ctx, authToken.User.UserContact.UserUUID)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	httpx.HttpResponseWithETag(w, r, http.StatusOK, history)
	return nil
}
