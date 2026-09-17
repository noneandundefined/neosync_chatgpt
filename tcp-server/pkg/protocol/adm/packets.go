package adm

import (
	"neomatica/neosync-tcp/infra/constants"
	"neomatica/neosync-tcp/infra/logger"
	"neomatica/neosync-tcp/infra/tcperrors"
	"net"
)

func (adm *ADMDevice) parseHelloPacket(buffer []byte) {
	packet, err := adm.protocol.Parse_WelcomePacket(buffer)
	if err != nil {
		logger.Error("ParseHelloPacket ip={%s}: Failed parse welcome packet: %s", adm.Conn.RemoteAddr().(*net.TCPAddr).IP, err.Error())
		adm.Emitter.Emit("error", err)
		return
	}

	if packet != nil && packet.Imei != nil {
		adm.Imei = packet.Imei
	}

	adm.Emitter.Emit("hello-pack-received", packet)
}

func (adm *ADMDevice) parseCommandPacket(buffer []byte) {
	answerBuffer := adm.protocol.Parse_CommandPacket(buffer, uint(len(buffer)))

	if n := len(answerBuffer); n > 0 && answerBuffer[n-1] == 0 {
		answerBuffer = answerBuffer[:n-1]
	}

	adm.Emitter.Emit("command-pack-received", answerBuffer)
}

func (adm *ADMDevice) parseSyncPacket(buffer []byte) {
	if len(buffer) != constants.ADM_RC_SYNC_PACK_SIZE_V2 {
		logger.Error("ParseSyncPacket ip={%s} imei={%s}: Incorrect package size: %d", adm.Conn.RemoteAddr().(*net.TCPAddr).IP, *adm.Imei, len(buffer))
		adm.Emitter.Emit("error", tcperrors.Err_IncorrectPackageSize)
		return
	}

	packet, err := adm.protocol.Parse_SyncPacket(buffer)
	if err != nil {
		logger.Error("ParseSyncPacket ip={%s} imei={%s}: Failed to parse sync packet: %s", adm.Conn.RemoteAddr().(*net.TCPAddr).IP, *adm.Imei, err.Error())
		adm.Emitter.Emit("error", err)
		return
	}

	adm.Emitter.Emit("sync-pack-received", packet)
}
