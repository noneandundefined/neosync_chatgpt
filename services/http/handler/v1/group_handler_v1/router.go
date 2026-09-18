package group_handler_v1

import (
	"neomatica/neosync/middleware"
	"neomatica/neosync/pkg/httpx"
	"net/http"

	"github.com/gorilla/mux"
)

/* Neosync HTTPx V1 */
/* RegisterRoutes: авторизация всех путей */

func (h *Handler) RegisterRoutes(router *mux.Router) {
	groupRouter := router.PathPrefix("/groups").Subrouter()

	/* Middleware - доступ могут получить только авторизованные пользователи */
	groupRouter.Use(middleware.IsAuthenticatedMiddleware(h.BaseHandler))

	/* Access: ALL */
	groupRouter.Handle("", httpx.ErrorHandler(h.GroupCreateHandler_V1)).Methods(http.MethodPost)

	/* Access: ALL */
	groupRouter.Handle("", httpx.ErrorHandler(h.GetGroupsHandler_V1)).Methods(http.MethodGet)

	/* Access: SUPERADMIN | SUPPORT */
	groupRouter.Handle("/full-list", httpx.ErrorHandler(h.GetGroupsFullHandler_V1)).Methods(http.MethodGet)

	/* Access: ALL */
	groupRouter.Handle("/fields/name", httpx.ErrorHandler(h.GetGroupsFieldsNameHandler_V1)).Methods(http.MethodGet)

	/* Access: ALL */
	groupRouter.Handle("/delete", middleware.PowDDos()(
		httpx.ErrorHandler(h.GroupDeleteMassiveHandler_V1),
	)).Methods(http.MethodPost)

	/* Access: ALL */
	groupRouter.Handle("/{id:[0-9]+}/members", middleware.PowDDos()(
		httpx.ErrorHandler(h.GroupMembersListHandler_V1),
	)).Methods(http.MethodGet)

	/* Access: ALL */
	groupRouter.Handle("/{id:[0-9]+}/members/upsert", middleware.PowDDos()(
		httpx.ErrorHandler(h.GroupMembersUpsertMassiveHandler_V1),
	)).Methods(http.MethodPost)

	/* Access: ALL */
	groupRouter.Handle("/{id:[0-9]+}/members/remove", middleware.PowDDos()(
		httpx.ErrorHandler(h.GroupMembersDeleteMassiveHandler_V1),
	)).Methods(http.MethodPost)

	/* Access: ALL */
	groupRouter.Handle("/{id:[0-9]+}", httpx.ErrorHandler(h.GetGroupWithDevicesHandler_V1)).Methods(http.MethodGet)

	/* Access: ALL */
	groupRouter.Handle("/{id:[0-9]+}", httpx.ErrorHandler(h.GroupUpdateByIdHandler_V1)).Methods(http.MethodPut)

	/* Access: ALL */
	groupRouter.Handle("/{id:[0-9]+}", httpx.ErrorHandler(h.GroupDeleteByIdHandler_V1)).Methods(http.MethodDelete)
}
