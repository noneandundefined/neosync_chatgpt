package telemetry

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"neomatica/neosync-tcp/infra/logger"
)

// Fixed header (ADM-5 / ADM-6), bytes 0..33 — protocol ADM 1.08.
const (
	adm6HeaderSize     = 34
	adm6StatusOffset   = 7
	adm6LatOffset      = 9
	adm6LonOffset      = 13
	adm6SatCountOffset = 25
	adm6DateTimeOffset = 26
	adm6VPowerOffset   = 30
	adm6VBatteryOffset = 32
	adm6BlocksOffset   = 34

	adm6StatusInvalidCoordsMask = 1 << 5

	adm6TypeBitAccel   = 2
	adm6TypeBitAnalog  = 3
	adm6TypeBitDigital = 4

	adm6MinDateTime = 946684800  // 2000-01-01
	adm6MaxDateTime = 4102444800 // 2100-01-01
)

// NormalizeADM6Packet unwraps COM0 payload when device returns ASCII-hex inside ADM RC frame.
func NormalizeADM6Packet(packet []byte) ([]byte, error) {
	if len(packet) < 3 {
		return nil, fmt.Errorf("packet too short")
	}

	if len(packet) == int(packet[2]) {
		return packet, nil
	}

	if len(packet) >= 4 && int(binary.LittleEndian.Uint16(packet[0:2])) == len(packet) {
		payload := packet[3:]
		for len(payload) > 0 && payload[len(payload)-1] == 0 {
			payload = payload[:len(payload)-1]
		}

		inner, err := hex.DecodeString(string(payload))
		if err != nil {
			return nil, fmt.Errorf("decode wrapped COM0 payload: %w", err)
		}

		if len(inner) < 3 {
			return nil, fmt.Errorf("wrapped COM0 payload too short")
		}

		return inner, nil
	}

	return packet, nil
}

func hasADM6Block(packetType uint8, bit uint8) bool {
	return packetType&(1<<bit) != 0
}

func isPlausibleDateTime(dt uint32) bool {
	return dt >= adm6MinDateTime && dt <= adm6MaxDateTime
}

// ParseStatCountTelemetry parses ADM-6 main data packet (COM0 response).
func ParseStatCountTelemetry(tl *Telemetry, packet []byte) error {
	if tl == nil {
		return fmt.Errorf("telemetry is nil")
	}

	normalized, err := NormalizeADM6Packet(packet)
	if err != nil {
		return err
	}

	if len(normalized) < adm6HeaderSize {
		logger.Error("ParseStatCountTelemetry: packet length %d is smaller than expected %d", len(normalized), adm6HeaderSize)
		return fmt.Errorf("packet length %d is smaller than expected %d", len(normalized), adm6HeaderSize)
	}

	if int(normalized[2]) != len(normalized) {
		logger.Error("ParseStatCountTelemetry: size mismatch: header=%d actual=%d", normalized[2], len(normalized))
		return fmt.Errorf("size mismatch: header=%d actual=%d", normalized[2], len(normalized))
	}

	satByte := normalized[adm6SatCountOffset]
	tl.GlonassSatellites = int(satByte >> 4)
	tl.GPSSatellites = int(satByte & 0x0F)

	dt := binary.LittleEndian.Uint32(normalized[adm6DateTimeOffset : adm6DateTimeOffset+4])
	if isPlausibleDateTime(dt) {
		tl.DateTime = dt
	} else if dt != 0 {
		logger.Warning("ParseStatCountTelemetry: implausible DATE_TIME in stat-count packet: dt=%d", dt)
	}

	tl.VPower = int(binary.LittleEndian.Uint16(normalized[adm6VPowerOffset : adm6VPowerOffset+2]))
	tl.VBattery = int(binary.LittleEndian.Uint16(normalized[adm6VBatteryOffset : adm6VBatteryOffset+2]))

	packetType := normalized[3]
	pos := adm6BlocksOffset

	if hasADM6Block(packetType, adm6TypeBitAccel) {
		if len(normalized) < pos+4 {
			return fmt.Errorf("packet too short for accel block")
		}
		pos += 4
	}

	tl.InA0 = 0
	tl.InA1 = 0
	if hasADM6Block(packetType, adm6TypeBitAnalog) {
		if len(normalized) < pos+12 {
			return fmt.Errorf("packet too short for analog block")
		}

		tl.InA0 = int(binary.LittleEndian.Uint16(normalized[pos : pos+2]))
		tl.InA1 = int(binary.LittleEndian.Uint16(normalized[pos+2 : pos+4]))
		pos += 12
	}

	tl.InD0 = 0
	tl.InD1 = 0
	if hasADM6Block(packetType, adm6TypeBitDigital) {
		if len(normalized) < pos+8 {
			return fmt.Errorf("packet too short for digital block")
		}

		tl.InD0 = binary.LittleEndian.Uint32(normalized[pos : pos+4])
		tl.InD1 = binary.LittleEndian.Uint32(normalized[pos+4 : pos+8])
	}

	return nil
}
