package group_handler_v1

import (
	"neomatica/neosync/infra/constants"
	"neomatica/neosync/infra/logger"
	"neomatica/neosync/middleware"
	"neomatica/neosync/pkg/httpx"
	"neomatica/neosync/pkg/httpx/httperr"
	"neomatica/neosync/types"
	"net/http"

	"github.com/go-playground/validator"
)

/* Neosync HTTPx V1 */
/* Handler: удаление массива групп у пользователя */

func (h *Handler) GroupDeleteMassiveHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*types.AuthToken)

	/* Check access */
	if !authToken.User.Has(constants.AccessGroupManage) {
		return httperr.Forbidden(tr.TErr("access-group-manage"))
	}

	var payload *DeleteGroupsPayload

	if err := httpx.HttpParse(r, &payload); err != nil {
		return httperr.BadRequest(err.Error())
	}

	if err := httpx.Validate.Struct(payload); err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			return httperr.BadRequest(httpx.ValidateMsg(tr, err))
		}

		return httperr.BadRequest(tr.TErr("fields-not-filled"))
	}

	for _, id := range payload.IDs {
		/* The user has group access */
		group, err := h.Store.Groups.Get_GroupById(ctx, id)
		if err != nil {
			logger.Error("GroupDeleteMassiveHandler_V1 req={%s}: %s", ctx.Value("XREQID").(string), err.Error())
			return httperr.Db(ctx, err)
		}

		if group == nil {
			return httperr.NotFound(tr.TErr("group-not-found"))
		}

		if group.UserUuid != authToken.User.UserContact.UserUUID {
			return h.groupForbiddenFor(ctx, group, authToken.User.UserContact.UserUUID, groupPermDelete)
		}
	}

	if err := h.Store.Groups.Delete_GroupsById(ctx, payload.IDs); err != nil {
		logger.Error("GroupDeleteMassiveHandler_V1 req={%s}: %s", ctx.Value("XREQID").(string), err.Error())
		return httperr.Db(ctx, err)
	}

	httpx.HttpResponse(w, r, http.StatusOK, tr.T("groups-deleted"))
	return nil
}
