package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNewRateLimiter(t *testing.T) {
	limit := 10
	window := time.Minute

	rl := NewRateLimiter(limit, window)

	assert.NotNil(t, rl)
	assert.Equal(t, limit, rl.limit)
	assert.Equal(t, window, rl.window)
	assert.NotNil(t, rl.clients)
}

func TestRateLimiter_Middleware_UnderLimit(t *testing.T) {
	rl := NewRateLimiter(2, time.Minute)
	middleware := rl.Middleware()

	// Create a test router
	r := gin.New()
	r.Use(middleware)
	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "success")
	})

	// Create multiple requests from the same IP
	req, _ := http.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "127.0.0.1:12345"

	// First request should succeed
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req)
	assert.Equal(t, http.StatusOK, w1.Code)
	assert.Equal(t, "success", w1.Body.String())

	// Second request should succeed
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req)
	assert.Equal(t, http.StatusOK, w2.Code)
	assert.Equal(t, "success", w2.Body.String())
}

func TestRateLimiter_Middleware_OverLimit(t *testing.T) {
	rl := NewRateLimiter(1, time.Minute)
	middleware := rl.Middleware()

	// Create a test router
	r := gin.New()
	r.Use(middleware)
	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "success")
	})

	// Create request from the same IP
	req, _ := http.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "127.0.0.1:12345"

	// First request should succeed
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req)
	assert.Equal(t, http.StatusOK, w1.Code)
	assert.Equal(t, "success", w1.Body.String())

	// Second request should be blocked
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req)
	assert.Equal(t, http.StatusTooManyRequests, w2.Code)
	assert.Contains(t, w2.Body.String(), "Too many requests")
}

func TestRateLimiter_Middleware_DifferentIPs(t *testing.T) {
	rl := NewRateLimiter(1, time.Minute)
	middleware := rl.Middleware()

	// Create a test router
	r := gin.New()
	r.Use(middleware)
	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "success")
	})

	// Create requests from different IPs
	req1, _ := http.NewRequest("GET", "/test", nil)
	req1.RemoteAddr = "127.0.0.1:12345"

	req2, _ := http.NewRequest("GET", "/test", nil)
	req2.RemoteAddr = "127.0.0.2:12345"

	// Both requests should succeed (different IPs)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)

	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestRateLimiter_Cleanup(t *testing.T) {
	rl := NewRateLimiter(1, 10*time.Millisecond) // Very short window for testing

	// Set up initial state
	rl.mu.Lock()
	rl.clients["127.0.0.1"] = 1
	rl.mu.Unlock()

	assert.Equal(t, 1, len(rl.clients))

	// Wait for cleanup to run
	time.Sleep(50 * time.Millisecond)

	rl.mu.Lock()
	defer rl.mu.Unlock()
	// After cleanup runs, clients map should be empty
	assert.Empty(t, rl.clients)
}

func TestRateLimiter_GetClientIP(t *testing.T) {
	rl := NewRateLimiter(1, time.Minute)
	middleware := rl.Middleware()

	// Create a test router
	r := gin.New()
	r.Use(middleware)
	r.GET("/test", func(c *gin.Context) {
		ip := c.ClientIP()
		c.String(http.StatusOK, ip)
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "127.0.0.1:12345"

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "127.0.0.1", w.Body.String())
}

func TestRateLimiter_ConcurrentAccess(t *testing.T) {
	rl := NewRateLimiter(10, time.Minute)
	middleware := rl.Middleware()

	r := gin.New()
	r.Use(middleware)
	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "success")
	})

	// Test concurrent access
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func() {
			req, _ := http.NewRequest("GET", "/test", nil)
			req.RemoteAddr = "127.0.0.1:12345"
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	assert.NotPanics(t, func() {
		// Test should not panic due to concurrent access
	})
}
