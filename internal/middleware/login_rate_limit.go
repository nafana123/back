package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"
)

type loginAttempt struct {
	count       int
	windowStart time.Time
}

func RateLimitByIP(maxAttempts int, window time.Duration, message string) func(http.Handler) http.Handler {
	var (
		mu       sync.Mutex
		attempts = make(map[string]loginAttempt)
	)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := clientIP(r.RemoteAddr)
			now := time.Now()

			mu.Lock()
			entry := attempts[ip]
			if entry.windowStart.IsZero() || now.Sub(entry.windowStart) > window {
				entry = loginAttempt{
					count:       0,
					windowStart: now,
				}
			}

			entry.count++
			attempts[ip] = entry
			remaining := maxAttempts - entry.count
			mu.Unlock()

			if remaining < 0 {
				http.Error(w, message, http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func AuthRateLimit(maxAttempts int, window time.Duration, message string) func(http.Handler) http.Handler {
	return RateLimitByIP(maxAttempts, window, message)
}

func clientIP(remoteAddr string) string {
	host, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		return remoteAddr
	}
	return host
}
