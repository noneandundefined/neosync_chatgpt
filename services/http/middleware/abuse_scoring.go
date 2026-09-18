package middleware

import (
	"bufio"
	"fmt"
	"neomatica/neosync/pkg/clientip"
	"neomatica/neosync/pkg/httpx"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

type abuseState struct {
	score           float64
	lastSeen        time.Time
	blockedUntil    time.Time
	recentEndpoints map[string]time.Time
}

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rw *statusRecorder) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *statusRecorder) Flush() {
	if f, ok := rw.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

func (rw *statusRecorder) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if h, ok := rw.ResponseWriter.(http.Hijacker); ok {
		return h.Hijack()
	}
	return nil, nil, fmt.Errorf("hijacker unsupported")
}

var (
	abuseMu    sync.Mutex
	abuseByIP  = make(map[string]*abuseState)
	abuseTTL   = 5 * time.Minute
	scoreDecay = 0.25
)

func getAbuseState(ip string, now time.Time) *abuseState {
	st, ok := abuseByIP[ip]
	if !ok {
		st = &abuseState{
			score:           0,
			lastSeen:        now,
			recentEndpoints: map[string]time.Time{},
		}
		abuseByIP[ip] = st

		return st
	}

	elapsedSec := now.Sub(st.lastSeen).Seconds()
	st.score -= elapsedSec * scoreDecay
	if st.score < 0 {
		st.score = 0
	}
	st.lastSeen = now

	for ep, ts := range st.recentEndpoints {
		if now.Sub(ts) > 10*time.Second {
			delete(st.recentEndpoints, ep)
		}
	}

	return st
}

func addPenalty(st *abuseState, score float64, now time.Time) {
	st.score += score

	if st.score >= 12 {
		st.blockedUntil = now.Add(30 * time.Second)
		return
	}

	if st.score >= 8 {
		st.blockedUntil = now.Add(10 * time.Second)
	}
}

func AbuseScoringMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/api/v1/devices?page") {
				next.ServeHTTP(w, r)
				return
			}

			now := time.Now()
			ip := clientip.IP(r)
			endpointKey := r.Method + ":" + r.URL.Path

			abuseMu.Lock()
			for key, state := range abuseByIP {
				if now.Sub(state.lastSeen) > abuseTTL {
					delete(abuseByIP, key)
				}
			}

			state := getAbuseState(ip, now)
			if now.Before(state.blockedUntil) {
				retryAfter := int(state.blockedUntil.Sub(now).Seconds())
				abuseMu.Unlock()

				w.Header().Set("Retry-After", "1")
				if retryAfter > 1 {
					w.Header().Set("Retry-After", "2")
				}

				tr := TranslatorFromContext(r.Context())

				httpx.HttpResponse(w, r, http.StatusTooManyRequests, tr.TErr("rate-limiter-exceeded"))
				return
			}

			state.recentEndpoints[endpointKey] = now
			if len(state.recentEndpoints) >= 7 {
				addPenalty(state, 1.5, now)
			}
			abuseMu.Unlock()

			recorder := &statusRecorder{ResponseWriter: w, statusCode: http.StatusOK}
			next.ServeHTTP(recorder, r)

			abuseMu.Lock()
			state = getAbuseState(ip, time.Now())
			switch {
			case recorder.statusCode == http.StatusTooManyRequests:
				addPenalty(state, 4, time.Now())
			case recorder.statusCode >= 500:
				addPenalty(state, 3, time.Now())
			case recorder.statusCode == http.StatusBadRequest || recorder.statusCode == http.StatusUnprocessableEntity:
				addPenalty(state, 1, time.Now())
			}
			abuseMu.Unlock()
		})
	}
}
