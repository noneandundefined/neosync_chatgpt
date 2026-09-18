package constants

var (
	REBOOT_COMMAND                   string = "RESET"
	UPDATE_COMMAND                   string = "UPDATE"
	BLEAUTOCATCH_FIND_NEARBY_COMMAND string = "BLEAUTOCATCH 65535,-60"
	BLEAUTOCATCH_FIND_ALL_COMMAND    string = "BLEAUTOCATCH 65535,0"
	BLEAUTOCATCH_STOP_FIND_COMMAND   string = "BLEAUTOCATCH 0"
	ONEWIRE_FIND_COMMAND             string = "OWCONFIG"
	BLESENSORINFO_COMMAND            string = "BLESENSORINFO"
	FUELINFO_COMMAND                 string = "FUELINFO"
	BLESENSOR_COMMAND                string = "BLESENSOR"
	FUEL_COMMAND                     string = "FUEL"
	COM0_COMMAND                     string = "COM0"
	ADM20INFO_COMMAND                string = "ADM20INFO"
	ERASE_EEPROM_COMMAND             string = "ERASE EEPROM" // reset terminal and reboot
	ERASE_FLASH_COMMAND              string = "ERASE FLASH"  // reset memory terminal
)

var TELEMATRY_COMMANDS = []string{
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
