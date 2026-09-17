package auth_handler_v1

import (
	"neomatica/neosync/middleware"
	"neomatica/neosync/pkg/httpx"
	"net/http"

	"github.com/gorilla/mux"
)

/* Neosync HTTPx V1 */
/* RegisterRoutes: авторизация всех путей */

func (h *Handler) RegisterRoutes(router *mux.Router) {
	authRouter := router.PathPrefix("/auth").Subrouter()

	/* Access: ALL */
	authRouter.Handle("/signin", middleware.PowDDos()(
		httpx.ErrorHandler(h.AuthSigninHandler_V1),
	)).Methods(http.MethodPost)

	/* Access: ALL */
	authRouter.Handle("/check", middleware.IsAuthenticatedMiddleware(h.BaseHandler, true)(
		httpx.ErrorHandler(h.AuthCheckHandler_V1),
	)).Methods(http.MethodGet)

	/* Access: ALL */
	/* /auth/password/reset/request?email=test@test.com */
	authRouter.Handle("/password/reset/request", middleware.PowDDos()(
		httpx.ErrorHandler(h.AuthRequestPasswordResetHandler_V1),
	)).Methods(http.MethodPost)

	/* Access: ALL */
	/* /auth/password/reset/confirm?uuid=&exp=&sig=, {password: "new"} */
	authRouter.Handle("/password/reset/confirm", middleware.PowDDos()(
		httpx.ErrorHandler(h.AuthPasswordResetHandler),
	)).Methods(http.MethodPost)

	/* Access: ALL */
	authRouter.Handle("/signout", httpx.ErrorHandler(h.AuthSignoutHandler_V1)).Methods(http.MethodPost)
}
