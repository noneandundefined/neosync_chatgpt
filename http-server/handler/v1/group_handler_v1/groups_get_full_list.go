package group_handler_v1

import (
	"neomatica/neosync/pkg/httpx"
	"neomatica/neosync/pkg/httpx/httperr"
	"net/http"
)

/* Neosync HTTPx V1 */
/* Handler: получение ВСЕХ групп у пользователя */

func (h *Handler) GetGroupsFullHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	groups, err := h.Store.Groups.Get_Groups(ctx)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	httpx.HttpResponseWithETag(w, r, http.StatusOK, groups)
	return nil
}
