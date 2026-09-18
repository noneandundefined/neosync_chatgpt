package device_handler_v1

import (
	"neomatica/neosync/infra/constants"
	"neomatica/neosync/infra/store/postgres/models"
	"neomatica/neosync/middleware"
	"neomatica/neosync/pkg/httpx"
	"neomatica/neosync/pkg/httpx/httperr"
	"neomatica/neosync/pkg/permissions"
	"neomatica/neosync/types"
	"neomatica/neosync/util"
	"net/http"
	"strings"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
)

/* Neosync HTTPx V1 */
/* Handler: обновление данных устройства */

func (h *Handler) DeviceUpdateByImeiHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*types.AuthToken)

	/* Check access */
	if !authToken.User.Has(constants.AccessTrekerEdit) {
		return httperr.Forbidden(tr.TErr("access-treker-edit"))
	}

	imei := mux.Vars(r)["imei"]
	if imei == "" {
		return httperr.NotFound(tr.TErr("device-imei-not-found"))
	}

	if len(imei) != 15 {
		return httperr.BadRequest(tr.TErr("invalid-or-nonexistent-imei"))
	}

	var payload *EditDevicePayload

	if err := httpx.HttpParse(r, &payload); err != nil {
		return httperr.BadRequest(err.Error())
	}

	if err := httpx.Validate.Struct(payload); err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			return httperr.BadRequest(httpx.ValidateMsg(tr, err))
		}

		return httperr.BadRequest(tr.TErr("fields-not-filled"))
	}

	payload.Name = optionalNonEmpty(payload.Name)
	payload.Phone = optionalNonEmpty(payload.Phone)
	payload.NameOrganization = optionalNonEmpty(payload.NameOrganization)

	if payload.Name != nil && len([]rune(*payload.Name)) > 30 {
		return httperr.BadRequest(tr.TErr("device-name-too-long"))
	}

	if payload.NameOrganization != nil && len([]rune(*payload.NameOrganization)) > 150 {
		return httperr.BadRequest(tr.TErr("organization-name-too-long"))
	}

	/* Device ownership verification */
	device, err := h.Store.Devices.Get_DeviceByImei(ctx, imei)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	if device == nil {
		return httperr.NotFound(tr.TErr("device-not-found"))
	}

	if !permissions.IsMainRole(authToken.RoleCode) && !permissions.DeviceBelongsToUser(device, &authToken.User) {
		return httperr.Forbidden(tr.TErr("device-not-owned"))
	}

	/* Check NameOrganization */
	if payload.NameOrganization != nil && len(*payload.NameOrganization) < 3 {
		return httperr.BadRequest(tr.TErr("organization-name-too-short"))
	}

	/* Check Phone */
	if payload.Phone != nil {
		phone := strings.TrimSpace(*payload.Phone)
		if phone != "" {
			phone = strings.TrimPrefix(phone, "+")
			phone = strings.ReplaceAll(phone, " ", "")
			if len(phone) < 8 {
				return httperr.BadRequest(tr.TErr("invalid-phone-number"))
			}
			payload.Phone = &phone
		}
	}

	var deviceModel *string

	if payload.Model != nil {
		deviceModel = util.ValidDeviceModel(*payload.Model)
	}

	tx, err := h.Db.BeginTx(ctx, nil)
	if err != nil {
		return httperr.Db(ctx, httperr.Err_DbNetwork)
	}

	defer func() {
		_ = tx.Rollback()
	}()

	deviceUpd := &models.UpdateDevice{
		Imei:                          payload.Imei,
		DeviceModel:                   deviceModel,
		Name:                          payload.Name,
		Phone:                         payload.Phone,
		NameOrganization:              payload.NameOrganization,
		RequestConfigurationOnConnect: payload.RequestConfigurationOnConnect,
	}

	if err := h.Store.Devices.Update_DeviceByImei(ctx, tx, imei, deviceUpd); err != nil {
		return httperr.Db(ctx, err)
	}

	if err := h.Store.Devices.Update_DeviceConfByDeviceId(ctx, tx, device.ID, deviceUpd); err != nil {
		return httperr.Db(ctx, err)
	}

	if err := tx.Commit(); err != nil {
		return httperr.Db(ctx, err)
	}

	httpx.HttpResponse(w, r, http.StatusOK, tr.T("device-data-updated"))
	return nil
}
