package gin

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	goiprequest "github.com/ihsanardanto-djoin/go-ip-request"
	goratelimit "github.com/ihsanardanto/go-ratelimit"
	"golang.org/x/time/rate"
)

func RateLimitMiddleware(rl *goratelimit.RateLimiter) gin.HandlerFunc {
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

	return func(c *gin.Context) {

		ip, _ := goiprequest.GetClientIP(c.Request)
		rl.Mu.Lock()
		if _, found := rl.Clients[ip]; !found {
			rl.Clients[ip] = &goratelimit.Client{Limiter: rate.NewLimiter(rate.Limit(rl.Limit), rl.Burst)}
		}
		rl.Clients[ip].LastSeen = time.Now()

		if !rl.AllowRequest(ip) {
			rl.Mu.Unlock()
			c.AbortWithStatusJSON(http.StatusTooManyRequests, map[string]string{"error": "Rate limit exceeded"})
			return
		}

		rl.Mu.Unlock()
		c.Next()
	}
}
