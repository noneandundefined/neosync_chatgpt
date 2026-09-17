package adm_v2

import (
	"neomatica/neosync-tcp/infra/constants"
	"neomatica/neosync-tcp/pkg/protocol/common"
)

func (adm *ADM_V2) Encode_CommandPacket(buffer []byte) []byte {
	if buffer == nil {
		return nil
	}

	bufferBytes := append(buffer, 0x00)

	sizeOfPacket := len(bufferBytes) + constants.NRC_V2_PACKET_HEADER_SIZE + 1
	sendBuffer := make([]byte, sizeOfPacket)

	common.AdmRcCreateHeaderForPacket(sendBuffer, uint(sizeOfPacket), constants.ADM_RC_PACK_TYPE_COMMAND_V2)
	sendBuffer[constants.ADM_RC_TYPE_OFFSET] = constants.ADM_RC_TYPE_STRING

	copy(sendBuffer[constants.ADM_RC_TYPE_OFFSET+1:], bufferBytes)

	return sendBuffer
}

func (adm *ADM_V2) Encode_GetSyncPacket() []byte {
	sendBuffer := make([]byte, constants.ADM_RC_COMMAND_GET_SYNC_PACK_SIZE_V2)

	common.AdmRcCreateHeaderForPacket(sendBuffer, constants.ADM_RC_COMMAND_GET_SYNC_PACK_SIZE_V2, constants.ADM_RC_PACK_TYPE_COMMAND_V2)
	sendBuffer[constants.ADM_RC_TYPE_OFFSET] = constants.ADM_RC_TYPE_GET_SYNC

	return sendBuffer
}

func (adm *ADM_V2) Encode_GetCfgPacket() []byte {
	sendBuffer := make([]byte, constants.ADM_RC_COMMAND_GET_CFG_PACK_SIZE_V2)

	common.AdmRcCreateHeaderForPacket(sendBuffer, constants.ADM_RC_COMMAND_GET_CFG_PACK_SIZE_V2, constants.ADM_RC_PACK_TYPE_COMMAND_V2)
	sendBuffer[constants.ADM_RC_TYPE_OFFSET] = constants.ADM_RC_TYPE_GET_CFG

	return sendBuffer
}

func (adm *ADM_V2) Encode_SetCfgPacket(buffer []byte, id uint64, cfgVersion uint8, firmwareVersion uint16) []byte {
	if buffer == nil || id == 0 {
		return nil
	}

	return buffer
}
