package goratelimit

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type Client struct {
	Limiter  *rate.Limiter
	LastSeen time.Time
}

type RateLimiter struct {
	Limit   int
	Burst   int
	Mu      sync.Mutex
	Clients map[string]*Client
}

func NewRateLimiter(rateLimit, burstLimit int) *RateLimiter {
	return &RateLimiter{
		Limit:   rateLimit,
		Burst:   burstLimit,
		Clients: make(map[string]*Client),
	}
}

func (rl *RateLimiter) AllowRequest(ip string) bool {
	return rl.Clients[ip].Limiter.Allow()
}
