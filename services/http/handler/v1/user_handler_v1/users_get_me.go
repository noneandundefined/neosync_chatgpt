package user_handler_v1

import (
	"neomatica/neosync/middleware"
	"neomatica/neosync/pkg/httpx"
	"neomatica/neosync/pkg/httpx/httperr"
	"neomatica/neosync/types"
	"net/http"
)

/* Neosync HTTPx V1 */
/* Handler: получение профиля */

func (h *Handler) GetUserMeHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*types.AuthToken)

	user, err := h.Store.Users.Get_UserByUuid(ctx, authToken.User.UserContact.UserUUID)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	if user == nil {
		return httperr.NotFound(tr.TErr("user-not-found"))
	}

	// httpx.HttpCache(w, 300) // 5 min.
	httpx.HttpResponseWithETag(w, r, http.StatusOK, user)
	return nil
}
