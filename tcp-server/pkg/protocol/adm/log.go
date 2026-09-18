package adm

import (
	"fmt"
	"neomatica/neosync-tcp/infra/constants"
	"neomatica/neosync-tcp/infra/logger"
)

func logIncomingPacket(ip, imei string, packet []byte) {
	size := len(packet)
	if size == 0 {
		return
	}

	logger.Packet("RX %s ip={%s} imei={%s} bytes={%d}", incomingPacketLabel(packet), ip, imei, size)
}

func logOutgoingTransit(ip, imei string, type_p uint8, payload []byte, sent int) {
	label, details := outgoingTransitLabel(type_p, payload)
	if details != "" {
		logger.Packet("TX %s ip={%s} imei={%s} bytes={%d} %s", label, ip, imei, sent, details)
		return
	}

	logger.Packet("TX %s ip={%s} imei={%s} bytes={%d}", label, ip, imei, sent)
}

func incomingPacketLabel(packet []byte) string {
	if len(packet) < constants.NRC_V2_PACKET_HEADER_SIZE {
		return "packet"
	}

	switch packet[2] {
	case constants.ADM_RC_PACK_TYPE_HELLO_V2:
		return "hello"
	case constants.ADM_RC_PACK_TYPE_SYNC_V2:
		return "sync"
	case constants.ADM_RC_PACK_TYPE_STRING_COMMAND_ANSWER_V2:
		return "command-response"
	case constants.ADM_RC_PACK_TYPE_CFG_V2:
		return "configuration"
	default:
		return fmt.Sprintf("packet-type-0x%02X", packet[2])
	}
}

func outgoingTransitLabel(type_p uint8, payload []byte) (label string, details string) {
	switch type_p {
	case constants.ADM_RC_TYPE_STRING:
		cmd := truncateASCII(payload, 48)
		if cmd == "" {
			return "command", ""
		}

		return "command", fmt.Sprintf("cmd={%s}", cmd)

	case constants.ADM_RC_TYPE_GET_SYNC:
		return "sync-request", ""

	case constants.ADM_RC_TYPE_GET_CFG:
		return "get-configuration-request", ""

	case constants.ADM_RC_TYPE_SET_CFG:
		return "configuration", ""

	default:
		return fmt.Sprintf("packet-type-0x%02X", type_p), ""
	}
}

func truncateASCII(payload []byte, max int) string {
	if len(payload) == 0 {
		return ""
	}

	end := len(payload)
	if end > max {
		end = max
	}

	return string(payload[:end])
}
