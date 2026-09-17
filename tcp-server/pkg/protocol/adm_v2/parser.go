package adm_v2

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"neomatica/neosync-tcp/infra/constants"
	"neomatica/neosync-tcp/infra/tcperrors"
	"neomatica/neosync-tcp/pkg/protocol"
	"neomatica/neosync-tcp/pkg/protocol/common"
	"neomatica/neosync-tcp/util"
)

func (adm *ADM_V2) Parse_WelcomePacket(buffer []byte) (*protocol.WelcomePacket, error) {
	if len(buffer) < constants.ADM_RC_HELLO_PACK_SIZE_V2 {
		return nil, tcperrors.Err_FirmwareOutdatedNotFound
	}

	imeiSize := int(buffer[4])
	if imeiSize > constants.IMEI_MAX_SIZE {
		return nil, fmt.Errorf("imei size is too large")
	}

	if len(buffer) < 5+imeiSize {
		return nil, fmt.Errorf("packet truncated (imei)")
	}

	parsedIMEI := common.DecodeIMEI(buffer[5 : 5+imeiSize])
	if parsedIMEI == "" {
		return nil, fmt.Errorf("response from tracker not recognized")
	}

	if len(buffer) < 14 {
		return nil, fmt.Errorf("packet truncated (password offset)")
	}

	passwordSize := int(buffer[13])
	if passwordSize > constants.PASSWORD_MAX_SIZE {
		return nil, fmt.Errorf("password is too large")
	}

	var password string
	if passwordSize > 0 {
		if len(buffer) < 14+passwordSize {
			return nil, fmt.Errorf("packet truncated (password)")
		}

		password = fmt.Sprintf("%x", buffer[14:14+passwordSize])
	}

	if len(buffer) < constants.ADM_RC_HELLO_PACK_SIZE_V2 {
		return nil, fmt.Errorf("packet truncated (fixed fields)")
	}

	welcomePacket := &protocol.WelcomePacket{
		Imei:            &parsedIMEI,
		Pass:            &password,
		Version:         buffer[3],
		FirmwareVersion: uint16(buffer[22]),
		CfgVersion:      buffer[23],
		LastModTime:     binary.LittleEndian.Uint32(buffer[24:28]),
		CfgHash:         binary.LittleEndian.Uint32(buffer[28:32]),
	}

	adm.Imei = util.PrepareImei(parsedIMEI)

	return welcomePacket, nil
}

func (adm *ADM_V2) Parse_CommandPacket(buffer []byte, receiveBufDataSize uint) string {
	ctx := context.Background()

	raw := buffer[constants.NRC_V2_PACKET_HEADER_SIZE:receiveBufDataSize]

	end := bytes.IndexByte(raw, 0x00)
	if end == -1 {
		end = len(raw)
	}
	raw = raw[:end]

	result := string(raw[:end])

	if adm.Imei != "" && !adm.deviceModelLookupDone.Load() {
		device, err := adm.Store.Devices.Get_DeviceByImei(ctx, adm.Imei)

		if err == nil && device != nil {
			hasModel := device.DeviceModel != nil && *device.DeviceModel != ""
			hasExtendedModel := device.DeviceExtendedModel != nil && *device.DeviceExtendedModel != ""

			if hasModel && hasExtendedModel {
				adm.deviceModelLookupDone.Store(true)
			} else {
				model, extendedModel := common.GetModelsByWho(raw)

				if (model != nil || extendedModel != nil) &&
					adm.Store.Devices.Update_DeviceModelsByImei(ctx, model, extendedModel, adm.Imei) == nil {
					adm.deviceModelLookupDone.Store(true)
				}
			}
		}
	}

	return result
}

func (adm *ADM_V2) Parse_SyncPacket(buffer []byte) (*protocol.SyncPacket, error) {
	lastModTime := uint32(buffer[8])<<24 | uint32(buffer[7])<<16 | uint32(buffer[6])<<8 | uint32(buffer[5])
	firmwareVersion := uint16(buffer[3])
	cfgVersion := buffer[4]

	if firmwareVersion < 0x0001 {
		return nil, fmt.Errorf("invalid firmware version in sync packet: %d", firmwareVersion)
	}

	return &protocol.SyncPacket{
		FirmwareVersion: firmwareVersion,
		CfgVersion:      cfgVersion,
		LastModTime:     lastModTime,
		CfgHash:         uint32(buffer[12])<<24 | uint32(buffer[11])<<16 | uint32(buffer[10])<<8 | uint32(buffer[9]),
	}, nil
}

func (adm *ADM_V2) Parse_ConfigurationPacket(buffer []byte) (*protocol.ConfigurationPacket, error) {
	return &protocol.ConfigurationPacket{
		CfgHash: binary.LittleEndian.Uint32(buffer[len(buffer)-4:]),
	}, nil
}

func (adm *ADM_V2) Parse_KeepAlivePacket(buffer []byte, imei *string) {
	/* Not parsed */
}
