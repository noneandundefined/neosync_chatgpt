package constants

import "strings"

const (
	FW_UPDATE_STATUS_IDLE    = "idle"
	FW_UPDATE_STATUS_PENDING = "pending"
)

func IsFirmwareUpdateCommand(command string) bool {
	cmd := strings.ToUpper(strings.TrimSpace(command))
	return cmd == UPDATE_COMMAND || strings.HasPrefix(cmd, UPDATE_COMMAND+" ")
}
