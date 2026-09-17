package analytic_handler_v1

import (
	"neomatica/neosync/infra/constants"
	"neomatica/neosync/middleware"
	"neomatica/neosync/pkg/httpx"
	"neomatica/neosync/pkg/httpx/httperr"
	"net/http"

	"github.com/go-playground/validator"
)

/* Neosync HTTPx V1 */
/* Handler: выполнение query запросов для аналитики */

func (h *Handler) AnalyticQueryExecHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)

	var payload *QueryPayload

	if err := httpx.HttpParse(r, &payload); err != nil {
		return httperr.BadRequest(err.Error())
	}

	if err := httpx.Validate.Struct(payload); err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			return httperr.BadRequest(httpx.ValidateMsg(tr, err))
		}

		return httperr.BadRequest(tr.TErr("fields-not-filled"))
	}

	if payload.Query == "" {
		return httperr.BadRequest(tr.TErr("fields-not-filled"))
	}

	var groupID *int
	if payload.GroupID != nil && *payload.GroupID > 0 {
		groupID = payload.GroupID
	}

	query := constants.PrepareAnalyticQuery(payload.Type, payload.Query, payload.From, payload.To, payload.UserUUID, payload.Model, groupID)

	result, err := h.Store.Analytics.Exec_SelectQuery(ctx, query)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	httpx.HttpResponseWithETag(w, r, http.StatusOK, result)
	return nil
}
