package main

import (
	"encoding/binary"
	"fmt"

	"neomatica/neosync-tcp/infra/constants"
	"neomatica/neosync-tcp/pkg/protocol/common"
)

func imeiToBCD(imei string) ([]byte, error) {
	if len(imei) > constants.IMEI_MAX_SIZE*2 {
		return nil, fmt.Errorf("imei too long")
	}

	padded := imei
	if len(padded)%2 != 0 {
		padded += "F"
	}

	out := make([]byte, len(padded)/2)
	for i := 0; i < len(out); i++ {
		hi := padded[i*2]
		lo := padded[i*2+1]

		if hi < '0' || hi > '9' || (lo != 'F' && lo != 'f' && (lo < '0' || lo > '9')) {
			return nil, fmt.Errorf("invalid imei digit pair %q", padded[i*2:i*2+2])
		}

		out[i] = (hi-'0')<<4 | imeiLoNibble(lo)
	}

	return out, nil
}

func imeiLoNibble(c byte) byte {
	if c == 'F' || c == 'f' {
		return 0x0F
	}

	return c - '0'
}

func buildHelloPacket(state *DeviceState) ([]byte, error) {
	imeiBCD, err := imeiToBCD(state.IMEI)
	if err != nil {
		return nil, err
	}

	pass := state.Password
	if pass == "" {
		pass = "0"
	}

	packet := make([]byte, constants.ADM_RC_HELLO_PACK_SIZE_V2)
	packet[0] = byte(constants.ADM_RC_HELLO_PACK_SIZE_V2)
	packet[1] = byte(constants.ADM_RC_HELLO_PACK_SIZE_V2 >> 8)
	packet[2] = constants.ADM_RC_PACK_TYPE_HELLO_V2
	packet[3] = 0x02
	packet[4] = byte(len(imeiBCD))
	copy(packet[5:5+len(imeiBCD)], imeiBCD)

	offset := 5 + len(imeiBCD)
	if offset+1+len(pass) > 22 {
		return nil, fmt.Errorf("hello packet layout overflow")
	}

	packet[offset] = byte(len(pass))
	copy(packet[offset+1:offset+1+len(pass)], []byte(pass))

	packet[22] = byte(state.FirmwareVersion)
	packet[23] = state.CfgVersion
	binary.LittleEndian.PutUint32(packet[24:28], state.LastModTime)
	binary.LittleEndian.PutUint32(packet[28:32], state.CfgHash)

	return packet, nil
}

func buildSyncPacket(state *DeviceState) []byte {
	packet := make([]byte, constants.ADM_RC_SYNC_PACK_SIZE_V2)
	packet[0] = constants.ADM_RC_SYNC_PACK_SIZE_V2
	packet[1] = 0x00
	packet[2] = constants.ADM_RC_PACK_TYPE_SYNC_V2
	packet[3] = byte(state.FirmwareVersion)
	packet[4] = state.CfgVersion

	// big-endian, как Parse_SyncPacket на сервере
	putUint32BE(packet[5:9], state.LastModTime)
	putUint32BE(packet[9:13], state.CfgHash)

	return packet
}

func putUint32BE(dst []byte, v uint32) {
	dst[0] = byte(v >> 24)
	dst[1] = byte(v >> 16)
	dst[2] = byte(v >> 8)
	dst[3] = byte(v)
}

func ensureCfgPacket(raw []byte) ([]byte, error) {
	if len(raw) < 3 {
		return nil, fmt.Errorf("config too short")
	}

	if raw[2] == constants.ADM_RC_PACK_TYPE_CFG_V2 {
		size := int(binary.LittleEndian.Uint16(raw[:2]))
		if size == len(raw) {
			return raw, nil
		}
	}

	// uid/value block без ADM-заголовка
	body := raw
	packet := make([]byte, len(body)+3)
	size := len(packet)
	binary.LittleEndian.PutUint16(packet[0:2], uint16(size))
	packet[2] = constants.ADM_RC_PACK_TYPE_CFG_V2
	copy(packet[3:], body)

	return packet, nil
}

func applyConfigPacket(state *DeviceState, packet []byte) error {
	if len(packet) < 4 || packet[2] != constants.ADM_RC_PACK_TYPE_CFG_V2 {
		return fmt.Errorf("not a cfg packet")
	}

	if hash := common.GetCfgHash(packet); hash != 0 {
		state.CfgHash = hash
	}
	if ts := common.GetCfgLastModTime(packet); ts != 0 {
		state.LastModTime = ts
	}

	return nil
}

func buildStringAnswer(command string) []byte {
	body := append([]byte(command), 0x00)
	size := len(body) + constants.NRC_V2_PACKET_HEADER_SIZE + 1
	packet := make([]byte, size)
	common.AdmRcCreateHeaderForPacket(packet, uint(size), constants.ADM_RC_PACK_TYPE_STRING_COMMAND_ANSWER_V2)
	copy(packet[constants.NRC_V2_PACKET_HEADER_SIZE:], body)
	return packet
}
