package device_handler_v1

import (
	"database/sql"
	"neomatica/neosync/infra/constants"
	"neomatica/neosync/infra/store/postgres/models"
	"neomatica/neosync/middleware"
	"neomatica/neosync/pkg/httpx"
	"neomatica/neosync/pkg/httpx/httperr"
	"neomatica/neosync/pkg/permissions"
	"neomatica/neosync/types"
	"neomatica/neosync/util"
	"net/http"

	"github.com/go-playground/validator"
)

/* Neosync HTTPx V1 */
/* Handler: импортирование устройств макс. 50 устройств */

func (h *Handler) DevicesImportHandler_V1(w http.ResponseWriter, r *http.Request) error { //nolint
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*types.AuthToken)

	/* Check access */
	if !authToken.User.Has(constants.AccessTrekerCreate) {
		return httperr.Forbidden(tr.TErr("access-treker-create"))
	}

	var currentUserUuid string = authToken.User.UserContact.UserUUID
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

	var payload []ImportDevicePayload

	if err := httpx.HttpParse(r, &payload); err != nil {
		return httperr.BadRequest(err.Error())
	}

	if len(payload) == 0 {
		return httperr.BadRequest(tr.TErr("fields-not-filled"))
	}

	for _, device := range payload {
		if err := httpx.Validate.Struct(device); err != nil {
			if _, ok := err.(validator.ValidationErrors); ok {
				return httperr.BadRequest(httpx.ValidateMsg(tr, err))
			}

			return httperr.BadRequest(tr.TErr("fields-not-filled"))
		}
	}

	var allowedOwnersByEmail map[string]string
	if authToken.RoleCode != constants.Role_User {
		role := constants.Role_AdminL2
		parentUuid := sql.NullString{Valid: false}

		switch authToken.RoleCode {
		case constants.Role_AdminL2:
			role = constants.Role_User
			parentUuid = sql.NullString{String: currentUserUuid, Valid: true}
		case constants.Role_AdminL2Support:
			role = constants.Role_User
			parentUuid = sql.NullString{String: *authToken.User.ParentUUID, Valid: true}
		}

		owners, _, err := h.Store.Users.Get_UsersToOwner(ctx, role, currentUserUuid, parentUuid, "", 0, 0, true)
		if err != nil {
			return httperr.Db(ctx, err)
		}

		allowedOwnersByEmail = make(map[string]string, len(owners))
		for _, user := range owners {
			allowedOwnersByEmail[user.Email] = user.UserUUID
		}
	}

	groups, err := h.Store.Groups.Get_GroupNamesByUuid(ctx, currentUserUuid)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	allowedGroupsByName := make(map[string]uint64, len(groups))
	for _, group := range groups {
		allowedGroupsByName[group.Name] = group.ID
	}

	imeis := make([]string, len(payload))
	for idx, device := range payload {
		imeis[idx] = device.Imei
	}

	existings, err := h.Store.Devices.Get_DevicesByImeis(ctx, imeis)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	existsMap := make(map[string]models.Device)
	for _, device := range existings {
		existsMap[device.IMEI] = device
	}

	type importAssignment struct {
		userUuid  *string
		ownerUuid *string
		groupId   *uint64
	}

	assignments := make(map[string]importAssignment, len(payload))

	var importErrors []ImportError

	for _, d := range payload {
		resolvedUserUuid := &userUuid
		resolvedOwnerUuid := ownerUuid
		if authToken.RoleCode != constants.Role_User && d.OwnerEmail != nil && *d.OwnerEmail != "" {
			ownerUserUuid, ok := allowedOwnersByEmail[*d.OwnerEmail]
			if !ok {
				importErrors = append(importErrors, ImportError{
					Imei:  d.Imei,
					Error: tr.TErr("user-not-found"),
				})
				continue
			}

			if permissions.IsMainRole(authToken.RoleCode) {
				resolvedUserUuid = &ownerUserUuid
				resolvedOwnerUuid = nil
			} else {
				resolvedOwnerUuid = &ownerUserUuid
			}
		}

		var resolvedGroupId *uint64
		if d.GroupName != nil && *d.GroupName != "" {
			groupId, ok := allowedGroupsByName[*d.GroupName]
			if !ok {
				importErrors = append(importErrors, ImportError{
					Imei:  d.Imei,
					Error: tr.TErr("group-not-found"),
				})
				continue
			}

			resolvedGroupId = &groupId
		}

		assignments[d.Imei] = importAssignment{
			userUuid:  resolvedUserUuid,
			ownerUuid: resolvedOwnerUuid,
			groupId:   resolvedGroupId,
		}
	}

	if len(importErrors) > 0 {
		httpx.HttpResponseWithETag(w, r, http.StatusConflict, importErrors)
		return nil
	}

	tx, err := h.Db.BeginTx(ctx, nil)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	defer func() {
		_ = tx.Rollback()
	}()

	syncData := &models.SyncData{
		FirmwareVersion: 0x00,
		CfgVersion:      0x00,
		LastModTime:     0x00,
		CfgHash:         0x00,
	}

	var devicesInsert []*models.Device
	var deviceConfsInsert []*models.DeviceConf
	var configurationsInsert []*models.Configuration
	var syncsInsert []*models.Sync

	groupDevices := make(map[uint64][]uint64)

	for _, d := range payload {
		assignment := assignments[d.Imei]

		if exist, ok := existsMap[d.Imei]; ok {
			if exist.UserUUID != nil {
				importErrors = append(importErrors, ImportError{
					Imei:  d.Imei,
					Error: linkedDeviceError(tr, exist, currentUserUuid, userUuid, authToken.RoleCode),
				})
				continue
			}

			if err := h.Store.Devices.Update_DeviceUuidByDeviceIdTx(ctx, tx, exist.ID, assignment.userUuid, assignment.ownerUuid); err != nil {
				return httperr.Db(ctx, err)
			}

			if assignment.groupId != nil {
				groupDevices[*assignment.groupId] = append(groupDevices[*assignment.groupId], exist.ID)
			}

			continue
		}

		device := &models.Device{
			UserUUID:    assignment.userUuid,
			OwnerUUID:   assignment.ownerUuid,
			IMEI:        d.Imei,
			Status:      false,
			DeviceModel: util.ValidDeviceModel(d.Model),
		}

		devicesInsert = append(devicesInsert, device)
	}

	/* Check errors and return HTTP response */
	if len(importErrors) > 0 {
		httpx.HttpResponseWithETag(w, r, http.StatusConflict, importErrors)
		return nil
	}

	insertedDevices, err := h.Store.Devices.Create_DevicesBatch(ctx, tx, devicesInsert)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	for _, device := range insertedDevices {
		assignment := assignments[device.IMEI]

		if assignment.groupId != nil {
			groupDevices[*assignment.groupId] = append(groupDevices[*assignment.groupId], device.ID)
		}

		deviceConfsInsert = append(deviceConfsInsert, &models.DeviceConf{
			DeviceID:                      device.ID,
			RequestConfigurationOnConnect: false,
		})

		syncsInsert = append(syncsInsert, &models.Sync{
			DeviceID:        device.ID,
			FirmwareVersion: syncData.FirmwareVersion,
			CfgVersion:      syncData.CfgVersion,
			LastModTime:     syncData.LastModTime,
			CfgHash:         syncData.CfgHash,
		})

		configurationsInsert = append(configurationsInsert, &models.Configuration{
			DeviceID: device.ID,
			CfgHash:  syncData.CfgHash,
		})
	}

	if err := h.Store.Devices.Create_DeviceConfBatch(ctx, tx, deviceConfsInsert); err != nil {
		return httperr.Db(ctx, err)
	}

	if err := h.Store.Syncs.Create_SyncsBatch(ctx, tx, syncsInsert); err != nil {
		return httperr.Db(ctx, err)
	}

	if err := h.Store.Configurations.Create_ConfigurationsBatch(ctx, tx, configurationsInsert); err != nil {
		return httperr.Db(ctx, err)
	}

	for groupId, deviceIds := range groupDevices {
		if len(deviceIds) == 0 {
			continue
		}

		if err := h.Store.Groups.Update_AddDevicesToGroupTx(ctx, tx, &models.GroupDevices{
			UserUuid:  currentUserUuid,
			GroupId:   groupId,
			DevicesId: deviceIds,
		}); err != nil {
			return httperr.Db(ctx, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return httperr.Conflict(tr.TErr("failed-to-save-data"))
	}

	httpx.HttpResponse(w, r, http.StatusCreated, tr.T("successful-device-import"))
	return nil
}
