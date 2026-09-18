package adm

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"neomatica/neosync/infra/locale"
	"neomatica/neosync/pkg/adm/help"
	"reflect"
	"strconv"
	"strings"
)

func (adm *Adm) ByteUint8(tr locale.Translator, uidName string, size uint8, val any) ([]byte, error) {
	rv := reflect.ValueOf(val)

	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return []byte{byte(rv.Int())}, nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return []byte{byte(rv.Uint())}, nil

	case reflect.Float32, reflect.Float64:
		return []byte{byte(rv.Float())}, nil

	case reflect.Slice, reflect.Array:
		if rv.Len() == 1 {
			return adm.ByteUint8(tr, uidName, size, rv.Index(0).Interface())
		}
	}

	return nil, fmt.Errorf("%s %s", tr.TErr("invalid-type"), uidName)
}

func (adm *Adm) ByteTwoUint8(tr locale.Translator, uidName string, size uint8, val any) ([]byte, error) {
	var arr []interface{}

	rv := reflect.ValueOf(val)
	switch rv.Kind() {
	case reflect.Slice, reflect.Array:
		if rv.Len() != 2 {
			return nil, fmt.Errorf("%s %s", tr.TErr("invalid-type"), uidName)
		}

		arr = make([]interface{}, 2)

		for i := 0; i < 2; i++ {
			arr[i] = rv.Index(i).Interface()
		}

	default:
		return nil, fmt.Errorf("%s %s", tr.TErr("invalid-type"), uidName)
	}

	result := make([]byte, 2)
	for i, v := range arr {
		rvElem := reflect.ValueOf(v)
		switch rvElem.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			result[i] = byte(rvElem.Int())

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			result[i] = byte(rvElem.Uint())

		case reflect.Float32, reflect.Float64:
			result[i] = byte(rvElem.Float())

		default:
			return nil, fmt.Errorf("%s %s", tr.TErr("invalid-type"), uidName)
		}
	}

	return result, nil
}

func (adm *Adm) ByteUInt16(tr locale.Translator, uidName string, size uint8, val any) ([]byte, error) {
	var n uint16
	rv := reflect.ValueOf(val)

	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n = uint16(rv.Int())

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		n = uint16(rv.Uint())

	case reflect.Float32, reflect.Float64:
		n = uint16(rv.Float())

	case reflect.Slice, reflect.Array:
		if rv.Len() != 1 {
			return nil, fmt.Errorf("%s %s", tr.TErr("invalid-type"), uidName)
		}

		b, err := adm.ByteUInt16(tr, uidName, size, rv.Index(0).Interface())
		if err != nil {
			return nil, err
		}

		return b, nil

	default:
		return nil, fmt.Errorf("%s %s", tr.TErr("invalid-type"), uidName)
	}

	return help.HelpUint16ToBytes(n), nil
}

func (adm *Adm) ByteTwoUint16(tr locale.Translator, uidName string, size uint8, val any) ([]byte, error) {
	var arrs []interface{}

	switch v := val.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, string:
		arrs = []interface{}{v}

	case []uint8:
		for _, n := range v {
			arrs = append(arrs, uint32(n))
		}

	case []uint16:
		for _, n := range v {
			arrs = append(arrs, uint32(n))
		}

	case []uint32:
		for _, n := range v {
			arrs = append(arrs, n)
		}

	case []int:
		for _, n := range v {
			arrs = append(arrs, uint32(n))
		}

	case []float32:
		for _, n := range v {
			arrs = append(arrs, uint32(n))
		}

	case []float64:
		for _, n := range v {
			arrs = append(arrs, uint32(n))
		}

	case []interface{}:
		arrs = v

	default:
		rv := reflect.ValueOf(val)

		if rv.Kind() == reflect.Array {
			arrs = make([]interface{}, rv.Len())
			for i := 0; i < rv.Len(); i++ {
				arrs[i] = rv.Index(i).Interface()
			}
		} else {
			return nil, fmt.Errorf("%s %s", tr.TErr("invalid-type"), uidName)
		}
	}

	if len(arrs) < 1 || len(arrs) > 2 {
		return nil, fmt.Errorf("%s %s", tr.TErr("invalid-type"), uidName)
	}

	buf := new(bytes.Buffer)

	for _, elem := range arrs {
		var num uint16

		switch n := elem.(type) {
		case nil:
			num = 0
		case int:
			num = uint16(n)
		case int8:
			num = uint16(n)
		case int16:
			num = uint16(n)
		case int32:
			num = uint16(n)
		case int64:
			num = uint16(n)
		case uint:
			num = uint16(n)
		case uint8:
			num = uint16(n)
		case uint16:
			num = n
		case uint32:
			num = uint16(n)
		case uint64:
			num = uint16(n)
		case float32:
			num = uint16(n)
		case float64:
			num = uint16(n)
		case []interface{}:
			if len(n) == 1 {
				b, err := adm.ByteUInt16(tr, uidName, size, n[0])
				if err != nil {
					return nil, err
				}

				buf.Write(b)
				continue
			} else {
				return nil, fmt.Errorf("%s %s", tr.TErr("invalid-type"), uidName)
			}

		default:
			return nil, fmt.Errorf("%s %s", tr.TErr("invalid-type"), uidName)
		}

		buf.Write(help.HelpUint16ToBytes(num))
	}

	result := buf.Bytes()

	if uint16(len(result)) < uint16(size) {
		result = append(result, make([]byte, int(size)-len(result))...)
	}

	if uint16(len(result)) > uint16(size) {
		return nil, fmt.Errorf("%s %s", tr.TErr("invalid-type"), uidName)
	}

	return result, nil
}

