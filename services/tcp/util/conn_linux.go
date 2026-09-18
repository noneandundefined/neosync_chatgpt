//go:build linux
// +build linux

package util

import (
	"fmt"
	"net"
	"time"

	"golang.org/x/sys/unix"
)

func SetQuickTCPDetect(tcpConn *net.TCPConn) error {
	if err := tcpConn.SetKeepAlive(true); err != nil {
		return fmt.Errorf("error set keep-alive: %w", err)
	}

	if err := tcpConn.SetKeepAlivePeriod(10 * time.Second); err != nil {
		return fmt.Errorf("error set keep-alive period: %w", err)
	}

	raw, err := tcpConn.SyscallConn()
	if err != nil {
		return err
	}

	var sockErr error
	raw.Control(func(fd uintptr) {
		const clientTimeoutMs = 5000
		sockErr = unix.SetsockoptInt(int(fd), unix.IPPROTO_TCP, unix.TCP_USER_TIMEOUT, clientTimeoutMs)
		if sockErr != nil {
			return
		}
	})

	return sockErr
}
