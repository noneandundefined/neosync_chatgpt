package main

import (
	"net/http"
	"os"

	"github.com/gorilla/handlers"
)

func (s *httpServer) cors(handler http.Handler) http.Handler {
	env := os.Getenv("GO_ENV")

	var origins []string
	if env == "DEV" {
		origins = []string{"http://localhost:5173"}
	} else {
		origins = []string{
			"https://neosync.neomatica.ru",
			"https://neosync.neomatica.com",
		}
	}

	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"})
	headers := handlers.AllowedHeaders([]string{"Content-Type", "X-Requested-With", "Authorization", "X-Captcha-Token", "X-Request-ID", "Pow-Challenge", "Pow-Nonce", "X-Idempotency-Key"})
	exposed := handlers.ExposedHeaders([]string{"X-Captcha-Required", "X-Captcha-Token", "X-Request-ID"})

	allowCredentials := handlers.AllowCredentials()

	return handlers.CORS(handlers.AllowedOrigins(origins), methods, headers, exposed, allowCredentials)(handler)
}
