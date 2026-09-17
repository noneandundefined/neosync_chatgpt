package main

import (
	"neomatica/neosync-tcp/infra/logger"
	"net"
)

func (tcp *tcpServer) handleConnectionCallback(conn net.Conn) {
	defer conn.Close()

	if err := tcp.Dispatch(conn); err != nil {
		logger.Error("HandleConnectionCallback ip={%s}: %s", conn.RemoteAddr().(*net.TCPAddr).IP, err.Error())
	}
}
