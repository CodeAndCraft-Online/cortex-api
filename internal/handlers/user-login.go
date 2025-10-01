package handlers

import (
	"net/http"
	"os"
	"time"

	db "github.com/CodeAndCraft-Online/cortex-api/internal/database"
	models "github.com/CodeAndCraft-Online/cortex-api/internal/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// UserLoginRequest represents the login request data
type UserLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// @Summary User Login
// @Description Authenticates a user with username and password, returns a JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body UserLoginRequest true "Login credentials - username and password"
// @Success 200 {object} map[string]string "token: JWT access token"
// @Failure 400 {object} map[string]string "error: Bad request - username and password required"
// @Failure 401 {object} map[string]string "error: Invalid credentials"
// @Failure 500 {object} map[string]string "error: Internal server error"
// @Router /auth/login [post]
func Login(c *gin.Context) {
	var req UserLoginRequest
	var foundUser models.User

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username is required"})
		return
	}

	if req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password is required"})
		return
	}

	db.DB.Where("username = ?", req.Username).First(&foundUser)
	if foundUser.ID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(req.Password))
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
