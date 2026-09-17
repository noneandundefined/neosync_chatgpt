package constants

import "fmt"

type FieldType string

const (
	UInt8           FieldType = "UInt8"
	TwoUint8        FieldType = "TwoUint8"
	UInt16          FieldType = "UInt16"
	UInt32          FieldType = "UInt32"
	TwoUint16       FieldType = "TwoUint16"
	MacAddressArray FieldType = "MacAddressArray"
	Hex             FieldType = "Hex"
	SNSArray        FieldType = "SNSArray"
	Uint32Array     FieldType = "Uint32Array"
	Uint64Array     FieldType = "Uint64Array"
	Phones          FieldType = "Phones"
	Ascii           FieldType = "Ascii"
	ByteArray       FieldType = "ByteArray"
	Raw             FieldType = "Raw"
)

type FieldSchema struct {
	UID        uint16    `json:"uid"`
	ResolveUID uint16    `json:"resolve_uid"`
	Name       string    `json:"name"`
	Type       FieldType `json:"type"`
	Default    any       `json:"default"`
	Min        *int      `json:"min,omitempty"`
	Max        *int      `json:"max,omitempty"`
	MinLength  *int      `json:"min_length,omitempty"`
	MaxLength  *int      `json:"max_length,omitempty"`
	MaxItems   *int      `json:"max_items,omitempty"`
	IsDigits   *bool     `json:"is_digits,omitempty"`
	Section    string    `json:"section"`
}

