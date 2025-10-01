package votes

import (
	"github.com/CodeAndCraft-Online/cortex-api/internal/handlers"
	middleware "github.com/CodeAndCraft-Online/cortex-api/pkg"
	"github.com/gin-gonic/gin"
)

func RegisterVotesRoutes(router *gin.RouterGroup) {
	votesRoutes := router.Group("/vote")
	votesRoutes.Use(middleware.AuthMiddleware())
	{
		votesRoutes.POST("/upvote", handlers.UpvotePost)
		votesRoutes.POST("/downvote", handlers.DownvotePost)
	}
}
