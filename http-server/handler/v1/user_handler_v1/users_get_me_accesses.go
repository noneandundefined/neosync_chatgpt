package user_handler_v1

import (
	"neomatica/neosync/middleware"
	"neomatica/neosync/pkg/httpx"
	"neomatica/neosync/pkg/httpx/httperr"
	"neomatica/neosync/types"
	"net/http"
)

/* Neosync HTTPx V1 */
/* Handler: получение доступов пользователя */

func (h *Handler) GetUserMeAccessesHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*types.AuthToken)

	userAccesses, err := h.Store.Users.Get_UserAccessByUuid(ctx, authToken.User.UserContact.UserUUID)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	if userAccesses == nil {
		return httperr.NotFound(tr.TErr("user-not-found"))
	}

	httpx.HttpResponseWithETag(w, r, http.StatusOK, userAccesses)
	return nil
}
