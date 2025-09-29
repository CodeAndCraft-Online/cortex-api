package auth

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRegisterAuthRoutes(t *testing.T) {
	// Create a router group
	router := gin.New()
	api := router.Group("/api")

	// Test that RegisterAuthRoutes doesn't panic
	assert.NotPanics(t, func() {
		RegisterAuthRoutes(api)
	})

	assert.NotNil(t, api)
}
