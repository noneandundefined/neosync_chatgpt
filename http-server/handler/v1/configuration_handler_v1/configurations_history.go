package configuration_handler_v1

import (
	"encoding/hex"
	"fmt"
	"neomatica/neosync/infra/constants"
	"neomatica/neosync/middleware"
	"neomatica/neosync/pkg/adm"
	"neomatica/neosync/pkg/httpx"
	"neomatica/neosync/pkg/httpx/httperr"
	"neomatica/neosync/pkg/permissions"
	gtypes "neomatica/neosync/types"
	"net/http"
	"reflect"
	"sort"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type ConfigurationHistoryChange struct {
	UID      uint16 `json:"uid"`
	Name     string `json:"name"`
	Section  string `json:"section,omitempty"`
	OldValue any    `json:"old_value,omitempty"`
	NewValue any    `json:"new_value,omitempty"`
}

type ConfigurationHistoryView struct {
	ID            uint64                       `json:"id"`
	ApplyAt       time.Time                    `json:"apply_at"`
	CfgHash       uint32                       `json:"cfg_hash"`
	CfgSyncStatus *string                      `json:"cfg_sync_status,omitempty"`
	CfgSyncError  *string                      `json:"cfg_sync_error,omitempty"`
	Origin        string                       `json:"origin"`
	ChangeCount   int                          `json:"change_count"`
	Changes       []ConfigurationHistoryChange `json:"changes"`
	ParseError    *string                      `json:"parse_error,omitempty"`
}

func historyFieldValue(field adm.FieldConfigurationParsed) any {
	if field.Unknown {
		return hex.EncodeToString(field.RawValue)
	}

	return field.Value
}

func buildHistoryFieldMap(fields []adm.FieldConfigurationParsed) map[uint16]adm.FieldConfigurationParsed {
	result := make(map[uint16]adm.FieldConfigurationParsed, len(fields))
	for _, field := range fields {
		result[field.RawUID] = field
	}
	return result
}

func configurationHistoryDiff(older, newer map[uint16]adm.FieldConfigurationParsed) []ConfigurationHistoryChange {
	uids := make(map[uint16]struct{}, len(older)+len(newer))
	for uid := range older {
		uids[uid] = struct{}{}
	}
	for uid := range newer {
		uids[uid] = struct{}{}
	}

	ordered := make([]int, 0, len(uids))
	for uid := range uids {
		ordered = append(ordered, int(uid))
	}
	sort.Ints(ordered)

	changes := make([]ConfigurationHistoryChange, 0)
	for _, rawUID := range ordered {
		uid := uint16(rawUID)
		oldField, oldOK := older[uid]
		newField, newOK := newer[uid]

		var oldValue any
		var newValue any
		if oldOK {
			oldValue = historyFieldValue(oldField)
		}
		if newOK {
			newValue = historyFieldValue(newField)
		}

		if oldOK == newOK && reflect.DeepEqual(oldValue, newValue) {
			continue
		}

		name := fmt.Sprintf("unknown_%d", uid)
		section := ""
		if schema, err := constants.GetCfgSchema(uid); err == nil {
			name = schema.Name
			section = schema.Section
		} else if newOK && newField.UID != "" {
			name = newField.UID
		} else if oldOK && oldField.UID != "" {
			name = oldField.UID
		}

		changes = append(changes, ConfigurationHistoryChange{
			UID:      uid,
			Name:     name,
			Section:  section,
			OldValue: oldValue,
			NewValue: newValue,
		})
	}

	return changes
}

func (h *Handler) GetConfigurationHistoriesHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*gtypes.AuthToken)
	imei := mux.Vars(r)["imei"]

	if !authToken.User.Has(constants.AccessConfigurationRead) {
		return httperr.Forbidden(tr.TErr("config-access-restricted"))
	}

	device, err := h.Store.Devices.Get_DeviceByImei(ctx, imei)
	if err != nil {
		return httperr.Db(ctx, err)
	}
	if device == nil {
		return httperr.NotFound(tr.TErr("device-not-found"))
	}
	if !permissions.CanAccessDevice(ctx, h.Store.GroupMembers, device, device.ID, &authToken.User, authToken.RoleCode, permissions.ShareReadConfig) {
		return httperr.Forbidden(tr.TErr("device-not-owned"))
	}

	if device.Activated == nil || !*device.Activated {
		return httperr.New(tr.TErr("device-not-activated"), http.StatusUnprocessableEntity)
	}

	histories, err := h.Store.Configurations.Get_ConfigurationHistoriesByDeviceId(ctx, device.ID)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	parsed := make([]map[uint16]adm.FieldConfigurationParsed, len(histories))
	parseErrors := make([]*string, len(histories))

	for i := range histories {
		fields, parseErr := getConfiguration(tr, histories[i].CfgData)
		if parseErr != nil {
			message := parseErr.Error()
			parseErrors[i] = &message
			continue
		}
		if fields == nil {
			message := tr.TErr("cfg-parse-error-device")
			parseErrors[i] = &message
			continue
		}
		parsed[i] = buildHistoryFieldMap(fields)
	}

	count := len(histories)
	if count > 100 {
		count = 100
	}

	result := make([]ConfigurationHistoryView, 0, count)
	for i := 0; i < count; i++ {
		item := histories[i]
		changes := make([]ConfigurationHistoryChange, 0)

		if i+1 < len(histories) && parsed[i] != nil && parsed[i+1] != nil {
			changes = configurationHistoryDiff(parsed[i+1], parsed[i])
		}

		result = append(result, ConfigurationHistoryView{
			ID:            item.ID,
			ApplyAt:       item.ApplyAt,
			CfgHash:       item.CfgHash,
			CfgSyncStatus: item.CfgSyncStatus,
			CfgSyncError:  item.CfgSyncError,
			Origin:        item.Origin,
			ChangeCount:   len(changes),
			Changes:       changes,
			ParseError:    parseErrors[i],
		})
	}

	httpx.HttpResponseWithETag(w, r, http.StatusOK, result)
	return nil
}

func (h *Handler) ExportConfigurationHistoryHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*gtypes.AuthToken)
	imei := mux.Vars(r)["imei"]

	if !authToken.User.Has(constants.AccessConfigurationRead) {
		return httperr.Forbidden(tr.TErr("config-access-restricted"))
	}

	historyID, err := strconv.ParseUint(mux.Vars(r)["historyId"], 10, 64)
	if err != nil || historyID == 0 {
		return httperr.BadRequest(tr.TErr("fields-not-filled"))
	}

	device, err := h.Store.Devices.Get_DeviceByImei(ctx, imei)
	if err != nil {
		return httperr.Db(ctx, err)
	}
	if device == nil {
		return httperr.NotFound(tr.TErr("device-not-found"))
	}
	if !permissions.CanAccessDevice(ctx, h.Store.GroupMembers, device, device.ID, &authToken.User, authToken.RoleCode, permissions.ShareReadConfig) {
		return httperr.Forbidden(tr.TErr("device-not-owned"))
	}

	if device.Activated == nil || !*device.Activated {
		return httperr.New(tr.TErr("device-not-activated"), http.StatusUnprocessableEntity)
	}

	history, err := h.Store.Configurations.Get_ConfigurationHistoryByIDAndDeviceId(ctx, uint64(historyID), device.ID)
	if err != nil {
		return httperr.Db(ctx, err)
	}
	if history == nil {
		return httperr.NotFound(tr.TErr("config-fetch-error"))
	}
	if len(history.CfgData) < 14 {
		return httperr.New(tr.TErr("error-file-too-short"), http.StatusUnprocessableEntity)
	}

	filename := fmt.Sprintf("AdmConfiguration_%s_%d_%s.bin", imei, history.CfgHash, history.ApplyAt.UTC().Format("20060102T150405Z"))
	httpx.HttpFileResponse(w, r, filename, history.CfgData, "application/octet-stream")
	return nil
}
