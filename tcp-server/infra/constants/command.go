package constants

var (
	WHO_COMMAND                      string = "WHO"
	REBOOT_COMMAND                   string = "RESET"
	UPDATE_COMMAND                   string = "UPDATE"
	BLEAUTOCATCH_FIND_NEARBY_COMMAND string = "BLEAUTOCATCH 65535,-60"
	BLEAUTOCATCH_FIND_ALL_COMMAND    string = "BLEAUTOCATCH 65535,0"
	BLEAUTOCATCH_STOP_FIND_COMMAND   string = "BLEAUTOCATCH 0"
	ONEWIRE_FIND_COMMAND             string = "OWCONFIG"
	ADM20INFO_COMMAND                string = "ADM20INFO"
	BLESENSORINFO_COMMAND            string = "BLESENSORINFO"
	FUELINFO_COMMAND                 string = "FUELINFO"
	COM0_COMMAND                     string = "COM0"
)

var TelemetryCommands = []string{
	COM0_COMMAND,
	ADM20INFO_COMMAND,
	BLESENSORINFO_COMMAND,
	FUELINFO_COMMAND,
}

var (
	CMD_STATUS_PENDING        string = "pending"
	CMD_STATUS_INPROGRESS     string = "inprogress"
	CMD_STATUS_EXECUTIONERROR string = "executionerror"
	CMD_STATUS_NOTCOMPLETED   string = "notcompleted"
	CMD_STATUS_COMPLETED      string = "completed"
)

var (
	CMD_SEND_MODE_INSTANT    string = "instant"
	CMD_SEND_MODE_ON_CONNECT string = "on_connect"
)
