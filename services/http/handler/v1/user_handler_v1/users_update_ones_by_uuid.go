package user_handler_v1

import (
	"neomatica/neosync/encryption"
	"neomatica/neosync/infra/constants"
	"neomatica/neosync/infra/logger"
	"neomatica/neosync/infra/store/postgres/models"
	"neomatica/neosync/middleware"
	"neomatica/neosync/pkg/httpx"
	"neomatica/neosync/pkg/httpx/httperr"
	"neomatica/neosync/pkg/permissions"
	"neomatica/neosync/types"
	"net/http"
	"strings"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
)

/* Neosync HTTPx V1 */
/* Handler: обновление данных у пользователя */

func (h *Handler) UserUpdateOnesByUuidHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*types.AuthToken)

	uuid := mux.Vars(r)["uuid"]
	if uuid == "" {
		return httperr.NotFound(tr.TErr("user-not-found"))
	}

	var payload *UpdateUserOnesPayload

	if err := httpx.HttpParse(r, &payload); err != nil {
		return httperr.BadRequest(err.Error())
	}

	if err := httpx.Validate.Struct(payload); err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			return httperr.BadRequest(httpx.ValidateMsg(tr, err))
		}

		return httperr.BadRequest(tr.TErr("fields-not-filled"))
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
			logger.Error("UserUpdateOnesByUuidHandler_V1 req={%s}: %s", ctx.Value("XREQID").(string), err.Error())
			return httperr.InternalServerError(err.Error())
		}

		passwordPtr = &pass
	}

	if payload.Phone != nil {
		if phone := strings.TrimSpace(*payload.Phone); phone != "" {
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
		Accesses: models.PtrAccessFields{
			AccessTrekerCreate:         payload.Accesses.AccessTrekerCreate,
			AccessTrekerEdit:           payload.Accesses.AccessTrekerEdit,
			AccessTrekerDelete:         payload.Accesses.AccessTrekerDelete,
			AccessGroupManage:          payload.Accesses.AccessGroupManage,
			AccessConfigurationRead:    payload.Accesses.AccessConfigurationRead,
			AccessConfigurationApply:   payload.Accesses.AccessConfigurationApply,
			AccessConfigurationHistory: payload.Accesses.AccessConfigurationHistory,
			AccessCommandSend:          payload.Accesses.AccessCommandSend,
			AccessLogRead:              payload.Accesses.AccessLogRead,
		},
	}

	if err := h.Store.Users.Update_UserCoreOnesByUuid(ctx, tx, uuid, UserUpdate); err != nil {
		return httperr.Db(ctx, err)
	}

	if err := h.Store.Users.Update_UserContactOnesByUuid(ctx, tx, uuid, UserUpdate); err != nil {
		return httperr.Db(ctx, err)
	}

	if err := h.Store.Users.Update_UserAccessOnesByUuid(ctx, tx, uuid, UserUpdate); err != nil {
		return httperr.Db(ctx, err)
	}

	if err := tx.Commit(); err != nil {
		return httperr.Db(ctx, err)
	}

	httpx.HttpResponse(w, r, http.StatusOK, tr.T("user-data-updated"))
	return nil
}
