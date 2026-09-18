package meta_handler_v1

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
	metaRouter := router.PathPrefix("/meta").Subrouter()

	/* Access: SUPERADMIN | SUPPORT */
	// ?table=configurations
	metaRouter.Handle("/schema", middleware.RequireRoles(constants.Role_SuperAdmin, constants.Role_Support)(
		httpx.ErrorHandler(h.MetaSchemaColumnsHandler_V1),
	)).Methods(http.MethodGet)

	/* Access: ALL */
	metaRouter.Handle("/health/tcp", httpx.ErrorHandler(h.MetaHealthTcpHandler_V1)).Methods(http.MethodGet)

	/* Access: ALL */
	metaRouter.Handle("/s/version", httpx.ErrorHandler(h.MetaGetVersionHandler_V1)).Methods(http.MethodGet)

	/* Access: ALL */
	metaRouter.Handle("/s/releases", httpx.ErrorHandler(h.MetaGetReleasesHandler_V1)).Methods(http.MethodGet)
}
