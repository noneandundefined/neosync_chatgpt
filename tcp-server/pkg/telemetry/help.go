package telemetry

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var packetRegexp = regexp.MustCompile(`(\[\d+(?:,[A-Za-z0-9]+)?\]:)`)

func tryNormalizeMac(raw string) string {
	cleaned := strings.ToUpper(strings.ReplaceAll(raw, ":", ""))

	if len(cleaned) != 12 {
		return raw
	}

	if _, err := strconv.ParseUint(cleaned, 16, 64); err != nil {
		return raw
	}

	return fmt.Sprintf("%s:%s:%s:%s:%s:%s", cleaned[0:2], cleaned[2:4], cleaned[4:6], cleaned[6:8], cleaned[8:10], cleaned[10:12])
}

func splitPackets(big string) []string {
	matches := packetRegexp.FindAllStringIndex(big, -1)
	if len(matches) == 0 {
		return nil
	}

	var parts []string

	for i := 0; i < len(matches); i++ {
		start := matches[i][0]

		var end int
		if i+1 < len(matches) {
			end = matches[i+1][0]
		} else {
			end = len(big)
		}

		segment := strings.TrimSpace(big[start:end])
		if segment != "" {
			parts = append(parts, segment)
		}
	}

	return parts
}
