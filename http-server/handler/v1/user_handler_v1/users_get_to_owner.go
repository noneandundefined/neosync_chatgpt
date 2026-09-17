package user_handler_v1

import (
	"database/sql"
	"neomatica/neosync/infra/constants"
	"neomatica/neosync/pkg/httpx"
	"neomatica/neosync/pkg/httpx/httperr"
	"neomatica/neosync/types"
	"neomatica/neosync/util"
	"net/http"
	"strings"
)

/* Neosync HTTPx V1 */
/* Handler: получения пользователей для трансфера терминалов */
/* Параметры - ?search=поиск ?page=страница ?limit=размер ?all=true */

func (h *Handler) GetUsersToOwnerHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	authToken := ctx.Value("identity").(*types.AuthToken)

	if authToken.RoleCode == constants.Role_User {
		httpx.HttpResponse(w, r, http.StatusNoContent, nil)
		return nil
	}

	search := strings.TrimSpace(r.URL.Query().Get("search"))
	all := r.URL.Query().Get("all") == "true"

	pagination := util.GetPagination(r, 25)

	if pagination.Limit > 100 {
		pagination.Limit = 100
		pagination.Offset = (pagination.Page - 1) * pagination.Limit
	}

	var role string = constants.Role_AdminL2
	var parentUuid sql.NullString = sql.NullString{
		Valid: false,
	}

	switch authToken.RoleCode {
	case constants.Role_AdminL2:
		role = constants.Role_User
		parentUuid = sql.NullString{
			String: authToken.User.UserContact.UserUUID,
			Valid:  true,
		}

	case constants.Role_AdminL2Support:
		role = constants.Role_User
		parentUuid = sql.NullString{
			String: *authToken.User.ParentUUID,
			Valid:  true,
		}
	}

	users, total, err := h.Store.Users.Get_UsersToOwner(ctx, role, authToken.User.UserContact.UserUUID, parentUuid, search, pagination.Limit, pagination.Offset, all)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	limit := pagination.Limit
	loaded := pagination.Offset + len(users)
	if all {
		limit = total
		loaded = total
	}

	httpx.HttpResponseWithETag(w, r, http.StatusOK, UsersToOwnerPageResponse{
		Items:   users,
		Page:    pagination.Page,
		Limit:   limit,
		Total:   total,
		HasMore: loaded < total,
	})
	return nil
}
