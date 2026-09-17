package common

import (
	"bytes"
	"errors"
	"neomatica/neosync-tcp/util"
	"strings"
	"sync"
)

func GetVersion_ByBinPacket(packet []byte) (byte, error) {
	if len(packet) < 4 {
		return 0, errors.New("packet too short to contain 'version'")
	}

	return packet[3], nil
}

func GetCleanBufFromPool(pool *sync.Pool) []byte {
	return pool.Get().([]byte)[:0]
}

func AdmRcCreateHeaderForPacket(packetBuf []byte, sizeOfPacket uint, type_x uint8) {
	packetBuf[0] = byte(sizeOfPacket)
	packetBuf[1] = byte(sizeOfPacket >> 8)
	packetBuf[2] = type_x
}

func DecodeBCD(bcd []byte) string {
	return DecodeIMEI(bcd)
}

func DecodeIMEI(bcd []byte) string {
	out := make([]byte, 0, len(bcd)*2)

	for _, b := range bcd {
		hi := b >> 4
		lo := b & 0x0F

		if hi <= 9 {
			out = append(out, '0'+hi)
		}

		if lo <= 9 {
			out = append(out, '0'+lo)
		}
	}

	imei := string(out)
	if len(imei) == 16 && imei[0] == '0' {
		imei = imei[1:]
	}

	return imei
}

func GetModelByWho(raw []byte) *string {
	end := bytes.IndexByte(raw, 0x00)
	if end == -1 {
		end = len(raw)
	}

	raw = raw[:end]

	first := bytes.IndexByte(raw, ',')
	if first == -1 {
		return nil
	}

	second := bytes.IndexByte(raw[first+1:], ',')
	if second == -1 {
		return nil
	}

	modelRaw := string(raw[first+1 : first+1+second])

	if len(modelRaw) >= 3 && strings.EqualFold(modelRaw[:3], "ADM") {
		return util.ValidDeviceModel(NormalizeModel(modelRaw))
	}

	return nil
}

func GetModelsByWho(raw []byte) (model *string, extendedModel *string) {
	end := bytes.IndexByte(raw, 0x00)
	if end == -1 {
		end = len(raw)
	}

	raw = raw[:end]
	parts := bytes.Split(raw, []byte{','})
	if len(parts) < 4 {
		return nil, nil
	}

	extendedRaw := strings.TrimSpace(string(parts[1]))
	if extendedRaw != "" {
		extendedModel = &extendedRaw
	}

	modelRaw := strings.TrimSpace(string(parts[3]))
	if len(modelRaw) >= 3 && strings.EqualFold(modelRaw[:3], "ADM") {
		model = util.ValidDeviceModel(NormalizeModel(modelRaw))
	}

	return model, extendedModel
}

func NormalizeModel(raw string) string {
	r := strings.ToUpper(raw)

	switch {
	case strings.HasPrefix(r, "ADM007BLE"):
		return "ADM007BLE"

	case strings.HasPrefix(r, "ADM007STD"):
		return "ADM007"

	case strings.HasPrefix(r, "ADM333STD"):
		return "ADM333"

	case strings.HasPrefix(r, "ADM333BLE"):
		return "ADM333BLE"

	case strings.HasPrefix(r, "ADM333V2"):
		return "ADM333V2"

	case strings.HasPrefix(r, "ADM333V3"):
		return "ADM333V2"

	case strings.HasPrefix(r, "ADM500"):
		return "ADM500"

	case strings.HasPrefix(r, "ADMP50 LTE"):
		return "ADMP50LTE"

	case strings.HasPrefix(r, "ADMP50"):
		return "ADMP50"

	default:
		return raw
	}
}
