package command_handler_v1

import (
	"context"
	"neomatica/neosync/infra/constants"
	"neomatica/neosync/infra/logger"
	"neomatica/neosync/middleware"
	"neomatica/neosync/pkg/httpx"
	"neomatica/neosync/pkg/httpx/httperr"
	gtypes "neomatica/neosync/types"
	"net/http"
	"sync"
	"time"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
)

var bleTasks = struct {
	m map[string]*BleTask
	sync.Mutex
}{m: make(map[string]*BleTask)}

func (h *Handler) CommandFindBleSensorsHandler_V1(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	tr := middleware.TranslatorFromContext(ctx)

	var payload *CommandWaitPayload

	imei := mux.Vars(r)["imei"]
	if imei == "" {
		return httperr.NotFound(tr.TErr("device-imei-not-found"))
	}

	if err := httpx.HttpParse(r, &payload); err != nil {
		return httperr.BadRequest(err.Error())
	}

	if err := httpx.Validate.Struct(payload); err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			return httperr.BadRequest(httpx.ValidateMsg(tr, err))
		}

		return httperr.BadRequest(tr.TErr("fields-not-filled"))
	}

	switch payload.Command {
	case constants.BLEAUTOCATCH_FIND_ALL_COMMAND, constants.BLEAUTOCATCH_FIND_NEARBY_COMMAND:
		bleTasks.Lock()
		if _, exists := bleTasks.m[imei]; exists {
			bleTasks.Unlock()
			return httperr.BadRequest(tr.TErr("scan-already-running-for-imei"))
		}

		ctx, cancel := context.WithCancel(context.Background())
		bleTasks.m[imei] = &BleTask{ctx: ctx, cancel: cancel}
		bleTasks.Unlock()

		go func() {
			const (
				totalTimeout = 3 * time.Minute
				retryDelay   = 20 * time.Second
			)

			ctxTask := ctx
			deadline := time.Now().Add(totalTimeout)

			rabbitData := gtypes.RabbitMQ_TransitBinary{
				Type:  constants.ADM_RC_TYPE_STRING,
				Imei:  imei,
				NResp: true,
				Data:  []byte(payload.Command),
			}

			for time.Now().Before(deadline) {
				select {
				case <-ctxTask.Done():
					logger.Info("Search BLE for (%s) stopped manually", imei)
					bleTasks.Lock()
					delete(bleTasks.m, imei)
					bleTasks.Unlock()
					return
				default:
				}

				_, err := h.RMQ.SendToRabbitAndWait(ctxTask, rabbitData, h.Session)
				if err == nil {
					logger.Info("Searching BLE sensors for (%s) still running...", imei)
				} else if err.Error() == tr.TErr("response-timeout-13s") {
					time.Sleep(retryDelay)
					continue
				} else {
					logger.Error("CommandFindBleSensorsHandler_V1: error by searching BLE sensors for (%s): %v", imei, err)
					return
				}

				time.Sleep(retryDelay)
			}

			logger.Info("Time search BLE sensors for (%s) is out", imei)
			bleTasks.Lock()
			delete(bleTasks.m, imei)
			bleTasks.Unlock()
		}()

		httpx.HttpResponse(w, r, http.StatusOK, tr.T("ble-scan-started"))

	case constants.BLEAUTOCATCH_STOP_FIND_COMMAND:
		bleTasks.Lock()
		if task, exists := bleTasks.m[imei]; exists {
			task.cancel()
			delete(bleTasks.m, imei)
			bleTasks.Unlock()
			httpx.HttpResponse(w, r, http.StatusOK, tr.T("ble-scan-stopped"))
		} else {
			bleTasks.Unlock()
			httpx.HttpResponse(w, r, http.StatusOK, tr.T("ble-scan-not-running-for-imei"))
		}

	default:
		return httperr.BadRequest(tr.TErr("unknown-command"))
	}

	return nil
}
