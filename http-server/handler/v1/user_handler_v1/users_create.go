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
	"github.com/google/uuid"
)

/* Neosync HTTPx V1 */
/* Handler: создание пользователя */

func (h *Handler) UserCreateUserHandler_V1(w http.ResponseWriter, r *http.Request) error { //nolint
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*types.AuthToken)

	var payload *UserCreatePayload

	if err := httpx.HttpParse(r, &payload); err != nil {
		return httperr.BadRequest(err.Error())
	}

	if err := httpx.Validate.Struct(payload); err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			return httperr.BadRequest(httpx.ValidateMsg(tr, err))
		}

		return httperr.BadRequest(tr.TErr("fields-not-filled"))
	}

	roleUniqueIndex, can := constants.CanCreateRole(authToken.RoleCode, payload.RoleCode)
	if !can {
		return httperr.Forbidden(tr.TErr("insufficient-permissions-create-user"))
	}

	if payload.NameOrganization != "" && len(payload.NameOrganization) < 3 {
		return httperr.BadRequest(tr.TErr("organization-name-too-short"))
	}

	if phone := strings.TrimSpace(payload.Phone); phone != "" {
		phone = strings.TrimPrefix(phone, "+")
		phone = strings.ReplaceAll(phone, " ", "")
		if len(phone) < 8 {
			return httperr.BadRequest(tr.TErr("invalid-phone-number"))
		}

		payload.Phone = phone
	}

	password := strings.TrimSpace(payload.Password)
	if _, chPass := constants.CheckSimplePasswords[strings.ToLower(password)]; chPass {
		return httperr.BadRequest(tr.TErr("simple-password"))
	}

	passwordEncrypted, err := encryption.Encrypt(tr, password)
	if err != nil {
		logger.Error("UserCreateUserHandler_V1 req={%s}: Failed encrypt password: %s", ctx.Value("XREQID").(string), err.Error())
		return httperr.InternalServerError(err.Error())
	}

	tx, err := h.Db.BeginTx(ctx, nil)
	if err != nil {
		return httperr.Db(ctx, httperr.Err_DbNetwork)
	}

	defer func() {
		_ = tx.Rollback()
	}()

	uuid := uuid.NewString()

	var parentUuid string = authToken.User.UserContact.UserUUID
	if authToken.RoleCode == constants.Role_AdminL2Support {
		if authToken.User.ParentUUID != nil {
			parentUuid = *authToken.User.ParentUUID
		}
	}

	userCore := &models.UserCore{
		UserUUID:   uuid,
		ParentUUID: &parentUuid,
		Email:      payload.Email,
		Password:   passwordEncrypted,
	}

	if err := h.Store.Users.Create_UserCore(ctx, tx, userCore); err != nil {
		return httperr.Db(ctx, err)
	}

	userContact := &models.UserContact{
		UserUUID:         uuid,
		NameOrganization: &payload.NameOrganization,
		Locality:         &payload.Locality,
		Phone:            &payload.Phone,
		Language:         payload.Language,
	}

	if err := h.Store.Users.Create_UserContact(ctx, tx, userContact); err != nil {
		return httperr.Db(ctx, err)
	}

	userRoleModel := &models.UserRole{
		UserUUID: uuid,
		RoleCode: roleUniqueIndex,
	}

	if err := h.Store.Users.Create_UserRole(ctx, tx, userRoleModel); err != nil {
		return httperr.Db(ctx, err)
	}

	userAccess := &models.UserAccess{
		UserUUID: uuid,
		AccessFields: models.AccessFields{
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

	if err := h.Store.Users.Create_UserAccess(ctx, tx, userAccess); err != nil {
		return httperr.Db(ctx, err)
	}

	if err := tx.Commit(); err != nil {
		return httperr.Conflict(tr.TErr("failed-to-save-data"))
	}

	httpx.HttpResponse(w, r, http.StatusCreated, tr.T("user-created-successfully"))
	return nil
}
