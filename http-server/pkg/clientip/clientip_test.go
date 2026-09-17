package clientip

import (
	"net/http/httptest"
	"testing"
)

func TestClientIPTrustBoundary(t *testing.T) {
	resolver, err := New("127.0.0.1/32,10.0.0.0/8")
	if err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct{ name, peer, xff, real, cf, want string }{
		{"direct spoof", "198.51.100.4:1234", "203.0.113.9", "203.0.113.8", "203.0.113.7", "198.51.100.4"},
		{"new port", "198.51.100.4:9999", "", "", "", "198.51.100.4"},
		{"trusted proxy", "127.0.0.1:9000", "198.51.100.4", "", "", "198.51.100.4"},
		{"untrusted left prefix", "127.0.0.1:9000", "203.0.113.9, 198.51.100.4", "", "", "198.51.100.4"},
		{"trusted chain", "127.0.0.1:9000", "198.51.100.4, 10.1.1.1", "", "", "198.51.100.4"},
		{"malformed chain", "127.0.0.1:9000", "garbage", "", "", "127.0.0.1"},
		{"ipv6", "[2001:db8::1]:1234", "", "", "", "2001:db8::1"},
		{"cloudflare header alone ignored", "127.0.0.1:9000", "", "", "203.0.113.7", "127.0.0.1"},
	} {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			req.RemoteAddr = test.peer
			req.Header.Set("X-Forwarded-For", test.xff)
			req.Header.Set("X-Real-IP", test.real)
			req.Header.Set("CF-Connecting-IP", test.cf)
			if got := resolver.IP(req); got != test.want {
				t.Fatalf("got %s, want %s", got, test.want)
			}
		})
	}
	if _, err := New("not-a-proxy"); err == nil {
		t.Fatal("invalid proxy accepted")
	}
}
