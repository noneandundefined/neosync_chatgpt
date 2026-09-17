package tests

import (
	"context"
	"neomatica/neosync/infra/constants"
	"neomatica/neosync/middleware"
	"neomatica/neosync/types"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestRequireRoles_ForbiddenForUser(t *testing.T) {
	tr := &MockTranslator{}

	routes := []struct {
		method string
		path   string
		roles  []string
	}{
		{"GET", "/", []string{constants.Role_SuperAdmin, constants.Role_Support, constants.Role_AdminL2}},
		{"GET", "/owners", []string{constants.Role_SuperAdmin, constants.Role_Support, constants.Role_AdminL2}},
		{"DELETE", "/", []string{constants.Role_SuperAdmin, constants.Role_Support, constants.Role_AdminL2}},
		{"DELETE", "/123e4567-e89b-12d3-a456-426614174000", []string{constants.Role_SuperAdmin, constants.Role_Support, constants.Role_AdminL2}},
		{"POST", "/123e4567-e89b-12d3-a456-426614174000/send/email", []string{constants.Role_SuperAdmin, constants.Role_Support, constants.Role_AdminL2}},
	}

	router := mux.NewRouter()
	for _, r := range routes {
		mw := middleware.RequireRoles(r.roles...)
		h := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		router.Handle(r.path, h).Methods(r.method)
	}

	for _, r := range routes {
		req := httptest.NewRequest(r.method, r.path, nil)
		//nolint
		ctx := context.WithValue(req.Context(), "translator", tr)
		ctx = context.WithValue(ctx, "identity", &types.AuthToken{RoleCode: constants.Role_User})
		req = req.WithContext(ctx)

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if rec.Result().StatusCode != http.StatusForbidden {
			t.Errorf("expected status 403 Forbidden, got %d", rec.Result().StatusCode)
		}
	}
}

func TestRequireRolesRejectsUnverifiedCookie(t *testing.T) {
	handler := middleware.RequireRoles(constants.Role_SuperAdmin)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { t.Fatal("unverified cookie authorized") }))
	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{Name: "access_token", Value: "unverified"})
	req = req.WithContext(context.WithValue(req.Context(), "translator", &MockTranslator{}))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("got %d", rec.Code)
	}
}
