package posts

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRegisterPostRoutes(t *testing.T) {
	router := gin.New()
	api := router.Group("/api")

	assert.NotPanics(t, func() {
		RegisterPostRoutes(api)
	})

	assert.NotNil(t, api)
}