func (adm *Adm) ByteUInt32(tr locale.Translator, uidName string, size uint8, val any) ([]byte, error) {
	var n uint32
	rv := reflect.ValueOf(val)

	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n = uint32(rv.Int())

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		n = uint32(rv.Uint())

	case reflect.Float32, reflect.Float64:
		n = uint32(rv.Float())

	case reflect.String:
		if val, err := strconv.ParseUint(rv.String(), 10, 32); err == nil {
			n = uint32(val)
		} else {
			return nil, fmt.Errorf("%s %s", tr.TErr("invalid-type"), uidName)
		}

	case reflect.Slice, reflect.Array:
		buf := new(bytes.Buffer)
		for i := 0; i < rv.Len(); i++ {
			b, err := adm.ByteUInt32(tr, uidName, size, rv.Index(i).Interface())
			if err != nil {
				return nil, err
			}

			buf.Write(b)
		}

		return buf.Bytes(), nil
	default:
		return nil, fmt.Errorf("%s %s", tr.TErr("invalid-type"), uidName)
	}

	return help.HelpUint32ToBytes(n), nil
}

func (adm *Adm) ByteUint32Array(tr locale.Translator, uidName string, size uint8, val any) ([]byte, error) {
	var arrs []interface{}

	rv := reflect.ValueOf(val)

	switch rv.Kind() {
	case reflect.Slice:
		arrs = make([]interface{}, rv.Len())
		for i := 0; i < rv.Len(); i++ {
			arrs[i] = rv.Index(i).Interface()
		}

	case reflect.Array:
		arrs = make([]interface{}, rv.Len())
		for i := 0; i < rv.Len(); i++ {
			arrs[i] = rv.Index(i).Interface()
		}

	default:
		return nil, fmt.Errorf("%s %s", tr.TErr("invalid-type"), uidName)
	}

	if len(arrs) == 0 {
		return nil, fmt.Errorf("%s %s", tr.TErr("invalid-type"), uidName)
	}

	buf := new(bytes.Buffer)
	for _, elem := range arrs {
		var num uint32

		rvElem := reflect.ValueOf(elem)

		switch rvElem.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			num = uint32(rvElem.Int())

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			num = uint32(rvElem.Uint())

		case reflect.Float32, reflect.Float64:
			num = uint32(math.Round(rvElem.Float()))

		case reflect.Slice, reflect.Array:
			b, err := adm.ByteUint32Array(tr, uidName, size, rvElem.Interface())
			if err != nil {
				return nil, err
			}

			buf.Write(b)
			continue

		default:
			return nil, fmt.Errorf("%s %s", tr.TErr("invalid-type"), uidName)
		}

		buf.Write(help.HelpUint32ToBytes(num))
	}

	if size > 0 {
		diff := int(size) - buf.Len()
		if diff > 0 {
			buf.Write(make([]byte, diff))
		}
	}

	return buf.Bytes(), nil
}

