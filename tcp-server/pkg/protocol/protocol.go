package protocol

type WelcomePacket struct {
	Imei            *string
	Pass            *string
	Version         uint8
	FirmwareVersion uint16
	CfgVersion      uint8
	LastModTime     uint32
	CfgHash         uint32
}

type SyncPacket struct {
	FirmwareVersion uint16
	CfgVersion      uint8
	LastModTime     uint32
	CfgHash         uint32
}

type ConfigurationPacket struct {
	CfgHash uint32
}

type Protocol interface {
	/* Parser */
	Parse_WelcomePacket(buffer []byte) (*WelcomePacket, error)
	Parse_CommandPacket(buffer []byte, receiveBufDataSize uint) string
	Parse_SyncPacket(buffer []byte) (*SyncPacket, error)
	Parse_ConfigurationPacket(buffer []byte) (*ConfigurationPacket, error)
	Parse_KeepAlivePacket(buffer []byte, imei *string)

	/* Encoded */
	Encode_CommandPacket(buffer []byte) []byte
	Encode_GetSyncPacket() []byte
	Encode_GetCfgPacket() []byte
	Encode_SetCfgPacket(buffer []byte, id uint64, cfgVersion uint8, firmwareVersion uint16) []byte
}
