package middleware

import (
	"neomatica/neosync/infra/constants"
	"neomatica/neosync/pkg/clientip"
	"neomatica/neosync/pkg/httpx"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/time/rate"
)

type clientLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	limiters = make(map[string]*clientLimiter)
	mu       sync.Mutex
)

func getLimiter(ip string, rps float64, burst int) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	if cl, exists := limiters[ip]; exists {
		cl.lastSeen = time.Now()
		return cl.limiter
	}

	limiter := rate.NewLimiter(rate.Limit(rps), burst)
	limiters[ip] = &clientLimiter{
		limiter:  limiter,
		lastSeen: time.Now(),
	}

	return limiter
}

/* Ограничивает кол-во запросов от клиентов */
func RateLimiterMiddleware(rps float64, burst int) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tr := TranslatorFromContext(r.Context())
			ip := clientip.IP(r)

			limiter := getLimiter(ip, rps, burst)

			if !limiter.Allow() {
				httpx.HttpResponse(w, r, http.StatusTooManyRequests, tr.TErr("rate-limiter-exceeded"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func RateLimiterCleanup() {
	go func() {
		for {
			time.Sleep(constants.RateLimiter_StateTTL)

			mu.Lock()
			for ip, cl := range limiters {
				if time.Since(cl.lastSeen) > constants.RateLimiter_StateTTL {
					delete(limiters, ip)
				}
			}
			mu.Unlock()
		}
	}()
}
