package users

import (
	handlers "github.com/CodeAndCraft-Online/cortex-api/internal/handlers"
	middleware "github.com/CodeAndCraft-Online/cortex-api/pkg"
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(router *gin.RouterGroup) {
	userRoutes := router.Group("/user")
	{
		userRoutes.POST("/invite/:inviteID/accept", handlers.AcceptInvite)
	}

	// Public user profile routes
	userRoutes.GET("/:username", handlers.GetUserProfile)

	// Protected user profile management routes (require authentication)
	protectedUserRoutes := router.Group("/user")
	protectedUserRoutes.Use(middleware.AuthMiddleware())
	{
		protectedUserRoutes.GET("/profile", handlers.GetCurrentUserProfile)
		protectedUserRoutes.PUT("/profile", handlers.UpdateUserProfile)
		protectedUserRoutes.DELETE("/profile", handlers.DeleteUserAccount)
	}
}
