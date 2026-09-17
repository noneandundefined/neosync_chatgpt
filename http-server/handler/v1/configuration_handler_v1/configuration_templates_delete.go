package configuration_handler_v1

import (
	"database/sql"
	"errors"
	"neomatica/neosync/infra/constants"
	"neomatica/neosync/middleware"
	"neomatica/neosync/pkg/httpx"
	"neomatica/neosync/pkg/httpx/httperr"
	gtypes "neomatica/neosync/types"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

/* Neosync HTTPx V1 */
/* Handler: удаление шаблона конфигурации */

func (h *Handler) DeleteConfigurationTemplateHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*gtypes.AuthToken)

	if !authToken.User.Has(constants.AccessConfigurationApply) {
		return httperr.Forbidden(tr.TErr("config-edit-access-restricted"))
	}

	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id == 0 {
		return httperr.BadRequest(tr.TErr("fields-not-filled"))
	}

	if err := h.Store.ConfigurationTemplates.Delete_ConfigurationTemplateByIdAndUserUuid(ctx, id, authToken.User.UserContact.UserUUID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return httperr.NotFound(tr.TErr("configuration-template-not-found"))
		}

		return httperr.Db(ctx, err)
	}

	httpx.HttpResponse(w, r, http.StatusOK, tr.T("configuration-template-deleted"))
	return nil
}
