package main

import (
	"net/http"

	"github.com/go-chi/chi"
	goratelimit "github.com/ihsanardanto/go-ratelimit"
	ratelimit "github.com/ihsanardanto/go-ratelimit/wrapper/http"
)

func main() {
	r := chi.NewRouter()

	// Create a new rate limiter
	rl := goratelimit.NewRateLimiter(2, 6)

	// Use the middleware with the list of allowed IP addresses
	r.Use(ratelimit.RateLimitMiddleware(rl))

	// Define your routes and handlers here
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, world!\n"))
	})

	http.ListenAndServe(":8080", r)
}
