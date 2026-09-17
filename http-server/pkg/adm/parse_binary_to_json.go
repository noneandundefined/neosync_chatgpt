package adm

import (
	"encoding/binary"
	"errors"
	"fmt"
	"neomatica/neosync/infra/constants"
	"neomatica/neosync/infra/locale"
	"neomatica/neosync/infra/logger"
	"neomatica/neosync/pkg/adm/help"
	"sort"
)

func (adm *Adm) ParseBinaryToJson(tr locale.Translator, data []byte) ([]FieldConfigurationParsed, error) {
	fieldConfiguration := make(map[uint16]FieldConfiguration)

	if len(data) < 4 {
		return nil, errors.New(tr.TErr("cfg-parse-error-device"))
	}

	size := int(binary.LittleEndian.Uint16(data[0:2]))
	if size > len(data) {
		return nil, errors.New(tr.TErr("cfg-parse-error-device"))
	}

	offset := 2 // SIZE
	if len(data) < offset+1 {
		return nil, errors.New(tr.TErr("cfg-parse-error-device"))
	}
	offset++ // TYPE

	offset += 8 // MAGIC + LAST_MODIFICATION_TIME + VERSION

	cfgSize := int(binary.LittleEndian.Uint16(data[offset : offset+2]))
	offset += 2

	/* cfg_size is the inner packet: magic + last_mod + version + this field + body + hash.
	   Fields start here; hash is the last 4 bytes of that inner packet. */
	innerStart := 3
	end := innerStart + cfgSize - 4
	if end > len(data) {
		end = len(data)
	}
	if end < offset {
		end = offset
	}

	for offset+3 <= end {
		uid := binary.LittleEndian.Uint16(data[offset : offset+2])
		length := data[offset+2]
		offset += 3

		if offset+int(length) > end {
			break
		}

		value := data[offset : offset+int(length)]
		offset += int(length)

		/* Padding / truncated tail: uid 0 or size 0 is not a real field.
		   Oversized cfg_size turns leftover bytes into fake UIDs
		   (e.g. ADM P50 has no Modbus, but 4c 00 00 becomes uid=76 size=0). */
		if uid == 0 || length == 0 {
			if uid != 0 {
				logger.Warning("ParseBinaryToJson: skip padding field uid=%d size=%d", uid, length)
			}
			continue
		}

		fieldConfiguration[uid] = FieldConfiguration{
			UID:   uid,
			Size:  length,
			Value: value,
		}
	}

	/* Build json view */
	uids := make([]uint16, 0, len(fieldConfiguration))
	for uid := range fieldConfiguration {
		uids = append(uids, uid)
	}

	sort.Slice(uids, func(i, j int) bool {
		return uids[i] < uids[j]
	})

	var configurationParsed []FieldConfigurationParsed
	for _, uid := range uids {
		cfgItem := fieldConfiguration[uid]

		value := adm.getValueByUid(cfgItem.UID, cfgItem)

		uidName, err := constants.GetCfgName(cfgItem.UID)
		if err != nil {
			/* Unknown uid in conf. */
			configurationParsed = append(configurationParsed, FieldConfigurationParsed{
				UID:      fmt.Sprintf("unknown_%d", cfgItem.UID),
				RawUID:   cfgItem.UID,
				Size:     cfgItem.Size,
				RawValue: cfgItem.Value,
				Unknown:  true,
			})
			continue
		}

		// fmt.Printf("%s - %d - %d: %s\n", uidName, cfgItem.UID, cfgItem.Size, hex.EncodeToString(cfgItem.Value))
		// fmt.Printf("%d: %v\n", cfgItem.UID, value)

		configurationParsed = append(configurationParsed, FieldConfigurationParsed{
			UID:    uidName,
			RawUID: cfgItem.UID,
			Size:   cfgItem.Size,
			Value:  value,
		})
	}

	return configurationParsed, nil
}

func (adm *Adm) getValueByUid(uid uint16, cfgItem FieldConfiguration) any {
	raw := cfgItem.Value

	fieldType, err := constants.GetCfgType(uid)
	if err != nil {
		return nil
	}

	switch uid {
	case 154, 155, 156, 157, 4: // ascii array
		return help.HelpConvertToFixedOrdered(raw, 2)

	case 5: // two uint16
		return help.HelpConvertTwoArrayUint16(raw)

	case 7, 124: // two uint8
		return help.HelpConvertTwoArrayUint8(raw)

	case 98: // beacon time sec
		bytes := cfgItem.Value
		if len(bytes)%4 != 0 {
			return nil
		}

		vals := make([]uint32, len(bytes)/4)
		for i := 0; i < len(bytes); i += 4 {
			vals[i/4] = binary.LittleEndian.Uint32(bytes[i : i+4])
		}

		return vals
	default:
		switch fieldType {
		case constants.ByteArray, constants.Raw:
			bytes := cfgItem.Value
			arr := make([]int, len(bytes))
			for i, b := range bytes {
				arr[i] = int(b)
			}
			return arr
		default:
			return cfgItem.GetValue()
		}
	}
}
