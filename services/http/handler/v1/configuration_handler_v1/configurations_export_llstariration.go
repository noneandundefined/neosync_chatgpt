package configuration_handler_v1

import (
	"fmt"
	"neomatica/neosync/infra/constants"
	"neomatica/neosync/middleware"
	"neomatica/neosync/pkg/httpx"
	"neomatica/neosync/pkg/httpx/httperr"
	"neomatica/neosync/pkg/permissions"
	"neomatica/neosync/types"
	"net/http"
	"strings"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
)

/* Neosync HTTPx V1 */
/* Handler: экспорт в файл LLS датчиков */

func (h *Handler) ConfigurationExportLLSTarirationHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*types.AuthToken)

	/* Check access */
	if !authToken.User.Has(constants.AccessConfigurationRead) {
		return httperr.Forbidden(tr.TErr("config-edit-access-restricted"))
	}

	imei := mux.Vars(r)["imei"]
	if imei == "" {
		return httperr.NotFound(tr.TErr("device-imei-not-found"))
	}

	/* Device ownership verification */
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

	if !*device.Activated {
		return httperr.New(tr.TErr("device-not-activated"), http.StatusUnprocessableEntity)
	}

	var payload *ConfigurationExportLLSTarirationPayload

	if err := httpx.HttpParse(r, &payload); err != nil {
		return httperr.BadRequest(err.Error())
	}

	if err := httpx.Validate.Struct(payload); err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			return httperr.BadRequest(httpx.ValidateMsg(tr, err))
		}

		return httperr.BadRequest(tr.TErr("fields-not-filled"))
	}

	var CsvRu = CsvLLSTarirationLabels{
		SensorHeader: "Датчик",
		ColIndex:     "№",
		ColValue:     "Значение датчика",
		ColLevel:     "Уровень топлива",
	}
	var labels CsvLLSTarirationLabels

	labels = CsvRu

	var sb strings.Builder

	for sensorId, rows := range payload.Tariration {
		sb.WriteString(fmt.Sprintf("%s %d\n", labels.SensorHeader, sensorId))

		sb.WriteString(fmt.Sprintf("%s;%s;%s\n",
			labels.ColIndex,
			labels.ColValue,
			labels.ColLevel,
		))

		for i, row := range rows {
			sb.WriteString(fmt.Sprintf("%d;%.2f;%.2f\n",
				i, row.Value, row.Level,
			))
		}

		sb.WriteString("\n")
	}

	csvData := sb.String()

	filename := "Calibration table.csv"
	if tr.GetLang() == "ru" {
		filename = "Тарировочная таблица.csv"
	}

	w.Header().Set("Content-Disposition", "attachment; filename=\""+filename+"\"")
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Write([]byte("\xEF\xBB\xBF"))
	w.Write([]byte(csvData))
	return nil
}
