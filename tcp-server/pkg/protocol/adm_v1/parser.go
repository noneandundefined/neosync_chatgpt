package adm_v1

import (
	"neomatica/neosync-tcp/infra/constants"
	"neomatica/neosync-tcp/infra/tcperrors"
	"neomatica/neosync-tcp/pkg/protocol"
	"neomatica/neosync-tcp/pkg/protocol/common"
)

func (adm *ADM_V1) Parse_WelcomePacket(buffer []byte) (*protocol.WelcomePacket, error) {
	if len(buffer) < constants.ADM_RC_HELLO_PACK_SIZE_V1 {
		return nil, tcperrors.Err_FirmwareOutdatedNotFound
	}

	if buffer[4] != 0x03 {
		return nil, tcperrors.Err_FirmwareOutdatedNotFound
	}

	welcomePacket := &protocol.WelcomePacket{}

	/* Base value v1 */
	welcomePacket.FirmwareVersion = 0x00
	welcomePacket.CfgVersion = 0x00
	welcomePacket.LastModTime = 0x00
	welcomePacket.CfgHash = 0x00

	welcomePacket.Version = buffer[5]

	id1Size := buffer[15]
	if int(16+id1Size) > len(buffer) {
		return nil, tcperrors.Err_IncorrectPackageSize
	}

	imeiBSD := common.DecodeBCD(buffer[7 : 7+8])
	welcomePacket.Imei = &imeiBSD

	pass := string(buffer[16 : 16+id1Size])
	welcomePacket.Pass = &pass

	return welcomePacket, nil
}

func (adm *ADM_V1) Parse_CommandPacket(buffer []byte, receiveBufDataSize uint) string {
	answerBuffer := buffer[5 : len(buffer)-1]
	return string(answerBuffer)
}

func (adm *ADM_V1) Parse_SyncPacket(buffer []byte) (*protocol.SyncPacket, error) {
	return nil, nil
}

func (adm *ADM_V1) Parse_ConfigurationPacket(buffer []byte) (*protocol.ConfigurationPacket, error) {
	return nil, nil
}

func (adm *ADM_V1) Parse_KeepAlivePacket(buffer []byte, imei *string) {
	/* Not parsed */
}