func (adm *Adm) ByteUint64Array(tr locale.Translator, uidName string, size uint8, val any) ([]byte, error) {
	var arrs []interface{}

	rv := reflect.ValueOf(val)

	switch rv.Kind() {
	case reflect.Slice:
		arrs = make([]interface{}, rv.Len())
		for i := 0; i < rv.Len(); i++ {
			arrs[i] = rv.Index(i).Interface()
		}

	case reflect.Array:
		arrs = make([]interface{}, rv.Len())
		for i := 0; i < rv.Len(); i++ {
			arrs[i] = rv.Index(i).Interface()
		}
	}

	if len(arrs) == 0 {
		return nil, fmt.Errorf("%s %s", tr.TErr("invalid-type"), uidName)
	}

	buf := new(bytes.Buffer)
	for _, elem := range arrs {
		var num uint64

		rvElem := reflect.ValueOf(elem)

		switch rvElem.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			num = uint64(rvElem.Int())

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			num = uint64(rvElem.Uint())

		case reflect.Float32, reflect.Float64:
			num = uint64(math.Round(rvElem.Float()))

		case reflect.Slice, reflect.Array:
			b, err := adm.ByteUint64Array(tr, uidName, size, rvElem.Interface())
			if err != nil {
				return nil, err
			}

			buf.Write(b)
			continue

		default:
			return nil, fmt.Errorf("%s %s", tr.TErr("invalid-type"), uidName)
		}

		buf.Write(help.HelpUint64ToBytes(num))
	}

	if size > 0 {
		diff := int(size) - buf.Len()
		if diff > 0 {
			buf.Write(make([]byte, diff))
		}
	}

	return buf.Bytes(), nil
}

func (adm *Adm) BytePhones(tr locale.Translator, uidName string, size uint8, val any) ([]byte, error) {
	buffer := make([]byte, size)

	segmentSize := int(size) / 4

	var values []string
	rv := reflect.ValueOf(val)

	if rv.Kind() == reflect.String {
		if err := json.Unmarshal([]byte(rv.String()), &values); err != nil {
			values = []string{rv.String()}
		}
	} else if rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array {
		for i := 0; i < rv.Len(); i++ {
			if str, ok := rv.Index(i).Interface().(string); ok {
				values = append(values, str)
			}
		}
	}

	if len(values) == 0 {
		return buffer, nil
	}

	for i := 0; i < 4; i++ {
		startIndex := i * segmentSize
		endIndex := startIndex + segmentSize
		if endIndex > int(size) {
			endIndex = int(size)
		}

		if i < len(values) {
			str := values[i]
			copy(buffer[startIndex:endIndex], []byte(str))
		}
	}

	return buffer, nil
}

func (adm *Adm) ByteHex(tr locale.Translator, uidName string, size uint8, val any) ([]byte, error) {
	if b, ok := val.(string); ok {
		bb, err := hex.DecodeString(b)
		if err != nil {
			return nil, err
		}

		return bb, nil
	}

	return nil, fmt.Errorf("%s %s", tr.TErr("invalid-type"), uidName)
}

func (adm *Adm) ByteAscii(tr locale.Translator, uidName string, size uint8, val any) ([]byte, error) {
	b := make([]byte, size)

	var arr []string
	rv := reflect.ValueOf(val)

	if rv.Kind() == reflect.String {
		if err := json.Unmarshal([]byte(rv.String()), &arr); err == nil && len(arr) > 0 {
		} else {
			copy(b, []byte(rv.String()))
			return b, nil
		}
	} else if rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array {
		for i := 0; i < rv.Len(); i++ {
			if str, ok := rv.Index(i).Interface().(string); ok {
				arr = append(arr, str)
			} else if x, ok := rv.Index(i).Interface().(float64); ok {
				arr = append(arr, fmt.Sprintf("%v", byte(x)))
			}
		}
	}

	if len(arr) == 0 {
		return b, nil
	}

	chunk := int(size) / len(arr)

	for i, str := range arr {
		start := i * chunk
		end := start + chunk
		if end > int(size) {
			end = int(size)
		}

		copy(b[start:end], []byte(str))
	}

	return b, nil
}

