package main

import (
	"fmt"
	"neomatica/neosync-tcp/infra/constants"
	"neomatica/neosync-tcp/infra/logger"
	"neomatica/neosync-tcp/pkg/protocol/common"
	"neomatica/neosync-tcp/util"
	"net"
	"sync"
)

var bufPool = sync.Pool{
	New: func() any {
		b := make([]byte, constants.ADM_RC_MAX_PACKET_SIZE)
		return &b
	},
}

func (tcp *tcpServer) Dispatch(conn net.Conn) error {
	buf := bufPool.Get().(*[]byte)
	buffer := *buf

	defer bufPool.Put(buf)

	n, err := util.ReadAtLast(conn, buffer, 4)
	if err != nil {
		return fmt.Errorf("error reading data: %s", err.Error())
	}

	version, err := common.GetVersion_ByBinPacket(buffer[:n])
	if err != nil {
		return fmt.Errorf("error reading data: %s", err.Error())
	}

	if version != 0x02 {
		logger.Info("RX rejected-version ip={%s} bytes={%d} version={%d}", conn.RemoteAddr().(*net.TCPAddr).IP, len(buffer[:n]), version)
		return fmt.Errorf("server get version={%d}, reject connection", version)
	}

	return tcp.HandleAdm(conn, buffer[:n])
}
