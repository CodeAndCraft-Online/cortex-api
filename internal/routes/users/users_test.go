package users

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRegisterUserRoutes(t *testing.T) {
	router := gin.New()
	api := router.Group("/api")

	assert.NotPanics(t, func() {
		RegisterUserRoutes(api)
	})

	assert.NotNil(t, api)
}
