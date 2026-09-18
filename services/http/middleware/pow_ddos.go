package middleware

import (
	"neomatica/neosync/pkg/clientip"
	"neomatica/neosync/pkg/httpx"
	"neomatica/neosync/pkg/pow"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

const defaultDifficulty = 3

func PowDDos() mux.MiddlewareFunc {
	guard, err := pow.New()
	if err != nil {
		panic("cannot initialize proof of work")
	}
	difficulty := defaultDifficulty
	if value := os.Getenv("POW_DDOS_DIFFICULTY"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed < 1 || parsed > pow.MaxDifficulty {
			panic("invalid POW_DDOS_DIFFICULTY (expected 1..6)")
		}
		difficulty = parsed
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			binding := r.Method + "|" + r.URL.Path + "|" + clientip.IP(r)
			if !guard.Verify(r.Header.Get("Pow-Challenge"), r.Header.Get("Pow-Nonce"), binding, difficulty) {
				challenge, err := guard.Challenge(binding, difficulty)
				if err != nil {
					httpx.HttpResponse(w, r, http.StatusServiceUnavailable, "proof of work unavailable")
					return
				}
				w.Header().Set("Cache-Control", "no-store")
				httpx.HttpResponse(w, r, http.StatusTooManyRequests, map[string]string{"challenge": challenge, "difficulty": strconv.Itoa(difficulty)})
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
