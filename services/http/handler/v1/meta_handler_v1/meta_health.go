package meta_handler_v1

import (
	"fmt"
	"neomatica/neosync/pkg/httpx"
	"net"
	"net/http"
	"os"
	"time"
)

func isTcpAlive() bool {
	tcp_addr := os.Getenv("REDIS_HOST")

	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:12346", tcp_addr), time.Second)
	if err != nil {
		return false
	}
	conn.Close()

	return true
}

func (h *Handler) MetaHealthTcpHandler_V1(w http.ResponseWriter, r *http.Request) error {
	tcpHealth := isTcpAlive()

	if !tcpHealth {
		httpx.HttpResponse(w, r, http.StatusAccepted, map[string]bool{"status": false})
		return nil
	}

	httpx.HttpResponse(w, r, http.StatusAccepted, map[string]bool{"status": true})
	return nil
}
