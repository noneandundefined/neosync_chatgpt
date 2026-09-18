package configuration_handler_v1

import (
	"fmt"
	"neomatica/neosync/infra/constants"
	"neomatica/neosync/infra/locale"
	"neomatica/neosync/infra/logger"
	"neomatica/neosync/infra/store/postgres/models"
	"neomatica/neosync/infra/store/redis"
	"neomatica/neosync/middleware"
	"neomatica/neosync/pkg/adm"
	"neomatica/neosync/pkg/adm/admparser"
	"neomatica/neosync/pkg/httpx"
	"neomatica/neosync/pkg/httpx/httperr"
	"neomatica/neosync/pkg/permissions"
	gtypes "neomatica/neosync/types"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func convertChangesToParsedWithSize(changes map[string]any, cfgParsed []adm.FieldConfigurationParsed) []adm.PacketFieldConfiguration {
	draft := &redis.ConfigurationDraftInsert{
		Changes: changes,
	}

	return convertDraftToParsedWithSize(draft, cfgParsed)
}

func buildConfigurationFromChanges(tr locale.Translator, changes map[string]any, cfgParsed []adm.FieldConfigurationParsed) ([]byte, error) {
	cfgConverted := convertChangesToParsedWithSize(changes, cfgParsed)
	admInst := adm.Adm{}

	return admInst.ParseJsonToBinary(tr, cfgConverted)
}

func convertDraftToParsedWithSize(draft *redis.ConfigurationDraftInsert, cfgParsed []adm.FieldConfigurationParsed) []adm.PacketFieldConfiguration {
	fullConfig := make([]adm.PacketFieldConfiguration, 0, len(cfgParsed))

	changes := make(map[string]interface{})
	for k, v := range draft.Changes {
		changes[k] = v
	}

	for _, field := range cfgParsed {
		/* If unknown UID */
		if field.Unknown {
			fullConfig = append(fullConfig, adm.PacketFieldConfiguration{
				UID:      field.RawUID,
				RawUID:   field.RawUID,
				Size:     field.Size,
				RawValue: field.RawValue,
				Unknown:  true,
			})
			continue
		}

		uidStr := field.UID
		var value interface{}

		if newVal, ok := changes[uidStr]; ok {
			value = newVal
		} else {
			value = field.Value
		}

		var uid uint16
		found := false
		for _, s := range constants.CfgSchema {
			if s.Name == uidStr {
				uid = s.UID
				found = true
				break
			}
		}
		if !found {
			continue
		}

		fullConfig = append(fullConfig, adm.PacketFieldConfiguration{
			UID:   uid,
			Size:  field.Size,
			Value: value,
		})
	}

	return fullConfig
}

/* Neosync HTTPx V1 */
/* Handler: build/отправка новой конфигурации */

func (h *Handler) ApplyConfigurationHandler_V1(w http.ResponseWriter, r *http.Request) error { //nolint
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*gtypes.AuthToken)

	/* Check access */
	if !authToken.User.Has(constants.AccessConfigurationApply) {
		return httperr.Forbidden(tr.TErr("config-edit-access-restricted"))
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

	if !permissions.HasDeviceAccount(device.UserUUID) {
		return httperr.New(tr.TErr("device-has-no-account"), http.StatusUnprocessableEntity)
	}

	var cfgParsed []adm.FieldConfigurationParsed
	var errCfg error
	var hash uint32

	createAnalyticRecord := func(cfgHash uint32, rmqPublished bool, errReason *string) {
		analytic := &models.AnalyticRecord{
			UserUUID:     authToken.User.UserContact.UserUUID,
			DeviceID:     device.ID,
			CfgHash:      cfgHash,
			RMQPublished: rmqPublished,
			DeviceOnline: device.Status,
			ErrorReason:  errReason,
		}

		if err := h.Store.Analytics.Create_AnalyticsRecord(ctx, analytic); err != nil {
			logger.Error("ApplyConfigurationHandler_V1 req={%s}: Failed create analytics record: %s", ctx.Value("XREQID").(string), err.Error())
		}
	}

	writeAnalyticError := func(cfgHash uint32, message string) {
		reason := message
		createAnalyticRecord(cfgHash, false, &reason)
	}

	cfgRedisBytes, err := redis.GetCacheConfiguration(imei)
	if err == nil && cfgRedisBytes != nil {
		cfgParsed, errCfg = getConfiguration(tr, cfgRedisBytes)
		hash = admparser.GetCfgHash(cfgRedisBytes)
	}

	/* RabbitMQ */
	rabbitData := gtypes.RabbitMQ_TransitBinary{
		Type: constants.ADM_RC_TYPE_GET_CFG,
		Imei: imei,
		Data: nil,
	}

	if cfgParsed == nil || errCfg != nil {
		response, err := h.RMQ.SendToRabbitAndWait(ctx, rabbitData, h.Session)
		if err == nil && response.Data != nil {
			cfgParsed, errCfg = getConfiguration(tr, response.Data)
			hash = admparser.GetCfgHash(response.Data)
		}
	}

	if cfgParsed == nil || errCfg != nil {
		configuration, err := h.Store.Configurations.Get_ConfigurationByDeviceId(ctx, device.ID)
		if err != nil {
			message := fmt.Sprintf("ApplyConfigurationHandler_V1 req={%s}: Failed get configuration from db: %s", ctx.Value("XREQID").(string), err.Error())
			logger.Error(message)
			writeAnalyticError(hash, message)
			return httperr.Conflict(tr.TErr("config-fetch-error"))
		}
		if configuration.CfgData == nil {
			message := fmt.Sprintf("ApplyConfigurationHandler_V1 req={%s}: Configuration NULL", ctx.Value("XREQID").(string))
			logger.Error(message)
			writeAnalyticError(hash, message)
			return httperr.Conflict(tr.TErr("config-fetch-error"))
		}

		cfgParsed, errCfg = getConfiguration(tr, configuration.CfgData)
		if errCfg != nil {
			message := fmt.Sprintf("ApplyConfigurationHandler_V1 req={%s}: Failed get configuration parsed: %s", ctx.Value("XREQID").(string), errCfg.Error())
			logger.Error(message)
			writeAnalyticError(hash, message)
			return httperr.BadRequest(tr.TErr("cfg-parse-error-device"))
		}
		if cfgParsed == nil {
			message := fmt.Sprintf("ApplyConfigurationHandler_V1 req={%s}: Configuration parsed is NULL", ctx.Value("XREQID").(string))
			logger.Error(message)
			writeAnalyticError(hash, message)
			return httperr.BadRequest(tr.TErr("cfg-parse-error-device"))
		}

		hash = configuration.CfgHash
	}

	draft, err := redis.GetMergedDraft(authToken.User.UserContact.UserUUID, imei)
	if err != nil {
		message := fmt.Sprintf("ApplyConfigurationHandler_V1 req={%s}: %s", ctx.Value("XREQID").(string), err.Error())
		logger.Error(message)
		writeAnalyticError(hash, message)
		return httperr.InternalServerError(tr.TErr("failed-to-get-draft"))
	}

	if draft == nil {
		httpx.HttpResponse(w, r, http.StatusNoContent, nil)
		return nil
	}

	if hash != draft.CfgHash {
		configuration, err := h.Store.Configurations.Get_ConfigurationByDeviceId(ctx, device.ID)
		if err != nil {
			message := fmt.Sprintf("ApplyConfigurationHandler_V1 req={%s}: Failed get configuration from db: %s", ctx.Value("XREQID").(string), err.Error())
			logger.Error(message)
			writeAnalyticError(hash, message)
			return httperr.Conflict(tr.TErr("config-fetch-error"))
		}

		draftCreatedAt := time.Unix(draft.Timestamp, 0)
		if configuration.UpdatedAt.After(draftCreatedAt) {
			logger.Info("ApplyConfigurationHandler_V1 req={%s}: Configuration was modified by another user after draft creation: draft timestamp %v, config updated at %v (user: %s, imei: %s)", ctx.Value("XREQID").(string), draftCreatedAt, configuration.UpdatedAt, authToken.User.UserContact.UserUUID, imei)
			message := fmt.Sprintf("ApplyConfigurationHandler_V1 req={%s}: Configuration was modified by another user after draft creation: draft timestamp %v, config updated at %v (user: %s, imei: %s)", ctx.Value("XREQID").(string), draftCreatedAt, configuration.UpdatedAt, authToken.User.UserContact.UserUUID, imei)
			logger.Error(message)
			writeAnalyticError(hash, message)
			return httperr.Conflict(tr.TErr("error-cfg-changed-refresh"))
		}

		logger.Info("ApplyConfigurationHandler_V1 req={%s}: Updating draft hash for same user: old %d -> new %d (user: %s, imei: %s)", ctx.Value("XREQID").(string), draft.CfgHash, hash, authToken.User.UserContact.UserUUID, imei)
		/* Обновляем хеш в памяти для текущего применения */
		/* Draft будет удален после применения конфигурации, поэтому обновление в Redis не нужно */
		draft.CfgHash = hash
	}

	cfgConverted := convertDraftToParsedWithSize(draft, cfgParsed)
	adm := adm.Adm{}

	/* Build binary configuration */
	newCfg, err := adm.ParseJsonToBinary(tr, cfgConverted)
	if err != nil {
		message := fmt.Sprintf("ApplyConfigurationHandler_V1 req={%s}: Failed build configuration: %s", ctx.Value("XREQID").(string), err.Error())
		logger.Error(message)
		writeAnalyticError(hash, message)
		return httperr.BadRequest(err.Error())
	}

	if newCfg == nil {
		message := fmt.Sprintf("ApplyConfigurationHandler_V1 req={%s}: New builded configuration is NULL", ctx.Value("XREQID").(string))
		logger.Error(message)
		writeAnalyticError(hash, message)
		return httperr.BadRequest(tr.TErr("cfg-parse-error-device"))
	}

	/* Upd. database -> only configuration */
	if err := h.Store.Configurations.Update_ConfigurationByDeviceId(ctx, device.ID, newCfg); err != nil {
		message := fmt.Sprintf("ApplyConfigurationHandler_V1 req={%s}: Failed update configuration in db: %s", ctx.Value("XREQID").(string), err.Error())
		logger.Error(message)
		writeAnalyticError(hash, message)
		return httperr.Db(ctx, err)
	}

	if err := h.Store.Configurations.Update_ConfigurationSyncStatusByDeviceId(ctx, device.ID, constants.CFG_SYNC_STATUS_PENDING, nil); err != nil {
		message := fmt.Sprintf("ApplyConfigurationHandler_V1 req={%s}: Failed update configuration sync status: %s", ctx.Value("XREQID").(string), err.Error())
		logger.Error(message)
		writeAnalyticError(hash, message)
		return httperr.Db(ctx, err)
	}

	/* Upd. database -> insert parsed configurations */
	cfgHashApplied := admparser.GetCfgHash(newCfg)
	if profValues, err := adm.ParseProfValues(cfgConverted); err != nil {
		logger.Error("ApplyConfigurationHandler_V1 req={%s}: Failed build parsed configuration fields: %s", ctx.Value("XREQID").(string), err.Error())
	} else if err := h.Store.Configurations.Replace_ConfigurationProf(ctx, device.ID, cfgHashApplied, profValues); err != nil {
		logger.Error("ApplyConfigurationHandler_V1 req={%s}: Failed save parsed configuration: %s", ctx.Value("XREQID").(string), err.Error())
	}

	_ = redis.SetCacheConfiguration(imei, constants.Redis_ConfigurationTTL, newCfg)

	/* RabbitMQ */
	rabbitData2 := gtypes.RabbitMQ_TransitBinary{
		Type: constants.ADM_RC_TYPE_SET_CFG,
		Imei: device.IMEI,
		Data: newCfg,
	}

	if err := h.RMQ.SendToRabbitAsync(ctx, rabbitData2); err != nil {
		message := fmt.Sprintf("ApplyConfigurationHandler_V1 req={%s}: Failed publish configuration in rabbitmq: %s", ctx.Value("XREQID").(string), err.Error())
		logger.Error(message)
		reason := message
		createAnalyticRecord(admparser.GetCfgHash(newCfg), false, &reason)

		if statusErr := h.Store.Configurations.Update_ConfigurationSyncStatusByDeviceId(ctx, device.ID, constants.CFG_SYNC_STATUS_FAILED, &reason); statusErr != nil {
			logger.Error("ApplyConfigurationHandler_V1 req={%s}: Failed to persist configuration delivery error: %s", ctx.Value("XREQID").(string), statusErr.Error())
		}

		return httperr.New(tr.TErr("rabbitmq-publish-failed"), http.StatusConflict)
	}

	createAnalyticRecord(cfgHashApplied, true, nil)

	/* Reset all in user draft and reset configurations */
	if err := redis.DeleteAllDrafts(authToken.User.UserContact.UserUUID, imei); err != nil {
		message := fmt.Sprintf("ApplyConfigurationHandler_V1 req={%s}: Failed to delete drafts after applying configuration (user: %s, imei: %s): %s", ctx.Value("XREQID").(string), authToken.User.UserContact.UserUUID, imei, err.Error())
		logger.Error(message)
		reason := message
		createAnalyticRecord(admparser.GetCfgHash(newCfg), true, &reason)
		return httperr.InternalServerError(tr.TErr("failed-to-reset-draft"))
	}

	var resp string

	if device.Status {
		resp = tr.T("config-sent-successfully")
	} else {
		resp = tr.T("config-sent-status-false-successfully")
	}

	httpx.HttpResponse(w, r, http.StatusAccepted, resp)
	return nil
}
