package group_handler_v1

import (
	"neomatica/neosync/infra/constants"
	"neomatica/neosync/infra/logger"
	"neomatica/neosync/infra/store/postgres/models"
	"neomatica/neosync/middleware"
	"neomatica/neosync/pkg/httpx"
	"neomatica/neosync/pkg/httpx/httperr"
	"neomatica/neosync/types"
	"net/http"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
)

/* Neosync HTTPx V1 */
/* Handler: обновление данных группы */

func (h *Handler) GroupUpdateByIdHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*types.AuthToken)

	/* Check access */
	if !authToken.User.Has(constants.AccessGroupManage) {
		return httperr.Forbidden(tr.TErr("access-group-manage"))
	}

	groupIdParam := mux.Vars(r)["id"]
	if groupIdParam == "" {
		return httperr.NotFound(tr.TErr("group-id-not-found"))
	}

	groupId, err := strconv.Atoi(groupIdParam)
	if err != nil {
		logger.Error("GroupUpdateByIdHandler_V1 req={%s}: %s", ctx.Value("XREQID").(string), err.Error())
		return httperr.BadRequest(tr.TErr("invalid-group-id"))
	}

	var payload *UpdateGroupPayload

	if err := httpx.HttpParse(r, &payload); err != nil {
		return httperr.BadRequest(err.Error())
	}

	if err := httpx.Validate.Struct(payload); err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			return httperr.BadRequest(httpx.ValidateMsg(tr, err))
		}

		return httperr.BadRequest(tr.TErr("fields-not-filled"))
	}

	/* The user has group access */
	group, err := h.Store.Groups.Get_GroupById(ctx, uint64(groupId))
	if err != nil {
		return httperr.Db(ctx, err)
	}

	if group == nil {
		return httperr.NotFound(tr.TErr("group-not-found"))
	}

	can, err := h.groupCanEdit(ctx, group, authToken.User.UserContact.UserUUID)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	if !can {
		return h.groupForbiddenFor(ctx, group, authToken.User.UserContact.UserUUID, groupPermEditGroup)
	}

	if !h.groupIsOwner(group, authToken.User.UserContact.UserUUID) && h.groupShareRightsChanged(payload, group) {
		return h.groupForbiddenFor(ctx, group, authToken.User.UserContact.UserUUID, groupPermOwnerOnly)
	}

	if len(payload.DevicesAdd) > 0 || len(payload.DevicesRemove) > 0 {
		canManage, err := h.groupCanManageDevices(ctx, group, authToken.User.UserContact.UserUUID)
		if err != nil {
			return httperr.Db(ctx, err)
		}

		if !canManage {
			return h.groupForbiddenFor(ctx, group, authToken.User.UserContact.UserUUID, groupPermManageDevices)
		}
	}

	groupUpd := &models.Group{
		Name:        payload.Name,
		Description: payload.Description,
	}

	if err := h.Store.Groups.Update_GroupById(ctx, uint64(groupId), groupUpd); err != nil {
		return httperr.Db(ctx, err)
	}

	if h.groupIsOwner(group, authToken.User.UserContact.UserUUID) {
		if err := h.Store.GroupMembers.Update_GroupShareDefaults(ctx, uint64(groupId), payload.CanEditGroup, payload.CanManageDevices, payload.CanReadConfig, payload.CanEditConfig, payload.CanSendCommands); err != nil {
			return httperr.Db(ctx, err)
		}
	}

	userUuid := authToken.User.UserContact.UserUUID

	if len(payload.DevicesRemove) > 0 {
		rm := &models.GroupDevices{
			UserUuid:  userUuid,
			GroupId:   uint64(groupId),
			DevicesId: payload.DevicesRemove,
		}

		if err := h.Store.Groups.Update_RemoveDevicesFromGroup(ctx, rm); err != nil {
			return httperr.Db(ctx, err)
		}
	}

	if len(payload.DevicesAdd) > 0 {
		add := &models.GroupDevices{
			UserUuid:  userUuid,
			GroupId:   uint64(groupId),
			DevicesId: payload.DevicesAdd,
		}

		if err := h.Store.Groups.Update_AddDevicesToGroup(ctx, add); err != nil {
			return httperr.Db(ctx, err)
		}
	}

	httpx.HttpResponse(w, r, http.StatusOK, tr.T("group-data-updated"))
	return nil
}
