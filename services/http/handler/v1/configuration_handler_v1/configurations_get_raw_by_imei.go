package configuration_handler_v1

import (
	"encoding/hex"
	"neomatica/neosync/infra/constants"
	"neomatica/neosync/middleware"
	"neomatica/neosync/pkg/httpx"
	"neomatica/neosync/pkg/httpx/httperr"
	"neomatica/neosync/pkg/permissions"
	"neomatica/neosync/types"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

/* Neosync HTTPx V1 */
/* Handler: получение полной конфигурации в HEX формате */

func (h *Handler) GetConfigurationRawByImeiHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*types.AuthToken)

	/* Check access */
	if !authToken.User.Has(constants.AccessConfigurationRead) {
		return httperr.Forbidden(tr.TErr("config-access-restricted"))
	}

	imei := mux.Vars(r)["imei"]

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

	configuration, err := h.Store.Configurations.Get_ConfigurationByDeviceId(ctx, device.ID)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	configurationRawModel := ConfigurationRaw{
		DeviceImei: imei,
		CfgHash:    configuration.CfgHash,
		ConfigHex:  hex.EncodeToString(configuration.CfgData),
		Timestamp:  time.Now().Unix(),
	}

	httpx.HttpResponseWithETag(w, r, http.StatusOK, configurationRawModel)
	return nil
}
