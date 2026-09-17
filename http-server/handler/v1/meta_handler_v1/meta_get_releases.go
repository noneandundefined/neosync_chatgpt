package meta_handler_v1

import (
	"neomatica/neosync/infra/logger"
	"neomatica/neosync/middleware"
	"neomatica/neosync/pkg/httpx"
	"neomatica/neosync/pkg/httpx/httperr"
	"net/http"
	"os"

	"gopkg.in/yaml.v2"
)

func (h *Handler) MetaGetReleasesHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)

	lang := tr.GetLang()
	if lang == "" {
		lang = "en"
	}

	data, err := os.ReadFile("config/releases.yml")
	if err != nil {
		logger.Error("MetaGetReleasesHandler_V1: %s", err.Error())
		httperr.InternalServerError(err.Error())
	}

	var releases ReleaseFile
	if err := yaml.Unmarshal(data, &releases); err != nil {
		logger.Error("MetaGetReleasesHandler_V1: %s", err.Error())
		httperr.InternalServerError(err.Error())
	}

	if len(releases.Releases) > 100 {
		releases.Releases = releases.Releases[len(releases.Releases)-100:]
	}

	httpx.HttpCache(w, 21600) // 6 hour.
	httpx.HttpResponseWithETag(w, r, http.StatusOK, releases)
	return nil
}
