package user_handler_v1

import (
	"math"
	"neomatica/neosync/infra/constants"
	"neomatica/neosync/middleware"
	"neomatica/neosync/pkg/httpx"
	"neomatica/neosync/pkg/httpx/httperr"
	"neomatica/neosync/pkg/permissions"
	"neomatica/neosync/types"
	"neomatica/neosync/util"
	"net/http"
)

/* Neosync HTTPx V1 */
/* Handler: получение пользователей */

func (h *Handler) GetUsersHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*types.AuthToken)

	pagination := util.GetPagination(r, 5)

	roleStr := r.URL.Query().Get("role")

	if permissions.IsMainRole(authToken.RoleCode) {
		users, total, totalAll, err := h.Store.Users.Get_UsersWithParams(ctx, pagination.Limit, pagination.Offset, pagination.Search, roleStr, pagination.ColumnSortKey, pagination.ColumnSortDir, authToken.User.UserContact.UserUUID)
		if err != nil {
			return httperr.Db(ctx, err)
		}

		totalPages := int(math.Ceil(float64(total) / float64(pagination.Limit)))
		usersResp := UsersWPResponse{
			Items:      users,
			Page:       pagination.Page,
			Limit:      pagination.Limit,
			Total:      total,
			TotalAll:   totalAll,
			TotalPages: totalPages,
		}

		httpx.HttpResponseWithETag(w, r, http.StatusOK, usersResp)
		return nil
	}

	_, ok := constants.RoleToSubordinateRole[authToken.RoleCode]
	if !ok {
		return httperr.Forbidden(tr.TErr("insufficient-permissions-get-users"))
	}

	dealerUuid := authToken.User.UserContact.UserUUID
	currentUserUuid := authToken.User.UserContact.UserUUID

	if authToken.RoleCode == constants.Role_AdminL2Support && authToken.User.ParentUUID != nil {
		dealerUuid = *authToken.User.ParentUUID
	}

	users, total, totalAll, err := h.Store.Users.Get_UsersByParentUuidWithParams(ctx, pagination.Limit, pagination.Offset, pagination.Search, roleStr, pagination.ColumnSortKey, pagination.ColumnSortDir, dealerUuid, currentUserUuid)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	totalPages := int(math.Ceil(float64(total) / float64(pagination.Limit)))
	usersResp := UsersWPResponse{
		Items:      users,
		Page:       pagination.Page,
		Limit:      pagination.Limit,
		Total:      total,
		TotalAll:   totalAll,
		TotalPages: totalPages,
	}

	httpx.HttpResponseWithETag(w, r, http.StatusOK, usersResp)
	return nil
}
