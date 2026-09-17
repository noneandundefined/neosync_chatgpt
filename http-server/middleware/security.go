package middleware

import (
	"io"
	"neomatica/neosync/infra/logger"
	"neomatica/neosync/pkg/clientip"
	"neomatica/neosync/pkg/httpx"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

const (
	MaxBodySize        = 1 << 20
	MaxViolationsPerIP = 3
	IpBlockDuration    = 10 * time.Minute
)

var (
	BlockedIPs   = make(map[string]time.Time)
	IpViolations = make(map[string]int)
	Mutex        sync.Mutex
)

/* Проверока данных от пользователей (XSS, SQLInjection, ...) */
func SecurityMiddleware() mux.MiddlewareFunc { //nolint
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tr := TranslatorFromContext(r.Context())
			ip := getIP(r)

			// logger.Ip(ip)

			if isBlocked(ip) {
				logger.Warning("SecurityMiddleware ip={%s}: Ip has been blocked", ip)
				httpx.HttpResponse(w, r, http.StatusForbidden, tr.TErr("access-denied"))
				return
			}

			/* Limit size */
			r.Body = http.MaxBytesReader(w, r.Body, MaxBodySize)

			var bodyContent string
			if r.Method == http.MethodPost || r.Method == http.MethodPut {
				data, err := io.ReadAll(r.Body)
				if err != nil {
					httpx.HttpResponse(w, r, http.StatusRequestEntityTooLarge, tr.TErr("request-text-toobig"))
					registerViolation(ip)
					return
				}

				bodyContent = string(data)
				r.Body = io.NopCloser(strings.NewReader(bodyContent))
			}

			next.ServeHTTP(w, r)
		})
	}
}

func getIP(r *http.Request) string { return clientip.IP(r) }

func isBlocked(ip string) bool {
	Mutex.Lock()
	defer Mutex.Unlock()

	expiry, exists := BlockedIPs[ip]
	if !exists {
		return false
	}

	if time.Now().After(expiry) {
		delete(BlockedIPs, ip)
		delete(IpViolations, ip)

		return false
	}

	return true
}

func registerViolation(ip string) {
	Mutex.Lock()
	defer Mutex.Unlock()

	IpViolations[ip]++
	if IpViolations[ip] >= MaxViolationsPerIP {
		BlockedIPs[ip] = time.Now().Add(IpBlockDuration)
	}
}
