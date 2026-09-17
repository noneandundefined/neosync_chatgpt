package company_handler_v1

import (
	"neomatica/neosync/middleware"
	"neomatica/neosync/pkg/httpx"
	"net/http"

	"github.com/gorilla/mux"
)

/* Neosync HTTPx V1 */
/* RegisterRoutes: авторизация всех путей */

func (h *Handler) RegisterRoutes(router *mux.Router) {
	companyRouter := router.PathPrefix("/companies").Subrouter()

	/* Middleware - доступ могут получить только авторизованные пользователи */
	companyRouter.Use(middleware.IsAuthenticatedMiddleware(h.BaseHandler))

	/* Access: ALL */
	companyRouter.Handle("", httpx.ErrorHandler(h.GetCompaniesHandlerV1)).Methods(http.MethodGet)

	/* Access: ALL */
	companyRouter.Handle("", httpx.ErrorHandler(h.CompanyCreateHandler_V1)).Methods(http.MethodPost)

	/* Access: ALL */
	companyRouter.Handle("/{id}", httpx.ErrorHandler(h.GetCompanyByIdHandler_V1)).Methods(http.MethodGet)

	/* Access: ALL */
	companyRouter.Handle("/{id}/retry-failed", middleware.PowDDos()(
		httpx.ErrorHandler(h.CompanyRetryFailedHandler_V1),
	)).Methods(http.MethodPost)

	/* Access: ALL */
	companyRouter.Handle("/{id}/cancel-pending", middleware.PowDDos()(
		httpx.ErrorHandler(h.CompanyCancelPendingHandler_V1),
	)).Methods(http.MethodPost)
}
