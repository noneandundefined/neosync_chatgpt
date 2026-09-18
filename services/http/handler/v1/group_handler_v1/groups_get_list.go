package group_handler_v1

import (
	"math"
	"neomatica/neosync/pkg/httpx"
	"neomatica/neosync/pkg/httpx/httperr"
	"neomatica/neosync/types"
	"neomatica/neosync/util"
	"net/http"
)

/* Neosync HTTPx V1 */
/* Handler: получение групп у пользователя */
/* Параметры - ?page=страница ?limit=макс. значение элементов ?search=поиск по названию и описанию группы */

func (h *Handler) GetGroupsHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	authToken := ctx.Value("identity").(*types.AuthToken)

	pagination := util.GetPagination(r, 10)

	groups, total, totalAll, err := h.Store.Groups.Get_GroupsByUuidWithParams(ctx, pagination.Limit, pagination.Offset, pagination.Search, pagination.ColumnSortKey, pagination.ColumnSortDir, authToken.User.UserContact.UserUUID)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	totalPages := int(math.Ceil(float64(total) / float64(pagination.Limit)))
	groupsResp := GroupsWPResponse{
		Items:      groups,
		Page:       pagination.Page,
		Limit:      pagination.Limit,
		Total:      total,
		TotalAll:   totalAll,
		TotalPages: totalPages,
	}

	httpx.HttpResponseWithETag(w, r, http.StatusOK, groupsResp)
	return nil
}
