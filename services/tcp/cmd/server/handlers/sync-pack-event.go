package handlers

import (
	"context"
	"neomatica/neosync-tcp/infra/constants"
	"neomatica/neosync-tcp/infra/logger"
	"neomatica/neosync-tcp/infra/store/postgres/models"
	"neomatica/neosync-tcp/internal/usecase"
	"neomatica/neosync-tcp/infra/store/redis"
	"neomatica/neosync-tcp/pkg/protocol"
	"neomatica/neosync-tcp/pkg/protocol/adm"
	"neomatica/neosync-tcp/types"
	"neomatica/neosync-tcp/util"
	"time"

	"github.com/google/uuid"
)

func (h *BasePackEventHandler) SyncPackEventHandler(device *adm.ADMDevice, uc *usecase.UseCase, received any) {
	packet, ok := received.(*protocol.SyncPacket)
	if !ok || device == nil || device.Imei == nil {
		logger.Warning("Device unavailable: cannot fetch events (device not connected)")
		return
	}

	currentSession, exists := h.Session.GetDeviceSession(*device.Imei)
	if !exists || currentSession.Device != device {
		return
	}
	currentSession.ConfigurationMu.Lock()
	defer currentSession.ConfigurationMu.Unlock()

	/* Sets base information device */
	device.FirmwareVersion = packet.FirmwareVersion
	device.CfgVersion = packet.CfgVersion
	if packet.LastModTime == 0 {
		device.LastModTime = uint32(time.Now().Unix())
	} else {
		device.LastModTime = packet.LastModTime
	}

	ip := util.RemoteIP(device.Conn)
	imei := util.PrepareImei(*device.Imei)

	var d *models.Device_DeviceConf_Sync
	var err error
	var bulkObserved bool

	h.limiterDatabase(func() {
		d, err = h.Store.Devices.Get_DeviceFullByImei(context.Background(), imei)
	})

	if err != nil {
		logger.Error("SyncPackEventHandler ip={%s} imei={%s}: Failed to get device: %s", ip, imei, err.Error())
		return
	}

	if d == nil {
		return
	}

	device.CfgHash = packet.CfgHash

	// check update device
	// -------------------
	upd, err := redis.GetUpdDevice(*device.Imei)
	if err != nil {
		logger.Error("SyncPackEventHandler ip={%s} imei={%s}: Failed to get UPDATE mode for device: %s", ip, imei, err.Error())
	}

	if upd != nil {
		h.Session.EnqueueCommand(types.RabbitMQ_TransitBinary{
			Type:      constants.ADM_RC_TYPE_STRING,
			RequestID: uuid.NewString(),
			Imei:      *device.Imei,
			NResp:     false,
			Data:      []byte(*upd),
		})

		if err := redis.DeleteUpdDevice(*device.Imei); err != nil {
			logger.Error("SyncPackEventHandler ip={%s} imei={%s}: Failed to delete Redis key for device: %s", ip, imei, err.Error())
		}
	}
	// -------------------
	// check update device

	_ = redis.WriteAdmLog(constants.EVENT_DEVICE_SYNC, *device.Imei, map[string]any{
		"device_id":        device.DeviceId,
		"firmware_version": device.FirmwareVersion,
		"cfg_version":      device.CfgVersion,
		"last_mod_time":    device.LastModTime,
		"cfg_hash":         device.CfgHash,
	})

	h.limiterDatabase(func() {
		if syncErr := h.Store.Syncs.Update_Sync(context.Background(), &models.Sync{
			DeviceID:        d.ID,
			FirmwareVersion: device.FirmwareVersion,
			CfgVersion:      device.CfgVersion,
			LastModTime:     device.LastModTime,
			CfgHash:         device.CfgHash,
		}); syncErr != nil {
			logger.Error("SyncPackEventHandler ip={%s} imei={%s}: %s", ip, imei, syncErr.Error())
		}

		if upd == nil {
			if fwErr := h.Store.Syncs.Complete_FirmwareUpdateByDeviceId(context.Background(), d.ID); fwErr != nil {
				logger.Error("SyncPackEventHandler ip={%s} imei={%s}: Failed to complete firmware update: %s", ip, imei, fwErr.Error())
			}
		}

		bulkActive, bulkErr := uc.CompanyTasks.ObserveDeviceConfiguration(context.Background(), d.ID, packet.CfgHash)
		if bulkErr != nil {
			err = bulkErr
			return
		}

		bulkObserved = bulkActive
		if !bulkActive {
			if syncErr := uc.Syncs.ResolveConfigurationSyncStatus(d, device.CfgHash, true); syncErr != nil {
				logger.Error("SyncPackEventHandler imei={%s}: Failed to resolve configuration sync status: %s", imei, syncErr.Error())
			}
		}

		// if clearErr := h.Store.DeviceCommands.Update_DeviceCommandClearProgressStatus(context.Background(), imei); clearErr != nil {
		// 	logger.Error("SyncPackEventHandler ip={%s} imei={%s}: timeout cleanup error: %s", device.Conn.RemoteAddr().(*net.TCPAddr).IP, imei, clearErr.Error())
		// }

		// var listErr error
		// commands, listErr = h.Store.DeviceCommands.Get_DeviceCommandsByImeiAndSendMode(context.Background(), imei, constants.CMD_SEND_MODE_ON_CONNECT)
		// if listErr != nil {
		// 	err = listErr
		// }
	})

	if err != nil {
		logger.Error("SyncPackEventHandler ip={%s} imei={%s}: get commands error: %s", ip, imei, err.Error())
		return
	}

	/* Telemetry commands — low priority, does not block user commands */
	h.SendTelemetryCommands(device)

	/* synchronization */
	/* --------------- */
	func() {
		h.limiterDatabase(func() {
			if device.Conn == nil {
				return
			}

			bulkHandled, bulkErr := uc.CompanyTasks.TryApplyPendingTask(context.Background(), h.Session, d)
			if bulkErr != nil {
				logger.Error("SyncPackEventHandler ip={%s} imei={%s}: Failed to process bulk configuration: %s", ip, imei, bulkErr.Error())
				return
			}

			if bulkHandled || bulkObserved {
				return
			}

			if device.CfgHash == 0 {
				return
			}

			if usecase.IsConfigurationMissing(d) {
				h.Session.TransitData(imei, constants.ADM_RC_TYPE_GET_CFG, nil)
				return
			}

			if device.CfgHash == usecase.EffectiveConfigurationHash(d) {
				return
			}

			if syncErr := uc.Syncs.Synchronization(h.Session, d); syncErr != nil {
				logger.Error("SyncPackEventHandler ip={%s} imei={%s}: Failed to pull/send configuration: %s", ip, imei, syncErr.Error())
			}
		})
	}()
	/* --------------- */
	/* synchronization */
}

func (h *BasePackEventHandler) SendTelemetryCommands(device *adm.ADMDevice) {
	if device == nil || device.Imei == nil {
		return
	}

	h.Session.ScheduleTelemetryCommands(*device.Imei, constants.TelemetryCommands)
}
