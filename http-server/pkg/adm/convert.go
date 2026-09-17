package adm

import (
	"bytes"
	"encoding/binary"
	"math"
	"neomatica/neosync/infra/constants"
	"neomatica/neosync/pkg/adm/help"
)

func (cfgItem FieldConfiguration) GetValue() any {
	fieldType, err := constants.GetCfgType(cfgItem.UID)
	if err != nil {
		return nil
	}

	switch fieldType {
	case constants.UInt8:
		if len(cfgItem.Value) < 1 {
			return 0
		}
		return cfgItem.Value[0]

	case constants.UInt16:
		if len(cfgItem.Value) < 2 {
			return uint16(0)
		}
		return binary.LittleEndian.Uint16(cfgItem.Value)

	case constants.UInt32:
		if len(cfgItem.Value) < 4 {
			return uint32(0)
		}
		return binary.LittleEndian.Uint32(cfgItem.Value)

	case constants.TwoUint8:
		arr := make([]int, len(cfgItem.Value))
		for i, b := range cfgItem.Value {
			arr[i] = int(b)
		}
		return arr

	case constants.TwoUint16:
		arr := make([]uint16, 0, 2)
		if len(cfgItem.Value) >= 2 {
			arr = append(arr, binary.LittleEndian.Uint16(cfgItem.Value[:2]))
		}
		if len(cfgItem.Value) >= 4 {
			arr = append(arr, binary.LittleEndian.Uint16(cfgItem.Value[2:4]))
		}
		return arr

	case constants.Uint32Array:
		count := len(cfgItem.Value) / 4
		arr := make([]uint32, count)
		for i := 0; i < count; i++ {
			arr[i] = binary.LittleEndian.Uint32(cfgItem.Value[i*4:])
		}
		return arr

	case constants.Uint64Array:
		count := len(cfgItem.Value) / 8
		arr := make([]uint64, count)
		for i := 0; i < count; i++ {
			arr[i] = binary.LittleEndian.Uint64(cfgItem.Value[i*8:])
		}
		return arr

	case constants.Phones:
		return help.HelpConvertToPhones(cfgItem.Value, 4)

	case constants.MacAddressArray:
		return help.HelpConvertToMacArray(cfgItem.Value)

	case constants.SNSArray:
		return help.HelpConvertToSNS(cfgItem.Value)

	case constants.Hex:
		return help.HelpConvertToHex(cfgItem.Value)

	case constants.Ascii:
		return sanitizePostgresText(string(bytes.Trim(cfgItem.Value, "\x00")))

	case constants.ByteArray, constants.Raw:
		return cfgItem.Value

	default:
		return cfgItem.Value
	}
}

func (cfgItem FieldConfiguration) ToFloat32() float32 {
	if len(cfgItem.Value) < 4 {
		return 0
	}

	bits := binary.LittleEndian.Uint32(cfgItem.Value)
	return math.Float32frombits(bits)
}
