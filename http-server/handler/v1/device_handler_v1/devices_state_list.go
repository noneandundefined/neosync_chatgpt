package device_handler_v1

import (
	"context"
	"neomatica/neosync/infra/constants"
	"neomatica/neosync/infra/locale"
	"neomatica/neosync/infra/store/postgres/models"
	"neomatica/neosync/middleware"
	"neomatica/neosync/pkg/httpx"
	"neomatica/neosync/pkg/httpx/httperr"
	"neomatica/neosync/types"
	"net/http"
	"strings"
)

/* Neosync HTTPx V1 */
/* Handler: получение групп устройств для отправки команд */
/* Параметры - ?search=поиск по IMEI, имени, организации */

func (h *Handler) GetDevicesStateListHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*types.AuthToken)
	search := strings.TrimSpace(r.URL.Query().Get("search"))

	devices, err := h.loadDevicesStateGroups(ctx, authToken, tr, search)
	if err != nil {
		return err
	}

	httpx.HttpResponseWithETag(w, r, http.StatusOK, devices)
	return nil
}

func (h *Handler) loadDevicesStateGroups(ctx context.Context, authToken *types.AuthToken, tr locale.Translator, search string) ([]models.DevicesForSendCommand, error) {
	var (
		devices []models.DevicesForSendCommand
		err     error
	)

	switch authToken.RoleCode {
	case constants.Role_SuperAdmin, constants.Role_Support:
		devices, err = h.Store.DeviceCommands.Get_DevicesState(ctx, authToken.User.UserContact.UserUUID, search)

	case constants.Role_AdminL2, constants.Role_AdminL2Support:
		devices, err = h.Store.DeviceCommands.Get_DevicesStateByUserUuid(ctx, authToken, search)

	case constants.Role_User:
		if authToken.User.ParentUUID == nil {
			return nil, httperr.Forbidden(tr.TErr("user-not-linked-to-admin"))
		}

		devices, err = h.Store.DeviceCommands.Get_DevicesStateByUserUuid(ctx, authToken, search)

	default:
		return nil, httperr.Forbidden(tr.TErr("insufficient-permissions-view-devices"))
	}

	if err != nil {
		return nil, httperr.Db(ctx, err)
	}

	return devices, nil
}

func (h *Handler) loadDevicesStateDevices(ctx context.Context, authToken *types.AuthToken, tr locale.Translator, groupID uint64, search string, limit, offset int, all bool) ([]models.DeviceShort, int, error) {
	var (
		devices []models.DeviceShort
		total   int
		err     error
	)

	switch authToken.RoleCode {
	case constants.Role_SuperAdmin, constants.Role_Support:
		devices, total, err = h.Store.DeviceCommands.Get_DevicesStateDevices(ctx, authToken.User.UserContact.UserUUID, groupID, search, limit, offset, all)

	case constants.Role_AdminL2, constants.Role_AdminL2Support:
		devices, total, err = h.Store.DeviceCommands.Get_DevicesStateDevicesByUserUuid(ctx, authToken, groupID, search, limit, offset, all)

	case constants.Role_User:
		if authToken.User.ParentUUID == nil {
			return nil, 0, httperr.Forbidden(tr.TErr("user-not-linked-to-admin"))
		}

		devices, total, err = h.Store.DeviceCommands.Get_DevicesStateDevicesByUserUuid(ctx, authToken, groupID, search, limit, offset, all)

	default:
		return nil, 0, httperr.Forbidden(tr.TErr("insufficient-permissions-view-devices"))
	}

	if err != nil {
		return nil, 0, httperr.Db(ctx, err)
	}

	return devices, total, nil
}
