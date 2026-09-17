package device_handler_v1

import (
	"neomatica/neosync/middleware"
	"neomatica/neosync/pkg/httpx"
	"neomatica/neosync/pkg/httpx/httperr"
	"neomatica/neosync/types"
	"neomatica/neosync/util"
	"net/http"
	"strconv"
	"strings"
)

/* Neosync HTTPx V1 */
/* Handler: получение IMEI группы для отправки команд */
/* Параметры - ?group_id=id ?search=поиск ?page=страница ?limit=размер ?all=true */

func (h *Handler) GetDevicesStateDevicesHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*types.AuthToken)

	groupIDStr := r.URL.Query().Get("group_id")
	if groupIDStr == "" {
		return httperr.BadRequest(tr.TErr("fields-not-filled"))
	}

	groupID, err := strconv.ParseUint(groupIDStr, 10, 64)
	if err != nil {
		return httperr.BadRequest(tr.TErr("fields-not-filled"))
	}

	search := strings.TrimSpace(r.URL.Query().Get("search"))
	all := r.URL.Query().Get("all") == "true"
	pagination := util.GetPagination(r, 25)

	if pagination.Limit > 100 {
		pagination.Limit = 100
		pagination.Offset = (pagination.Page - 1) * pagination.Limit
	}

	devices, total, err := h.loadDevicesStateDevices(ctx, authToken, tr, groupID, search, pagination.Limit, pagination.Offset, all)
	if err != nil {
		return err
	}

	limit := pagination.Limit
	if all {
		limit = total
	}

	loaded := pagination.Offset + len(devices)
	if all {
		loaded = total
	}

	httpx.HttpResponseWithETag(w, r, http.StatusOK, DevicesStateDevicesResponse{
		Items:   devices,
		Page:    pagination.Page,
		Limit:   limit,
		Total:   total,
		HasMore: loaded < total,
	})
	return nil
}
