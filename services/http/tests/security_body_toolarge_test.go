package tests

import (
	"bytes"
	"context"
	"neomatica/neosync/middleware"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSecurityMiddleware_BodyTooLarge(t *testing.T) {
	handler := middleware.SecurityMiddleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))

	largeBody := bytes.Repeat([]byte("A"), 1024*1024+1)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(largeBody))
	//nolint
	ctx := context.WithValue(req.Context(), "translator", &MockTranslator{})
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Errorf("Expected status %d, got %d", http.StatusRequestEntityTooLarge, rec.Code)
	}
}
