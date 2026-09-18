package handlers

import (
	"encoding/hex"
	"neomatica/neosync-tcp/infra/constants"
	"neomatica/neosync-tcp/infra/logger"
	"neomatica/neosync-tcp/infra/store/redis"
	"neomatica/neosync-tcp/pkg"
	"neomatica/neosync-tcp/pkg/protocol/adm"
	tl "neomatica/neosync-tcp/pkg/telemetry"
	"neomatica/neosync-tcp/util"
	"net"
	"strings"
)

/* Initial command events workers */
func (h *BasePackEventHandler) CommandPackEventWorkers(n int) {
	for i := 0; i < n; i++ {
		pkg.RecoverGo("CommandPackEventJob", func() {
			for job := range h.QueueCommandEvent {
				h.CommandPackEventJob(job)
			}
		})
	}
}

/* Command event job */
func (h *BasePackEventHandler) CommandPackEventJob(job JobCommandEvent) {
	device := job.Device
	packet := job.Packet /* Answer type:string */

	if device == nil || device.Imei == nil {
		return
	}
	defer device.ReleaseTerminal()

	if packet == "" {
		return
	}

	imei := util.PrepareImei(*device.Imei)

	/* Get session by Imei */
	session, exists := h.Session.GetDeviceSession(*device.Imei)
	if session == nil || !exists {
		return
	}

	session.Mu.Lock()
	isTelemetry := session.InFlight != nil && session.InFlight.Cmd.Telemetry
	session.Mu.Unlock()

	// Periodic telemetry commands run on every sync and are intentionally not
	// part of the support event log. User/manual command answers are retained.
	if !isTelemetry {
		_ = redis.WriteAdmLog(constants.EVENT_DEVICE_RECEVCOMMAND, *device.Imei, map[string]any{
			"answer": packet,
		})
	}

	/* Check resp. telemetry */
	telemetrySave(imei, packet)

	h.Session.DeliverCommandAnswer(imei, packet)
}

func (h *BasePackEventHandler) CommandPackEventHandler(device *adm.ADMDevice, received any) {
	if device == nil || device.Imei == nil {
		logger.Warning("Device unavailable: cannot fetch events (device not connected)")
		return
	}

	answer, ok := received.(string)
	if !ok {
		logger.Error("CommandPackEventHandler ip={%s} imei={%s}: expected string payload, got %T", device.Conn.RemoteAddr().(*net.TCPAddr).IP, *device.Imei, received)
		return
	}

	/* Send packet to queue jobs */
	select {
	case h.QueueCommandEvent <- JobCommandEvent{Device: device, Packet: answer}:
	default:
		logger.Error("CommandPackEventHandler ip={%s} imei={%s}: Command events queue overflow", device.Conn.RemoteAddr().(*net.TCPAddr).IP, *device.Imei)
	}
}

func telemetrySave(imei, answer string) {
	tlStruct, _ := redis.GetTelemetry(imei)

	answer = strings.TrimSpace(answer)

	switch {
	case strings.HasPrefix(answer, "ADM20INFO:"):
		_ = tl.ParseADM20Info(tlStruct, answer)
	case strings.HasPrefix(answer, "BLESENSORINFO:"):
		_ = tl.ParseBLEString(tlStruct, answer)
	case strings.HasPrefix(answer, "FUELINFO:"):
		_ = tl.ParseFUELString(tlStruct, answer)
	default:
		packetBytes, err := hex.DecodeString(answer)
		if err != nil {
			break
		}

		normalized, normErr := tl.NormalizeADM6Packet(packetBytes)
		if normErr != nil {
			break
		}

		if len(normalized) < 3 || len(normalized) != int(normalized[2]) {
			break
		}

		if err := tl.ParseStatCountTelemetry(tlStruct, packetBytes); err != nil {
			logger.Warning("telemetrySave imei={%s}: failed to parse COM0 packet: %s", imei, err.Error())
		}
	}

	_ = redis.SaveTelemetry(imei, tlStruct)
}
