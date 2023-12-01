package main

import (
	"net/http"

	goratelimit "github.com/ihsanardanto/go-ratelimit"
	echowritelimit "github.com/ihsanardanto/go-ratelimit/wrapper/echo"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	// Create a new rate limiter
	rl := goratelimit.NewRateLimiter(2, 6)

	// Use the rate limiter as middleware
	e.Use(echowritelimit.RateLimitMiddleware(rl))

	// Define your routes and handlers
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!\n")
	})

	// Start the server
	e.Start(":8080")
}
