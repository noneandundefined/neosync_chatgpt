package user_handler_v1

import (
	"fmt"
	"html"
	"neomatica/neosync/encryption"
	"neomatica/neosync/infra/locale"
	"neomatica/neosync/infra/logger"
	"neomatica/neosync/middleware"
	"neomatica/neosync/pkg"
	"neomatica/neosync/pkg/httpx"
	"neomatica/neosync/pkg/httpx/httperr"
	"neomatica/neosync/pkg/permissions"
	"neomatica/neosync/pkg/smartcaptcha"
	"neomatica/neosync/types"
	"net/http"
	"os"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
)

/* Neosync HTTPx V1 */
/* Handler: отправка login, password на почту пользователя */

func normalizeEmailLanguage(language string) string {
	switch language {
	case "ru", "en", "es":
		return language
	default:
		return "en"
	}
}

func (h *Handler) UserSendEmailInfoHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*types.AuthToken)

	var payload *SendEmailInfoPayload

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
		logger.Error("UserSendEmailInfoHandler_V1 req={%s}: Failed validation smartcaptcha token: %s", ctx.Value("XREQID").(string), err.Error())
		return httperr.BadRequest(err.Error())
	}

	uuid := mux.Vars(r)["uuid"]
	if uuid == "" {
		return httperr.NotFound(tr.TErr("user-not-found"))
	}

	user, err := h.Store.Users.Get_UserCoreByUuid(ctx, uuid)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	if user == nil {
		return httperr.NotFound(tr.TErr("user-not-found"))
	}

	/* if this is adminL2 and user not created current adminL2 -> forbidden */
	if !permissions.IsMainRole(authToken.RoleCode) {
		if user.ParentUUID == nil || *user.ParentUUID != authToken.User.UserContact.UserUUID {
			return httperr.Forbidden(tr.TErr("user-not-owned"))
		}
	}

	emailLanguage := normalizeEmailLanguage(payload.Language)
	emailTr := locale.NewTranslator(emailLanguage)

	password, err := encryption.Decrypt(tr, user.Password)
	if err != nil {
		logger.Error("UserSendEmailInfoHandler_V1 req={%s}: Failed decrypt password: %s", ctx.Value("XREQID").(string), err.Error())
		return httperr.InternalServerError(err.Error())
	}

	go func() {
		body := fmt.Sprintf(`<p>%s:</p><div style="border-left:4px solid #395d95; padding:15px; margin:20px 0; color:#49525f;"><strong>%s:</strong> <span>%s</span><br><strong>%s:</strong> <span>%s</span></div><p>%s</p><a href="%s/sign_in" style="display:flex;">%s</a>`,
			emailTr.T("account-created-email"), emailTr.T("login"), html.EscapeString(user.Email),
			emailTr.T("password"), html.EscapeString(password), emailTr.T("login-instructions"),
			html.EscapeString(os.Getenv("CLIENT_URL")), emailTr.T("go-to-login"))
		if err := pkg.SendEmail(user.Email, emailTr.T("login-data-subject"), body, emailTr); err != nil {
			logger.Error("UserSendEmailInfoHandler_V1: Failed to send login data email")
		}
	}()

	httpx.HttpResponse(w, r, http.StatusOK, tr.T("login-data-sent-to-user"))
	return nil
}
