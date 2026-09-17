package device_handler_v1

import (
	"neomatica/neosync/infra/constants"
	"neomatica/neosync/middleware"
	"neomatica/neosync/pkg/httpx"
	"neomatica/neosync/pkg/httpx/httperr"
	"neomatica/neosync/types"
	"net/http"
)

const importCheckMaxImeis = 100

func (h *Handler) DevicesImportCheckHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*types.AuthToken)

	if !authToken.User.Has(constants.AccessTrekerCreate) {
		return httperr.Forbidden(tr.TErr("access-treker-create"))
	}

	currentUserUuid := authToken.User.UserContact.UserUUID
	accountUserUuid := currentUserUuid

	if authToken.RoleCode == constants.Role_AdminL2Support && authToken.User.ParentUUID != nil {
		accountUserUuid = *authToken.User.ParentUUID
	}

	if authToken.RoleCode == constants.Role_User && authToken.User.ParentUUID != nil {
		accountUserUuid = *authToken.User.ParentUUID
	}

	var payload ImportCheckPayload
	if err := httpx.HttpParse(r, &payload); err != nil {
		return httperr.BadRequest(err.Error())
	}

	imeis := payload.Imeis
	if len(imeis) > importCheckMaxImeis {
		imeis = imeis[:importCheckMaxImeis]
	}

	importErrors := make([]ImportError, 0)
	if len(imeis) == 0 {
		httpx.HttpResponse(w, r, http.StatusOK, importErrors)
		return nil
	}

	existings, err := h.Store.Devices.Get_DevicesByImeis(ctx, imeis)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	for _, exist := range existings {
		if exist.UserUUID == nil {
			continue
		}

		importErrors = append(importErrors, ImportError{
			Imei:  exist.IMEI,
			Error: linkedDeviceError(tr, exist, currentUserUuid, accountUserUuid, authToken.RoleCode),
		})
	}

	httpx.HttpResponse(w, r, http.StatusOK, importErrors)
	return nil
}
