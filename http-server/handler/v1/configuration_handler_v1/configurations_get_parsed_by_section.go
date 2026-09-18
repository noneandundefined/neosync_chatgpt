package configuration_handler_v1

import (
	"fmt"
	"neomatica/neosync/config"
	"neomatica/neosync/infra/constants"
	"neomatica/neosync/infra/locale"
	"neomatica/neosync/infra/logger"
	"neomatica/neosync/infra/store/redis"
	"neomatica/neosync/middleware"
	"neomatica/neosync/pkg/adm"
	"neomatica/neosync/pkg/adm/admparser"
	"neomatica/neosync/pkg/httpx/httperr"
	"neomatica/neosync/pkg/permissions"
	gtypes "neomatica/neosync/types"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

func getConfigurationBySection(section string, cfgParsed []adm.FieldConfigurationParsed) ([]adm.FieldConfigurationParsed, []constants.FieldSchema) {
	var schemaList []constants.FieldSchema
	uidSet := make(map[string]struct{})

	for _, schema := range constants.CfgSchema {
		if schema.Section == "" {
			continue
		}

		sections := strings.Split(schema.Section, ";")

		for _, s := range sections {
			if strings.EqualFold(strings.TrimSpace(s), section) {
				schemaList = append(schemaList, schema)
				uidSet[schema.Name] = struct{}{}
				break
			}
		}
	}

	var filteredCfg []adm.FieldConfigurationParsed
	for _, cfg := range cfgParsed {
		if _, ok := uidSet[cfg.UID]; ok {
			filteredCfg = append(filteredCfg, cfg)
		}
	}

	return filteredCfg, schemaList
}

func sendSSEMessage(w http.ResponseWriter, flusher http.Flusher, eventType string, message string) {
	jsonMessage, err := config.JSON.Marshal(map[string]string{"message": message})
	if err != nil {
		fmt.Printf("Error marshaling SSE message: %v\n", err)
		jsonMessage = []byte("{\"message\": \"\"}")
	}

	fmt.Fprintf(w, "event: %s\ndata: %s\n\n", eventType, jsonMessage)
	flusher.Flush()
}

/* Neosync HTTPx V1 */
/* Handler: получение конфигурации по секции, по Imei устройства */

func (h *Handler) GetConfigurationParsedBySectionHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*gtypes.AuthToken)

	/* Check access */
	if !authToken.User.Has(constants.AccessConfigurationRead) {
		return httperr.Forbidden(tr.TErr("config-access-restricted"))
	}

	lang := r.URL.Query().Get("lang")
	if lang != "" {
		tr = locale.NewTranslator(lang)
	}

	imei := mux.Vars(r)["imei"]
	section := mux.Vars(r)["section"]

	if imei == "" {
		return httperr.BadRequest(tr.TErr("device-imei-not-found"))
	}

	if section == "" {
		return httperr.BadRequest(tr.TErr("device-section-not-found"))
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

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	flusher, ok := w.(http.Flusher)
	if !ok {
		return httperr.Conflict("streaming unsupported")
	}

	/* SSE event. connected */
	sendSSEMessage(w, flusher, "open", tr.T("sse-connected-message"))

	select {
	case <-ctx.Done():
		return nil
	default:
	}

	var cfgParsed []adm.FieldConfigurationParsed
	var errCfg error
	var hash uint32

	/* SSE event. check cached configurations */
	sendSSEMessage(w, flusher, "progress", tr.T("checking-config-cache-message"))

	cfgRedisBytes, err := redis.GetCacheConfiguration(imei)
	if err == nil && cfgRedisBytes != nil {
		/* SSE event. success get cached configurations */
		sendSSEMessage(w, flusher, "progress", tr.T("config-found-in-cache-message"))

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
		/* SSE event. getting configuration in terminal */
		sendSSEMessage(w, flusher, "progress", tr.T("requesting-config-from-tracker-message"))

		response, err := h.RMQ.SendToRabbitAndWait(ctx, rabbitData, h.Session)
		if err == nil && response.Data != nil {
			/* SSE event. success get configuration in terminal */
			sendSSEMessage(w, flusher, "progress", tr.T("config-received-from-tracker-message"))

			_ = redis.SetCacheConfiguration(imei, constants.Redis_ConfigurationTTL, response.Data)

			cfgParsed, errCfg = getConfiguration(tr, response.Data)
			hash = admparser.GetCfgHash(response.Data)
		}
	}

	if cfgParsed == nil || errCfg != nil {
		/* SSE event. getting configuration in database */
		sendSSEMessage(w, flusher, "progress", tr.T("retrieving-config-from-db-message"))

		configuration, err := h.Store.Configurations.Get_ConfigurationByDeviceId(ctx, device.ID)
		if err != nil {
			logger.Error("GetConfigurationParsedBySectionHandler_V1 req={%s} imei={%s}: Failed get configuration from db: %s", ctx.Value("XREQID").(string), imei, err.Error())
			sendSSEMessage(w, flusher, "error", tr.TErr("config-fetch-error"))
			return nil
		}
		if configuration == nil || configuration.CfgData == nil {
			sendSSEMessage(w, flusher, "error", tr.TErr("config-fetch-error"))
			return nil
		}

		cfgParsed, errCfg = getConfiguration(tr, configuration.CfgData)
		if errCfg != nil {
			logger.Error("GetConfigurationParsedBySectionHandler_V1 req={%s} imei={%s}: Failed get configuration parsed: %s", ctx.Value("XREQID").(string), imei, errCfg.Error())
			sendSSEMessage(w, flusher, "error", tr.TErr("cfg-parse-error-device"))
			return nil
		}
		if cfgParsed == nil {
			logger.Error("GetConfigurationParsedBySectionHandler_V1 req={%s} imei={%s}: Configuration parsed is NULL", ctx.Value("XREQID").(string), imei)
			sendSSEMessage(w, flusher, "error", tr.TErr("cfg-parse-error-device"))
			return nil
		}

		/* SSE event. success get configuration in database */
		sendSSEMessage(w, flusher, "progress", tr.T("config-received-from-db-message"))

		_ = redis.SetCacheConfiguration(imei, constants.Redis_ConfigurationTmpTTL, configuration.CfgData)
		hash = configuration.CfgHash
	}

	/* SSE event. loading prepare configuration for view */
	sendSSEMessage(w, flusher, "progress", tr.T("loading-config-message"))

	filteredCfg, schemaList := getConfigurationBySection(section, cfgParsed)

	cfgResponse := ConfigurationSectionParsed{
		DeviceImei:   imei,
		CfgHash:      hash,
		Section:      section,
		ConfigParsed: filteredCfg,
		Schema:       schemaList,
		Timestamp:    time.Now().Unix(),
	}

	respJson, err := config.JSON.Marshal(cfgResponse)
	if err != nil {
		logger.Error("GetConfigurationParsedBySectionHandler_V1 req={%s} imei={%s}: Failed marshal JSON: %s", ctx.Value("XREQID").(string), imei, err.Error())
		sendSSEMessage(w, flusher, "error", tr.TErr("config-fetch-error"))
		return err
	}

	/* Event configuratons */
	fmt.Fprintf(w, "event: configuration\ndata: %s\n\n", respJson)
	flusher.Flush()

	/* Event ok. before close sse */
	sendSSEMessage(w, flusher, "done", tr.T("config-loaded-message"))

	<-ctx.Done()
	return nil
}
