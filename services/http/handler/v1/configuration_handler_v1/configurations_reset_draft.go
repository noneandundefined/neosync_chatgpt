package configuration_handler_v1

import (
	"neomatica/neosync/infra/constants"
	"neomatica/neosync/infra/logger"
	"neomatica/neosync/infra/store/redis"
	"neomatica/neosync/middleware"
	"neomatica/neosync/pkg/httpx"
	"neomatica/neosync/pkg/httpx/httperr"
	"neomatica/neosync/pkg/permissions"
	"neomatica/neosync/types"
	"net/http"

	"github.com/gorilla/mux"
)

/* Neosync HTTPx V1 */
/* Handler: сброс/очистка конфигурации */

func (h *Handler) ResetConfigurationDraftHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*types.AuthToken)

	/* Check access */
	if !authToken.User.Has(constants.AccessConfigurationApply) {
		return httperr.Forbidden(tr.TErr("config-access-restricted"))
	}

	imei := mux.Vars(r)["imei"]
	if imei == "" {
		return httperr.BadRequest(tr.TErr("device-imei-not-found"))
	}

	/* Device ownership verification */
	device, err := h.Store.Devices.Get_DeviceByImei(ctx, imei)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	if device == nil {
		return httperr.NotFound(tr.TErr("device-not-found"))
	}

	if !permissions.CanAccessDevice(ctx, h.Store.GroupMembers, device, device.ID, &authToken.User, authToken.RoleCode, permissions.ShareEditConfig) {
		return httperr.Forbidden(tr.TErr("device-not-owned"))
	}

	if !*device.Activated {
		return httperr.New(tr.TErr("device-not-activated"), http.StatusUnprocessableEntity)
	}

	if err := redis.ClearCacheConfiguration(imei); err != nil {
		return httperr.Db(ctx, err)
	}

	go func(uuid, _imei string) {
		if err := redis.DeleteAllDrafts(uuid, _imei); err != nil {
			logger.Error("ResetConfigurationDraftHandler_V1 req={%s} imei={%s}: Failed delete all drafts: %s", ctx.Value("XREQID").(string), _imei, err.Error())
		}
	}(authToken.User.UserContact.UserUUID, imei)

	httpx.HttpResponse(w, r, http.StatusNoContent, nil)
	return nil
}
