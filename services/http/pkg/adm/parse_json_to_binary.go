package adm

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"neomatica/neosync/infra/constants"
	"neomatica/neosync/infra/locale"
	"neomatica/neosync/util"
	"reflect"
	"time"
)

const (
	magic   uint16 = 35814
	version uint16 = 0x02
)

func (adm *Adm) ParseJsonToBinary(tr locale.Translator, cfg []PacketFieldConfiguration) ([]byte, error) {
	cfgFields := new(bytes.Buffer)
	cfgPacket := new(bytes.Buffer)

	write := func(w io.Writer, value interface{}) error {
		return binary.Write(w, binary.LittleEndian, value)
	}

	/* Magic */
	if err := write(cfgPacket, magic); err != nil {
		return nil, errors.New(tr.TErr("cfg-write-magic-error"))
	}

	/* last_mod_time */
	if err := write(cfgPacket, uint32(time.Now().Unix())); err != nil {
		return nil, errors.New(tr.TErr("cfg-write-last-mod-time-error"))
	}

	/* Vsersion */
	if err := write(cfgPacket, version); err != nil {
		return nil, errors.New(tr.TErr("cfg-write-version-error"))
	}

	for _, cfgItem := range cfg {
		var raw []byte
		var errRaw error

		if cfgItem.Size == 0 {
			continue
		}

		fieldType, _ := constants.GetCfgType(cfgItem.UID)
		if (fieldType == constants.ByteArray || fieldType == constants.Raw) && isEmptySlice(cfgItem.Value) {
			continue
		}

		/* Parse unknown raw value */
		if cfgItem.Unknown {
			raw = cfgItem.RawValue
			if len(raw) != int(cfgItem.Size) {
				return nil, errors.New("invalid raw size for unknown uid")
			}
		} else {
			raw, errRaw = adm.getValueByType(tr, cfgItem.UID, cfgItem.Value, cfgItem.Size)
			if errRaw != nil {
				return nil, errRaw
			}
		}

		if err := binary.Write(cfgFields, binary.LittleEndian, cfgItem.UID); err != nil {
			return nil, errors.New(tr.TErr("cfg-uid-to-bytes-error"))
		}

		if err := write(cfgFields, cfgItem.Size); err != nil {
			return nil, errors.New(tr.TErr("cfg-size-to-bytes-error"))
		}

		if _, err := cfgFields.Write(raw); err != nil {
			return nil, errors.New(tr.TErr("cfg-body-to-bytes-error"))
		}
	}

	sizePacket := uint16(2 + 4 + 2 + 2 + cfgFields.Len() + 4)
	if err := write(cfgPacket, sizePacket); err != nil {
		return nil, errors.New(tr.TErr("cfg-write-size-error"))
	}

	if _, err := cfgPacket.Write(cfgFields.Bytes()); err != nil {
		return nil, errors.New(tr.TErr("cfg-write-packet-body-error"))
	}

	hash := util.FNV1aHash(cfgPacket.Bytes())
	if err := write(cfgPacket, hash); err != nil {
		return nil, err
	}

	final := new(bytes.Buffer)

	totalSize := uint16(2 + 1 + cfgPacket.Len())
	if err := write(final, totalSize); err != nil {
		return nil, errors.New(tr.TErr("cfg-write-final-size-error"))
	}

	if err := write(final, uint8(0x05)); err != nil {
		return nil, errors.New(tr.TErr("cfg-write-protocol-version-error"))
	}

	if _, err := final.Write(cfgPacket.Bytes()); err != nil {
		return nil, errors.New(tr.TErr("cfg-write-final-error"))
	}

	return final.Bytes(), nil
}

func (adm *Adm) getValueByType(tr locale.Translator, uid uint16, val any, size uint8) ([]byte, error) {
	fieldType, err := constants.GetCfgType(uid)
	if err != nil {
		return nil, nil
	}

	uidName, _ := constants.GetCfgName(uid)

	switch fieldType {
	case constants.UInt8:
		return adm.ByteUint8(tr, uidName, size, val)

	case constants.TwoUint8:
		return adm.ByteTwoUint8(tr, uidName, size, val)

	case constants.UInt16:
		return adm.ByteUInt16(tr, uidName, size, val)

	case constants.TwoUint16:
		return adm.ByteTwoUint16(tr, uidName, size, val)

	case constants.UInt32:
		return adm.ByteUInt32(tr, uidName, size, val)

	case constants.Uint32Array:
		return adm.ByteUint32Array(tr, uidName, size, val)

	case constants.Uint64Array:
		return adm.ByteUint64Array(tr, uidName, size, val)

	case constants.Phones:
		return adm.BytePhones(tr, uidName, size, val)

	case constants.Hex:
		return adm.ByteHex(tr, uidName, size, val)

	case constants.Ascii:
		return adm.ByteAscii(tr, uidName, size, val)

	case constants.SNSArray:
		return adm.ByteSNSArray(tr, uidName, size, val)

	case constants.MacAddressArray:
		return adm.ByteMacAddressArray(tr, uidName, size, val)

	case constants.ByteArray, constants.Raw:
		return adm.ByteArray(tr, uidName, size, val)

	default:
		return make([]byte, size), nil
	}
}

func isEmptySlice(val any) bool {
	if val == nil {
		return true
	}

	rv := reflect.ValueOf(val)
	return (rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array) && rv.Len() == 0
}
