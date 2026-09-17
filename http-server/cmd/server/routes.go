package main

import (
	"neomatica/neosync/handler"
	"neomatica/neosync/infra/worker"

	"neomatica/neosync/handler/v1/analytic_handler_v1"
	"neomatica/neosync/handler/v1/auth_handler_v1"
	"neomatica/neosync/handler/v1/command_handler_v1"
	"neomatica/neosync/handler/v1/company_handler_v1"
	"neomatica/neosync/handler/v1/configuration_handler_v1"
	"neomatica/neosync/handler/v1/device_handler_v1"
	"neomatica/neosync/handler/v1/firmware_handler_v1"
	"neomatica/neosync/handler/v1/group_handler_v1"
	"neomatica/neosync/handler/v1/meta_handler_v1"
	"neomatica/neosync/handler/v1/user_handler_v1"

	"neomatica/neosync/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *httpServer) routes() http.Handler {
	router := mux.NewRouter()

	/* Middleware for logging API request */
	router.Use(middleware.NewLogger().LoggerMiddleware)
	/* Middleware for X-Request-Id */
	router.Use(middleware.XRequestIdMiddleware())
	/* Middleware for i18n language */
	router.Use(middleware.LanguageMiddleware())
	/* Middleware for get exception errors */
	router.Use(middleware.RecoveryMiddleware())
	/* Middleware for security API */
	router.Use(middleware.SecurityMiddleware())
	/* Middleware rate limiter */
	router.Use(middleware.RateLimiterMiddleware(6, 10))
	/* Middleware abuse scoring */
	// router.Use(middleware.AbuseScoringMiddleware())

	/* Middleware rate limiter cleanup */
	middleware.RateLimiterCleanup()

	subrouter := router.PathPrefix("/api/v1").Subrouter()

	baseHandler := &handler.BaseHandler{
		Db:        s.db,
		Store:     s.store,
		UseCase:   s.usecase,
		RMQ:       s.rmq,
		Session:   s.session,
		Analytics: s.analytics,
	}

	companyWorker := worker.NewCompanyTasksWorker(s.store)

	/* Configuration routes */
	configuration_handler_v1.NewHandler(baseHandler).RegisterRoutes(subrouter)
	/* Command routes */
	command_handler_v1.NewHandler(baseHandler).RegisterRoutes(subrouter)
	/* Device routes */
	device_handler_v1.NewHandler(baseHandler).RegisterRoutes(subrouter)
	/* Authenticate rotues */
	auth_handler_v1.NewHandler(baseHandler).RegisterRoutes(subrouter)
	/* User routes */
	user_handler_v1.NewHandler(baseHandler).RegisterRoutes(subrouter)
	/* Group routes */
	group_handler_v1.NewHandler(baseHandler).RegisterRoutes(subrouter)
	/* Analytic routes */
	company_handler_v1.NewHandler(baseHandler, companyWorker).RegisterRoutes(subrouter)
	/* Device firmware routes */
	firmware_handler_v1.NewHandler(baseHandler).RegisterRoutes(subrouter)
	/* MetaData routes */
	meta_handler_v1.NewHandler(baseHandler).RegisterRoutes(subrouter)
	/* Analytic routes */
	analytic_handler_v1.NewHandler(baseHandler).RegisterRoutes(subrouter)

	/* Doc routes */
	s.docs(subrouter)

	return s.cors(router)
}
