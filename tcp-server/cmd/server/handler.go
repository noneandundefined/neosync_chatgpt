package main

import (
	"errors"
	"fmt"
	"io"
	"neomatica/neosync-tcp/infra/constants"
	"neomatica/neosync-tcp/infra/event"
	"neomatica/neosync-tcp/infra/logger"
	"neomatica/neosync-tcp/infra/store/redis"
	"neomatica/neosync-tcp/pkg/protocol"
	"neomatica/neosync-tcp/pkg/protocol/adm"
	"neomatica/neosync-tcp/pkg/protocol/adm_v2"
	"neomatica/neosync-tcp/util"
	"net"
	"os"
)

func (tcp *tcpServer) HandleAdm(conn net.Conn, buffer []byte) error {
	var protocol protocol.Protocol = &adm_v2.ADM_V2{Store: tcp.store, Db: tcp.db}
	return tcp.HandleAdmConn(conn, buffer, protocol)
}

func (tcp *tcpServer) HandleAdmConn(conn net.Conn, buffer []byte, protocol protocol.Protocol) error {
	emitter := event.NewEventEmitter()

	device := adm.NewADMDevice(conn, emitter, protocol)

	defer func() {
		var imei string

		if device.Imei != nil {
			imei = util.PrepareImei(*device.Imei)
		}

		if imei != "" {
			_ = redis.WriteAdmLog(constants.EVENT_DEVICE_DISCONN, imei, nil)

			tcp.session.RemoveDeviceSessionForDevice(imei, device)
			tcp.cache.Delete(fmt.Sprintf("sync:%s", imei))
		}

		_ = conn.Close()
	}()

	// set keep-alive tcp conn
	// -----------------------
	if tcpConn, ok := conn.(*net.TCPConn); ok {
		if err := util.SetQuickTCPDetect(tcpConn); err != nil {
			logger.Warning("HandleAdmConn ip={%s}: Warning set keep-alive: %s", conn.RemoteAddr().(*net.TCPAddr).IP, err.Error())
		}
	}
	// -----------------------
	// set keep-alive tcp conn

	/* Listen events */
	tcp.ListenEvents(device)

	if len(buffer) > 0 {
		device.ParseData(buffer)
	}

	/* Parser bytes */
	for {
		buf := bufPool.Get().(*[]byte)

		n, err := conn.Read(*buf)
		if err != nil {
			bufPool.Put(buf)

			ip := conn.RemoteAddr().(*net.TCPAddr).IP

			if !errors.Is(err, io.EOF) && !readDisconnect(err) && !os.IsTimeout(err) {
				logger.Error("HandleAdmConn ip={%s}: Error read n bytes: %s", ip, err.Error())
			}

			break
		}

		data := (*buf)[:n]

		/* Parse ADM bytes */
		device.ParseData(data)

		/* Clear buffer */
		bufPool.Put(buf)
	}

	return nil
}

func readDisconnect(err error) bool {
	var opErr *net.OpError

	if !errors.As(err, &opErr) || opErr.Op != "read" {
		return false
	}

	return errors.Is(opErr.Err, net.ErrClosed) || opErr.Err.Error() == "use of closed network connection"
}
