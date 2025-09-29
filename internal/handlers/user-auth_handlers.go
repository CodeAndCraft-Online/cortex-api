package handlers

import (
	"net/http"

	db "github.com/CodeAndCraft-Online/cortex-api/internal/database"
	models "github.com/CodeAndCraft-Online/cortex-api/internal/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// @Summary Register a new user
// @Description Creates a new user account
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string "Success"
// @Failure 400 {object} map[string]string "Bad request"
// @Router /auth/register [post]
func Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)

	db.DB.Create(&user)
	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}
