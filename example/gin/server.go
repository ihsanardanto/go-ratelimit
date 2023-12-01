package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	goratelimit "github.com/ihsanardanto/go-ratelimit"
	ginratelimit "github.com/ihsanardanto/go-ratelimit/wrapper/gin"
)

func main() {
	r := gin.Default()

	// Create a new rate limiter
	rl := goratelimit.NewRateLimiter(2, 6)

	// Use the middleware with the list of allowed IP addresses
	r.Use(ginratelimit.RateLimitMiddleware(rl))

	// Define your routes and handlers here
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, World!\n")
	})

	r.Run()
}
