// nolint
package middleware

import (
	"context"
	"neomatica/neosync/handler"
	"neomatica/neosync/infra/logger"
	"neomatica/neosync/infra/store/postgres/models"
	"neomatica/neosync/infra/store/redis"
	"neomatica/neosync/pkg/authtoken"
	"neomatica/neosync/pkg/clientip"
	"neomatica/neosync/pkg/httpx"
	"neomatica/neosync/types"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

var (
	analyticsUserOnlineMu       sync.Mutex
	analyticsUserOnlineLastSeen = map[string]time.Time{}
	analyticsUserOnlineTTL      = 5 * time.Minute
	analyticsUserOnlineLastGC   time.Time
)

func trackUserOnlineAsync(_ context.Context, h *handler.BaseHandler, userUUID string) {
	now := time.Now()

	analyticsUserOnlineMu.Lock()

	// Prune at most once per TTL. This keeps the cache bounded without scanning
	// the whole map on every authenticated request.
	if analyticsUserOnlineLastGC.IsZero() || now.Sub(analyticsUserOnlineLastGC) >= analyticsUserOnlineTTL {
		for uuid, seenAt := range analyticsUserOnlineLastSeen {
			if now.Sub(seenAt) >= analyticsUserOnlineTTL {
				delete(analyticsUserOnlineLastSeen, uuid)
			}
		}
		analyticsUserOnlineLastGC = now
	}

	lastSeen, ok := analyticsUserOnlineLastSeen[userUUID]
	if ok && now.Sub(lastSeen) < analyticsUserOnlineTTL {
		analyticsUserOnlineMu.Unlock()
		return
	}

	analyticsUserOnlineLastSeen[userUUID] = now
	analyticsUserOnlineMu.Unlock()

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := h.Store.Analytics.Create_AnalyticsUser(ctx, &models.AnalyticsUser{
			UserUUID: userUUID,
			IsOnline: true,
		}); err != nil {
			logger.Error("IsAuthenticatedMiddleware: failed create analytics user activity: %s", err.Error())
		}
	}()
}

/* Проверка пользователя, IP и сессии в Redis */
/* allowRefresh используется для /auth/check, остальные запросы только проверяют access */
func IsAuthenticatedMiddleware(h *handler.BaseHandler, allowRefresh ...bool) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			tr := TranslatorFromContext(ctx)
			canRefresh := len(allowRefresh) > 0 && allowRefresh[0]

			unauthorized := func() {
				authtoken.ClearCookies(w)
				httpx.HttpResponse(w, r, http.StatusUnauthorized, tr.TErr("auth-login-required"))
			}

			unavailable := func(err error) {
				logger.Error("IsAuthenticatedMiddleware: %s", err.Error())
				httpx.HttpResponse(w, r, http.StatusServiceUnavailable, tr.TErr("auth-temporarily-unavailable"))
			}

			var authToken types.RequestAuthToken
			var validAccess bool

			if cookie, err := r.Cookie("access_token"); err == nil {
				authToken, err = authtoken.OpenAccess(cookie.Value)
				validAccess = err == nil
			}

			refreshToken := ""

			if cookie, err := r.Cookie("refresh_token"); err == nil {
				refreshToken = cookie.Value
			}

			sessionID := authToken.SessionID

			if !validAccess {
				sessionID = authtoken.RefreshSessionID(refreshToken)
			}

			if sessionID == "" {
				unauthorized()
				return
			}

			session, err := redis.GetAuthSession(ctx, sessionID)
			if err != nil {
				unavailable(err)
				return
			}

			if session == nil || (validAccess && authToken.UUID != session.UserUUID) {
				unauthorized()
				return
			}

			/* Обновление выполняется последовательно через /auth/check */
			if !validAccess && !canRefresh {
				w.Header().Set("Cache-Control", "no-store")
				httpx.HttpResponse(w, r, http.StatusUnauthorized, map[string]string{"code": "auth_refresh_required"})
				return
			}

			validRefresh := authtoken.RefreshSessionID(refreshToken) == session.ID && authtoken.Equal(session.RefreshHash, authtoken.RefreshHash(refreshToken))

			if !validAccess && !validRefresh {
				unauthorized()
				return
			}

			ipHash, err := authtoken.IPHash(clientip.IP(r))
			if err != nil {
				unavailable(err)
				return
			}

			if !authtoken.Equal(session.IPHash, ipHash) || (validAccess && !authtoken.Equal(authToken.IPAddress, ipHash)) {
				if err := redis.DeleteAuthSession(ctx, session.ID, session.UserUUID); err != nil {
					unavailable(err)
					return
				}

				unauthorized()
				return
			}

			refreshNeeded := !validAccess || !time.Now().Before(authToken.Timestamp.Add(authtoken.AccessTTL))

			if refreshNeeded && !canRefresh {
				w.Header().Set("Cache-Control", "no-store")
				httpx.HttpResponse(w, r, http.StatusUnauthorized, map[string]string{"code": "auth_refresh_required"})
				return
			}

			if refreshNeeded && !validRefresh {
				unauthorized()
				return
			}

			user, err := h.Store.Users.Get_UserAuthByUuid(ctx, session.UserUUID)
			if err != nil {
				unavailable(err)
				return
			}

			if user == nil {
				unauthorized()
				return
			}

			userRole, err := h.Store.Users.Get_UserRoleByUuid(ctx, session.UserUUID)
			if err != nil {
				unavailable(err)
				return
			}

			if userRole == nil {
				unauthorized()
				return
			}

			if refreshNeeded {
				now := time.Now()

				newAccess, err := authtoken.SealAccess(types.RequestAuthToken{
					UUID:      session.UserUUID,
					IPAddress: ipHash,
					SessionID: session.ID,
					RoleCode:  userRole.RoleCode,
					Timestamp: now,
				})
				if err != nil {
					unavailable(err)
					return
				}

				newRefresh, err := authtoken.NewRefresh(session.ID)
				if err != nil {
					unavailable(err)
					return
				}

				rotated, err := redis.RotateAuthSessionRefresh(ctx, session.ID, session.UserUUID, authtoken.RefreshHash(refreshToken), authtoken.RefreshHash(newRefresh))
				if err != nil {
					unavailable(err)
					return
				}

				if !rotated {
					unauthorized()
					return
				}

				accessExpires := now.Add(authtoken.AccessTTL)

				if session.ExpiresAt.Before(accessExpires) {
					accessExpires = session.ExpiresAt
				}

				authtoken.SetCookie(w, "access_token", newAccess, accessExpires)
				authtoken.SetCookie(w, "refresh_token", newRefresh, session.ExpiresAt)
			}

			trackUserOnlineAsync(ctx, h, user.UserContact.UserUUID)

			ctx = context.WithValue(ctx, "identity", &types.AuthToken{
				User:     *user,
				RoleCode: userRole.RoleCode,
			})

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetIdentity(ctx context.Context) *types.AuthToken {
	if value, ok := ctx.Value("identity").(*types.AuthToken); ok {
		return value
	}

	return nil
}
