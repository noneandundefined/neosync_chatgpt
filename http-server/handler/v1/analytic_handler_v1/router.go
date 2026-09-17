package analytic_handler_v1

import (
	"neomatica/neosync/infra/constants"
	"neomatica/neosync/middleware"
	"neomatica/neosync/pkg/httpx"
	"net/http"

	"github.com/gorilla/mux"
)

/* Neosync HTTPx V1 */
/* RegisterRoutes: авторизация всех путей */

func (h *Handler) RegisterRoutes(router *mux.Router) {
	/* Anonymous acquisition and authentication events */
	router.Handle("/analytics/public-events", httpx.ErrorHandler(h.CreateAnalyticsEventsHandler_V1)).Methods(http.MethodPost)

	analyticRouter := router.PathPrefix("/analytics").Subrouter()

	/* Middleware - доступ могут получить только авторизованные пользователи */
	analyticRouter.Use(middleware.IsAuthenticatedMiddleware(h.BaseHandler))

	/* Product analytics ingestion for authenticated clients */
	analyticRouter.Handle("/events", httpx.ErrorHandler(h.CreateAnalyticsEventsHandler_V1)).Methods(http.MethodPost)

	/* Access: SUPERADMIN | SUPPORT */
	analyticRouter.Handle("/widgets", middleware.RequireRoles(constants.Role_SuperAdmin, constants.Role_Support)(
		httpx.ErrorHandler(h.GetAnalyticWidgetsHandler_V1),
	)).Methods(http.MethodGet)

	/* Access: SUPERADMIN | SUPPORT */
	analyticRouter.Handle("/query", middleware.RequireRoles(constants.Role_SuperAdmin, constants.Role_Support)(
		httpx.ErrorHandler(h.AnalyticQueryExecHandler_V1),
	)).Methods(http.MethodPost)
}
