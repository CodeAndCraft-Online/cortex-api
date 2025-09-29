package subs

import (
	handlers "github.com/CodeAndCraft-Online/cortex-api/internal/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterSubRoutes(router *gin.RouterGroup) {
	subRoutes := router.Group("/sub")
	{
		subRoutes.GET("/", handlers.GetSubs)
		subRoutes.GET("/sub/:subID/postCount", handlers.GetPostCountPerSub)
		subRoutes.POST("/sub", handlers.CreateSub)
		subRoutes.POST("/sub/:subID/join", handlers.JoinSub)
		subRoutes.POST("/sub/:subID/leave", handlers.LeaveSub)
		subRoutes.GET("/sub/:subID/posts", handlers.ListSubPosts)
		subRoutes.POST("/sub/:subID/invite", handlers.InviteUser)
	}
}

// // Protected routes (require authentication)
// protectedRoutes := router.Group("/v1")
// protectedRoutes.Use(middleware.AuthMiddleware()) // Apply JWT middleware // ✅ Private
// {
// 	protectedRoutes.POST("/upvote", controllers.UpvotePost)
// 	protectedRoutes.POST("/downvote", controllers.DownvotePost)
// }
