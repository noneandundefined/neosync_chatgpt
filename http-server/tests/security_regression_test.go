package tests

import (
	"context"
	"fmt"
	"neomatica/neosync/middleware"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDebugHeaderDoesNotBypassPoW(t *testing.T) {
	t.Setenv("POW_DDOS_DIFFICULTY", "1")
	handler := middleware.PowDDos()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { t.Fatal("PoW bypassed") }))
	req := httptest.NewRequest("POST", "/signin", nil)
	req.Header.Set("X-Debug", "1")
	req = req.WithContext(context.WithValue(req.Context(), "translator", &MockTranslator{}))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("got %d", rec.Code)
	}
}

func TestRateLimitSurvivesPortAndHeaderChanges(t *testing.T) {
	handler := middleware.RateLimiterMiddleware(0.01, 1)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) }))
	for i, want := range []int{http.StatusOK, http.StatusTooManyRequests, http.StatusTooManyRequests} {
		req := httptest.NewRequest("GET", "/", nil)
		req.RemoteAddr = fmt.Sprintf("198.51.100.231:%d", 1000+i)
		req.Header.Set("X-Forwarded-For", fmt.Sprintf("203.0.113.%d", i+1))
		req.Header.Set("CF-Connecting-IP", fmt.Sprintf("203.0.113.%d", i+1))
		req = req.WithContext(context.WithValue(req.Context(), "translator", &MockTranslator{}))
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		if rec.Code != want {
			t.Fatalf("request %d: got %d want %d", i, rec.Code, want)
		}
	}
}
