package tests

import (
	"net/http"
	"net/http/httptest"
	"strings"
)

type MockTranslator struct{}

func (m *MockTranslator) T(key string) string {
	return key
}

func (m *MockTranslator) TErr(key string) string {
	return strings.ToLower(key)
}

func (m *MockTranslator) GetLang() string {
	return "ru"
}

func MakeReqWithIp(method, target, body, ip string) *http.Request {
	req := httptest.NewRequest(method, target, strings.NewReader(body))
	req.RemoteAddr = ip + ":12345"
	return req
}
