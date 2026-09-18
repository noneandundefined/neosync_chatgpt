//go:build !linux
// +build !linux

package util

import "net"

func SetQuickTCPDetect(_ *net.TCPConn) error {
	return nil
}
