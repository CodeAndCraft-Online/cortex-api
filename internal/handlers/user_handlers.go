// Package handlers provides HTTP handlers for user management
package handlers

import (
	"net/http"
	"strconv"

	"github.com/CodeAndCraft-Online/cortex-api/internal/models"
	"github.com/CodeAndCraft-Online/cortex-api/internal/repositories"
	"github.com/CodeAndCraft-Online/cortex-api/internal/services"
	"github.com/gin-gonic/gin"
)

// GetUserProfile retrieves the public profile of a user by username
// @Summary Get user profile
// @Description Get public profile information for a user
// @Tags users
// @Accept  json
// @Produce  json
// @Param username path string true "Username"
// @Success 200 {object} models.UserResponse
// @Failure 404 {object} map[string]string "error: User not found"
// @Router /user/{username} [get]
func GetUserProfile(c *gin.Context) {
	username := c.Param("username")

	// Get user profile via service
	userResponse, err := services.GetUserProfile(username, nil)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, userResponse)
}

// GetCurrentUserProfile retrieves the authenticated user's full profile
// @Summary Get current user profile
// @Description Get the full profile for the authenticated user
// @Tags users
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} models.UserProfileResponse
// @Failure 401 {object} map[string]string "error: Unauthorized"
// @Failure 404 {object} map[string]string "error: User not found"
// @Router /user/profile [get]
func GetCurrentUserProfile(c *gin.Context) {
	// Get username from JWT token (set by auth middleware)
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Look up user by username to get user ID
	repo := repositories.NewUserRepository()
	user, err := repo.GetUserByUsername(username.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User lookup failed"})
		return
	}

	// Get user profile via service
	userProfile, err := services.GetUserProfileInternal(user.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, userProfile)
}

// UpdateUserProfile updates the authenticated user's profile
// @Summary Update user profile
// @Description Update the authenticated user's profile information
// @Tags users
// @Accept  json
// @Produce  json
// @Param profile body models.UserUpdateRequest true "Profile update data"
// @Security ApiKeyAuth
// @Success 200 {object} models.UserProfileResponse
// @Failure 400 {object} map[string]string "error: Bad request or validation error"
// @Failure 401 {object} map[string]string "error: Unauthorized"
// @Failure 500 {object} map[string]string "error: Internal server error"
// @Router /user/profile [put]
func UpdateUserProfile(c *gin.Context) {
	// Get username from JWT token
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Look up user by username to get user ID
	repo := repositories.NewUserRepository()
	user, err := repo.GetUserByUsername(username.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User lookup failed"})
		return
	}

	// Parse request body
	var updateRequest models.UserUpdateRequest
	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update user profile via service
	updatedProfile, err := services.UpdateUserProfile(user.ID, updateRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedProfile)
}

// DeleteUserAccount deletes the authenticated user's account
// @Summary Delete user account
// @Description Delete the authenticated user's account permanently
// @Tags users
// @Accept  json
// @Produce  json
// @Param password body map[string]string true "Current password confirmation"
// @Security ApiKeyAuth
// @Success 200 {object} map[string]string "message: Account deleted successfully"
// @Failure 400 {object} map[string]string "error: Bad request or invalid password"
// @Failure 401 {object} map[string]string "error: Unauthorized"
// @Failure 500 {object} map[string]string "error: Internal server error"
// @Router /user/profile [delete]
func DeleteUserAccount(c *gin.Context) {
	// Get username from JWT token
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Look up user by username to get user ID
	repo := repositories.NewUserRepository()
	user, err := repo.GetUserByUsername(username.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User lookup failed"})
		return
	}

	// Parse request body for password confirmation
	var requestBody struct {
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password confirmation required"})
		return
	}

	// Delete user account via service
	err = services.DeleteUserAccount(user.ID, requestBody.Password)
	if err != nil {
		if err.Error() == "invalid password" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account deleted successfully"})
}

// Legacy function for backward compatibility
func GetUserByID(c *gin.Context) {
	// Extract user ID from URL parameter
	idStr := c.Param("id")
	userID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	repo := repositories.NewUserRepository()
	user, err := repo.GetUserByID(uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}
