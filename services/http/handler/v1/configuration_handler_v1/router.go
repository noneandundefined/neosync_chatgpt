package configuration_handler_v1

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
	configurationRouter := router.PathPrefix("/devices").Subrouter()

	/* Middleware - доступ могут получить только авторизованные пользователи */
	configurationRouter.Use(middleware.IsAuthenticatedMiddleware(h.BaseHandler))

	/* Access: ALL */
	configurationRouter.Handle("/configuration-templates", httpx.ErrorHandler(h.GetConfigurationTemplatesHandler_V1)).Methods(http.MethodGet)

	/* Access: ALL */
	configurationRouter.Handle("/configuration-templates", middleware.PowDDos()(
		httpx.ErrorHandler(h.CreateConfigurationTemplateHandler_V1),
	)).Methods(http.MethodPost)

	/* Access: ALL */
	configurationRouter.Handle("/configuration-templates/{id}", httpx.ErrorHandler(h.GetConfigurationTemplateByIdHandler_V1)).Methods(http.MethodGet)

	/* Access: ALL */
	configurationRouter.Handle("/configuration-templates/{id}", middleware.PowDDos()(
		httpx.ErrorHandler(h.UpdateConfigurationTemplateHandler_V1),
	)).Methods(http.MethodPut)

	/* Access: ALL */
	configurationRouter.Handle("/configuration-templates/{id}", middleware.PowDDos()(
		httpx.ErrorHandler(h.DeleteConfigurationTemplateHandler_V1),
	)).Methods(http.MethodDelete)

	/* Access: ALL */
	configurationRouter.Handle("/configuration-templates/{id}/configuration/parsed/{section}", httpx.ErrorHandler(h.GetConfigurationTemplateByIdParsedBySectionHandler_V1)).Methods(http.MethodGet)

	/* Access: ALL */
	configurationRouter.Handle("/{imei:[0-9]{15}}/configuration/parsed", httpx.ErrorHandler(h.GetConfigurationParsedByImeiHandler_V1)).Methods(http.MethodGet)

	/* Access: SUPERADMIN | SUPPORT */
	configurationRouter.Handle("/{imei:[0-9]{15}}/configuration/raw", middleware.RequireRoles(constants.Role_SuperAdmin, constants.Role_Support)(
		httpx.ErrorHandler(h.GetConfigurationRawByImeiHandler_V1),
	)).Methods(http.MethodGet)

	/* Access: ALL */
	configurationRouter.Handle("/{imei:[0-9]{15}}/configuration/parsed/{section}", httpx.ErrorHandler(h.GetConfigurationParsedBySectionHandler_V1)).Methods(http.MethodGet)

	/* Access: ALL */
	configurationRouter.Handle("/templates/configuration/parsed/{section}", httpx.ErrorHandler(h.GetConfigurationTemplateParsedBySectionHandler_V1)).Methods(http.MethodGet)

	/* Access: ALL */
	configurationRouter.Handle("/{imei:[0-9]{15}}/configuration/draft/reset", middleware.PowDDos()(
		httpx.ErrorHandler(h.ResetConfigurationDraftHandler_V1),
	)).Methods(http.MethodDelete)

	/* Access: ALL */
	configurationRouter.Handle("/{imei:[0-9]{15}}/configuration/draft/check", httpx.ErrorHandler(h.GetConfigurationDraftCheckHandler_V1)).Methods(http.MethodGet)

	/* Access: ALL */
	configurationRouter.Handle("/{imei:[0-9]{15}}/configuration/draft/{section}", httpx.ErrorHandler(h.GetConfigurationDraftHandler_V1)).Methods(http.MethodGet)

	/* Access: ALL */
	configurationRouter.Handle("/{imei:[0-9]{15}}/configuration/draft/{section}", httpx.ErrorHandler(h.SetConfigurationDraftHandler_V1)).Methods(http.MethodPost)

	/* Access: ALL */
	configurationRouter.Handle("/{imei:[0-9]{15}}/configuration/apply", middleware.PowDDos()(
		httpx.ErrorHandler(h.ApplyConfigurationHandler_V1),
	)).Methods(http.MethodPost)

	/* Access: ALL */
	configurationRouter.Handle("/{imei:[0-9]{15}}/configuration/apply-template/{id}", middleware.PowDDos()(
		httpx.ErrorHandler(h.ApplyConfigurationTemplateHandler_V1),
	)).Methods(http.MethodPost)

	/* Access: ALL */
	configurationRouter.Handle("/{imei:[0-9]{15}}/configuration/export", httpx.ErrorHandler(h.ConfigurationExportByImeiHandler_V1)).Methods(http.MethodGet)

	/* Access: ALL */
	configurationRouter.Handle("/{imei:[0-9]{15}}/configuration/import", httpx.ErrorHandler(h.ConfigurationImportByImeiHandler_V1)).Methods(http.MethodPost)

	/* Access: ALL */
	configurationRouter.Handle("/{imei:[0-9]{15}}/configuration/telemetry", httpx.ErrorHandler(h.GetConfigurationTelemetryHandler_V1)).Methods(http.MethodGet)

	/* Access: ALL */
	configurationRouter.Handle("/{imei:[0-9]{15}}/configuration/telemetry", httpx.ErrorHandler(h.RebootConfigurationTelemetryHandler_V1)).Methods(http.MethodPost)

	/* Access: ALL */
	configurationRouter.Handle("/{imei:[0-9]{15}}/configuration/pass", middleware.PowDDos()(
		httpx.ErrorHandler(h.ConfigurationChangePassHandler_V1),
	)).Methods(http.MethodPost)

	/* Access: ALL */
	configurationRouter.Handle("/{imei:[0-9]{15}}/configuration/export/llstariration", httpx.ErrorHandler(h.ConfigurationExportLLSTarirationHandler_V1)).Methods(http.MethodPost)
}
