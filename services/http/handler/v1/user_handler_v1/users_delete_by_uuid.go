package user_handler_v1

import (
	"neomatica/neosync/middleware"
	"neomatica/neosync/pkg/httpx"
	"neomatica/neosync/pkg/httpx/httperr"
	"neomatica/neosync/pkg/permissions"
	"neomatica/neosync/types"
	"net/http"

	"github.com/gorilla/mux"
)

/* Neosync HTTPx V1 */
/* Handler: удаление пользователя */

func (h *Handler) UserDeleteByUuidHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*types.AuthToken)

	uuid := mux.Vars(r)["uuid"]
	if uuid == "" {
		return httperr.NotFound(tr.TErr("user-not-found"))
	}

	user, err := h.Store.Users.Get_UserCoreByUuid(ctx, uuid)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	if user == nil {
		return httperr.NotFound(tr.TErr("user-not-found"))
	}

	/* if this is adminL2 and user not created current adminL2 -> forbidden */
	if !permissions.IsMainRole(authToken.RoleCode) {
		if user.ParentUUID == nil || *user.ParentUUID != authToken.User.UserContact.UserUUID {
			return httperr.Forbidden(tr.TErr("user-not-owned"))
		}
	}

	if err := h.Store.Users.Delete_UserByUuid(ctx, uuid); err != nil {
		return httperr.Db(ctx, err)
	}

	httpx.HttpResponse(w, r, http.StatusOK, tr.T("user-deleted-successfully"))
	return nil
}
