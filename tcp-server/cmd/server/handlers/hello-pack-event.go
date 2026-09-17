package handlers

import (
	"context"
	"neomatica/neosync-tcp/infra/constants"
	"neomatica/neosync-tcp/infra/logger"
	"neomatica/neosync-tcp/infra/store/postgres/models"
	"neomatica/neosync-tcp/infra/store/postgres/usecase"
	"neomatica/neosync-tcp/infra/store/redis"
	"neomatica/neosync-tcp/pkg"
	"neomatica/neosync-tcp/pkg/protocol"
	"neomatica/neosync-tcp/pkg/protocol/adm"
	"neomatica/neosync-tcp/types"
	"neomatica/neosync-tcp/util"
	"time"

	"github.com/google/uuid"
)

/* Initial hello events workers */
func (h *BasePackEventHandler) HelloPackEventWorkers(n int) {
	for i := 0; i < n; i++ {
		pkg.RecoverGo("HelloPackEventJob", func() {
			for job := range h.QueueHelloEvent {
				h.HelloPackEventJob(job)
			}
		})
	}
}

/* Hello event job */
func (h *BasePackEventHandler) HelloPackEventJob(job JobHelloEvent) {
	device := job.Device
	packet := job.Packet

	if device == nil || packet == nil || packet.Imei == nil {
		return
	}

	if device.Imei == nil {
		device.Imei = packet.Imei
	}

	imei := util.PrepareImei(*packet.Imei)
	ip := util.RemoteIP(device.Conn)
	imeiLog := util.ImeiString(packet.Imei)

	/* Sets base information device part.2 */
	device.Pass = packet.Pass
	device.FirmwareVersion = packet.FirmwareVersion
	device.CfgVersion = packet.CfgVersion
	if packet.LastModTime == 0 {
		device.LastModTime = uint32(time.Now().Unix())
	} else {
		device.LastModTime = packet.LastModTime
	}

	var d *models.Device_DeviceConf_Sync
	var err error

	h.limiterDatabase(func() {
		d, err = h.Store.Devices.Get_DeviceFullByImei(context.Background(), imei)
		if err != nil {
			return
		}

		device.CfgHash = packet.CfgHash

		/* Create device in database */
		if err = job.UseCase.Devices.Create_Device(); err != nil {
			return
		}

		/* if device(d) is NULL, get again Get_DeviceFullByImei */
		if d == nil {
			d, err = h.Store.Devices.Get_DeviceFullByImei(context.Background(), imei)
			if err != nil || d == nil {
				return
			}
		}

		/* Sets information device part.3 */
		device.DeviceId = d.ID
		device.UserUuid = d.UserUUID
		device.RCOnC = d.RequestConfigurationOnConnect

		/* Create session for device */
		h.Session.NewDeviceSession(imei, device)
	})

	if err != nil {
		logger.Error("HelloPackEventJob ip={%s} imei={%s}: Failed to prepare device: %s", ip, imeiLog, err.Error())
		return
	}

	if d == nil {
		logger.Error("HelloPackEventJob ip={%s} imei={%s}: Failed to retrieve device after creation", ip, imeiLog)
		return
	}

	_ = redis.WriteAdmLog(constants.EVENT_DEVICE_CONN, imei, map[string]any{
		"device_id":     d.ID,
		"cfg_hash":      packet.CfgHash,
		"last_mod_time": packet.LastModTime,
	})

	currentSession, exists := h.Session.GetDeviceSession(imei)
	if !exists || currentSession.Device != device {
		return
	}
	currentSession.ConfigurationMu.Lock()
	defer currentSession.ConfigurationMu.Unlock()

	hasModel := d.DeviceModel != nil && *d.DeviceModel != ""
	hasExtendedModel := d.DeviceExtendedModel != nil && *d.DeviceExtendedModel != ""

	if !hasModel || !hasExtendedModel {
		h.Session.EnqueueCommand(types.RabbitMQ_TransitBinary{
			Type:      constants.ADM_RC_TYPE_STRING,
			RequestID: uuid.NewString(),
			Imei:      d.IMEI,
			NResp:     true,
			Data:      []byte(constants.WHO_COMMAND),
		})
	} else {
		device.MarkV2DeviceModelLookupSatisfied()
	}

	/* Если в БД нет конфигурации — запросить с трекера. Иначе RCOnC — опциональный запрос. */
	bulkHandled := false

	h.limiterDatabase(func() {
		var bulkErr error

		bulkObserved, observeErr := job.UseCase.CompanyTasks.ObserveDeviceConfiguration(context.Background(), d.ID, packet.CfgHash)
		if observeErr != nil {
			err = observeErr
			bulkHandled = true
			return
		}

		bulkHandled, bulkErr = job.UseCase.CompanyTasks.TryApplyPendingTask(context.Background(), h.Session, d)
		bulkHandled = bulkHandled || bulkObserved
		if bulkErr != nil {
			err = bulkErr
		}
	})

	if err != nil {
		logger.Error("HelloPackEventJob ip={%s} imei={%s}: Failed to process bulk configuration: %s", ip, imeiLog, err.Error())
	}

	if !bulkHandled {
		configMissing := usecase.IsConfigurationMissing(d)

		if configMissing {
			if packet.CfgHash != 0 {
				h.Session.TransitData(d.IMEI, constants.ADM_RC_TYPE_GET_CFG, nil)
			}
		} else if !d.ForceNeosyncConfigurationPriority && device.RCOnC && device.LastModTime >= uint32(d.CfgUpdatedAt.Unix()) {
			h.Session.TransitData(d.IMEI, constants.ADM_RC_TYPE_GET_CFG, nil)
		}

		if !configMissing && packet.CfgHash != 0 {
			if syncErr := job.UseCase.Syncs.Synchronization(h.Session, d); syncErr != nil {
				logger.Error("HelloPackEventJob ip={%s} imei={%s}: Failed to pull/send configuration: %s", ip, imeiLog, syncErr.Error())
			}
		}
	}

	var commands []models.DeviceCommandWithExecutions
	var cmdErr error

	h.limiterDatabase(func() {
		if clearErr := h.Store.DeviceCommands.Update_DeviceCommandClearProgressStatus(context.Background(), imei); clearErr != nil {
			logger.Error("HelloPackEventJob imei={%s}: timeout cleanup error: %s", imei, clearErr.Error())
		}

		commands, cmdErr = h.Store.DeviceCommands.Get_DeviceCommandsByImeiAndSendMode(context.Background(), imei, constants.CMD_SEND_MODE_ON_CONNECT)
	})

	if cmdErr != nil {
		logger.Error("HelloPackEventJob imei={%s}: get commands error: %s", imei, cmdErr.Error())
		return
	}

	for _, cmd := range commands {
		h.Session.EnqueueCommand(types.RabbitMQ_TransitBinary{
			Type:        constants.ADM_RC_TYPE_STRING,
			RequestID:   uuid.NewString(),
			Imei:        cmd.IMEI,
			NResp:       false,
			Data:        []byte(cmd.Command),
			ExecutionID: cmd.ExecutionID,
		})
	}
}

func (h *BasePackEventHandler) HelloPackEventHandler(device *adm.ADMDevice, usecase *usecase.UseCase, received any) {
	packet, ok := received.(*protocol.WelcomePacket)
	if !ok || packet == nil || packet.Imei == nil {
		logger.Warning("Device unavailable: cannot fetch events (device not connected)")
		return
	}

	/* Sets base information device part.1 */
	device.Imei = packet.Imei

	/* Send packet to queue jobs */
	select {
	case h.QueueHelloEvent <- JobHelloEvent{Device: device, UseCase: usecase, Packet: packet}:
	default:
		logger.Error("HelloPackEventHandler ip={%s} imei={%s}: Hello events queue overflow", util.RemoteIP(device.Conn), util.ImeiString(packet.Imei))
	}
}
