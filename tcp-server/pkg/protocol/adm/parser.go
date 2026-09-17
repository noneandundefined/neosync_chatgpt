package adm

import (
	"neomatica/neosync-tcp/infra/constants"
	"neomatica/neosync-tcp/infra/logger"
	"net"
)

func (adm *ADMDevice) ParsePacket(buffer []byte) {
	/* Min\Max packet size */
	packSize := len(buffer)
	if packSize < constants.NRC_V2_PACKET_HEADER_SIZE || packSize > constants.ADM_RC_MAX_PACKET_SIZE {
		logger.Error("ParsePacket ip={%s}: Incorrect package size: %d", adm.Conn.RemoteAddr().(*net.TCPAddr).IP, packSize)
		return
	}

	/* Validate packet type */
	admType := buffer[2]
	if admType < constants.ADM_RC_PACK_TYPE_HELLO_V2 || admType > constants.ADM_RC_PACK_TYPE_CFG_V2 {
		logger.Error("ParsePacket ip={%s}: Unknown packet type: 0x%02X", adm.Conn.RemoteAddr().(*net.TCPAddr).IP, admType)
		return
	}

	switch admType {
	case constants.ADM_RC_PACK_TYPE_HELLO_V2:
		adm.parseHelloPacket(buffer)

	case constants.ADM_RC_PACK_TYPE_SYNC_V2:
		adm.parseSyncPacket(buffer)

	case constants.ADM_RC_PACK_TYPE_STRING_COMMAND_ANSWER_V2:
		adm.parseCommandPacket(buffer)

	case constants.ADM_RC_PACK_TYPE_CFG_V2:
		adm.Emitter.Emit("configuration-pack-received", buffer)

	default:
		logger.Error("ParsePacket ip={%s}: Unknown packet type: 0x%02X", adm.Conn.RemoteAddr().(*net.TCPAddr).IP, admType)
	}
}
