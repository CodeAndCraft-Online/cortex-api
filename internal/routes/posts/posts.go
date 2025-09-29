package posts

import (
	"github.com/CodeAndCraft-Online/cortex-api/internal/handlers"
	"github.com/gin-gonic/gin"
)

// RegisterPostRoutes sets up routes for posts
func RegisterPostRoutes(router *gin.RouterGroup) {
	posts := router.Group("/posts")
	{
		posts.GET("/:id", handlers.GetPostByID)
		posts.POST("/", handlers.CreatePost)
		posts.GET("/", handlers.GetPosts)
		posts.POST("/posts/:postID", handlers.GetPostByID)
		posts.GET("/posts/:postID/comments", handlers.GetCommentsByPostID)
		// posts.PUT("/:id", handlers.UpdatePost)
		// posts.DELETE("/:id", handlers.DeletePost)
	}
}
