package configuration_handler_v1

import (
	"fmt"
	"neomatica/neosync/infra/constants"
	"neomatica/neosync/infra/logger"
	"neomatica/neosync/middleware"
	"neomatica/neosync/pkg/httpx"
	"neomatica/neosync/pkg/httpx/httperr"
	"neomatica/neosync/pkg/permissions"
	"neomatica/neosync/types"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
)

func (h *Handler) ConfigurationChangePassHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*types.AuthToken)

	/* Check access */
	if !authToken.User.Has(constants.AccessConfigurationApply) {
		return httperr.Forbidden(tr.TErr("access-treker-create"))
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

	if !permissions.CanAccessDevice(ctx, h.Store.GroupMembers, device, device.ID, &authToken.User, authToken.RoleCode, permissions.ShareEditConfig) {
		return httperr.Forbidden(tr.TErr("device-not-owned"))
	}

	if !*device.Activated {
		return httperr.New(tr.TErr("device-not-activated"), http.StatusUnprocessableEntity)
	}

	if !device.Status {
		return httperr.New(tr.TErr("device-not-connected-to-server"), http.StatusUnprocessableEntity)
	}

	var payload *ConfigurationChangePassPayload

	if err := httpx.HttpParse(r, &payload); err != nil {
		return httperr.BadRequest(err.Error())
	}

	if err := httpx.Validate.Struct(payload); err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			return httperr.BadRequest(httpx.ValidateMsg(tr, err))
		}

		return httperr.BadRequest(tr.TErr("fields-not-filled"))
	}

	command := fmt.Sprintf("PASS %v %v", payload.OldPassword, payload.NewPassword)

	rabbitData := types.RabbitMQ_TransitBinary{
		Type:  constants.ADM_RC_TYPE_STRING,
		Imei:  imei,
		NResp: true,
		Data:  []byte(command),
	}

	if err := h.RMQ.SendToRabbitAsync(ctx, rabbitData); err != nil {
		logger.Error("ConfigurationChangePassHandler_V1 req={%s}: failed to send to RabbitMQ: %s", ctx.Value("XREQID").(string), err.Error())
		return httperr.Conflict(tr.TErr("rabbitmq-publish-failed"))
	}

	httpx.HttpResponseWithETag(w, r, http.StatusOK, tr.T("success-change-password"))
	return nil
}
