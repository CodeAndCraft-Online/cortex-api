package comments

import (
	"github.com/CodeAndCraft-Online/cortex-api/internal/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterCommentsRoutes(router *gin.RouterGroup) {
	comments := router.Group("/comments")
	{
		comments.POST("/comments", handlers.CreateComment)
	}
}
