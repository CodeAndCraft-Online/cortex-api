package subs

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRegisterSubRoutes(t *testing.T) {
	router := gin.New()
	api := router.Group("/api")

	assert.NotPanics(t, func() {
		RegisterSubRoutes(api)
	})

	assert.NotNil(t, api)
}
