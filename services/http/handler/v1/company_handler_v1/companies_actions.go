package company_handler_v1

import (
	"neomatica/neosync/infra/constants"
	"neomatica/neosync/middleware"
	"neomatica/neosync/pkg/httpx"
	"neomatica/neosync/pkg/httpx/httperr"
	gtypes "neomatica/neosync/types"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type CompanyActionResponse struct {
	Affected int    `json:"affected"`
	Message  string `json:"message"`
}

func (h *Handler) CompanyRetryFailedHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*gtypes.AuthToken)

	if !authToken.User.Has(constants.AccessConfigurationApply) {
		return httperr.Forbidden(tr.TErr("config-edit-access-restricted"))
	}

	companyID, err := parseCompanyID(mux.Vars(r)["id"])
	if err != nil {
		return httperr.BadRequest(tr.TErr("invalid-id"))
	}

	company, err := h.Store.Companies.Get_CompanyByIdAndUserUuid(ctx, companyID, authToken.User.UserContact.UserUUID)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	if company == nil {
		return httperr.NotFound(tr.TErr("company-not-found"))
	}

	affected, err := h.Store.Companies.Retry_FailedCompanyTasksByCompanyId(ctx, companyID)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	if affected == 0 {
		return httperr.BadRequest(tr.TErr("company-retry-failed-nothing"))
	}

	httpx.HttpResponse(w, r, http.StatusOK, CompanyActionResponse{
		Affected: int(affected),
		Message:  tr.T("company-retry-failed-success"),
	})

	return nil
}

func (h *Handler) CompanyCancelPendingHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*gtypes.AuthToken)

	if !authToken.User.Has(constants.AccessConfigurationApply) {
		return httperr.Forbidden(tr.TErr("config-edit-access-restricted"))
	}

	companyID, err := parseCompanyID(mux.Vars(r)["id"])
	if err != nil {
		return httperr.BadRequest(tr.TErr("invalid-id"))
	}

	company, err := h.Store.Companies.Get_CompanyByIdAndUserUuid(ctx, companyID, authToken.User.UserContact.UserUUID)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	if company == nil {
		return httperr.NotFound(tr.TErr("company-not-found"))
	}

	affected, err := h.Store.Companies.Cancel_PendingCompanyTasksByCompanyId(ctx, companyID)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	if affected == 0 {
		return httperr.BadRequest(tr.TErr("company-cancel-pending-nothing"))
	}

	httpx.HttpResponse(w, r, http.StatusOK, CompanyActionResponse{
		Affected: int(affected),
		Message:  tr.T("company-cancel-pending-success"),
	})

	return nil
}

func parseCompanyID(idParam string) (uint64, error) {
	companyID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil || companyID == 0 {
		return 0, err
	}

	return companyID, nil
}
