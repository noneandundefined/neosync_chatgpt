package firmware_handler_v1

import (
	"neomatica/neosync/middleware"
	"neomatica/neosync/pkg/httpx"
	"net/http"

	"github.com/gorilla/mux"
)

/* Neosync HTTPx V1 */
/* RegisterRoutes: авторизация всех путей */

func (h *Handler) RegisterRoutes(router *mux.Router) {
	firmwareRouter := router.PathPrefix("/firmwares").Subrouter()

	/* Middleware - доступ могут получить только авторизованные пользователи */
	firmwareRouter.Use(middleware.IsAuthenticatedMiddleware(h.BaseHandler))

	/* Access: ALL */
	firmwareRouter.Handle("/models", httpx.ErrorHandler(h.GetFirmwareModelsHandler_V1)).Methods(http.MethodGet)

	/* Access: ALL */
	firmwareRouter.Handle("/sources", httpx.ErrorHandler(h.GetFirmwareListHandler_V1)).Methods(http.MethodGet)

	/* Access: ALL */
	firmwareRouter.Handle("/devices/{imei:[0-9]{15}}/firmware", httpx.ErrorHandler(h.FirmwareUpdateFirmwareVersionByImeiHandler_V1)).Methods(http.MethodPost)
}
