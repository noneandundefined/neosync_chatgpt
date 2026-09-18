package configuration_handler_v1

import (
	"encoding/binary"
	"errors"
	"neomatica/neosync/infra/constants"
	"neomatica/neosync/infra/locale"
	"neomatica/neosync/infra/logger"
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

func getConfiguration(tr locale.Translator, configuration []byte) ([]adm.FieldConfigurationParsed, error) {
	adm := adm.Adm{}

	if len(configuration) < 4 {
		return nil, nil
	}

	size := binary.LittleEndian.Uint16(configuration[0:2])

	if int(size) != len(configuration) {
		return nil, errors.New(tr.TErr("cfg-parse-error-device"))
	}

	cfgParsed, errCfg := adm.ParseBinaryToJson(tr, configuration)
	if errCfg != nil {
		return nil, errCfg
	}

	return cfgParsed, nil
}

/* Neosync HTTPx V1 */
/* Handler: получение всей конфигурации по Imei устройства */

func (h *Handler) GetConfigurationParsedByImeiHandler_V1(w http.ResponseWriter, r *http.Request) error { //nolint
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*gtypes.AuthToken)

	/* Check access */
	if !authToken.User.Has(constants.AccessConfigurationRead) {
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

	if !permissions.CanAccessDevice(ctx, h.Store.GroupMembers, device, device.ID, &authToken.User, authToken.RoleCode, permissions.ShareReadConfig) {
		return httperr.Forbidden(tr.TErr("device-not-owned"))
	}

	if !*device.Activated {
		return httperr.New(tr.TErr("device-not-activated"), http.StatusUnprocessableEntity)
	}

	var cfgParsed []adm.FieldConfigurationParsed
	var errCfg error
	var hash uint32

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
			_ = redis.SetCacheConfiguration(imei, constants.Redis_ConfigurationTTL, response.Data)

			cfgParsed, errCfg = getConfiguration(tr, response.Data)
			hash = admparser.GetCfgHash(response.Data)
		}
	}

	if cfgParsed == nil || errCfg != nil {
		configuration, err := h.Store.Configurations.Get_ConfigurationByDeviceId(ctx, device.ID)
		if err != nil {
			logger.Error("ApplyConfigurationHandler_V1 req={%s}: Failed get configuration from db: %s", ctx.Value("XREQID").(string), err.Error())
			return httperr.Conflict(tr.TErr("config-fetch-error"))
		}
		if configuration.CfgData == nil {
			logger.Error("ApplyConfigurationHandler_V1 req={%s}: Configuration NULL", ctx.Value("XREQID").(string))
			return httperr.Conflict(tr.TErr("config-fetch-error"))
		}

		cfgParsed, errCfg = getConfiguration(tr, configuration.CfgData)
		if errCfg != nil {
			logger.Error("ApplyConfigurationHandler_V1 req={%s}: Failed get configuration parsed: %s", ctx.Value("XREQID").(string), errCfg.Error())
			return httperr.BadRequest(tr.TErr("cfg-parse-error-device"))
		}
		if cfgParsed == nil {
			logger.Error("ApplyConfigurationHandler_V1 req={%s}: Configuration parsed is NULL", ctx.Value("XREQID").(string))
			return httperr.BadRequest(tr.TErr("cfg-parse-error-device"))
		}

		_ = redis.SetCacheConfiguration(imei, constants.Redis_ConfigurationTmpTTL, configuration.CfgData)
		hash = configuration.CfgHash
	}

	var schemaList []constants.FieldSchema
	nameSet := make(map[string]struct{})

	for _, field := range cfgParsed {
		nameSet[field.UID] = struct{}{}
	}

	for _, schema := range constants.CfgSchema {
		if _, ok := nameSet[schema.Name]; ok {
			schemaList = append(schemaList, schema)
		}
	}

	cfgResponse := ConfigurationParsed{
		DeviceImei:   imei,
		CfgHash:      hash,
		ConfigParsed: cfgParsed,
		Schema:       schemaList,
		Timestamp:    time.Now().Unix(),
	}

	httpx.HttpResponseWithETag(w, r, http.StatusOK, cfgResponse)
	return nil
}
