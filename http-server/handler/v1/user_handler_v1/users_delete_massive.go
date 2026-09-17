package user_handler_v1

import (
	"neomatica/neosync/middleware"
	"neomatica/neosync/pkg/httpx"
	"neomatica/neosync/pkg/httpx/httperr"
	"neomatica/neosync/pkg/permissions"
	"neomatica/neosync/types"
	"net/http"

	"github.com/go-playground/validator"
)

/* Neosync HTTPx V1 */
/* Handler: удаление массива пользователей */

func (h *Handler) UserMassiveDeleteHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*types.AuthToken)

	var payload *DeleteUsersPayload

	if err := httpx.HttpParse(r, &payload); err != nil {
		return httperr.BadRequest(err.Error())
	}

	if err := httpx.Validate.Struct(payload); err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			return httperr.BadRequest(httpx.ValidateMsg(tr, err))
		}

		return httperr.BadRequest(tr.TErr("fields-not-filled"))
	}

	if len(payload.UUIDs) == 0 {
		return httperr.BadRequest(tr.TErr("user-list-empty"))
	}

	if !permissions.IsMainRole(authToken.RoleCode) {
		for _, uuid := range payload.UUIDs {
			user, err := h.Store.Users.Get_UserCoreByUuid(ctx, uuid)
			if err != nil {
				return httperr.Db(ctx, err)
			}

			if user == nil {
				return httperr.BadRequest(tr.TErr("user-not-found"))
			}

			if user.ParentUUID == nil || *user.ParentUUID != authToken.User.UserContact.UserUUID {
				return httperr.Forbidden(tr.TErr("user-not-owned"))
			}
		}
	}

	if err := h.Store.Users.Delete_UsersByUuid(ctx, payload.UUIDs); err != nil {
		return httperr.Db(ctx, err)
	}

	httpx.HttpResponse(w, r, http.StatusOK, tr.T("users-deleted"))
	return nil
}
