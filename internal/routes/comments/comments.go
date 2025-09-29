package comments

import (
	"github.com/CodeAndCraft-Online/cortex-api/internal/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterCommentsRoutes(router *gin.RouterGroup) {
	comments := router.Group("/comments")
	{
		// Individual comment CRUD operations
		comments.GET("/:id", handlers.GetCommentByID)   // Get comment by ID (public)
		comments.PUT("/:id", handlers.UpdateComment)    // Update comment (author only, handled in handler)
		comments.DELETE("/:id", handlers.DeleteComment) // Delete comment (author only, handled in handler)

		// Legacy routes (keep for backward compatibility)
		comments.POST("/comments", handlers.CreateComment) // Create comment via post endpoint
	}
}
