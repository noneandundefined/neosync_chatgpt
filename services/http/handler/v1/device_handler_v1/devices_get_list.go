package device_handler_v1

import (
	"math"
	"neomatica/neosync/infra/constants"
	"neomatica/neosync/infra/store/postgres/models"
	"neomatica/neosync/middleware"
	"neomatica/neosync/pkg/httpx"
	"neomatica/neosync/pkg/httpx/httperr"
	"neomatica/neosync/types"
	"neomatica/neosync/util"
	"net/http"
)

/* Neosync HTTPx V1 */
/* Handler: получение устройств у пользователя */
/* Параметры - ?page=страница ?limit=макс. значение элементов ?search=поиск */

func (h *Handler) GetDevicesHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*types.AuthToken)

	pagination := util.GetPagination(r, 10)

	var (
		devices  []models.Device_DeviceConf_Sync
		total    int
		totalAll int
		err      error
	)

	var (
		userUuid  *string
		ownerUuid *string
	)

	switch authToken.RoleCode {
	case constants.Role_SuperAdmin, constants.Role_Support:
		userUuid = &authToken.User.UserContact.UserUUID

		devices, total, totalAll, err = h.Store.Devices.Get_DevicesWithParams(ctx, pagination.Limit, pagination.Offset, pagination.Search, pagination.ColumnSortKey, pagination.ColumnSortDir, *userUuid)

	case constants.Role_AdminL2:
		userUuid = &authToken.User.UserContact.UserUUID

		devices, total, totalAll, err = h.Store.Devices.Get_DevicesByUuidsWithParams(ctx, pagination.Limit, pagination.Offset, pagination.Search, pagination.ColumnSortKey, pagination.ColumnSortDir, authToken.RoleCode, userUuid, ownerUuid)

	case constants.Role_AdminL2Support:
		userUuid = authToken.User.ParentUUID

		devices, total, totalAll, err = h.Store.Devices.Get_DevicesByUuidsWithParams(ctx, pagination.Limit, pagination.Offset, pagination.Search, pagination.ColumnSortKey, pagination.ColumnSortDir, authToken.RoleCode, userUuid, ownerUuid)

	case constants.Role_User:
		ownerUuid = &authToken.User.UserContact.UserUUID

		devices, total, totalAll, err = h.Store.Devices.Get_DevicesByUuidsWithParams(ctx, pagination.Limit, pagination.Offset, pagination.Search, pagination.ColumnSortKey, pagination.ColumnSortDir, authToken.RoleCode, userUuid, ownerUuid)

	default:
		return httperr.Forbidden(tr.TErr("insufficient-permissions-view-devices"))
	}

	if err != nil {
		return httperr.Db(ctx, err)
	}

	totalPages := int(math.Ceil(float64(total) / float64(pagination.Limit)))
	devicesResp := DevicesWPResponse{
		Items:      devices,
		Page:       pagination.Page,
		Limit:      pagination.Limit,
		Total:      total,
		TotalAll:   totalAll,
		TotalPages: totalPages,
	}

	w.Header().Set("Cache-Control", "no-store")
	httpx.HttpResponseWithETag(w, r, http.StatusOK, devicesResp)
	return nil
}

/* Neosync HTTPx V1 */
/* Handler: получение всех устройств пользователя */

// func (h *Handler) GetDevicesAlHandler_V1(w http.ResponseWriter, r *http.Request) error {
// 	ctx := r.Context()
// 	tr := middleware.TranslatorFromContext(ctx)
// 	authToken := ctx.Value("identity").(*types.AuthToken)

// 	var (
// 		devices []models.Device_DeviceConf_Sync
// 		err     error
// 	)

// 	switch authToken.RoleCode {
// 	case constants.Role_SuperAdmin, constants.Role_Support:
// 		devices, err = h.Store.Devices.Get_Devices(ctx, authToken.User.UserContact.UserUUID)

// 	case constants.Role_AdminL2:
// 		devices, err = h.Store.Devices.Get_DevicesByUuid(ctx, authToken)

// 	case constants.Role_User:
// 		if authToken.User.ParentUUID == nil {
// 			return httperr.Forbidden(tr.TErr("user-not-linked-to-admin"))
// 		}

// 		devices, err = h.Store.Devices.Get_DevicesByUuid(ctx, authToken)

// 	default:
// 		return httperr.Forbidden(tr.TErr("insufficient-permissions-view-devices"))
// 	}

// 	if err != nil {
// 		return httperr.Db(ctx, err)
// 	}

// 	httpx.HttpResponseWithETag(w, r, http.StatusOK, devices)
// 	return nil
// }
