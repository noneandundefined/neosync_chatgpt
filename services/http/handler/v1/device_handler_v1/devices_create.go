package device_handler_v1

import (
	"neomatica/neosync/infra/constants"
	"neomatica/neosync/infra/logger"
	"neomatica/neosync/infra/store/postgres/models"
	"neomatica/neosync/middleware"
	"neomatica/neosync/pkg/httpx"
	"neomatica/neosync/pkg/httpx/httperr"
	"neomatica/neosync/pkg/smartcaptcha"
	"neomatica/neosync/types"
	"neomatica/neosync/util"
	"neomatica/neosync/util/antispam"
	"net/http"
	"strings"
	"time"

	"github.com/go-playground/validator"
)

/* Anti spammer */
var spam = antispam.New()

/* Neosync HTTPx V1 */
/* Handler: создание устройства (также трансфер - если устройство уже создано) */

func (h *Handler) CreateDeviceHandler_V1(w http.ResponseWriter, r *http.Request) error { //nolint
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*types.AuthToken)

	/* Check access */
	if !authToken.User.Has(constants.AccessTrekerCreate) {
		return httperr.Forbidden(tr.TErr("access-treker-create"))
	}

	var userUuid string = authToken.User.UserContact.UserUUID
	var ownerUuid *string = nil

	/* Role AdminL2Support */
	if authToken.RoleCode == constants.Role_AdminL2Support {
		if authToken.User.ParentUUID != nil {
			userUuid = *authToken.User.ParentUUID
		}
	}

	/* Role User */
	if authToken.RoleCode == constants.Role_User {
		if authToken.User.ParentUUID != nil {
			userUuid = *authToken.User.ParentUUID
		}

		ownerUuid = &authToken.User.UserContact.UserUUID
	}

	var payload *CreateDevicePayload

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

	/* Yandex SmartCaptcha */
	if err := smartcaptcha.VerifyRequest(tr, payload.TurnstileToken, r); err != nil {
		logger.Error("CreateDeviceHandler_V1 req={%s}: Failed validation smartcaptcha token: %s", ctx.Value("XREQID").(string), err.Error())
		return httperr.BadRequest(err.Error())
	}

	exists, err := h.Store.Devices.Get_DeviceByImei(ctx, payload.Imei)
	if err != nil {
		return httperr.Db(ctx, err)
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

	if exists != nil {
		if exists.UserUUID != nil {
			spam.RegisterError(authToken.User.UserContact.UserUUID)

			if blocked, left := spam.IsBlocked(authToken.User.UserContact.UserUUID); blocked {
				return httperr.New(tr.TErr("device-add-blocked")+" "+left.Truncate(time.Second).String(), http.StatusTooManyRequests)
			}

			return httperr.BadRequest(linkedDeviceError(tr, *exists, authToken.User.UserContact.UserUUID, userUuid, authToken.RoleCode))
		}

		/* Check blocked */
		if blocked, left := spam.IsBlocked(authToken.User.UserContact.UserUUID); blocked {
			return httperr.New(tr.TErr("device-add-blocked")+" "+left.Truncate(time.Second).String(), http.StatusTooManyRequests)
		}

		tx, err := h.Db.BeginTx(ctx, nil)
		if err != nil {
			return httperr.Db(ctx, err)
		}

		defer func() {
			_ = tx.Rollback()
		}()

		if err := h.Store.Devices.Update_DeviceUuidByDeviceIdTx(ctx, tx, exists.ID, &userUuid, nil); err != nil {
			return httperr.Db(ctx, err)
		}

		deviceUpd := &models.UpdateDevice{
			Name:             optionalNonEmpty(payload.Name),
			Phone:            optionalNonEmpty(payload.Phone),
			NameOrganization: optionalNonEmpty(payload.NameOrganization),
		}

		if err := h.Store.Devices.Update_DeviceByImei(ctx, tx, payload.Imei, deviceUpd); err != nil {
			return httperr.Db(ctx, err)
		}

		if err := tx.Commit(); err != nil {
			return httperr.Conflict(tr.TErr("failed-to-save-data"))
		}

		/* Reset blocked */
		spam.Reset(authToken.User.UserContact.UserUUID)

		httpx.HttpResponse(w, r, http.StatusOK, tr.T("tracker-added-successfully"))
		return nil
	}

	/* Check blocked */
	if blocked, left := spam.IsBlocked(authToken.User.UserContact.UserUUID); blocked {
		return httperr.New(tr.TErr("device-add-blocked")+" "+left.Truncate(time.Second).String(), http.StatusTooManyRequests)
	}

	/* Reset blocked */
	spam.Reset(authToken.User.UserContact.UserUUID)

	syncData := &models.SyncData{
		FirmwareVersion: 0x00,
		CfgVersion:      0x00,
		LastModTime:     0x00,
		CfgHash:         0x00,
	}

	tx, err := h.Db.BeginTx(ctx, nil)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	defer func() {
		_ = tx.Rollback()
	}()

	device := &models.Device{
		UserUUID:         &userUuid,
		OwnerUUID:        ownerUuid,
		Name:             payload.Name,
		Phone:            payload.Phone,
		NameOrganization: payload.NameOrganization,
		IMEI:             payload.Imei,
		Status:           false,
		DeviceModel:      util.ValidDeviceModel(payload.Model),
	}

	deviceID, err := h.Store.Devices.Create_Device(ctx, tx, device)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	deviceConf := &models.DeviceConf{
		DeviceID:                      deviceID,
		RequestConfigurationOnConnect: payload.RequestConfigurationOnConnect,
	}

	if err := h.Store.Devices.Create_DeviceConf(ctx, tx, deviceConf); err != nil {
		return httperr.Db(ctx, err)
	}

	sync := &models.Sync{
		DeviceID:        deviceID,
		FirmwareVersion: syncData.FirmwareVersion,
		CfgVersion:      syncData.CfgVersion,
		LastModTime:     syncData.LastModTime,
		CfgHash:         syncData.CfgHash,
	}

	configuration := &models.Configuration{
		DeviceID: deviceID,
		CfgHash:  syncData.CfgHash,
	}

	if err := h.Store.Syncs.Create_Sync(ctx, tx, sync); err != nil {
		return httperr.Db(ctx, err)
	}

	if err := h.Store.Configurations.Create_Configuration(ctx, tx, configuration); err != nil {
		return httperr.Db(ctx, err)
	}

	if err := tx.Commit(); err != nil {
		return httperr.Conflict(tr.TErr("failed-to-save-data"))
	}

	httpx.HttpResponse(w, r, http.StatusCreated, tr.T("device-activates-on-server-connect"))
	return nil
}

func optionalNonEmpty(s *string) *string {
	if s == nil {
		return nil
	}

	v := strings.TrimSpace(*s)
	if v == "" || v == "+" {
		return nil
	}

	return &v
}
