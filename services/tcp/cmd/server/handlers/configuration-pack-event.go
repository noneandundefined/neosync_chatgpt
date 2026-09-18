package handlers

import (
	"context"
	"neomatica/neosync-tcp/infra/constants"
	"neomatica/neosync-tcp/infra/logger"
	"neomatica/neosync-tcp/infra/store/postgres/models"
	"neomatica/neosync-tcp/internal/usecase"
	"neomatica/neosync-tcp/infra/store/redis"
	"neomatica/neosync-tcp/pkg/protocol/adm"
	"neomatica/neosync-tcp/pkg/protocol/common"
	"neomatica/neosync-tcp/types"
	"neomatica/neosync-tcp/util"
	"net"
	"net/http"
)

func (h *BasePackEventHandler) ConfigurationPackEventHandler(device *adm.ADMDevice, uc *usecase.UseCase, received any) {
	if device == nil || device.Imei == nil {
		logger.Warning("Device unavailable: cannot fetch events (device not connected)")
		return
	}

	cfgBytes, ok := received.([]byte)
	if !ok || cfgBytes == nil {
		logger.Error("ConfigurationPackEventHandler ip={%s} imei={%s}: expected []byte payload, got %T", device.Conn.RemoteAddr().(*net.TCPAddr).IP, *device.Imei, received)
		return
	}

	currentSession, exists := h.Session.GetDeviceSession(*device.Imei)
	if !exists || currentSession.Device != device {
		return
	}
	currentSession.ConfigurationMu.Lock()
	defer currentSession.ConfigurationMu.Unlock()

	cfgHash := common.GetCfgHash(cfgBytes)
	failed := false

	h.limiterDatabase(func() {
		d, err := h.Store.Devices.Get_DeviceFullByImei(context.Background(), util.PrepareImei(*device.Imei))
		if err != nil || d == nil {
			logger.Error("ConfigurationPackEventHandler imei={%s}: cannot load configuration priority: %v", *device.Imei, err)
			failed = true
			return
		}

		task, err := h.Store.Companies.Get_ActiveCompanyTaskByDeviceId(context.Background(), device.DeviceId)
		if err != nil {
			failed = true
			return
		}

		if _, err := uc.CompanyTasks.ObserveDeviceConfiguration(context.Background(), device.DeviceId, cfgHash); err != nil {
			logger.Error("ConfigurationPackEventHandler imei={%s}: Failed to observe bulk configuration: %s", *device.Imei, err.Error())
			failed = true
			return
		}

		preserveTaskConfiguration := task != nil && common.GetCfgHash(task.CfgData) != cfgHash
		if !preserveTaskConfiguration && !usecase.ShouldPreserveNeosyncConfiguration(d, cfgHash) {
			if err := h.Store.Configurations.Update_ConfigurationByDeviceId(context.Background(), &models.Configuration{
				DeviceID: device.DeviceId,
				CfgHash:  cfgHash,
				CfgData:  cfgBytes,
			}); err != nil {
				logger.Error("ConfigurationPackEventHandler ip={%s} imei={%s}: %s", device.Conn.RemoteAddr().(*net.TCPAddr).IP, *device.Imei, err.Error())
				device.Emitter.Emit("error", err)

				failed = true

				return
			}

			/* Analytics */
			userUuid := "system"
			if device.UserUuid != nil {
				userUuid = *device.UserUuid
			}

			if err := h.Store.Analytics.Mark_AnalyticsConfigurationAppliedByDeviceHash(context.Background(), userUuid, device.DeviceId, cfgHash); err != nil {
				logger.Error("ConfigurationPackEventHandler ip={%s} imei={%s}: Failed to mark analytics apply: %s", device.Conn.RemoteAddr().(*net.TCPAddr).IP, *device.Imei, err.Error())
			}

			if err := h.Store.Configurations.Update_ConfigurationSyncStatusByDeviceId(context.Background(), device.DeviceId, constants.CFG_SYNC_STATUS_CONFIRMED, nil); err != nil {
				logger.Error("ConfigurationPackEventHandler ip={%s} imei={%s}: Failed to update configuration sync status: %s", device.Conn.RemoteAddr().(*net.TCPAddr).IP, *device.Imei, err.Error())
			}
		} else {
			logger.Info("ConfigurationPackEventHandler imei={%s}: keep NeoSync configuration, received hash={%d}, expected={%d}", *device.Imei, cfgHash, usecase.EffectiveConfigurationHash(d))
		}

		cfgLastModTime := common.GetCfgLastModTime(cfgBytes)
		if cfgLastModTime == 0 {
			cfgLastModTime = device.LastModTime
		}

		if err := h.Store.Syncs.Update_SyncByDeviceId(context.Background(), &models.Sync{
			DeviceID:        device.DeviceId,
			FirmwareVersion: device.FirmwareVersion,
			CfgVersion:      device.CfgVersion,
			LastModTime:     cfgLastModTime,
			CfgHash:         cfgHash,
		}); err != nil {
			logger.Error("ConfigurationPackEventHandler ip={%s} imei={%s}: Failed to update sync cfg_hash: %s", device.Conn.RemoteAddr().(*net.TCPAddr).IP, *device.Imei, err.Error())
		}

		device.CfgHash = cfgHash
		device.LastModTime = cfgLastModTime
	})

	if failed {
		return
	}

	/* ADM write log */
	_ = redis.WriteAdmLog(constants.EVENT_DEVICE_CFGRECEV, *device.Imei, map[string]any{
		"device_id": device.DeviceId,
		"cfg_hash":  common.GetCfgHash(cfgBytes),
	})

	if err := redis.ClearCacheConfiguration(*device.Imei); err != nil {
		logger.Error("ConfigurationPackEventHandler ip={%s} imei={%s}: Failed to clear cache configuration: %s", device.Conn.RemoteAddr().(*net.TCPAddr).IP, *device.Imei, err.Error())
		device.Emitter.Emit("error", err)
	}

	session, exists := h.Session.GetDeviceSession(*device.Imei)
	if session == nil || !exists {
		return
	}

	requestID := session.RequestId
	if requestID == nil {
		return
	}

	/* RabbitMQ */
	rabbitmqTransit := types.RabbitMQ_TransitBinary{
		Type:      0x00,
		RequestID: *requestID,
		Imei:      *device.Imei,
		StatCode:  http.StatusOK,
		Data:      cfgBytes,
	}

	if err := h.RMQ.SendToRabbitAsync(rabbitmqTransit); err != nil {
		logger.Error("ConfigurationPackEventHandler ip={%s}: %s", device.Conn.RemoteAddr().(*net.TCPAddr).IP, err.Error())
		device.Emitter.Emit("error", err)
	}
	h.Session.RemoveRequestIdSession(*requestID)
}
