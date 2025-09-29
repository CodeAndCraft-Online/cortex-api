package comments

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRegisterCommentsRoutes(t *testing.T) {
	router := gin.New()
	api := router.Group("/api")

	assert.NotPanics(t, func() {
		RegisterCommentsRoutes(api)
	})

	assert.NotNil(t, api)
}
