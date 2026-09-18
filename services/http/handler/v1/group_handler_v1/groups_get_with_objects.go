package group_handler_v1

import (
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
/* Handler: получение группы у пользователя */

func (h *Handler) GetGroupWithDevicesHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*types.AuthToken)

	groupIdParam := mux.Vars(r)["id"]
	if groupIdParam == "" {
		return httperr.NotFound(tr.TErr("group-id-not-found"))
	}

	groupId, err := strconv.Atoi(groupIdParam)
	if err != nil {
		logger.Error("GetGroupWithObjectsHandler_V1 req={%s}: %s", ctx.Value("XREQID").(string), err.Error())
		return httperr.BadRequest(tr.TErr("invalid-group-id"))
	}

	canAccess, err := h.Store.GroupMembers.Get_UserHasGroupAccess(ctx, uint64(groupId), authToken.User.UserContact.UserUUID)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	if !canAccess {
		return httperr.Forbidden(tr.TErr("group-access-denied"))
	}

	group, err := h.Store.Groups.Get_GroupWithDevices(ctx, uint64(groupId), authToken)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	if group == nil {
		httpx.HttpResponse(w, r, http.StatusNoContent, nil)
		return nil
	}

	httpx.HttpResponseWithETag(w, r, http.StatusOK, group)
	return nil
}
