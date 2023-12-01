package http

import (
	"net/http"
	"time"

	goiprequest "github.com/ihsanardanto-djoin/go-ip-request"
	goratelimit "github.com/ihsanardanto/go-ratelimit"
	"golang.org/x/time/rate"
)

func RateLimitMiddleware(rl *goratelimit.RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		go func() {
			for {
				time.Sleep(time.Minute)
				// Lock the mutex to protect this section from race conditions.
				rl.Mu.Lock()
				for ip, client := range rl.Clients {
					if time.Since(client.LastSeen) > 3*time.Minute {
						delete(rl.Clients, ip)
					}
				}
				rl.Mu.Unlock()
			}
		}()

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			ip, err := goiprequest.GetClientIP(r)
			if err != nil {
				http.Error(w, "Unauthorized, IP Client not detected\n", http.StatusUnauthorized)
				return
			}

			rl.Mu.Lock()
			if _, found := rl.Clients[ip]; !found {
				rl.Clients[ip] = &goratelimit.Client{Limiter: rate.NewLimiter(rate.Limit(rl.Limit), rl.Burst)}
			}
			rl.Clients[ip].LastSeen = time.Now()

			if !rl.AllowRequest(ip) {
				rl.Mu.Unlock()
				http.Error(w, "Rate limit exceeded", http.StatusForbidden)
				return
			}

			rl.Mu.Unlock()
			next.ServeHTTP(w, r)
		})
	}
}

// Custom responseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
