package group_handler_v1

import (
	"neomatica/neosync/infra/constants"
	"neomatica/neosync/infra/logger"
	"neomatica/neosync/infra/store/postgres/models"
	"neomatica/neosync/middleware"
	"neomatica/neosync/pkg/httpx"
	"neomatica/neosync/pkg/httpx/httperr"
	"neomatica/neosync/pkg/smartcaptcha"
	"neomatica/neosync/types"
	"net/http"

	"github.com/go-playground/validator"
)

/* Neosync HTTPx V1 */
/* Handler: создание группы */

func (h *Handler) GroupCreateHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*types.AuthToken)

	/* Check access */
	if !authToken.User.Has(constants.AccessGroupManage) {
		return httperr.Forbidden(tr.TErr("access-group-manage"))
	}

	var payload *CreateGroupPayload

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
		logger.Error("GroupCreateHandler_V1 req={%s}: Failed validation smartcaptcha token: %s", ctx.Value("XREQID").(string), err.Error())
		return httperr.BadRequest(err.Error())
	}

	if len(payload.Name) > 30 {
		return httperr.BadRequest(tr.TErr("group-name-max-30"))
	}

	if payload.Description != "" && len(payload.Description) > 300 {
		return httperr.BadRequest(tr.TErr("group-description-max-300"))
	}

	groupModel := &models.Group{
		UserUuid:         authToken.User.UserContact.UserUUID,
		Name:             payload.Name,
		Description:      payload.Description,
		CanEditGroup:     payload.CanEditGroup,
		CanManageDevices: payload.CanManageDevices,
		CanReadConfig:    payload.CanReadConfig,
		CanEditConfig:    payload.CanEditConfig,
		CanSendCommands:  payload.CanSendCommands,
	}

	if err := h.Store.Groups.Create_Group(ctx, groupModel); err != nil {
		return httperr.Db(ctx, err)
	}

	httpx.HttpResponse(w, r, http.StatusCreated, tr.T("group-created-successfully"))
	return nil
}
