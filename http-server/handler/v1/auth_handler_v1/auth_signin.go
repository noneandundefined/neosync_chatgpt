package auth_handler_v1

import (
	"crypto/subtle"
	"database/sql"
	"neomatica/neosync/encryption"
	"neomatica/neosync/infra/logger"
	"neomatica/neosync/infra/store/postgres/store"
	"neomatica/neosync/infra/store/redis"
	"neomatica/neosync/middleware"
	"neomatica/neosync/pkg/authtoken"
	"neomatica/neosync/pkg/clientip"
	"neomatica/neosync/pkg/httpx"
	"neomatica/neosync/pkg/httpx/httperr"
	"neomatica/neosync/pkg/smartcaptcha"
	"neomatica/neosync/types"
	"net/http"
	"strings"
	"time"

	"github.com/go-playground/validator"
)

/* Neosync HTTPx V1 */
/* Handler: вход пользователя. устанавливает access_token и, при remember_me, refresh_token в cookie */

func (h *Handler) AuthSigninHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)

	var payload *SigninPayload

	if err := httpx.HttpParse(r, &payload); err != nil {
		return httperr.BadRequest(err.Error())
	}

	if err := httpx.Validate.Struct(payload); err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			return httperr.BadRequest(httpx.ValidateMsg(tr, err))
		}

		return httperr.BadRequest(tr.TErr("fields-not-filled"))
	}

	/* Yandex SmartCaptcha */
	if err := smartcaptcha.VerifyRequest(tr, payload.TurnstileToken, r); err != nil {
		logger.Error("AuthSigninHandler_V1 req={%s}: Failed validation smartcaptcha token: %s", ctx.Value("XREQID").(string), err.Error())
		return httperr.BadRequest(err.Error())
	}

	user, err := h.Store.Users.Get_UserCoreByEmail(ctx, strings.TrimSpace(payload.Login))
	if err != nil {
		return httperr.Db(ctx, err)
	}

	if user == nil {
		return httperr.NotFound(tr.TErr(httperr.Err_UserNotFound.Error()))
	}

	userRole, err := h.Store.Users.Get_UserRoleByUuid(ctx, user.UserUUID)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	if userRole == nil {
		return httperr.NotFound(tr.TErr("user-not-found"))
	}

	/* Password func */
	time.Sleep(200 * time.Millisecond)

	passwordUserHashed, err := encryption.Decrypt(tr, strings.TrimSpace(user.Password))
	if err != nil {
		logger.Error("AuthSigninHandler_V1 req={%s}: Failed decrypt password: %s", ctx.Value("XREQID").(string), err.Error())
		return httperr.InternalServerError(err.Error())
	}

	if subtle.ConstantTimeCompare([]byte(strings.TrimSpace(payload.Password)), []byte(passwordUserHashed)) != 1 {
		time.Sleep(500 * time.Millisecond)
		return httperr.NotFound(tr.TErr("user-not-found"))
	}

	now := time.Now()

	sessionID, err := authtoken.NewID()
	if err != nil {
		logger.Error("AuthSigninHandler_V1: failed to create auth token: %s", err.Error())
		return httperr.InternalServerError(tr.TErr("auth-token-create-failed"))
	}

	ipHash, err := authtoken.IPHash(clientip.IP(r))
	if err != nil {
		logger.Error("AuthSigninHandler_V1: failed to create auth token: %s", err.Error())
		return httperr.InternalServerError(tr.TErr("auth-token-create-failed"))
	}

	session := &redis.AuthSession{
		ID:        sessionID,
		UserUUID:  user.UserUUID,
		IPHash:    ipHash,
		ExpiresAt: now.Add(authtoken.AccessTTL),
	}

	refresh := ""

	if payload.RememberMe {
		refresh, err = authtoken.NewRefresh(sessionID)
		if err != nil {
			logger.Error("AuthSigninHandler_V1: failed to create auth token: %s", err.Error())
			return httperr.InternalServerError(tr.TErr("auth-token-create-failed"))
		}

		session.RefreshHash = authtoken.RefreshHash(refresh)
		session.ExpiresAt = now.Add(authtoken.RefreshTTL)
	}

	access, err := authtoken.SealAccess(types.RequestAuthToken{
		UUID:      user.UserUUID,
		IPAddress: ipHash,
		SessionID: sessionID,
		RoleCode:  userRole.RoleCode,
		Timestamp: now,
	})
	if err != nil {
		logger.Error("AuthSigninHandler_V1: failed to create auth token: %s", err.Error())
		return httperr.InternalServerError(tr.TErr("auth-token-create-failed"))
	}

	/* Блокировка защищает от входа со старым паролем во время его сброса */
	if err := store.WithTx(h.Db, ctx, func(tx *sql.Tx) error {
		var password string

		query := `SELECT password FROM user_cores WHERE user_uuid = $1 FOR UPDATE`

		if err := tx.QueryRowContext(ctx, query, user.UserUUID).Scan(&password); err != nil {
			return err
		}

		if password != user.Password {
			return httperr.Unauthorized(tr.TErr("auth-password-changed"))
		}

		if err := redis.CreateAuthSession(ctx, session); err != nil {
			return httperr.Redis(ctx, err)
		}

		return nil
	}); err != nil {
		if _, ok := err.(httperr.HTTPError); ok {
			return err
		}

		return httperr.Db(ctx, err)
	}

	authtoken.SetCookie(w, "access_token", access, now.Add(authtoken.AccessTTL))

	if payload.RememberMe {
		authtoken.SetCookie(w, "refresh_token", refresh, session.ExpiresAt)
	} else {
		authtoken.ClearCookie(w, "refresh_token")
	}

	httpx.HttpResponse(w, r, http.StatusOK, map[string]interface{}{
		"login":     user.Email,
		"role_code": userRole.RoleCode,
	})
	return nil
}
