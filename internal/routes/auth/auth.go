package auth

import (
	handlers "github.com/CodeAndCraft-Online/cortex-api/internal/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(router *gin.RouterGroup) {
	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/register", handlers.Register)
		authRoutes.POST("/login", handlers.Login)
		authRoutes.POST("/password-reset/request", handlers.RequestPasswordReset)
		authRoutes.POST("/password-reset/reset", handlers.ResetPassword) // ✅ Public (requires reset token)
	}
}
