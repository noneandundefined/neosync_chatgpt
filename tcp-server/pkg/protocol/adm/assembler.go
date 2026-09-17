package adm

import (
	"encoding/binary"
	"encoding/hex"
	"neomatica/neosync-tcp/infra/constants"
	"neomatica/neosync-tcp/infra/logger"
)

const assemblerMaxBufferSize = constants.ADM_RC_MAX_PACKET_SIZE * 2

type Assembler struct {
	buffer []byte
}

func isValidIncomingHeader(size int, packetType byte) bool {
	if size < constants.NRC_V2_PACKET_HEADER_SIZE || size > constants.ADM_RC_MAX_PACKET_SIZE {
		return false
	}

	switch packetType {
	case constants.ADM_RC_PACK_TYPE_HELLO_V2:
		// Минимум текущего v2; новая прошивка может увеличить hello — size в header задаёт фактическую длину.
		return size >= constants.ADM_RC_HELLO_PACK_SIZE_V2
	case constants.ADM_RC_PACK_TYPE_SYNC_V2:
		return size == constants.ADM_RC_SYNC_PACK_SIZE_V2
	case constants.ADM_RC_PACK_TYPE_STRING_COMMAND_ANSWER_V2:
		return size > constants.NRC_V2_PACKET_HEADER_SIZE
	case constants.ADM_RC_PACK_TYPE_CFG_V2:
		return size >= constants.NRC_V2_PACKET_HEADER_SIZE+4
	default:
		return false
	}
}

func headerPreview(buf []byte) string {
	n := len(buf)
	if n > 16 {
		n = 16
	}

	return hex.EncodeToString(buf[:n])
}

/* Feed дожидается всего пакета по [size] */
func (a *Assembler) Feed(data []byte) [][]byte {
	a.buffer = append(a.buffer, data...)
	var packets [][]byte

	for {
		if len(a.buffer) > assemblerMaxBufferSize {
			drop := len(a.buffer) - assemblerMaxBufferSize
			logger.Warning("Feed: assembler buffer overflow, dropping %d leading bytes", drop)
			a.buffer = a.buffer[drop:]
		}

		if len(a.buffer) < constants.NRC_V2_PACKET_HEADER_SIZE {
			break
		}

		size := int(binary.LittleEndian.Uint16(a.buffer[:2]))
		packetType := a.buffer[2]

		if !isValidIncomingHeader(size, packetType) {
			skipped := 1
			a.buffer = a.buffer[1:]

			for len(a.buffer) >= constants.NRC_V2_PACKET_HEADER_SIZE {
				size = int(binary.LittleEndian.Uint16(a.buffer[:2]))
				packetType = a.buffer[2]

				if isValidIncomingHeader(size, packetType) {
					logger.Warning(
						"Feed: resync skipped %d bytes, next size=%d type=0x%02X prefix=%s",
						skipped,
						size,
						packetType,
						headerPreview(a.buffer),
					)
					break
				}

				a.buffer = a.buffer[1:]
				skipped++
			}

			if len(a.buffer) < constants.NRC_V2_PACKET_HEADER_SIZE {
				break
			}

			continue
		}

		if len(a.buffer) < size {
			break
		}

		pkt := make([]byte, size)
		copy(pkt, a.buffer[:size])

		packets = append(packets, pkt)

		a.buffer = a.buffer[size:]
	}

	return packets
}
