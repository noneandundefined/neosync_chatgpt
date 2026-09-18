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
/* Handler: список участников группы (только владелец) */

func (h *Handler) GroupMembersListHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*types.AuthToken)

	if !authToken.User.Has(constants.AccessGroupManage) {
		return httperr.Forbidden(tr.TErr("access-group-manage"))
	}

	groupIDStr := mux.Vars(r)["id"]
	if groupIDStr == "" {
		return httperr.NotFound(tr.TErr("group-id-not-found"))
	}

	groupID, err := strconv.ParseUint(groupIDStr, 10, 64)
	if err != nil {
		logger.Error("GroupMembersListHandler_V1 req={%s}: %s", ctx.Value("XREQID").(string), err.Error())
		return httperr.BadRequest(tr.TErr("invalid-group-id"))
	}

	group, err := h.Store.Groups.Get_GroupById(ctx, groupID)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	if group == nil {
		return httperr.NotFound(tr.TErr("group-not-found"))
	}

	if !h.groupIsOwner(group, authToken.User.UserContact.UserUUID) {
		return h.groupForbiddenFor(ctx, group, authToken.User.UserContact.UserUUID, groupPermOwnerOnly)
	}

	members, err := h.Store.GroupMembers.Get_GroupMembers(ctx, groupID)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	httpx.HttpResponseWithETag(w, r, http.StatusOK, members)
	return nil
}
