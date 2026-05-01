package security

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
)

const maxBodyBytes = 1 << 20 // 1 MB

// Headers applies security headers to every response (ADR-0003).
func Headers(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("X-Frame-Options", "DENY")
		h.Set("X-XSS-Protection", "1; mode=block")
		h.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		h.Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
		// HSTS only makes sense over TLS
		if r.TLS != nil {
			h.Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		}
		next.ServeHTTP(w, r)
	})
}

// BodyLimit caps the request body to 1 MB.
func BodyLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
		next.ServeHTTP(w, r)
	})
}

// TrustedProxies validates X-Forwarded-For against an allowlist.
// TRUSTED_PROXIES env var is a comma-separated list of CIDRs.
// If empty, remote addr is used directly.
func TrustedProxies(next http.Handler) http.Handler {
	cidrs := parseTrustedProxies()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(cidrs) > 0 {
			remoteIP, _, _ := net.SplitHostPort(r.RemoteAddr)
			if isTrusted(remoteIP, cidrs) {
				if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
					// Use first IP in chain as the real client IP
					parts := strings.Split(xff, ",")
					r.RemoteAddr = strings.TrimSpace(parts[0]) + ":0"
				}
			}
		}
		next.ServeHTTP(w, r)
	})
}

func parseTrustedProxies() []*net.IPNet {
	raw := os.Getenv("TRUSTED_PROXIES")
	if raw == "" {
		return nil
	}
	var cidrs []*net.IPNet
	for _, cidr := range strings.Split(raw, ",") {
		cidr = strings.TrimSpace(cidr)
		if cidr == "" {
			continue
		}
		if !strings.Contains(cidr, "/") {
			cidr = fmt.Sprintf("%s/32", cidr)
		}
		_, network, err := net.ParseCIDR(cidr)
		if err == nil {
			cidrs = append(cidrs, network)
		}
	}
	return cidrs
}

func isTrusted(ip string, cidrs []*net.IPNet) bool {
	parsed := net.ParseIP(ip)
	if parsed == nil {
		return false
	}
	for _, cidr := range cidrs {
		if cidr.Contains(parsed) {
			return true
		}
	}
	return false
}
