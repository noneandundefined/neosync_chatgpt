package main

import (
	"fmt"
	"neomatica/neosync-tcp/config"
	"neomatica/neosync-tcp/infra/logger"
	"net"
	"os"
	"runtime"
	"time"
)

const maxConnections = 10000

func (tcp *tcpServer) tcpStart() error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", config.DeviceRemotePort))
	if err != nil {
		return fmt.Errorf("Error starting server: %s", err.Error())
	}
	defer listener.Close()

	fmt.Printf("\n[%v] [INFO] Tcp server started :%d\n", time.Now().Format("2006-01-02 15:04:05"), config.DeviceRemotePort)
	fmt.Printf("[%v] [INFO] Proccess PID: %d, Version: %s\n", time.Now().Format("2006-01-02 15:04:05"), os.Getpid(), Version)
	fmt.Printf("[%v] [INFO] Golang version: %s\n\n", time.Now().Format("2006-01-02 15:04:05"), runtime.Version())

	semaphore := make(chan struct{}, maxConnections)

	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Error("TcpStart: Error accepting connection: %s", err.Error())
			continue
		}

		select {
		case semaphore <- struct{}{}:
			go func() {
				defer func() { <-semaphore }()
				tcp.handleConnectionCallback(conn)
			}()

		default:
			logger.Warning("TcpStart ip={%s}: Too many active connections total={%d}", conn.RemoteAddr().(*net.TCPAddr).IP, len(semaphore))
			_ = conn.Close()
		}
	}
}
