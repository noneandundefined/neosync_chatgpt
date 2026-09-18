package configuration_handler_v1

import (
	"fmt"
	"neomatica/neosync/config"
	"neomatica/neosync/infra/constants"
	"neomatica/neosync/infra/locale"
	"neomatica/neosync/infra/logger"
	"neomatica/neosync/middleware"
	"neomatica/neosync/pkg/httpx/httperr"
	gtypes "neomatica/neosync/types"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

/* Neosync HTTPx V1 */
/* Handler: получение конфигурации шаблона по секции и модели (значения Default из схемы) */

func (h *Handler) GetConfigurationTemplateParsedBySectionHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)
	authToken := ctx.Value("identity").(*gtypes.AuthToken)

	if !authToken.User.Has(constants.AccessConfigurationRead) {
		return httperr.Forbidden(tr.TErr("config-access-restricted"))
	}

	lang := r.URL.Query().Get("lang")
	if lang != "" {
		tr = locale.NewTranslator(lang)
	}

	section := mux.Vars(r)["section"]

	if section == "" {
		return httperr.BadRequest(tr.TErr("device-section-not-found"))
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	flusher, ok := w.(http.Flusher)
	if !ok {
		return httperr.Conflict("streaming unsupported")
	}

	sendSSEMessage(w, flusher, "open", tr.T("sse-connected-message"))

	select {
	case <-ctx.Done():
		return nil
	default:
	}

	sendSSEMessage(w, flusher, "progress", tr.T("loading-config-message"))

	filteredCfg, schemaList := getDefaultConfigurationBySection(section)

	cfgResponse := ConfigurationTemplateSectionParsed{
		CfgHash:      0,
		Section:      section,
		ConfigParsed: filteredCfg,
		Schema:       schemaList,
		Timestamp:    time.Now().Unix(),
	}

	respJson, err := config.JSON.Marshal(cfgResponse)
	if err != nil {
		logger.Error("GetConfigurationTemplateParsedBySectionHandler_V1 req={%s}: Failed marshal JSON: %s", ctx.Value("XREQID").(string), err.Error())
		sendSSEMessage(w, flusher, "error", tr.TErr("config-fetch-error"))
		return err
	}

	fmt.Fprintf(w, "event: configuration\ndata: %s\n\n", respJson)
	flusher.Flush()

	sendSSEMessage(w, flusher, "done", tr.T("config-loaded-message"))

	<-ctx.Done()
	return nil
}
