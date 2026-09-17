package auth_handler_v1

import (
	"neomatica/neosync/pkg/httpx"
	"net/http"
)

/* Handler: авторизация проверена в IsAuthenticatedMiddleware */
func (h *Handler) AuthCheckHandler_V1(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Cache-Control", "no-store")

	httpx.HttpResponse(w, r, http.StatusOK, map[string]bool{"authenticated": true})
	return nil
}