func (adm *Adm) ByteSNSArray(tr locale.Translator, uidName string, size uint8, val any) ([]byte, error) {
	var sns []string

	if val != nil {
		rv := reflect.ValueOf(val)
		if rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array {
			sns = make([]string, rv.Len())
			for i := 0; i < rv.Len(); i++ {
				if s, ok := rv.Index(i).Interface().(string); ok {
					sns[i] = s
				} else {
					return nil, fmt.Errorf("%s %s", tr.TErr("invalid-type"), uidName)
				}
			}
		} else {
			return nil, fmt.Errorf("%s %s", tr.TErr("invalid-type"), uidName)
		}
	}

	var result []byte

	for _, sn := range sns {
		if len(sn) != 16 {
			return nil, errors.New(tr.TErr("invalid-1wire-sn-size"))
		}

		b, err := hex.DecodeString(sn)
		if err != nil {
			return nil, fmt.Errorf(tr.TErr("invalid-sn-format"), sn)
		}

		result = append(result, b...)
	}

	if uint16(len(result)) < uint16(size) {
		padding := make([]byte, uint16(size)-uint16(len(result)))
		result = append(result, padding...)
	} else if uint16(len(result)) > uint16(size) {
		return nil, fmt.Errorf(tr.TErr("invalid-sn-format-size-exceeded"), size)
	}

	return result, nil
}

func (adm *Adm) ByteMacAddressArray(tr locale.Translator, uidName string, size uint8, val any) ([]byte, error) {
	result := make([]byte, size)

	if val == nil {
		return result, nil
	}

	var macs []string
	rv := reflect.ValueOf(val)

	switch rv.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < rv.Len(); i++ {
			if s, ok := rv.Index(i).Interface().(string); ok {
				macs = append(macs, s)
			} else {
				return result, fmt.Errorf("%s %s", tr.TErr("invalid-type"), uidName)
			}
		}
	default:
		return result, fmt.Errorf("%s %s", tr.TErr("invalid-type"), uidName)
	}

	if len(macs) == 0 {
		return result, nil
	}

	blockSize := int(size) / len(macs)
	if blockSize < 6 {
		return nil, errors.New(tr.TErr("size-too-small-for-mac-count"))
	}

	offset := 0
	for _, mac := range macs {
		clean := strings.ReplaceAll(mac, ":", "")
		if len(clean) != 12 {
			return nil, fmt.Errorf(tr.TErr("invalid-mac-format"), mac)
		}

		b, err := hex.DecodeString(clean)
		if err != nil {
			return nil, fmt.Errorf(tr.TErr("mac-decode-error"), mac, err)
		}

		copy(result[offset:], b)
		offset += blockSize
	}

	return result, nil
}

func (adm *Adm) ByteArray(tr locale.Translator, uidName string, size uint8, val any) ([]byte, error) {
	var arr []interface{}

	switch v := val.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, string:
		arr = []interface{}{v}
	case []int:
		for _, n := range v {
			arr = append(arr, n)
		}
	case []int8:
		for _, n := range v {
			arr = append(arr, n)
		}
	case []int16:
		for _, n := range v {
			arr = append(arr, n)
		}
	case []int32:
		for _, n := range v {
			arr = append(arr, n)
		}
	case []int64:
		for _, n := range v {
			arr = append(arr, n)
		}
	case []uint:
		for _, n := range v {
			arr = append(arr, n)
		}
	case []uint8:
		for _, n := range v {
			arr = append(arr, n)
		}
	case []uint16:
		for _, n := range v {
			arr = append(arr, n)
		}
	case []uint32:
		for _, n := range v {
			arr = append(arr, n)
		}
	case []uint64:
		for _, n := range v {
			arr = append(arr, n)
		}
	case []float32:
		for _, n := range v {
			arr = append(arr, n)
		}
	case []float64:
		for _, n := range v {
			arr = append(arr, n)
		}
	case []interface{}:
		arr = v

	default:
		rv := reflect.ValueOf(val)

		if rv.Kind() == reflect.Array {
			arr = make([]interface{}, rv.Len())
			for i := 0; i < rv.Len(); i++ {
				arr[i] = rv.Index(i).Interface()
			}
		} else {
			return nil, fmt.Errorf("%s %s", tr.TErr("invalid-type"), uidName)
		}
	}

	if len(arr) == 0 {
		return make([]byte, size), nil
	}

	b := make([]byte, size)
	for i := 0; i < len(arr) && i < int(size); i++ {
		switch v := arr[i].(type) {
		case int, int8, int16, int32, int64:
			b[i] = byte(reflect.ValueOf(v).Int())
		case uint, uint8, uint16, uint32, uint64:
			b[i] = byte(reflect.ValueOf(v).Uint())
		case float32:
			b[i] = byte(v)
		case float64:
			b[i] = byte(v)

		default:
			return nil, fmt.Errorf("%s %s", tr.TErr("invalid-type"), uidName)
		}
	}

	return b, nil
}
