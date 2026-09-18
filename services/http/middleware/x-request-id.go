package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func XRequestIdMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get("X-Request-ID")

			if requestID == "" {
				requestID = uuid.NewString()
			}

			w.Header().Set("X-Request-ID", requestID)

			//nolint
			ctx := context.WithValue(r.Context(), "XREQID", requestID)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
