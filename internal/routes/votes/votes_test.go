package votes

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRegisterVotesRoutes(t *testing.T) {
	router := gin.New()
	api := router.Group("/api")

	assert.NotPanics(t, func() {
		RegisterVotesRoutes(api)
	})

	assert.NotNil(t, api)
}
