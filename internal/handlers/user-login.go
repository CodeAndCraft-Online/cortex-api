package handlers

import (
	"net/http"
	"time"

	db "github.com/CodeAndCraft-Online/cortex-api/internal/database"
	models "github.com/CodeAndCraft-Online/cortex-api/internal/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("your-secret-key")

// @Summary User Login
// @Description Authenticates a user and returns a JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string "Success"
// @Failure 400 {object} map[string]string "Bad request"
// @Router /auth/login [post]
func Login(c *gin.Context) {
	var user models.User
	var foundUser models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.DB.Where("username = ?", user.Username).First(&foundUser)
	if foundUser.ID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(user.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": foundUser.Username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, _ := token.SignedString(jwtSecret)

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
