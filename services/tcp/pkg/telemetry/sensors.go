package telemetry

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func ParseLineIntoTelemetry(tl *Telemetry, line string) error {
	headerRegexp := regexp.MustCompile(`^\[(\d+)(?:,([A-Za-z0-9]+))?\]:\s*([A-F0-9]+|[0-9]+)`)

	header := headerRegexp.FindStringSubmatch(line)
	if header == nil {
		return fmt.Errorf("invalid line header: %s", line)
	}

	index := header[1]
	protocol := header[2]
	mainValue := header[3]

	data := BLESensor{}

	fields := strings.Split(line, ";")

	for _, f := range fields {
		parts := strings.Split(strings.TrimSpace(f), ":")
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])

		switch key {
		case "R": // RSSI
			data.RSSI, _ = strconv.Atoi(val)

		case "V": // sensor battery voltage, V
			data.Voltage, _ = strconv.ParseFloat(val, 64)

		case "LMT": // time since last message
			data.LMT, _ = strconv.Atoi(val)

		case "L": // fuel level
			normalized := strings.ReplaceAll(val, ",", ".")

			if levelInt, err := strconv.Atoi(normalized); err == nil {
				data.FuelLevel = levelInt
				break
			}

			if levelFloat, err := strconv.ParseFloat(normalized, 64); err == nil {
				data.FuelLevel = int(levelFloat)
			}

		case "T": // temperature
			data.Temp, _ = strconv.Atoi(val)
		}
	}

	if protocol == "" && len(mainValue) >= 10 {
		key := tryNormalizeMac(mainValue)
		tl.BLESENSORINFO[key] = data
		return nil
	}

	if protocol == "BLE" {
		key := tryNormalizeMac(mainValue)
		tl.FUELINFO[key] = data
		return nil
	}

	if protocol == "485" {
		key := index
		tl.FUELINFO[key] = data
		return nil
	}

	return fmt.Errorf("unknown packet format: %s", line)
}

func ParseBLEString(tl *Telemetry, big string) error {
	packets := splitPackets(big)

	for _, p := range packets {
		if !strings.Contains(p, ",") {
			if err := ParseLineIntoTelemetry(tl, p); err != nil {
				return err
			}
		}
	}

	return nil
}

func ParseFUELString(tl *Telemetry, big string) error {
	packets := splitPackets(big)

	for _, p := range packets {
		if strings.Contains(p, ",") {
			if err := ParseLineIntoTelemetry(tl, p); err != nil {
				return err
			}
		}
	}

	return nil
}
