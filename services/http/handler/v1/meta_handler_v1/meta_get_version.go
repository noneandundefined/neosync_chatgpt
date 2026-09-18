package meta_handler_v1

import (
	"neomatica/neosync/infra/logger"
	"neomatica/neosync/pkg/httpx"
	"neomatica/neosync/pkg/httpx/httperr"
	"net/http"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

func (h *Handler) MetaGetVersionHandler_V1(w http.ResponseWriter, r *http.Request) error {
	data, err := os.ReadFile("config/releases.yml")
	if err != nil {
		logger.Error("MetaGetVersionHandler_V1: %s", err.Error())
		httperr.InternalServerError(err.Error())
	}

	var releases ReleaseFile
	if err := yaml.Unmarshal(data, &releases); err != nil {
		logger.Error("MetaGetVersionHandler_V1: %s", err.Error())
		httperr.InternalServerError(err.Error())
	}

	if len(releases.Releases) == 0 {
		httpx.HttpResponse(w, r, http.StatusOK, "Neosync version")
		return nil
	}

	var latestRelease *Release
	for i := range releases.Releases {
		releaseDate, err := time.Parse("02.01.2006", releases.Releases[i].Date)
		if err != nil {
			continue
		}

		if latestRelease == nil {
			latestRelease = &releases.Releases[i]
			continue
		}

		latestDate, _ := time.Parse("02.01.2006", latestRelease.Date)
		if releaseDate.After(latestDate) {
			latestRelease = &releases.Releases[i]
		}
	}

	if latestRelease == nil {
		httpx.HttpResponse(w, r, http.StatusOK, "Neosync version")
		return nil
	}

	httpx.HttpCache(w, 21600) // 6 hour.
	httpx.HttpResponse(w, r, http.StatusOK, latestRelease.Version)
	return nil
}
