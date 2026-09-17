package clientip

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync/atomic"
)

type Resolver struct{ trusted []*net.IPNet }

var configured atomic.Pointer[Resolver]

func New(value string) (*Resolver, error) {
	result := &Resolver{}
	for _, item := range strings.Split(value, ",") {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		if ip := net.ParseIP(item); ip != nil {
			bits := 128
			if ip.To4() != nil {
				ip = ip.To4()
				bits = 32
			}
			result.trusted = append(result.trusted, &net.IPNet{IP: ip, Mask: net.CIDRMask(bits, bits)})
			continue
		}
		_, network, err := net.ParseCIDR(item)
		if err != nil {
			return nil, fmt.Errorf("invalid trusted proxy: %s", item)
		}
		result.trusted = append(result.trusted, network)
	}
	return result, nil
}

func Configure(value string) error {
	resolver, err := New(value)
	if err != nil {
		return err
	}
	configured.Store(resolver)
	return nil
}

func peer(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		host = r.RemoteAddr
	}
	if ip := net.ParseIP(host); ip != nil {
		return ip.String()
	}
	return host
}

func (r *Resolver) isTrusted(value string) bool {
	ip := net.ParseIP(value)
	for _, network := range r.trusted {
		if network.Contains(ip) {
			return true
		}
	}
	return false
}

func (resolver *Resolver) IP(r *http.Request) string {
	address := peer(r)
	if !resolver.isTrusted(address) {
		return address
	}
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		parts := strings.Split(forwarded, ",")
		for i := len(parts) - 1; i >= 0; i-- {
			if !resolver.isTrusted(address) {
				return address
			}
			ip := net.ParseIP(strings.TrimSpace(parts[i]))
			if ip == nil {
				return peer(r)
			}
			address = ip.String()
		}
		return address
	}
	if ip := net.ParseIP(strings.TrimSpace(r.Header.Get("X-Real-IP"))); ip != nil {
		return ip.String()
	}
	return address
}

func IP(r *http.Request) string {
	resolver := configured.Load()
	if resolver == nil {
		return peer(r)
	}
	return resolver.IP(r)
}
