package group_handler_v1

import (
	"neomatica/neosync/pkg/httpx"
	"neomatica/neosync/pkg/httpx/httperr"
	"neomatica/neosync/types"
	"net/http"
)

/* Neosync HTTPx V1 */
/* Handler: получение минимальной информации группы у пользователя */

func (h *Handler) GetGroupsFieldsNameHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	authToken := ctx.Value("identity").(*types.AuthToken)

	groups, err := h.Store.Groups.Get_GroupNamesByUuid(ctx, authToken.User.UserContact.UserUUID)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	httpx.HttpResponseWithETag(w, r, http.StatusOK, groups)
	return nil
}
