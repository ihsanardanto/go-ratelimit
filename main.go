package goratelimit

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type RateLimiter struct {
	limit   int
	burst   int
	mu      sync.Mutex
	clients map[string]*client
}

func NewRateLimiter(rateLimit, burstLimit int) *RateLimiter {
	return &RateLimiter{
		limit:   rateLimit,
		burst:   burstLimit,
		clients: make(map[string]*client),
	}
}

func (rl *RateLimiter) AllowRequest(ip string) bool {
	return rl.clients[ip].limiter.Allow()
}
