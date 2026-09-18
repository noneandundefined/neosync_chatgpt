package middleware

import (
	"neomatica/neosync/pkg/httpx"
	"net/http"

	"github.com/gorilla/mux"
)

/* Используются в handler, для проверок ролей пользователей */
func RequireRoles(allowedRoles ...string) mux.MiddlewareFunc {
	roleSet := map[string]struct{}{}
	for _, role := range allowedRoles {
		roleSet[role] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tr := TranslatorFromContext(r.Context())

			identity := GetIdentity(r.Context())
			if identity == nil {
				httpx.HttpResponse(w, r, http.StatusUnauthorized, tr.TErr("auth-login-required"))
				return
			}

			if _, ok := roleSet[identity.RoleCode]; !ok {
				httpx.HttpResponse(w, r, http.StatusForbidden, tr.TErr("auth-not-enough-rights"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
