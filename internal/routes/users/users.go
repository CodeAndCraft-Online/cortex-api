package users

import (
	handlers "github.com/CodeAndCraft-Online/cortex-api/internal/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(router *gin.RouterGroup) {
	userRoutes := router.Group("/user")
	{
		userRoutes.POST("/invite/:inviteID/accept", handlers.AcceptInvite)
	}
}
