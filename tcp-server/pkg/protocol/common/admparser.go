package common

import "encoding/binary"

func GetCfgHash(cfg []byte) uint32 {
	if len(cfg) < 7 {
		return 0
	}

	body := cfg[3:]
	return binary.LittleEndian.Uint32(body[len(body)-4:])
}

func GetCfgLastModTime(cfg []byte) uint32 {
	if len(cfg) < 9 {
		return 0
	}

	return binary.LittleEndian.Uint32(cfg[5:9])
}
