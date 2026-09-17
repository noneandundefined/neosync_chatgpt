package group_handler_v1

import (
	"neomatica/neosync/infra/constants"
	"neomatica/neosync/infra/logger"
	"neomatica/neosync/middleware"
	"neomatica/neosync/pkg/httpx"
	"neomatica/neosync/pkg/httpx/httperr"
	"neomatica/neosync/types"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

/* Neosync HTTPx V1 */
/* Handler: удаление группы по id */

func (h *Handler) GroupDeleteByIdHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*types.AuthToken)

	/* Check access */
	if !authToken.User.Has(constants.AccessGroupManage) {
		return httperr.Forbidden(tr.TErr("access-group-manage"))
	}

	groupIdStr := mux.Vars(r)["id"]
	if groupIdStr == "" {
		return httperr.NotFound(tr.TErr("group-id-not-found"))
	}

	groupId, err := strconv.Atoi(groupIdStr)
	if err != nil {
		logger.Error("GroupDeleteByIdHandler_V1 req={%s}: %s", ctx.Value("XREQID").(string), err.Error())
		return httperr.BadRequest(tr.TErr("invalid-group-id"))
	}

	/* The user has group access */
	group, err := h.Store.Groups.Get_GroupById(ctx, uint64(groupId))
	if err != nil {
		return httperr.Db(ctx, err)
	}

	if group == nil {
		return httperr.NotFound(tr.TErr("group-not-found"))
	}

	userOwnerGroup, err := h.Store.Users.Get_UserCoreByUuid(ctx, group.UserUuid)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	currentUser := authToken.User.UserContact.UserUUID

	isDealerDeletingChild := false
	isOwner := group.UserUuid == currentUser

	if authToken.RoleCode == constants.Role_AdminL2 {
		if userOwnerGroup.ParentUUID != nil && *userOwnerGroup.ParentUUID == currentUser {
			isDealerDeletingChild = true
		}
	}
	if !isOwner && !isDealerDeletingChild {
		return h.groupForbiddenFor(ctx, group, currentUser, groupPermDelete)
	}

	if err := h.Store.Groups.Delete_GroupById(ctx, group.ID); err != nil {
		logger.Error("GroupDeleteByIdHandler_V1 req={%s}: %s", ctx.Value("XREQID").(string), err.Error())
		return httperr.Db(ctx, err)
	}

	httpx.HttpResponse(w, r, http.StatusOK, tr.T("group-deleted-successfully"))
	return nil
}
