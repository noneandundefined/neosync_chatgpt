package device_handler_v1

import (
	"neomatica/neosync/infra/constants"
	"neomatica/neosync/infra/logger"
	"neomatica/neosync/middleware"
	"neomatica/neosync/pkg/httpx"
	"neomatica/neosync/pkg/httpx/httperr"
	"neomatica/neosync/pkg/smartcaptcha"
	"neomatica/neosync/types"
	"net/http"

	"github.com/go-playground/validator"
)

/* Neosync HTTPx V1 */
/* Handler: трансфер устройств в аккаунты дилеров */

func (h *Handler) DeviceTransferToAccountHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*types.AuthToken)

	var payload *TransferAccountPayload

	if err := httpx.HttpParse(r, &payload); err != nil {
		return httperr.BadRequest(err.Error())
	}

	if err := httpx.Validate.Struct(payload); err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			return httperr.BadRequest(httpx.ValidateMsg(tr, err))
		}

		return httperr.BadRequest(tr.TErr("fields-not-filled"))
	}

	/* Yandex SmartCaptcha */
	if err := smartcaptcha.VerifyRequest(tr, payload.TurnstileToken, r); err != nil {
		logger.Error("DeviceTransferToAccountHandler_V1 req={%s}: Failed validation smartcaptcha token: %s", ctx.Value("XREQID").(string), err.Error())
		return httperr.BadRequest(err.Error())
	}

	if len(payload.IMEIs) == 0 {
		return httperr.BadRequest(tr.TErr("fields-not-filled"))
	}

	var targetUUID *string
	if payload.UUID != nil && *payload.UUID != "" {
		targetUUID = payload.UUID
	}

	switch authToken.RoleCode {
	case constants.Role_SuperAdmin, constants.Role_Support:
		for _, imei := range payload.IMEIs {
			if len(imei) != 15 {
				return httperr.BadRequest(tr.TErr("invalid-or-nonexistent-imei"))
			}

			if err := h.Store.Devices.Update_DeviceUuidByImei(ctx, imei, targetUUID, nil); err != nil {
				return httperr.Db(ctx, err)
			}
		}

	case constants.Role_AdminL2, constants.Role_AdminL2Support:
		for _, imei := range payload.IMEIs {
			if len(imei) != 15 {
				return httperr.BadRequest(tr.TErr("invalid-or-nonexistent-imei"))
			}

			if err := h.Store.Devices.Update_DeviceOwnerUuidByImei(ctx, imei, targetUUID); err != nil {
				return httperr.Db(ctx, err)
			}
		}
	}

	httpx.HttpResponse(w, r, http.StatusOK, tr.T("device-transferred-to-account"))
	return nil
}
