package user_handler_v1

import (
	"neomatica/neosync/infra/constants"
	"neomatica/neosync/pkg/httpx"
	"neomatica/neosync/pkg/httpx/httperr"
	"neomatica/neosync/types"
	"neomatica/neosync/util"
	"net/http"
	"strings"
)

/* Neosync HTTPx V1 */
/* Handler: рекурсивное дерево пользователей по иерархии parent_uuid */
/* Параметры - ?search=поиск ?page=страница ?limit=размер ?all=true */

func (h *Handler) GetUsersTreeHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	authToken := ctx.Value("identity").(*types.AuthToken)

	viewerUUID := authToken.User.UserContact.UserUUID
	var dealerUUID *string

	switch authToken.RoleCode {
	case constants.Role_AdminL2Support, constants.Role_User:
		dealerUUID = authToken.User.ParentUUID
	default:
		dealerUUID = nil
	}

	search := strings.TrimSpace(r.URL.Query().Get("search"))
	all := r.URL.Query().Get("all") == "true"
	pagination := util.GetPagination(r, 25)

	if pagination.Limit > 100 {
		pagination.Limit = 100
		pagination.Offset = (pagination.Page - 1) * pagination.Limit
	}

	users, total, err := h.Store.Users.Get_UsersForTree(ctx, authToken.RoleCode, viewerUUID, dealerUUID, search, pagination.Limit, pagination.Offset, all)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	limit := pagination.Limit
	loaded := pagination.Offset + len(users)
	if all {
		limit = total
		loaded = total
	}

	httpx.HttpResponseWithETag(w, r, http.StatusOK, UsersTreePageResponse{
		Items:   users,
		Page:    pagination.Page,
		Limit:   limit,
		Total:   total,
		HasMore: loaded < total,
	})
	return nil
}
