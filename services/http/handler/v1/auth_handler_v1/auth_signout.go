package auth_handler_v1

import (
	"neomatica/neosync/infra/store/redis"
	"neomatica/neosync/middleware"
	"neomatica/neosync/pkg/authtoken"
	"neomatica/neosync/pkg/httpx"
	"neomatica/neosync/pkg/httpx/httperr"
	"net/http"
)

/* Handler: отзыв сессии в Redis и выход пользователя */
func (h *Handler) AuthSignoutHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)

	/* Access можно использовать для выхода и после окончания его срока */
	if cookie, err := r.Cookie("access_token"); err == nil {
		if authToken, err := authtoken.OpenAccess(cookie.Value); err == nil {
			if err := redis.DeleteAuthSession(ctx, authToken.SessionID, authToken.UUID); err != nil {
				return httperr.Redis(ctx, err)
			}
		}
	}

	if cookie, err := r.Cookie("refresh_token"); err == nil {
		sessionID := authtoken.RefreshSessionID(cookie.Value)

		if sessionID != "" {
			session, err := redis.GetAuthSession(ctx, sessionID)
			if err != nil {
				return httperr.Redis(ctx, err)
			}

			if session != nil && authtoken.Equal(session.RefreshHash, authtoken.RefreshHash(cookie.Value)) {
				if err := redis.DeleteAuthSession(ctx, session.ID, session.UserUUID); err != nil {
					return httperr.Redis(ctx, err)
				}
			}
		}
	}

	authtoken.ClearCookies(w)

	httpx.HttpResponse(w, r, http.StatusOK, tr.T("logged-out"))
	return nil
}
