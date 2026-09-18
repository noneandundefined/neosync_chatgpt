package device_handler_v1

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"neomatica/neosync/infra/constants"
	"neomatica/neosync/infra/store/redis"
	"neomatica/neosync/middleware"
	"neomatica/neosync/pkg/httpx"
	"neomatica/neosync/pkg/httpx/httperr"
	"neomatica/neosync/pkg/permissions"
	"neomatica/neosync/types"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func (h *Handler) ExportDeviceLogsHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*types.AuthToken)

	if !authToken.User.Has(constants.AccessLogRead) {
		return httperr.Forbidden(tr.TErr("log-access-restricted"))
	}

	imei := mux.Vars(r)["imei"]
	if len(imei) != 15 {
		return httperr.BadRequest(tr.TErr("invalid-or-nonexistent-imei"))
	}

	device, err := h.Store.Devices.Get_DeviceAndConfByImei(ctx, imei)
	if err != nil {
		return httperr.Db(ctx, err)
	}
	if device == nil {
		return httperr.NotFound(tr.TErr("device-not-found"))
	}

	if !permissions.CanAccessDevice(ctx, h.Store.GroupMembers, device, device.ID, &authToken.User, authToken.RoleCode, permissions.ShareViewDevice) {
		return httperr.Forbidden(tr.TErr("device-not-owned"))
	}

	logs, err := redis.ReadAdmLog(imei)
	if err != nil {
		return httperr.Redis(ctx, err)
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	if err := writer.Write([]string{"time_utc", "timestamp", "event", "payload"}); err != nil {
		return httperr.InternalServerError(err.Error())
	}

	for _, item := range logs {
		payload := ""
		if item.Payload != nil {
			if encoded, marshalErr := json.Marshal(item.Payload); marshalErr == nil {
				payload = string(encoded)
			} else {
				payload = fmt.Sprint(item.Payload)
			}
		}

		if err := writer.Write([]string{
			time.Unix(item.TS, 0).UTC().Format(time.RFC3339),
			strconv.FormatInt(item.TS, 10),
			item.Event,
			payload,
		}); err != nil {
			return httperr.InternalServerError(err.Error())
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return httperr.InternalServerError(err.Error())
	}

	httpx.HttpFileResponse(w, r, fmt.Sprintf("device-%s-events.csv", imei), buf.Bytes(), "text/csv; charset=utf-8")
	return nil
}
