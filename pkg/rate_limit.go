package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter stores request counts
type RateLimiter struct {
	mu      sync.Mutex
	clients map[string]int
	limit   int
	window  time.Duration
}

// NewRateLimiter initializes a rate limiter
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		clients: make(map[string]int),
		limit:   limit,
		window:  window,
	}
	// Reset counts periodically
	go rl.cleanup()
	return rl
}

// Middleware applies rate limiting
func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		rl.mu.Lock()
		defer rl.mu.Unlock()

		if rl.clients[ip] >= rl.limit {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			c.Abort()
			return
		}

		rl.clients[ip]++
		c.Next()
	}
}

// cleanup resets request counts periodically
func (rl *RateLimiter) cleanup() {
	for {
		time.Sleep(rl.window)
		rl.mu.Lock()
		rl.clients = make(map[string]int)
		rl.mu.Unlock()
	}
}
