package handlers

import (
	"net/http"
	"strings"

	db "github.com/CodeAndCraft-Online/cortex-api/internal/database"
	"github.com/CodeAndCraft-Online/cortex-api/internal/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// @Summary Register a new user
// @Description Creates a new user account with username and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.User true "User registration details"
// @Success 200 {object} map[string]string "message: User registered successfully"
// @Failure 400 {object} map[string]string "error: Bad request - username and password required"
// @Failure 409 {object} map[string]string "error: Username already taken"
// @Failure 500 {object} map[string]string "error: Internal server error"
// @Router /auth/register [post]
func Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if user.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username is required"})
		return
	}

	if user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password is required"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = string(hashedPassword)

	if err := db.DB.Create(&user).Error; err != nil {
		// Check if it's a duplicate username error
		if err.Error() == "UNIQUE constraint failed: users.username" ||
			err.Error() == "duplicate key value violates unique constraint \"uni_users_username\"" ||
			strings.Contains(err.Error(), "duplicate") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "username already taken"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}
