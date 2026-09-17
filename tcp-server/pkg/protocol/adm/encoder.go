package adm

import (
	"io"
	"neomatica/neosync-tcp/config"
	"neomatica/neosync-tcp/infra/constants"
	"neomatica/neosync-tcp/infra/logger"
	"neomatica/neosync-tcp/infra/store/redis"
	"neomatica/neosync-tcp/infra/tcperrors"
	"neomatica/neosync-tcp/pkg/protocol/common"
	"neomatica/neosync-tcp/util"
	"time"
)

func (adm *ADMDevice) TransitData(type_p uint8, buffer []byte) error {
	ip := util.RemoteIP(adm.Conn)
	imei := adm.imeiForLog()

	if adm.Conn == nil {
		logger.Error("Encoder imei={%s}: Socket is NULL", imei)
		adm.Emitter.Emit("error", tcperrors.Err_DeviceNotConnected)
		return tcperrors.Err_DeviceNotConnected
	}

	if type_p < constants.ADM_RC_TYPE_STRING || type_p > constants.ADM_RC_TYPE_SET_CFG || (type_p == constants.ADM_RC_TYPE_STRING && buffer == nil) {
		logger.Error("Encoder ip={%s} imei={%s}: Wrong command type", ip, imei)
		adm.Emitter.Emit("error", tcperrors.Err_UnknownTypeOfCommand)
		return tcperrors.Err_UnknownTypeOfCommand
	}

	var sendBuffer []byte
	switch type_p {
	case constants.ADM_RC_TYPE_STRING:
		sendBuffer = adm.protocol.Encode_CommandPacket(buffer)

	case constants.ADM_RC_TYPE_GET_SYNC:
		sendBuffer = adm.protocol.Encode_GetSyncPacket()

	case constants.ADM_RC_TYPE_GET_CFG:
		sendBuffer = adm.protocol.Encode_GetCfgPacket()

		go func(imeiValue string, deviceID uint64) {
			_ = redis.WriteAdmLog(constants.EVENT_DEVICE_CFGREQUEST, imeiValue, map[string]any{
				"device_id": deviceID,
			})
		}(imei, adm.DeviceId)

	case constants.ADM_RC_TYPE_SET_CFG:
		sendBuffer = adm.protocol.Encode_SetCfgPacket(buffer, adm.DeviceId, adm.CfgVersion, adm.FirmwareVersion)

		cfgHash := common.GetCfgHash(buffer)
		if cfgHash == 0 {
			cfgHash = adm.CfgHash
		}

		// adm write log
		go func(imeiValue string, hash uint32) {
			_ = redis.WriteAdmLog(constants.EVENT_DEVICE_CFGSEND, imeiValue, map[string]any{
				"device_id": adm.DeviceId,
				"cfg_hash":  hash,
			})
		}(imei, cfgHash)

	default:
		logger.Warning("Encoder ip={%s} imei={%s}: Firmware is outdated or not found on the server: %d", ip, imei, type_p)
		adm.Emitter.Emit("error", tcperrors.Err_FirmwareOutdatedNotFound)
		return tcperrors.Err_FirmwareOutdatedNotFound
	}

	if sendBuffer == nil {
		adm.Emitter.Emit("error", tcperrors.Err_SendPacketToDevice)
		return tcperrors.Err_SendPacketToDevice
	}

	holdsTerminal := type_p == constants.ADM_RC_TYPE_STRING
	if holdsTerminal {
		adm.AcquireTerminal()
	}

	if err := adm.Conn.SetWriteDeadline(time.Now().Add(config.WriteDeadline)); err != nil {
		logger.Error("Encoder ip={%s} imei={%s}: Failed to set time_out (write): %s", ip, imei, err.Error())

		adm.Emitter.Emit("error", tcperrors.Err_SendPacketToDevice)

		if holdsTerminal {
			adm.ReleaseTerminal()
		}

		return err
	}

	n, err := adm.Conn.Write(sendBuffer)
	if err == nil && n != len(sendBuffer) {
		err = io.ErrShortWrite
	}
	if err != nil {
		logger.Error("Encoder ip={%s} imei={%s}: Failed sending command to device: %s", ip, imei, err.Error())

		adm.Emitter.Emit("error", tcperrors.Err_SendPacketToDevice)

		if holdsTerminal {
			adm.ReleaseTerminal()
		}

		adm.OnDisconnect()

		return err
	}

	logOutgoingTransit(ip, imei, type_p, buffer, n)
	return nil
}
