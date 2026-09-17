package configuration_handler_v1

import (
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
/* Handler: получение шаблона конфигурации по id */

func (h *Handler) GetConfigurationTemplateByIdHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*gtypes.AuthToken)

	if !authToken.User.Has(constants.AccessConfigurationRead) {
		return httperr.Forbidden(tr.TErr("config-access-restricted"))
	}

	id, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil || id == 0 {
		return httperr.BadRequest(tr.TErr("fields-not-filled"))
	}

	template, err := h.Store.ConfigurationTemplates.Get_ConfigurationTemplateByIdAndUserUuid(ctx, id, authToken.User.UserContact.UserUUID)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	if template == nil {
		return httperr.NotFound(tr.TErr("configuration-template-not-found"))
	}

	httpx.HttpResponse(w, r, http.StatusOK, template)
	return nil
}
