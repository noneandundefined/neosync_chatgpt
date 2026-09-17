package admparser

import "encoding/binary"

func GetCfgHash(cfg []byte) uint32 {
	return binary.LittleEndian.Uint32(cfg[3:][len(cfg[3:])-4:])
}

func GetCfgLastModTime(cfg []byte) uint32 {
	return binary.LittleEndian.Uint32(cfg[5:9])
}

func GetCfgPacketSize(cfg []byte) uint16 {
	if len(cfg) < 2 {
		return 0
	}

	return binary.LittleEndian.Uint16(cfg[:2])
}
