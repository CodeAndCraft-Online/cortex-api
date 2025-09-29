package handlers

import (
	"fmt"
	"net/http"

	"github.com/CodeAndCraft-Online/cortex-api/internal/services"
	"github.com/gin-gonic/gin"
)

// RequestPasswordReset generates a reset token for the user
func RequestPasswordReset(c *gin.Context) {
	var request struct {
		Username string `json:"username"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if request.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username is required"})
		return
	}

	passwordResetToken, err := services.ResetPasswordRequest(request.Username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"err": err,
		})
		return
	}

	// **In real-world apps, send this token via email or SMS**
	fmt.Println("Password reset token for", passwordResetToken.UserID, ":", passwordResetToken.Token)

	c.JSON(http.StatusOK, gin.H{"message": "Reset token generated. Use it to reset your password.", "token": passwordResetToken.Token})
}

// ResetPassword allows users to reset their password using a valid token
func ResetPassword(c *gin.Context) {
	var request struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := services.ResetPassword(request.Token, request.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password has been reset successfully"})
}
