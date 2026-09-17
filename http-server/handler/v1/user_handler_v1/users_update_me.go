package user_handler_v1

import (
	"neomatica/neosync/encryption"
	"neomatica/neosync/infra/constants"
	"neomatica/neosync/infra/logger"
	"neomatica/neosync/infra/store/postgres/models"
	"neomatica/neosync/middleware"
	"neomatica/neosync/pkg/httpx"
	"neomatica/neosync/pkg/httpx/httperr"
	"neomatica/neosync/types"
	"net/http"
	"strings"

	"github.com/go-playground/validator"
)

func ptr[T any](v T) *T { return &v }

/* Neosync HTTPx V1 */
/* Handler: обновление профиля пользователя */

func (h *Handler) UserUpdateMeHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*types.AuthToken)

	var payload *UpdateUserMePayload

	if err := httpx.HttpParse(r, &payload); err != nil {
		return httperr.BadRequest(err.Error())
	}

	if err := httpx.Validate.Struct(payload); err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			return httperr.BadRequest(httpx.ValidateMsg(tr, err))
		}

		return httperr.BadRequest(tr.TErr("fields-not-filled"))
	}

	var emailPtr, passwordPtr *string

	if payload.Email != nil && strings.TrimSpace(*payload.Email) != "" {
		email := strings.ToLower(strings.TrimSpace(*payload.Email))
		emailPtr = &email
	}

	if payload.Password != nil && strings.TrimSpace(*payload.Password) != "" {
		password := strings.TrimSpace(*payload.Password)

		if _, chPass := constants.CheckSimplePasswords[strings.ToLower(password)]; chPass {
			return httperr.BadRequest(tr.TErr("simple-password"))
		}

		pass, err := encryption.Encrypt(tr, password)
		if err != nil {
			logger.Error("UserUpdateMeHandler_V1 req={%s}: %s", ctx.Value("XREQID").(string), err.Error())
			return httperr.InternalServerError(err.Error())
		}

		passwordPtr = &pass
	}

	if payload.Phone != nil {
		phone := strings.TrimSpace(*payload.Phone)

		if phone == "" {
			payload.Phone = ptr("")
		} else {
			phone = strings.TrimPrefix(phone, "+")
			phone = strings.ReplaceAll(phone, " ", "")

			if len(phone) < 8 {
				return httperr.BadRequest(tr.TErr("invalid-phone-number"))
			}

			payload.Phone = &phone
		}
	}

	tx, err := h.Db.BeginTx(ctx, nil)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	defer func() {
		_ = tx.Rollback()
	}()

	UserUpdate := &models.UserUpdate{
		Email:            emailPtr,
		Password:         passwordPtr,
		Phone:            payload.Phone,
		NameOrganization: payload.NameOrganization,
		Locality:         payload.Locality,
		Language:         payload.Language,
	}

	if err := h.Store.Users.Update_UserCoreOnesByUuid(ctx, tx, authToken.User.UserContact.UserUUID, UserUpdate); err != nil {
		return httperr.Db(ctx, err)
	}

	if err := h.Store.Users.Update_UserContactOnesByUuid(ctx, tx, authToken.User.UserContact.UserUUID, UserUpdate); err != nil {
		return httperr.Db(ctx, err)
	}

	if err := tx.Commit(); err != nil {
		return httperr.Db(ctx, err)
	}

	httpx.HttpResponse(w, r, http.StatusOK, tr.T("profile-updated"))
	return nil
}
