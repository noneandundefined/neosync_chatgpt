package company_handler_v1

import (
	"fmt"
	"neomatica/neosync/middleware"
	"neomatica/neosync/pkg/httpx"
	"neomatica/neosync/pkg/httpx/httperr"
	"neomatica/neosync/types"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

/* Neosync HTTPx V1 */
/* Handler: получение массовой настройки по ID */

func (h *Handler) GetCompanyByIdHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*types.AuthToken)

	idParam := mux.Vars(r)["id"]
	companyID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil || companyID == 0 {
		return httperr.BadRequest(fmt.Sprintf(tr.TErr("invalid-id"), idParam))
	}

	company, err := h.Store.Companies.Get_CompanyByIdAndUserUuid(ctx, companyID, authToken.User.UserContact.UserUUID)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	if company == nil {
		return httperr.NotFound(tr.TErr("company-not-found"))
	}

	response := buildCompanyDetailsResponse(company)

	httpx.HttpResponseWithETag(w, r, http.StatusOK, response)
	return nil
}
