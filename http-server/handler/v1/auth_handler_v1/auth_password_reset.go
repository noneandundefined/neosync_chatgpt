package auth_handler_v1

import (
	"database/sql"
	"fmt"
	"neomatica/neosync/encryption"
	"neomatica/neosync/infra/constants"
	"neomatica/neosync/infra/logger"
	"neomatica/neosync/infra/store/postgres/store"
	"neomatica/neosync/infra/store/redis"
	"neomatica/neosync/middleware"
	"neomatica/neosync/pkg"
	"neomatica/neosync/pkg/authtoken"
	"neomatica/neosync/pkg/httpx"
	"neomatica/neosync/pkg/httpx/httperr"
	"neomatica/neosync/pkg/smartcaptcha"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator"
)

/* Neosync HTTPx V1 */
/* Handler: отправляет ссылку для сброса пароля на указанный email */

func (h *Handler) AuthRequestPasswordResetHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)

	email := strings.TrimSpace(r.URL.Query().Get("email"))
	if email != "" {
		user, err := h.Store.Users.Get_UserCoreByEmail(ctx, email)
		if err != nil {
			return httperr.Db(ctx, err)
		}

		if user != nil {
			expiresAt := time.Now().Add(24 * time.Hour)
			exp := expiresAt.Unix()

			sig, err := authtoken.NewID()
			if err != nil {
				logger.Error("AuthRequestPasswordResetHandler_V1: failed to create reset token: %s", err.Error())
				return httperr.InternalServerError(tr.TErr("auth-token-create-failed"))
			}

			if err := redis.SetPasswordResetToken(ctx, user.UserUUID, sig, expiresAt); err != nil {
				return httperr.Redis(ctx, err)
			}

			link := fmt.Sprintf("%s/password/new?uuid=%s&exp=%d&sig=%s", os.Getenv("CLIENT_URL"), user.UserUUID, exp, sig)

			logger.Info("AuthRequestPasswordResetHandler_V1 req={%s} email={%s}: Send password reset email", ctx.Value("XREQID").(string), user.Email)

			go func() {
				if err := pkg.SendEmail(user.Email, tr.T("password-reset-data-subject"), fmt.Sprintf(`
		<p style="margin-top:0;">%s</p>

		<p>%s:<br>%s</p>

		<p>%s</p>
		`, tr.T("password-reset-request-received"), tr.T("go-to-link-reset-password"), link, tr.T("exp-limit-reset-password")), tr); err != nil {
					logger.Error("AuthRequestPasswordResetHandler_V1 req={%s} email={%s}: Failed sent reset link to email: %s", ctx.Value("XREQID").(string), user.Email, err.Error())
					return
				}
			}()
		}
	}

	httpx.HttpResponse(w, r, http.StatusOK, tr.T("password-reset-sent-to-email"))
	return nil
}

/* Neosync HTTPx V1 */
/* Handler: сбрасывает пароль пользователя по ссылке с uuid, exp и sig. */

func (h *Handler) AuthPasswordResetHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)

	uuid := r.URL.Query().Get("uuid")
	expStr := r.URL.Query().Get("exp")
	sig := r.URL.Query().Get("sig")

	if uuid == "" || expStr == "" || sig == "" {
		return httperr.BadRequest(tr.TErr("incorrect-password-reset-link"))
	}

	exp, err := strconv.ParseInt(expStr, 10, 64)
	if err != nil {
		return httperr.BadRequest(tr.TErr("incorrect-password-reset-link"))
	}

	var payload *PasswordResetPayload

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
		logger.Error("AuthPasswordResetHandler req={%s}: Failed validation smartcaptcha token: %s", ctx.Value("XREQID").(string), err.Error())
		return httperr.BadRequest(err.Error())
	}

	password := strings.TrimSpace(payload.Password)
	if _, chPass := constants.CheckSimplePasswords[strings.ToLower(password)]; chPass {
		return httperr.BadRequest(tr.TErr("simple-password"))
	}

	user, err := h.Store.Users.Get_UserByUuid(ctx, uuid)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	if user == nil {
		return httperr.NotFound(tr.TErr(httperr.Err_UserNotFound.Error()))
	}

	if time.Now().Unix() >= exp {
		return httperr.BadRequest(tr.TErr("incorrect-password-reset-link"))
	}

	passwordHashed, err := encryption.Encrypt(tr, password)
	if err != nil {
		logger.Error("AuthPasswordResetHandler req={%s}: Failed encrypt password: %s", ctx.Value("XREQID").(string), err.Error())
		return httperr.InternalServerError(err.Error())
	}

	/* Блокируем пользователя, обновляем пароль и однократно погашаем ссылку */
	if err := store.WithTx(h.Db, ctx, func(tx *sql.Tx) error {
		var currentPassword string

		query := `SELECT password FROM user_cores WHERE user_uuid = $1 FOR UPDATE`

		if err := tx.QueryRowContext(ctx, query, uuid).Scan(&currentPassword); err != nil {
			return err
		}

		query = `UPDATE user_cores SET password = $1, refresh_token = NULL WHERE user_uuid = $2`

		if _, err := tx.ExecContext(ctx, query, passwordHashed, uuid); err != nil {
			return err
		}

		/* Ошибка или использованная ссылка откатывают изменение пароля */
		consumed, err := redis.ConsumePasswordResetToken(ctx, uuid, sig)
		if err != nil {
			return httperr.Redis(ctx, err)
		}

		if !consumed {
			return httperr.BadRequest(tr.TErr("incorrect-password-reset-link"))
		}

		return nil
	}); err != nil {
		if _, ok := err.(httperr.HTTPError); ok {
			return err
		}

		return httperr.Db(ctx, err)
	}

	httpx.HttpResponse(w, r, http.StatusOK, tr.T("success-password-reset"))
	return nil
}
