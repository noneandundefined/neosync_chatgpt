package user_handler_v1

import (
	"neomatica/neosync/infra/constants"
	"neomatica/neosync/pkg/httpx"
	"net/http"

	"neomatica/neosync/middleware"

	"github.com/gorilla/mux"
)

/* Neosync HTTPx V1 */
/* RegisterRoutes: авторизация всех путей */

func (h *Handler) RegisterRoutes(router *mux.Router) {
	userRouter := router.PathPrefix("/users").Subrouter()

	/* Middleware - доступ могут получить только авторизованные пользователи */
	userRouter.Use(middleware.IsAuthenticatedMiddleware(h.BaseHandler))

	/* Access: SUPERADMIN | SUPPORT | DEALER | DEALER_SUPPORT */
	userRouter.Handle("",
		middleware.RequireRoles(constants.Role_SuperAdmin, constants.Role_Support, constants.Role_AdminL2, constants.Role_AdminL2Support)(
			httpx.ErrorHandler(h.GetUsersHandler_V1),
		),
	).Methods(http.MethodGet)

	/* Access: ALL */
	userRouter.Handle("/tree", httpx.ErrorHandler(h.GetUsersTreeHandler_V1)).Methods(http.MethodGet)

	/* Access: SUPERADMIN | SUPPORT */
	userRouter.Handle("/full-list",
		middleware.RequireRoles(constants.Role_SuperAdmin, constants.Role_Support)(
			httpx.ErrorHandler(h.GetUsersFullHandler_V1),
		),
	).Methods(http.MethodGet)

	/* Access: ALL */
	userRouter.Handle("/owners", httpx.ErrorHandler(h.GetUsersToOwnerHandler_V1)).Methods(http.MethodGet)

	/* Access: ALL */
	userRouter.Handle("/me", httpx.ErrorHandler(h.GetUserMeHandler_V1)).Methods(http.MethodGet)

	/* Access: ALL */
	userRouter.Handle("/me", httpx.ErrorHandler(h.UserUpdateMeHandler_V1)).Methods(http.MethodPatch)

	/* Access: ALL */
	userRouter.Handle("/me/configuration-priority", httpx.ErrorHandler(h.UserUpdateConfigurationPriorityHandler_V1)).Methods(http.MethodPatch)

	/* Access: DEALER */
	userRouter.Handle("/me/cvcg",
		middleware.RequireRoles(constants.Role_AdminL2)(
			httpx.ErrorHandler(h.UserUpdateCVCGByUuidHandler_V1),
		),
	).Methods(http.MethodPatch)

	/* Access: ALL */
	userRouter.Handle("/me/accesses", httpx.ErrorHandler(h.GetUserMeAccessesHandler_V1)).Methods(http.MethodGet)

	/* Access: ALL */
	userRouter.Handle("/me/login-with-role", httpx.ErrorHandler(h.GetUserMeLoginWithRoleCodeHandler_V1)).Methods(http.MethodGet)

	/* Access: SUPERADMIN | SUPPORT | DEALER | DEALER_SUPPORT */
	userRouter.Handle("/create",
		middleware.PowDDos()(
			middleware.RequireRoles(constants.Role_SuperAdmin, constants.Role_Support, constants.Role_AdminL2, constants.Role_AdminL2Support)(
				httpx.ErrorHandler(h.UserCreateUserHandler_V1),
			),
		),
	).Methods(http.MethodPost)

	/* Access: SUPERADMIN | SUPPORT | DEALER | DEALER_SUPPORT */
	userRouter.Handle("/delete",
		middleware.PowDDos()(
			middleware.RequireRoles(constants.Role_SuperAdmin, constants.Role_Support, constants.Role_AdminL2, constants.Role_AdminL2Support)(
				httpx.ErrorHandler(h.UserMassiveDeleteHandler_V1),
			),
		),
	).Methods(http.MethodPost)

	/* Access: SUPERADMIN | SUPPORT | DEALER | DEALER_SUPPORT */
	userRouter.Handle("/{uuid:[0-9a-fA-F-]{36}}",
		middleware.PowDDos()(
			middleware.RequireRoles(constants.Role_SuperAdmin, constants.Role_Support, constants.Role_AdminL2, constants.Role_AdminL2Support)(
				httpx.ErrorHandler(h.UserGetOnesByUuidHandler_V1),
			),
		),
	).Methods(http.MethodGet)

	/* Access: SUPERADMIN | SUPPORT | DEALER | DEALER_SUPPORT */
	userRouter.Handle("/{uuid:[0-9a-fA-F-]{36}}",
		middleware.RequireRoles(constants.Role_SuperAdmin, constants.Role_Support, constants.Role_AdminL2, constants.Role_AdminL2Support)(
			httpx.ErrorHandler(h.UserUpdateOnesByUuidHandler_V1),
		),
	).Methods(http.MethodPatch)

	/* Access: SUPERADMIN | SUPPORT | DEALER | DEALER_SUPPORT */
	userRouter.Handle("/{uuid:[0-9a-fA-F-]{36}}",
		middleware.RequireRoles(constants.Role_SuperAdmin, constants.Role_Support, constants.Role_AdminL2, constants.Role_AdminL2Support)(
			httpx.ErrorHandler(h.UserDeleteByUuidHandler_V1),
		),
	).Methods(http.MethodDelete)

	/* Access: SUPERADMIN | SUPPORT | DEALER | DEALER_SUPPORT */
	userRouter.Handle("/{uuid:[0-9a-fA-F-]{36}}/send/email",
		middleware.PowDDos()(
			middleware.RequireRoles(constants.Role_SuperAdmin, constants.Role_Support, constants.Role_AdminL2, constants.Role_AdminL2Support)(
				httpx.ErrorHandler(h.UserSendEmailInfoHandler_V1),
			),
		),
	).Methods(http.MethodPost)
}
