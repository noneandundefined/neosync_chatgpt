package user_handler_v1

import (
	"neomatica/neosync/middleware"
	"neomatica/neosync/pkg/httpx"
	"neomatica/neosync/pkg/httpx/httperr"
	"neomatica/neosync/types"
	"net/http"
)

/* Neosync HTTPx V1 */
/* Handler: получение минимальной информации о пользователе */

func (h *Handler) GetUserMeLoginWithRoleCodeHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*types.AuthToken)

	user, err := h.Store.Users.Get_UserLoginWithRoleCodeByUuid(ctx, authToken.User.UserContact.UserUUID)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	if user == nil {
		return httperr.NotFound(tr.TErr("user-not-found"))
	}

	httpx.HttpResponseWithETag(w, r, http.StatusOK, user)
	return nil
}
