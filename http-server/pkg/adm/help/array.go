package help

import (
	"encoding/binary"
	"strings"
	"unicode/utf8"
)

func HelpConvertToFixedOrdered(raw []byte, count int) []string {
	chunk := len(raw) / count
	result := make([]string, count)

	for i := 0; i < count; i++ {
		start := i * chunk
		end := start + chunk

		if end > len(raw) {
			end = len(raw)
		}

		part := sanitizeString(string(raw[start:end]))
		if part != "" {
			result[i] = part
		} else {
			result[i] = ""
		}
	}

	return result
}

func HelpConvertToPhones(raw []byte, count int) []string {
	chunk := len(raw) / count
	result := make([]string, 0, count)

	allEmpty := true

	for i := 0; i < count; i++ {
		start := i * chunk
		end := start + chunk

		if end > len(raw) {
			end = len(raw)
		}

		part := sanitizeString(string(raw[start:end]))
		if part != "" {
			result = append(result, part)
			allEmpty = false
		}
	}

	if allEmpty {
		return []string{}
	}

	return result
}

func HelpConvertTwoArrayUint8(raw []byte) []int {
	nums := make([]int, 2)

	if len(raw) >= 1 {
		nums[0] = int(raw[0])
	}
	if len(raw) >= 2 {
		nums[1] = int(raw[1])
	}

	return nums
}

func HelpConvertTwoArrayUint16(raw []byte) [2]any {
	var ports [2]any
	if len(raw) >= 4 {
		first := binary.LittleEndian.Uint16(raw[:2])
		second := binary.LittleEndian.Uint16(raw[2:4])

		if first != 0 {
			ports[0] = first
		} else {
			ports[0] = nil
		}

		if second != 0 {
			ports[1] = second
		} else {
			ports[1] = nil
		}
	}

	return ports
}

func sanitizeString(s string) string {
	s = strings.TrimRight(s, "\x00")
	s = strings.ReplaceAll(s, "\x00", "")

	if utf8.ValidString(s) {
		return s
	}

	return strings.ToValidUTF8(s, "")
}
