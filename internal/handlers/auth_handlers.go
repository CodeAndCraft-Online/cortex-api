package handlers

import (
	"fmt"
	"net/http"

	"github.com/CodeAndCraft-Online/cortex-api/internal/services"
	"github.com/gin-gonic/gin"
)

// @Summary Request Password Reset
// @Description Generates a password reset token for the specified username
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body map[string]string true "Username for password reset"
// @Success 200 {object} map[string]interface{} "message: Reset token generated. Use it to reset your password., token: Reset token"
// @Failure 400 {object} map[string]string "error: Bad request or username required"
// @Failure 404 {object} map[string]string "error: User not found"
// @Router /auth/password-reset/request [post]
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

// @Summary Reset Password
// @Description Resets a user's password using a valid reset token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body map[string]string true "Reset token and new password"
// @Success 200 {object} map[string]string "message: Password has been reset successfully"
// @Failure 400 {object} map[string]string "error: Bad request, invalid token, or weak password"
// @Failure 404 {object} map[string]string "error: Token not found or expired"
// @Router /auth/password-reset/reset [post]
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
