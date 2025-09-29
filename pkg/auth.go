package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// ✅ Read values from environment variables
var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// AuthMiddleware extracts and verifies the JWT token
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			fmt.Println("Authorization header missing")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			c.Abort()
			return
		}

		// Extract the token from "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			fmt.Println("Invalid authorization format")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		tokenString := tokenParts[1]

		// Parse and verify the JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Replace this with your secret key
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			fmt.Println("Invalid token:", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Extract username from claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			fmt.Println("Failed to parse JWT claims")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		username, exists := claims["username"].(string)
		if !exists {
			fmt.Println("Username claim missing in token")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("username", username) // ✅ Store the username in context
		c.Next()
	}
}
