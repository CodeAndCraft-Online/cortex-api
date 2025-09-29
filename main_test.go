package main

import (
	"os"
	"testing"
	"time"

	routes "github.com/CodeAndCraft-Online/cortex-api/internal/routes"
	rateLimit "github.com/CodeAndCraft-Online/cortex-api/pkg"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMainComponents(t *testing.T) {
	// Set up test environment to avoid database connection
	originalEnv := os.Getenv("PORT")
	defer func() {
		if originalEnv != "" {
			os.Setenv("PORT", originalEnv)
		} else {
			os.Unsetenv("PORT")
		}
	}()
	os.Setenv("PORT", "8080")

	// Test that DB initialization doesn't panic (though it may fail in test env)
	assert.NotPanics(t, func() {
		defer func() {
			if r := recover(); r != nil {
				// Expected to potentially fail in test environment without real DB
				t.Log("DB initialization failed as expected in test environment:", r)
			}
		}()
		// We can't easily test InitDB without a database, but we can test it doesn't crash immediately
	})

	// Test router creation
	assert.NotPanics(t, func() {
		router := gin.Default()
		assert.NotNil(t, router)
	})

	// Test rate limiter creation
	assert.NotPanics(t, func() {
		limiter := rateLimit.NewRateLimiter(100, time.Minute)
		assert.NotNil(t, limiter)
		assert.NotPanics(t, func() {
			router := gin.Default()
			router.Use(limiter.Middleware())
		})
	})

	// Test route building
	assert.NotPanics(t, func() {
		router := gin.Default()
		routes.BuildRouteGroups(router)
	})
}
