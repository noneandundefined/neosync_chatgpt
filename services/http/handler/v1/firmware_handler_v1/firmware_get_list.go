package firmware_handler_v1

import (
	"bufio"
	"fmt"
	"math"
	"neomatica/neosync/infra/locale"
	"neomatica/neosync/infra/store/postgres/models"
	"neomatica/neosync/middleware"
	"neomatica/neosync/pkg/httpx"
	"neomatica/neosync/pkg/httpx/httperr"
	"neomatica/neosync/util"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

func fetchFotaForModels(tr locale.Translator, sources []models.DeviceModelSource) []Fota {
	var wg sync.WaitGroup
	var mu sync.Mutex
	latestFotas := make(map[string]Fota)

	for _, src := range sources {
		wg.Add(1)
		go func(m models.DeviceModelSource) {
			defer wg.Done()

			url := fmt.Sprintf("https://neomatica.com/upload/docs/changelog/%s/%s.txt", tr.GetLang(), m.DeviceModel)
			resp, err := http.Get(url)
			if err != nil {
				return
			}
			defer resp.Body.Close()

			scanner := bufio.NewScanner(resp.Body)
			re := regexp.MustCompile(`^0x([0-9A-Fa-f]+)\s*-\s*(\d{2}\.\d{2}\.\d{4})$`)

			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())
				match := re.FindStringSubmatch(line)
				if match != nil {
					releaseDate, err := time.Parse("02.01.2006", match[2])
					if err != nil {
						continue
					}

					f := Fota{
						CreatedAt:   time.Now(),
						DeviceModel: m.DeviceModel,
						Firmware:    "0x" + strings.ToUpper(match[1]),
						ReleaseDate: match[2],
						Type:        "release",
					}

					mu.Lock()
					if existing, ok := latestFotas[f.DeviceModel]; !ok {
						latestFotas[f.DeviceModel] = f
					} else {
						existingDate, _ := time.Parse("02.01.2006", existing.ReleaseDate)
						if releaseDate.After(existingDate) {
							latestFotas[f.DeviceModel] = f
						}
					}
					mu.Unlock()
				}
			}
		}(src)
	}

	wg.Wait()

	result := make([]Fota, 0, len(latestFotas))
	for _, f := range latestFotas {
		result = append(result, f)
	}

	return result
}

/* Neosync HTTPx V1 */
/* Handler: получение истории прошивок устройств */

func (h *Handler) GetFirmwareListHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)

	pagination := util.GetPagination(r, 10)

	sources, err := h.Store.DeviceModelSources.Get_Sources(ctx)
	if err != nil {
		return httperr.Db(ctx, err)
	}

	if len(sources) == 0 {
		httpx.HttpResponse(w, r, http.StatusOK, nil)
		return nil
	}

	fotas := fetchFotaForModels(tr, sources)

	start := pagination.Offset
	end := start + pagination.Limit
	if start > len(fotas) {
		start = len(fotas)
	}
	if end > len(fotas) {
		end = len(fotas)
	}
	pagedFotas := fotas[start:end]

	total := len(fotas)
	totalPages := int(math.Ceil(float64(total) / float64(pagination.Limit)))

	resp := FotasWPResponse{
		Items:      pagedFotas,
		Page:       pagination.Page,
		Limit:      pagination.Limit,
		Total:      total,
		TotalPages: totalPages,
	}

	httpx.HttpResponseWithETag(w, r, http.StatusOK, resp)
	return nil
}
