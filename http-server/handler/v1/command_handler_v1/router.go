package command_handler_v1

import (
	"neomatica/neosync/middleware"
	"neomatica/neosync/pkg/httpx"
	"net/http"

	"github.com/gorilla/mux"
)

/* Neosync HTTPx V1 */
/* RegisterRoutes: авторизация всех путей */

func (h *Handler) RegisterRoutes(router *mux.Router) {
	commandRouter := router.PathPrefix("/devices").Subrouter()

	/* Middleware - доступ могут получить только авторизованные пользователи */
	commandRouter.Use(middleware.IsAuthenticatedMiddleware(h.BaseHandler))

	/* Access: ALL */
	commandRouter.Handle("/cmd", httpx.ErrorHandler(h.CommandCreateAndSendHandler_V1)).Methods(http.MethodPost)

	/* Access: ALL */
	commandRouter.Handle("/cmd", httpx.ErrorHandler(h.GetCommandHistoryHandler_V1)).Methods(http.MethodGet)

	/* Access: ALL */
	/* Получение по сессии - /cmd?session_id=123 */
	commandRouter.Handle("/cmd/sessions/{session_id}", httpx.ErrorHandler(h.GetCommandRecentHistoryHandler_V1)).Methods(http.MethodGet)

	/* Access: ALL */
	commandRouter.Handle("/cmd/{id:[0-9]+}", httpx.ErrorHandler(h.GetCommandByCommandIdHandler_V1)).Methods(http.MethodGet)

	/* Access: ALL */
	commandRouter.Handle("/cmd/{id:[0-9]+}/a/cancel", middleware.PowDDos()(
		httpx.ErrorHandler(h.CommandCancelHandler_V1),
	)).Methods(http.MethodPatch)

	/* Access: ALL */
	commandRouter.Handle("/cmd/{id:[0-9]+}/a/delete", middleware.PowDDos()(
		httpx.ErrorHandler(h.DeleteCommandHandler_V1),
	)).Methods(http.MethodDelete)

	/* Access: ALL */
	commandRouter.Handle("/{imei:[0-9]{15}}/cmd/send-and-wait", httpx.ErrorHandler(h.CommandSendAndWaitHandler_V1)).Methods(http.MethodPost)

	/* Access: ALL */
	commandRouter.Handle("/{imei:[0-9]{15}}/cmd/reboot", middleware.PowDDos()(
		httpx.ErrorHandler(h.CommandRebootHandler_V1),
	)).Methods(http.MethodPost)

	/* Access: ALL */
	commandRouter.Handle("/{imei:[0-9]{15}}/cmd/erase_eeprom", middleware.PowDDos()(
		httpx.ErrorHandler(h.CommandEraseEepromHandler_V1),
	)).Methods(http.MethodPost)

	/* Access: ALL */
	commandRouter.Handle("/{imei:[0-9]{15}}/cmd/erase_flash", middleware.PowDDos()(
		httpx.ErrorHandler(h.CommandEraseFlashHandler_V1),
	)).Methods(http.MethodPost)

	/* Access: ALL */
	commandRouter.Handle("/{imei:[0-9]{15}}/cmd/find_ble_sensors", middleware.PowDDos()(
		httpx.ErrorHandler(h.CommandFindBleSensorsHandler_V1),
	)).Methods(http.MethodPost)

	/* Access: ALL */
	commandRouter.Handle("/{imei:[0-9]{15}}/cmd/find_owire_sensors", middleware.PowDDos()(
		httpx.ErrorHandler(h.CommandOWireHandler_V1),
	)).Methods(http.MethodPost)
}
