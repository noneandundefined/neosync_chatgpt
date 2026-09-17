package company_handler_v1

import (
	"math"
	"neomatica/neosync/pkg/httpx"
	"neomatica/neosync/pkg/httpx/httperr"
	"neomatica/neosync/types"
	"neomatica/neosync/util"
	"net/http"
)

/* Neosync HTTPx V1 */
/* Handler: получение массовых настроек пользователя */

func (h *Handler) GetCompaniesHandlerV1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	authToken := ctx.Value("identity").(*types.AuthToken)

	pagination := util.GetPagination(r, 10)

	companies, total, totalAll, err := h.Store.Companies.Get_CompaniesByUserUuid(
		ctx,
		pagination.Limit,
		pagination.Offset,
		pagination.Search,
		pagination.ColumnSortKey,
		pagination.ColumnSortDir,
		authToken.User.UserContact.UserUUID,
	)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	totalPages := int(math.Ceil(float64(total) / float64(pagination.Limit)))
	response := CompaniesWPResponse{
		Items:      companies,
		Page:       pagination.Page,
		Limit:      pagination.Limit,
		Total:      total,
		TotalAll:   totalAll,
		TotalPages: totalPages,
	}

	httpx.HttpResponseWithETag(w, r, http.StatusOK, response)
	return nil
}