var CfgSchema = map[uint16]FieldSchema{
	4:   {UID: 4, ResolveUID: 4, Name: "server_host", Type: Ascii, Default: []string{"", ""}, Section: "server;sim"},
	5:   {UID: 5, ResolveUID: 5, Name: "server_port", Type: TwoUint16, Default: []int{0, 0}, Min: intPtr(0), Max: intPtr(65535), IsDigits: boolPtr(true), Section: "server;sim"},
	6:   {UID: 6, ResolveUID: 6, Name: "device_id", Type: UInt32, Default: 1, IsDigits: boolPtr(true), Section: "server"},
	7:   {UID: 7, ResolveUID: 7, Name: "traffic_protocol", Type: TwoUint8, Default: []int{0, 0}, Min: intPtr(0), Max: intPtr(255), IsDigits: boolPtr(true), Section: "server"},
	8:   {UID: 8, ResolveUID: 8, Name: "black_box_structure", Type: UInt32, Default: 0, Min: intPtr(0), Max: intPtr(4_294_967_295), IsDigits: boolPtr(true), Section: "server"},
	9:   {UID: 9, ResolveUID: 9, Name: "auth_phone", Type: Phones, Default: []string{}, MaxItems: intPtr(4), IsDigits: boolPtr(true), Section: "device"},
	10:  {UID: 10, ResolveUID: 10, Name: "auth_global_pass", Type: Ascii, Default: "0", MaxLength: intPtr(8), Section: "device"},
	11:  {UID: 11, ResolveUID: 11, Name: "sim_status", Type: UInt8, Default: 0, Min: intPtr(0), Max: intPtr(255), IsDigits: boolPtr(true), Section: "sim"},
	12:  {UID: 12, ResolveUID: 12, Name: "ops_white_list_deprecated", Type: ByteArray, Default: nil, IsDigits: boolPtr(true), Section: ""},
	13:  {UID: 13, ResolveUID: 13, Name: "device_name", Type: Ascii, Default: nil, MaxLength: intPtr(15), Section: "device"},
	14:  {UID: 14, ResolveUID: 14, Name: "static_mode", Type: UInt8, Default: 4, Min: intPtr(0), Max: intPtr(255), IsDigits: boolPtr(true), Section: "track"},
	15:  {UID: 15, ResolveUID: 15, Name: "static_ain_number", Type: UInt8, Default: 255, Min: intPtr(0), Max: intPtr(255), IsDigits: boolPtr(true), Section: "track"},
	16:  {UID: 16, ResolveUID: 16, Name: "accel_instatic_timeout_sec", Type: UInt8, Default: 30, IsDigits: boolPtr(true), Section: ""},
	17:  {UID: 17, ResolveUID: 17, Name: "accel_outstatic_threshold_mg", Type: UInt8, Default: 6, IsDigits: boolPtr(true), Section: ""},
	18:  {UID: 18, ResolveUID: 18, Name: "accel_outstatic_timeout_samples", Type: UInt8, Default: 1, IsDigits: boolPtr(true), Section: ""},
	19:  {UID: 19, ResolveUID: 19, Name: "ain_false_high", Type: TwoUint16, Default: []int{4000, 4000}, Min: intPtr(0), Max: intPtr(60000), IsDigits: boolPtr(true), Section: "inputs"},
	20:  {UID: 20, ResolveUID: 20, Name: "ain_true_low", Type: TwoUint16, Default: []int{7000, 7000}, Min: intPtr(0), Max: intPtr(60000), IsDigits: boolPtr(true), Section: "inputs"},
	21:  {UID: 21, ResolveUID: 21, Name: "output", Type: UInt8, Default: 0, Min: intPtr(0), Max: intPtr(255), IsDigits: boolPtr(true), Section: "outputs"},
	22:  {UID: 22, ResolveUID: 22, Name: "points_period_move", Type: UInt16, Default: 30, Min: intPtr(0), Max: intPtr(65535), IsDigits: boolPtr(true), Section: "track"},
	23:  {UID: 23, ResolveUID: 23, Name: "points_period_hold", Type: UInt16, Default: 300, Min: intPtr(0), Max: intPtr(65535), IsDigits: boolPtr(true), Section: "track"},
	24:  {UID: 24, ResolveUID: 24, Name: "points_period_alarm", Type: UInt8, Default: 255, IsDigits: boolPtr(true), Section: ""},
	25:  {UID: 25, ResolveUID: 25, Name: "points_track_speed_stop", Type: UInt8, Default: 4, Min: intPtr(0), Max: intPtr(255), IsDigits: boolPtr(true), Section: "track"},
	26:  {UID: 26, ResolveUID: 26, Name: "points_track_acceleration", Type: UInt8, Default: 10, Min: intPtr(0), Max: intPtr(255), IsDigits: boolPtr(true), Section: "track"},
	27:  {UID: 27, ResolveUID: 27, Name: "points_track_distance", Type: UInt16, Default: 1000, Min: intPtr(0), Max: intPtr(65535), IsDigits: boolPtr(true), Section: "track"},
	28:  {UID: 28, ResolveUID: 28, Name: "points_track_course", Type: ByteArray, Default: []int{20, 10, 5}, Min: intPtr(0), Max: intPtr(255), IsDigits: boolPtr(true), Section: "track"},
	29:  {UID: 29, ResolveUID: 29, Name: "points_track_crosstrack", Type: ByteArray, Default: []int{6, 6, 6}, Min: intPtr(0), Max: intPtr(255), IsDigits: boolPtr(true), Section: "track"},
	31:  {UID: 31, ResolveUID: 31, Name: "device_function_2", Type: UInt8, Default: 1, Min: intPtr(0), Max: intPtr(255), IsDigits: boolPtr(true), Section: "track"},
	32:  {UID: 32, ResolveUID: 32, Name: "opactive_in", Type: UInt8, Default: 255, IsDigits: boolPtr(true), Section: ""},
	33:  {UID: 33, ResolveUID: 33, Name: "gsmsignal_in", Type: UInt8, Default: 2, IsDigits: boolPtr(true), Section: ""},
	34:  {UID: 34, ResolveUID: 34, Name: "diagnostic_channels", Type: UInt8, Default: 0, IsDigits: boolPtr(true), Section: ""},
	35:  {UID: 35, ResolveUID: 35, Name: "beacon_mode", Type: UInt8, Default: nil, IsDigits: boolPtr(true), Section: ""},
	36:  {UID: 36, ResolveUID: 36, Name: "alarm_phone", Type: Phones, Default: []string{}, MaxItems: intPtr(4), IsDigits: boolPtr(true), Section: "events"},
	37:  {UID: 37, ResolveUID: 37, Name: "ble_adm_sensor_address_list", Type: MacAddressArray, Default: []string{}, Section: "bluetooth"},
	38:  {UID: 38, ResolveUID: 38, Name: "temp_min_value", Type: ByteArray, Default: []int{0, 0, 0, 0, 0}, IsDigits: boolPtr(true), Section: ""},
	39:  {UID: 39, ResolveUID: 39, Name: "temp_max_value", Type: ByteArray, Default: []int{0, 0, 0, 0, 0}, IsDigits: boolPtr(true), Section: ""},
	40:  {UID: 40, ResolveUID: 40, Name: "light_min_value", Type: ByteArray, Default: []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, IsDigits: boolPtr(true), Section: ""},
	41:  {UID: 41, ResolveUID: 41, Name: "light_max_value", Type: ByteArray, Default: []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, IsDigits: boolPtr(true), Section: ""},
	42:  {UID: 42, ResolveUID: 42, Name: "humid_min_value", Type: ByteArray, Default: []int{0, 0, 0, 0, 0}, IsDigits: boolPtr(true), Section: ""},
	43:  {UID: 43, ResolveUID: 43, Name: "humid_max_value", Type: ByteArray, Default: []int{0, 0, 0, 0, 0}, IsDigits: boolPtr(true), Section: ""},
	44:  {UID: 44, ResolveUID: 44, Name: "ble_fuel_sensor_address_list", Type: MacAddressArray, Default: nil, Section: "bluetooth"},
	45:  {UID: 45, ResolveUID: 45, Name: "scan_interval", Type: UInt8, Default: nil, IsDigits: boolPtr(true), Section: ""},
	46:  {UID: 46, ResolveUID: 46, Name: "scan_window", Type: UInt8, Default: nil, IsDigits: boolPtr(true), Section: ""},
	47:  {UID: 47, ResolveUID: 47, Name: "scan_timeout", Type: UInt8, Default: nil, IsDigits: boolPtr(true), Section: ""},
	48:  {UID: 48, ResolveUID: 48, Name: "ble_enabled", Type: UInt8, Default: nil, IsDigits: boolPtr(true), Section: ""},
	49:  {UID: 49, ResolveUID: 49, Name: "send_light_data_mode", Type: UInt8, Default: 0, IsDigits: boolPtr(true), Section: ""},
	50:  {UID: 50, ResolveUID: 50, Name: "bleprotocol", Type: UInt8, Default: 255, IsDigits: boolPtr(true), Section: ""},
	51:  {UID: 51, ResolveUID: 51, Name: "ble_sensor_type", Type: ByteArray, Default: []int{255, 255, 255, 255, 255, 255, 255}, IsDigits: boolPtr(true), Section: "bluetooth"},
	52:  {UID: 52, ResolveUID: 52, Name: "ble_angle_false_min", Type: UInt8, Default: 0, IsDigits: boolPtr(true), Section: ""},
	53:  {UID: 53, ResolveUID: 53, Name: "ble_angle_false_max", Type: UInt8, Default: 0, IsDigits: boolPtr(true), Section: ""},
	54:  {UID: 54, ResolveUID: 54, Name: "ble_angle_true_min", Type: UInt8, Default: 0, IsDigits: boolPtr(true), Section: ""},
	55:  {UID: 55, ResolveUID: 55, Name: "ble_angle_true_max", Type: UInt8, Default: 0, IsDigits: boolPtr(true), Section: ""},
	56:  {UID: 56, ResolveUID: 56, Name: "traffic_sec", Type: UInt8, Default: 6, IsDigits: boolPtr(true), Section: ""},
	57:  {UID: 57, ResolveUID: 57, Name: "ble_output", Type: UInt8, Default: 255, IsDigits: boolPtr(true), Min: intPtr(0), Max: intPtr(255), Section: "bluetooth"},
	58:  {UID: 58, ResolveUID: 58, Name: "ble_broadcast_silent_period_sec", Type: UInt16, Default: 15, IsDigits: boolPtr(true), Min: intPtr(0), Max: intPtr(65535), Section: "bluetooth"},
	59:  {UID: 59, ResolveUID: 59, Name: "ble_broadcast_active_period_sec", Type: UInt8, Default: 5, IsDigits: boolPtr(true), Min: intPtr(0), Max: intPtr(255), Section: "bluetooth"},
	60:  {UID: 60, ResolveUID: 60, Name: "alarming_sms_timeout_period_min", Type: ByteArray, Default: []int{60, 60, 60}, IsDigits: boolPtr(true), Section: ""},
	61:  {UID: 61, ResolveUID: 61, Name: "ble_output_key", Type: Hex, Default: "00010203000102030001020300010203", Section: "bluetooth"},
	62:  {UID: 62, ResolveUID: 62, Name: "ble_output_addr", Type: Hex, Default: "000000000000", Section: "bluetooth"},
	63:  {UID: 63, ResolveUID: 63, Name: "ble_output_security_mode", Type: UInt8, Default: nil, IsDigits: boolPtr(true), Section: ""},
	64:  {UID: 64, ResolveUID: 64, Name: "speedalarm_max_speed_value", Type: UInt16, Default: 0, IsDigits: boolPtr(true), Section: ""},
	65:  {UID: 65, ResolveUID: 65, Name: "beacon", Type: UInt8, Default: 0, IsDigits: boolPtr(true), Section: ""},
	66:  {UID: 66, ResolveUID: 66, Name: "beacon_point_placement_mode", Type: UInt8, Default: 0, IsDigits: boolPtr(true), Section: ""},
	67:  {UID: 67, ResolveUID: 67, Name: "lls_fuel_sensors_addr", Type: ByteArray, Default: nil, IsDigits: boolPtr(true), Section: "rs485"},
	68:  {UID: 68, ResolveUID: 68, Name: "ow_temp_sn", Type: SNSArray, Default: []string{}, Section: "onewire"},
	69:  {UID: 69, ResolveUID: 69, Name: "ibutton_enabled", Type: UInt8, Default: 0, Min: intPtr(0), Max: intPtr(255), IsDigits: boolPtr(true), Section: "onewire"},
	70:  {UID: 70, ResolveUID: 70, Name: "ibutton_output_mode", Type: UInt8, Default: 0, IsDigits: boolPtr(true), Section: ""},
	71:  {UID: 71, ResolveUID: 71, Name: "adm20_addr", Type: UInt8, Default: 255, Min: intPtr(0), Max: intPtr(255), IsDigits: boolPtr(true), Section: "rs485"},
	72:  {UID: 72, ResolveUID: 72, Name: "adm20_mode", Type: UInt8, Default: 0, Min: intPtr(0), Max: intPtr(255), IsDigits: boolPtr(true), Section: "rs485"},
	73:  {UID: 73, ResolveUID: 73, Name: "adm20_output_state_mask", Type: UInt8, Default: 0, Min: intPtr(0), Max: intPtr(255), IsDigits: boolPtr(true), Section: "rs485"},
	74:  {UID: 74, ResolveUID: 74, Name: "adm20_alarm_state", Type: UInt8, Default: 0, Min: intPtr(0), Max: intPtr(255), IsDigits: boolPtr(true), Section: "rs485"},
	75:  {UID: 75, ResolveUID: 75, Name: "adm20_id_mode", Type: UInt8, Default: 0, Min: intPtr(0), Max: intPtr(255), IsDigits: boolPtr(true), Section: "rs485"},
	76:  {UID: 76, ResolveUID: 76, Name: "modbus_devices_addr", Type: ByteArray, Default: []int{255, 255, 255, 255, 255}, Min: intPtr(0), Max: intPtr(255), IsDigits: boolPtr(true), Section: "rs485"},
	77:  {UID: 77, ResolveUID: 77, Name: "modbus_devices_bd_rate", Type: Uint32Array, Default: []int{19200, 19200, 19200, 19200, 19200}, Min: intPtr(0), Max: intPtr(4_294_967_295), IsDigits: boolPtr(true), Section: "rs485"},
	78:  {UID: 78, ResolveUID: 78, Name: "modbus_registers_addr_deprecated", Type: ByteArray, Default: nil, IsDigits: boolPtr(true), Section: ""},
	79:  {UID: 79, ResolveUID: 79, Name: "modbus_bytes_count", Type: ByteArray, Default: []int{2, 2, 2, 2, 2}, Min: intPtr(0), Max: intPtr(255), IsDigits: boolPtr(true), Section: "rs485"},
	80:  {UID: 80, ResolveUID: 80, Name: "modbus_bytes_order", Type: ByteArray, Default: []int{255, 255, 255, 255, 255}, Min: intPtr(0), Max: intPtr(255), IsDigits: boolPtr(true), Section: "rs485"},
	81:  {UID: 81, ResolveUID: 81, Name: "modbus_registers_type", Type: ByteArray, Default: []int{3, 3, 3, 3, 3}, Min: intPtr(0), Max: intPtr(255), IsDigits: boolPtr(true), Section: "rs485"},
	82:  {UID: 82, ResolveUID: 82, Name: "modbus_reset_timeout_min", Type: ByteArray, Default: []int{0, 0, 0, 0, 0}, IsDigits: boolPtr(true), Section: ""},
	83:  {UID: 83, ResolveUID: 83, Name: "canlog", Type: UInt32, Default: 0, IsDigits: boolPtr(true), Section: ""},
	84:  {UID: 84, ResolveUID: 84, Name: "speedalarm_out_alarm_state", Type: UInt8, Default: 1, IsDigits: boolPtr(true), Section: ""},
	85:  {UID: 85, ResolveUID: 85, Name: "intrueout_ain_num", Type: UInt8, Default: 255, IsDigits: boolPtr(true), Section: ""},
	86:  {UID: 86, ResolveUID: 86, Name: "intrueout_out_alarm_state", Type: UInt8, Default: 1, IsDigits: boolPtr(true), Section: ""},
	87:  {UID: 87, ResolveUID: 87, Name: "intrueout_speed_threshold_kmh", Type: UInt8, Default: 255, Min: intPtr(0), Max: intPtr(255), IsDigits: boolPtr(true), Section: ""},
	88:  {UID: 88, ResolveUID: 88, Name: "beacon_whitelist", Type: MacAddressArray, Default: []string{}, IsDigits: boolPtr(true), Section: ""},
	89:  {UID: 89, ResolveUID: 89, Name: "beacon_sending_group_timeout", Type: UInt8, Default: 60, IsDigits: boolPtr(true), Section: ""},
	90:  {UID: 90, ResolveUID: 90, Name: "beacon_sending_group_size", Type: UInt8, Default: 5, IsDigits: boolPtr(true), Section: ""},
	91:  {UID: 91, ResolveUID: 91, Name: "fuel_sensor_address_list", Type: ByteArray, Default: []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, Min: intPtr(0), Max: intPtr(255), IsDigits: boolPtr(true), Section: "rs485;bluetooth"},
	92:  {UID: 92, ResolveUID: 92, Name: "fuel_sensor_type", Type: UInt8, Default: 255, Min: intPtr(0), Max: intPtr(255), IsDigits: boolPtr(true), Section: "rs485;bluetooth"},
	93:  {UID: 93, ResolveUID: 93, Name: "tcp_keep_alive_idle_period_min", Type: UInt8, Default: 0, IsDigits: boolPtr(true), Section: ""},
	94:  {UID: 94, ResolveUID: 94, Name: "ble_scan_timeout_lock_flag", Type: UInt8, Default: 1, IsDigits: boolPtr(true), Section: ""},
	95:  {UID: 95, ResolveUID: 95, Name: "static_program_speed", Type: UInt8, Default: 3, Min: intPtr(0), Max: intPtr(255), IsDigits: boolPtr(true), Section: "track"},
	96:  {UID: 96, ResolveUID: 96, Name: "stat_mask", Type: UInt32, Default: 4294966783, IsDigits: boolPtr(true), Section: "track"},
	97:  {UID: 97, ResolveUID: 97, Name: "device_functions", Type: UInt32, Default: 65537, Min: intPtr(0), Max: intPtr(4_294_967_295), IsDigits: boolPtr(true), Section: "device"},
	98:  {UID: 98, ResolveUID: 98, Name: "beacon_time_sec", Type: Uint32Array, Default: []int{300, 600, 3600, 1200}, Min: intPtr(60), Max: intPtr(2_147_483_647), IsDigits: boolPtr(true), Section: "device"},
	99:  {UID: 99, ResolveUID: 99, Name: "static_power_mode", Type: UInt8, Default: 0, Min: intPtr(0), Max: intPtr(255), IsDigits: boolPtr(true), Section: "device"},
	100: {UID: 100, ResolveUID: 100, Name: "static_power_delay_min", Type: UInt8, Default: 5, Min: intPtr(0), Max: intPtr(255), IsDigits: boolPtr(true), Section: "device"},
	101: {UID: 101, ResolveUID: 101, Name: "fuel_extend_data_mode", Type: UInt8, Default: 0, IsDigits: boolPtr(true), Section: ""},
	103: {UID: 103, ResolveUID: 103, Name: "ed_max_acceleration_kmh_s", Type: UInt8, Default: 255, IsDigits: boolPtr(true), Section: ""},
	104: {UID: 104, ResolveUID: 104, Name: "ed_acc_out_alarm_state", Type: UInt8, Default: 1, IsDigits: boolPtr(true), Section: ""},
	105: {UID: 105, ResolveUID: 105, Name: "ed_max_deceleration_kmh_s", Type: UInt8, Default: 255, IsDigits: boolPtr(true), Section: ""},
	106: {UID: 106, ResolveUID: 106, Name: "serial_uploading", Type: UInt8, Default: 0, IsDigits: boolPtr(true), Section: ""},
	107: {UID: 107, ResolveUID: 107, Name: "simtime_min", Type: UInt16, Default: nil, IsDigits: boolPtr(true), Section: ""},
	108: {UID: 108, ResolveUID: 108, Name: "adm6_data_confirm", Type: TwoUint8, Default: []int{0, 0}, IsDigits: boolPtr(true), Section: ""},
	109: {UID: 109, ResolveUID: 109, Name: "nav_mode", Type: UInt8, Default: 0, IsDigits: boolPtr(true), Section: ""},
	110: {UID: 110, ResolveUID: 110, Name: "ibutton_output_time_sec", Type: UInt8, Default: 5, IsDigits: boolPtr(true), Section: ""},
	112: {UID: 112, ResolveUID: 112, Name: "adm34_tagout", Type: UInt8, Default: 64, IsDigits: boolPtr(true), Section: ""},
	114: {UID: 114, ResolveUID: 114, Name: "nav_filter_sats_count", Type: UInt8, Default: 6, Min: intPtr(0), Max: intPtr(255), IsDigits: boolPtr(true), Section: "track"},
	115: {UID: 115, ResolveUID: 115, Name: "nav_filter_hdop_tenths", Type: UInt8, Default: 20, IsDigits: boolPtr(true), Section: "track"},
	116: {UID: 116, ResolveUID: 116, Name: "nav_filter_hdop_max", Type: UInt8, Default: 25, Min: intPtr(0), Max: intPtr(255), IsDigits: boolPtr(true), Section: "track"},
	117: {UID: 117, ResolveUID: 117, Name: "custom_mask_default_false_1", Type: UInt32, Default: 515, IsDigits: boolPtr(true), Section: "device;events;bluetooth"},
	118: {UID: 118, ResolveUID: 118, Name: "custom_mask_default_false_2", Type: UInt32, Default: 0, IsDigits: boolPtr(true), Section: ""},
	119: {UID: 119, ResolveUID: 119, Name: "custom_mask_default_true_1", Type: UInt32, Default: 4294967293, IsDigits: boolPtr(true), Section: "events"},
	120: {UID: 120, ResolveUID: 120, Name: "custom_mask_default_true_2", Type: UInt32, Default: 4294967295, IsDigits: boolPtr(true), Section: ""},
	121: {UID: 121, ResolveUID: 121, Name: "beacon_missing_timeout_sec", Type: UInt8, Default: 120, IsDigits: boolPtr(true), Section: ""},
	122: {UID: 122, ResolveUID: 122, Name: "sos_config", Type: UInt8, Default: 0, IsDigits: boolPtr(true), Section: "events"},
	123: {UID: 123, ResolveUID: 123, Name: "device_function_3", Type: UInt8, Default: 129, Min: intPtr(0), Max: intPtr(255), IsDigits: boolPtr(true), Section: "sim"},
	124: {UID: 124, ResolveUID: 124, Name: "impulse_inputs", Type: TwoUint8, Default: nil, Min: intPtr(0), Max: intPtr(255), IsDigits: boolPtr(true), Section: "inputs"},
	125: {UID: 125, ResolveUID: 125, Name: "intrueout_out_num", Type: UInt8, Default: nil, IsDigits: boolPtr(true), Section: ""},
	126: {UID: 126, ResolveUID: 126, Name: "ed_acc_out_num", Type: UInt8, Default: nil, IsDigits: boolPtr(true), Section: ""},
	127: {UID: 127, ResolveUID: 127, Name: "speedalarm_out_num_mask", Type: UInt8, Default: nil, IsDigits: boolPtr(true), Section: ""},
	128: {UID: 128, ResolveUID: 128, Name: "adm34_tagout_num", Type: UInt8, Default: nil, IsDigits: boolPtr(true), Section: ""},
	129: {UID: 129, ResolveUID: 129, Name: "fota_update_mode", Type: UInt8, Default: nil, IsDigits: boolPtr(true), Section: ""},
	130: {UID: 130, ResolveUID: 130, Name: "fota_update_sim_index", Type: UInt8, Default: nil, IsDigits: boolPtr(true), Section: ""},
	131: {UID: 131, ResolveUID: 131, Name: "fota_update_url", Type: Ascii, Default: nil, MaxLength: intPtr(128), Section: ""},
	132: {UID: 132, ResolveUID: 132, Name: "net_mode", Type: TwoUint8, Default: []int{0, 0}, Min: intPtr(0), Max: intPtr(255), IsDigits: boolPtr(true), Section: "sim"},
	133: {UID: 133, ResolveUID: 133, Name: "beacon_mode_gsm_gnss", Type: ByteArray, Default: []int{0, 0, 3, 0}, IsDigits: boolPtr(true), Section: "device"},
	134: {UID: 134, ResolveUID: 134, Name: "low_battery_level", Type: UInt16, Default: 3550, Min: intPtr(0), Max: intPtr(65535), IsDigits: boolPtr(true), Section: "events"},
	135: {UID: 135, ResolveUID: 135, Name: "imei", Type: Ascii, Default: nil, Section: ""},
	136: {UID: 136, ResolveUID: 136, Name: "device_mode", Type: UInt8, Default: 0, IsDigits: boolPtr(true), Section: "device;sim;server;track"},
	137: {UID: 137, ResolveUID: 137, Name: "vibro_mode", Type: UInt8, Default: 1, IsDigits: boolPtr(true), Section: ""},
	138: {UID: 138, ResolveUID: 138, Name: "bleoutain", Type: UInt8, Default: 0, IsDigits: boolPtr(true), Section: ""},
	139: {UID: 139, ResolveUID: 139, Name: "no_nav_reset_time_min", Type: UInt8, Default: 10, IsDigits: boolPtr(true), Section: ""},
	140: {UID: 140, ResolveUID: 140, Name: "bt_ble_pwr_was_set_flags_mask", Type: UInt8, Default: nil, IsDigits: boolPtr(true), Section: ""},
	141: {UID: 141, ResolveUID: 141, Name: "nav_filter_invalid_height_decameters", Type: UInt16, Default: nil, IsDigits: boolPtr(true), Section: ""},
	142: {UID: 142, ResolveUID: 142, Name: "battery_safe_mode", Type: UInt8, Default: 4, IsDigits: boolPtr(true), Section: "device"},
	143: {UID: 143, ResolveUID: 143, Name: "jamout_status_list", Type: ByteArray, Default: []int{0, 255}, IsDigits: boolPtr(true), Section: ""},
	144: {UID: 144, ResolveUID: 144, Name: "modbus_registers_addr", Type: Uint32Array, Default: []int{0, 0}, Min: intPtr(0), Max: intPtr(4_294_967_295), IsDigits: boolPtr(true), Section: "rs485"},
	145: {UID: 145, ResolveUID: 145, Name: "tmode_value", Type: UInt8, Default: 3, IsDigits: boolPtr(true), Section: ""},
	146: {UID: 146, ResolveUID: 146, Name: "ble_beacon_time_extra_point", Type: UInt8, Default: 5, IsDigits: boolPtr(true), Section: ""},
	147: {UID: 147, ResolveUID: 147, Name: "bleout_master_id", Type: ByteArray, Default: []int{0, 0, 0, 0, 0, 0}, IsDigits: boolPtr(true), Section: ""},
	148: {UID: 148, ResolveUID: 148, Name: "vibstart_time", Type: ByteArray, Default: nil, IsDigits: boolPtr(true), Section: ""},
	149: {UID: 149, ResolveUID: 149, Name: "gsm_led_brightness", Type: UInt8, Default: 50, Section: ""},
	150: {UID: 150, ResolveUID: 150, Name: "gps_led_brightness", Type: UInt8, Default: 50, IsDigits: boolPtr(true), Section: ""},
	151: {UID: 151, ResolveUID: 151, Name: "pwr_led_brightness", Type: UInt8, Default: 50, IsDigits: boolPtr(true), Section: ""},
	152: {UID: 152, ResolveUID: 152, Name: "beacon_middle_timing", Type: Uint32Array, Default: []int{0, 60}, IsDigits: boolPtr(true), Section: ""},
	153: {UID: 153, ResolveUID: 153, Name: "temperature_calib_offset", Type: UInt8, Default: 0, Min: intPtr(0), Max: intPtr(255), IsDigits: boolPtr(true), Section: "rs485"},
	154: {UID: 154, ResolveUID: 154, Name: "sim_pin", Type: Ascii, Default: []string{"", ""}, Section: "sim"},
	155: {UID: 155, ResolveUID: 155, Name: "apn_name", Type: Ascii, Default: []string{"m2m.beeline.ru", "m2m.beeline.ru"}, Section: "sim"},
	156: {UID: 156, ResolveUID: 156, Name: "apn_user", Type: Ascii, Default: []string{"beeline", "beeline"}, Section: "sim"},
	157: {UID: 157, ResolveUID: 157, Name: "apn_pass", Type: Ascii, Default: []string{"beeline", "beeline"}, Section: "sim"},
	158: {UID: 158, ResolveUID: 158, Name: "ops_white_list", Type: Uint32Array, Default: []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, IsDigits: boolPtr(true), Min: intPtr(0), Max: intPtr(1_000_000_000), Section: "sim"},
	159: {UID: 159, ResolveUID: 159, Name: "ops_black_list", Type: Uint32Array, Default: []int{0, 0, 0, 0}, IsDigits: boolPtr(true), Min: intPtr(0), Max: intPtr(1_000_000_000), Section: "sim"},
	160: {UID: 160, ResolveUID: 160, Name: "allowed_iccid", Type: ByteArray, Default: []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, IsDigits: boolPtr(true), Section: ""},
	161: {UID: 161, ResolveUID: 161, Name: "simtime_min_conf", Type: TwoUint16, Default: []int{30, 30}, Min: intPtr(0), Max: intPtr(65535), IsDigits: boolPtr(true), Section: "sim"},
	162: {UID: 162, ResolveUID: 162, Name: "geofence_points", Type: Ascii, Default: nil, IsDigits: boolPtr(true), Section: ""},
	163: {UID: 163, ResolveUID: 163, Name: "geofence_sms_period_min", Type: UInt16, Default: nil, IsDigits: boolPtr(true), Section: ""},
	165: {UID: 165, ResolveUID: 165, Name: "navtimesync", Type: UInt8, Default: 0, IsDigits: boolPtr(true), Section: "track"},
	166: {UID: 166, ResolveUID: 166, Name: "navtimesync_max_time_diff_sec", Type: UInt16, Default: 1800, IsDigits: boolPtr(true), Section: "track"},
	167: {UID: 167, ResolveUID: 167, Name: "nav_filter_valid_height_deca_hecto_meters", Default: 0, Type: TwoUint8, IsDigits: boolPtr(true), Section: "track"},
	168: {UID: 168, ResolveUID: 168, Name: "beacontype", Type: UInt8, Default: 0, IsDigits: boolPtr(true), Section: ""},
	170: {UID: 170, ResolveUID: 170, Name: "adm40_addr", Type: Hex, Default: "000000000000", Section: "track"},
	249: {UID: 249, ResolveUID: 249, Name: "remote_server_host_address", Type: Ascii, Default: "www.neosync.neomatica.ru", Section: ""},
	250: {UID: 250, ResolveUID: 250, Name: "remote_server_host_port", Type: UInt16, Default: 12346, IsDigits: boolPtr(true), Section: ""},
	251: {UID: 251, ResolveUID: 251, Name: "remote_server_sync_period_min", Type: UInt8, Default: 5, IsDigits: boolPtr(true), Section: ""},
	252: {UID: 252, ResolveUID: 252, Name: "majorfilter", Type: UInt16, Default: 0, IsDigits: boolPtr(true), Section: ""},
	253: {UID: 253, ResolveUID: 253, Name: "alternative_sources_positioning_query_period", Type: UInt16, Default: 120, IsDigits: boolPtr(true), Min: intPtr(10), Max: intPtr(65535), Section: "track"},
	254: {UID: 254, ResolveUID: 254, Name: "alternative_sources_positioning_token", Type: Ascii, Default: "0", MaxLength: intPtr(41), Section: "track"},
	255: {UID: 255, ResolveUID: 255, Name: "coordinate_substitution_alternative", Type: UInt8, Default: 0, IsDigits: boolPtr(true), Min: intPtr(0), Max: intPtr(8), Section: "track"},
	256: {UID: 256, ResolveUID: 256, Name: "multiplier_coordinate_validity_area", Type: UInt16, Default: 1, IsDigits: boolPtr(true), Min: intPtr(1), Max: intPtr(65535), Section: "track"},
	257: {UID: 257, ResolveUID: 257, Name: "discarded_points_during_ejection", Type: UInt8, Default: 0, IsDigits: boolPtr(true), Min: intPtr(0), Max: intPtr(50), Section: "track"},
	259: {UID: 259, ResolveUID: 259, Name: "distance_between_points", Type: UInt16, Default: 1000, IsDigits: boolPtr(true), Min: intPtr(10), Max: intPtr(65535), Section: "track"},
	266: {UID: 266, ResolveUID: 266, Name: "alternative_sources_positioning_operating_mode", Type: UInt8, Default: 0, IsDigits: boolPtr(true), Section: "track"},
	268: {UID: 268, ResolveUID: 268, Name: "alternative_sources_positioning_min_time_requests", Type: UInt8, Default: 10, IsDigits: boolPtr(true), Section: "track"},
	269: {UID: 269, ResolveUID: 269, Name: "alternative_sources_positioning_accuracy_coordinates", Type: UInt16, Default: 100, IsDigits: boolPtr(true), Min: intPtr(10), Max: intPtr(1000), Section: "track"},
	271: {UID: 271, ResolveUID: 271, Name: "cooling_period_after_transition", Type: UInt8, Default: 0, IsDigits: boolPtr(true), Section: "track"},
	275: {UID: 275, ResolveUID: 275, Name: "sos_mode", Type: UInt8, Default: 0, IsDigits: boolPtr(true), Section: "events"},
}

func intPtr(i int) *int {
	return &i
}

func boolPtr(i bool) *bool {
	return &i
}

func GetCfgUid(uidName string) (uint16, error) {
	for uid, schema := range CfgSchema {
		if schema.Name == uidName {
			return uid, nil
		}
	}

	return 0, fmt.Errorf("unknown field name: %s", uidName)
}

func GetCfgName(uid uint16) (string, error) {
	s, ok := CfgSchema[uid]
	if !ok {
		return "", fmt.Errorf("unknown uid: %d", uid)
	}

	return s.Name, nil
}

func GetCfgType(uid uint16) (FieldType, error) {
	s, ok := CfgSchema[uid]
	if !ok {
		return "", fmt.Errorf("unknown uid: %d", uid)
	}

	return s.Type, nil
}

func GetCfgSchema(uid uint16) (FieldSchema, error) {
	s, ok := CfgSchema[uid]
	if !ok {
		return FieldSchema{}, fmt.Errorf("unknown uid: %d", uid)
	}

	return s, nil
}
