package middleware

import (
	"log/slog"
	"net/http"
	"sync"
	"time"
)

// RateLimit is a simple in-memory token bucket to reduce brute force traffic.
func RateLimit(maxRequests int, window time.Duration) func(http.Handler) http.Handler {
	if maxRequests <= 0 {
		maxRequests = 60
	}
	if window <= 0 {
		window = time.Minute
	}

	type bucket struct {
		mu     sync.Mutex
		hits   []time.Time
	}

	store := make(map[string]*bucket)
	var mu sync.Mutex

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.Header.Get("X-Forwarded-For")
			if ip == "" {
				ip = r.RemoteAddr
			}

			mu.Lock()
			b, ok := store[ip]
			if !ok {
				b = &bucket{}
				store[ip] = b
			}
			mu.Unlock()

			b.mu.Lock()
			cutoff := time.Now().Add(-window)
			filtered := b.hits[:0]
			for _, hit := range b.hits {
				if hit.After(cutoff) {
					filtered = append(filtered, hit)
				}
			}
			b.hits = filtered
			if len(b.hits) >= maxRequests {
				b.mu.Unlock()
				slog.Warn("rate limit hit", "ip", ip, "limit", maxRequests)
				http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
				return
			}
			b.hits = append(b.hits, time.Now())
			b.mu.Unlock()

			next.ServeHTTP(w, r)
		})
	}
}
