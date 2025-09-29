package votes

import (
	handlers "github.com/CodeAndCraft-Online/cortex-api/internal/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterVotesRoutes(router *gin.RouterGroup) {
	votesRoutes := router.Group("/vote")
	{
		votesRoutes.POST("/upvote", handlers.UpvotePost)
		votesRoutes.POST("/downvote", handlers.DownvotePost)
	}
}
