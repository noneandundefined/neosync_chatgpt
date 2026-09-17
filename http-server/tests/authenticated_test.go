package tests

import (
	"context"
	"neomatica/neosync/handler"
	"neomatica/neosync/infra/store/postgres/store"
	"neomatica/neosync/middleware"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIsAuthenticatedMiddleware(t *testing.T) {
	baseHandler := &handler.BaseHandler{
		Store: store.Storage{},
	}

	mw := middleware.IsAuthenticatedMiddleware(baseHandler)
	handlerFunc := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called when no auth-token")
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	//nolint
	ctx := context.WithValue(req.Context(), "translator", &MockTranslator{})
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()

	handlerFunc.ServeHTTP(rec, req)

	if rec.Result().StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401 Unauthorized, got %d", rec.Result().StatusCode)
	}
}
