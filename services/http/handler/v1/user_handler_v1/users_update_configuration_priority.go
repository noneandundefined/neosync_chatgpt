package user_handler_v1

import (
	"neomatica/neosync/middleware"
	"neomatica/neosync/pkg/httpx"
	"neomatica/neosync/pkg/httpx/httperr"
	"neomatica/neosync/types"
	"net/http"

	"github.com/go-playground/validator"
)

func (h *Handler) UserUpdateConfigurationPriorityHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*types.AuthToken)

	var payload *UpdateConfigurationPriority

	if err := httpx.HttpParse(r, &payload); err != nil {
		return httperr.BadRequest(err.Error())
	}

	if err := httpx.Validate.Struct(payload); err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			return httperr.BadRequest(httpx.ValidateMsg(tr, err))
		}

		return httperr.BadRequest(tr.TErr("fields-not-filled"))
	}

	if err := h.Store.Users.Update_ForceNeosyncConfigurationPriority(ctx, authToken.User.UserContact.UserUUID, *payload.Enabled); err != nil {
		return httperr.Db(ctx, err)
	}

	httpx.HttpResponse(w, r, http.StatusNoContent, nil)
	return nil
}
