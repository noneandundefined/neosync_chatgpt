package util

import (
	"net"
)

func ReadAtLast(conn net.Conn, buffer []byte, min int) (int, error) {
	total := 0
	for total < min {
		n, err := conn.Read(buffer[total:])
		if err != nil {
			return total, err
		}

		total += n
	}

	return total, nil
}

func RemoteIP(conn net.Conn) string {
	if conn == nil {
		return ""
	}

	addr := conn.RemoteAddr()
	if addr == nil {
		return ""
	}

	if tcpAddr, ok := addr.(*net.TCPAddr); ok {
		return tcpAddr.IP.String()
	}

	return addr.String()
}

func ImeiString(imei *string) string {
	if imei == nil {
		return ""
	}

	return *imei
}
