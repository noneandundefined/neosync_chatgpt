package adm

import (
	"neomatica/neosync-tcp/infra/constants"
	"neomatica/neosync-tcp/infra/event"
	"neomatica/neosync-tcp/infra/logger"
	"neomatica/neosync-tcp/infra/tcperrors"
	"neomatica/neosync-tcp/pkg/protocol"
	"neomatica/neosync-tcp/pkg/protocol/adm_v2"
	"net"
	"sync/atomic"
	"time"
)

type ADMDevice struct {
	Conn            net.Conn
	DeviceId        uint64
	UserUuid        *string
	Imei            *string
	Pass            *string
	FirmwareVersion uint16
	CfgVersion      uint8
	LastModTime     uint32
	CfgHash         uint32
	RCOnC           bool
	Emitter         *event.EventEmitter
	protocol        protocol.Protocol
	assembler       Assembler
	terminalBusy    atomic.Bool
}

func NewADMDevice(conn net.Conn, emitter *event.EventEmitter, protocol protocol.Protocol) *ADMDevice {
	return &ADMDevice{
		Conn:            conn,
		DeviceId:        0,
		UserUuid:        nil,
		Imei:            nil,
		Pass:            nil,
		FirmwareVersion: 0x00,
		CfgVersion:      0x00,
		LastModTime:     0x00,
		CfgHash:         0x00,
		RCOnC:           false,
		Emitter:         emitter,
		protocol:        protocol,
	}
}

func (adm *ADMDevice) AcquireTerminal() {
	for !adm.terminalBusy.CompareAndSwap(false, true) {
		time.Sleep(2 * time.Millisecond)
	}
}

func (adm *ADMDevice) ReleaseTerminal() {
	adm.terminalBusy.Store(false)
}

func (adm *ADMDevice) TerminalBusy() bool {
	return adm.terminalBusy.Load()
}

func (adm *ADMDevice) imeiForLog() string {
	if adm.Imei == nil {
		return "unknown"
	}

	return *adm.Imei
}

func (adm *ADMDevice) ParseData(buffer []byte) {
	if len(buffer) == 0 {
		adm.Emitter.Emit("error", tcperrors.Err_PacketIsNull)
		return
	}

	/* Keep-alive */
	if len(buffer) == 1 && buffer[0] == constants.ADM_RC_KEEP_ALIVE_VALUE {
		adm.Emitter.Emit("keepalive-pack-received", "")
		return
	}

	packets := adm.assembler.Feed(buffer)

	for _, pkt := range packets {
		logIncomingPacket(adm.Conn.RemoteAddr().(*net.TCPAddr).IP.String(), adm.imeiForLog(), pkt)
		adm.ParsePacket(pkt)
	}
}

func (adm *ADMDevice) MarkV2DeviceModelLookupSatisfied() {
	if v2, ok := adm.protocol.(*adm_v2.ADM_V2); ok {
		v2.MarkDeviceModelLookupSatisfied()
	}
}

func (adm *ADMDevice) OnDisconnect() {
	if adm.Conn != nil {
		if err := adm.Conn.Close(); err != nil {
			if adm.Imei != nil {
				logger.Error("OnDisconnect ip={%s} imei={%s}: Failed to close connection: %s", adm.Conn.RemoteAddr().(*net.TCPAddr).IP, *adm.Imei, err.Error())
			} else {
				logger.Error("OnDisconnect ip={%s}: Failed to close connection: %s", adm.Conn.RemoteAddr().(*net.TCPAddr).IP, err.Error())
			}

			return
		}

		adm.Conn = nil
	}

	adm.DeviceId = 0x00
	adm.UserUuid, adm.Imei, adm.Pass = nil, nil, nil
	adm.FirmwareVersion = 0x00
	adm.CfgVersion, adm.LastModTime = 0x00, 0x00
	adm.CfgHash = 0x00
	adm.assembler.buffer = adm.assembler.buffer[:0]
	adm.RCOnC = false
	adm.terminalBusy.Store(false)
}
