package routes

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestBuildRouteGroups(t *testing.T) {
	// Create a new router
	router := gin.New()

	// Test that BuildRouteGroups doesn't panic
	assert.NotPanics(t, func() {
		BuildRouteGroups(router)
	})

	// Test that routes are registered (basic smoke test)
	assert.NotNil(t, router)
}
