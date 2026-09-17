package device_handler_v1

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
	deviceRouter := router.PathPrefix("/devices").Subrouter()

	/* Middleware - доступ могут получить только авторизованные пользователи */
	deviceRouter.Use(middleware.IsAuthenticatedMiddleware(h.BaseHandler))

	/* Access: ALL */
	deviceRouter.Handle("", httpx.ErrorHandler(h.CreateDeviceHandler_V1)).Methods(http.MethodPost)

	/* Access: ALL */
	deviceRouter.Handle("", httpx.ErrorHandler(h.GetDevicesHandler_V1)).Methods(http.MethodGet)

	/* Access: ALL */
	// deviceRouter.Handle("/full-list-by-user", middleware.PowDDos()(
	// 	httpx.ErrorHandler(h.GetDevicesByUserHandler_V1),
	// )).Methods(http.MethodGet)

	/* Access: SUPERADMIN | SUPPORT */
	deviceRouter.Handle("/full-list", httpx.ErrorHandler(h.GetDevicesFullHandler_V1)).Methods(http.MethodGet)

	/* Access: ALL */
	deviceRouter.Handle("/import", middleware.PowDDos()(
		httpx.ErrorHandler(h.DevicesImportHandler_V1),
	)).Methods(http.MethodPost)

	/* Access: ALL */
	deviceRouter.Handle("/import/check", httpx.ErrorHandler(h.DevicesImportCheckHandler_V1)).Methods(http.MethodPost)

	/* Access: ALL */
	deviceRouter.Handle("/state", httpx.ErrorHandler(h.GetDevicesStateListHandler_V1)).Methods(http.MethodGet)

	/* Access: ALL */
	deviceRouter.Handle("/state/devices", httpx.ErrorHandler(h.GetDevicesStateDevicesHandler_V1)).Methods(http.MethodGet)

	/* Access: ALL */
	deviceRouter.Handle("/status/detailed", httpx.ErrorHandler(h.GetDeviceDetailedStatusByImeisHandler_V1)).Methods(http.MethodPost)

	/* Access: ALL */
	deviceRouter.Handle("/delete", middleware.PowDDos()(
		httpx.ErrorHandler(h.DeviceMassiveDeleteHandler_V1),
	)).Methods(http.MethodPost)

	/* Access: SUPERADMIN | SUPPORT | DEALER | DEALER_SUPPORT */
	deviceRouter.Handle("/transfer",
		middleware.PowDDos()(
			middleware.RequireRoles(constants.Role_SuperAdmin, constants.Role_Support, constants.Role_AdminL2, constants.Role_AdminL2Support)(
				httpx.ErrorHandler(h.DeviceTransferToAccountHandler_V1),
			),
		),
	).Methods(http.MethodPatch)

	/* Access: ALL */
	deviceRouter.Handle("/{imei:[0-9]{15}}/logs", httpx.ErrorHandler(h.GetDeviceLogsHandler_V1)).Methods(http.MethodGet)

	/* Access: ALL */
	deviceRouter.Handle("/{imei:[0-9]{15}}/status", httpx.ErrorHandler(h.GetDeviceStatusByImeiHandler_V1)).Methods(http.MethodGet)

	/* Access: ALL */
	deviceRouter.Handle("/{imei:[0-9]{15}}", httpx.ErrorHandler(h.DeviceUpdateByImeiHandler_V1)).Methods(http.MethodPatch)

	/* Access: ALL */
	deviceRouter.Handle("/{imei:[0-9]{15}}", httpx.ErrorHandler(h.GetDeviceByImeiHandler_V1)).Methods(http.MethodGet)

	/* Access: ALL */
	deviceRouter.Handle("/{imei:[0-9]{15}}", httpx.ErrorHandler(h.DeleteDeviceHandler_V1)).Methods(http.MethodDelete)
}
