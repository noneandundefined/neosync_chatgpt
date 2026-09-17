package user_handler_v1

import (
	"neomatica/neosync/pkg/httpx"
	"neomatica/neosync/pkg/httpx/httperr"
	"net/http"
)

/* Neosync HTTPx V1 */
/* Handler: получение ВСЕХ пользователей */

func (h *Handler) GetUsersFullHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	users, err := h.Store.Users.Get_Users(ctx)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	httpx.HttpResponseWithETag(w, r, http.StatusOK, users)
	return nil
}
