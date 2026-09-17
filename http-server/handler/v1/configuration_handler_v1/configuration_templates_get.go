package configuration_handler_v1

import (
	"math"
	"neomatica/neosync/infra/constants"
	"neomatica/neosync/middleware"
	"neomatica/neosync/pkg/httpx"
	"neomatica/neosync/pkg/httpx/httperr"
	"neomatica/neosync/types"
	"neomatica/neosync/util"
	"net/http"
)

/* Neosync HTTPx V1 */
/* Handler: получение шаблонов конфигурации пользователя */

func (h *Handler) GetConfigurationTemplatesHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*types.AuthToken)

	if !authToken.User.Has(constants.AccessConfigurationRead) {
		return httperr.Forbidden(tr.TErr("config-access-restricted"))
	}

	pagination := util.GetPagination(r, 10)

	templates, total, totalAll, err := h.Store.ConfigurationTemplates.Get_ConfigurationTemplatesByUserUuid(
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
	response := ConfigurationTemplatesWPResponse{
		Items:      templates,
		Page:       pagination.Page,
		Limit:      pagination.Limit,
		Total:      total,
		TotalAll:   totalAll,
		TotalPages: totalPages,
	}

	httpx.HttpResponseWithETag(w, r, http.StatusOK, response)
	return nil
}
