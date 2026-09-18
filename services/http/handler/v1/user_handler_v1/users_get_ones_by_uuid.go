package user_handler_v1

import (
	"neomatica/neosync/encryption"
	"neomatica/neosync/infra/constants"
	"neomatica/neosync/middleware"
	"neomatica/neosync/pkg/httpx"
	"neomatica/neosync/pkg/httpx/httperr"
	"neomatica/neosync/pkg/permissions"
	"neomatica/neosync/types"
	"net/http"

	"github.com/gorilla/mux"
)

/* Neosync HTTPx V1 */
/* Handler: получения данных пользователя по uuid */

func (h *Handler) UserGetOnesByUuidHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*types.AuthToken)

	uuid := mux.Vars(r)["uuid"]
	if uuid == "" {
		return httperr.NotFound(tr.TErr("user-not-found"))
	}

	user, err := h.Store.Users.Get_UserByUuid(ctx, uuid)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	if user == nil {
		return httperr.NotFound(tr.TErr("user-not-found"))
	}

	var parentUuid string = authToken.User.UserContact.UserUUID
	if authToken.RoleCode == constants.Role_AdminL2Support {
		if authToken.User.ParentUUID != nil {
			parentUuid = *authToken.User.ParentUUID
		}
	}

	/* if this is adminL2 and user not created current adminL2 -> forbidden */
	if !permissions.IsMainRole(authToken.RoleCode) {
		if user.ParentUUID == nil || *user.ParentUUID != parentUuid {
			return httperr.Forbidden(tr.TErr("user-not-owned"))
		}
	}

	pass, err := encryption.Decrypt(tr, user.Password)
	if err != nil {
		return httperr.InternalServerError(err.Error())
	}
	user.Password = pass

	httpx.HttpCache(w, 300) // 5 min.
	httpx.HttpResponseWithETag(w, r, http.StatusOK, user)
	return nil
}
