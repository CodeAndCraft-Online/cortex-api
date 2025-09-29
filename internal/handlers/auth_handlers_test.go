package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/CodeAndCraft-Online/cortex-api/internal/database"
	"github.com/CodeAndCraft-Online/cortex-api/internal/models"
	"github.com/CodeAndCraft-Online/cortex-api/internal/testutils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	db, teardown, err := testutils.SetupTestDB()
	if err != nil {
		println("Docker not available, skipping handler integration tests:", err.Error())
		database.DB = nil // Ensure no stale database connection
		os.Exit(0)        // Skip all tests in this package
	}

	database.DB = db
	m.Run()
	teardown()
}

// Helper function to create test router
func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/auth/reset-password-request", RequestPasswordReset)
	router.POST("/auth/reset-password", ResetPassword)
	return router
}

func TestRequestPasswordResetHandler(t *testing.T) {
	router := setupTestRouter()

	// Create test user
	user := models.User{
		Username: "handleruser",
		Password: "password",
	}
	database.DB.Create(&user)

	// Create request payload
	requestBody := map[string]string{
		"username": "handleruser",
	}
	jsonData, _ := json.Marshal(requestBody)

	// Create test request
	req, _ := http.NewRequest("POST", "/auth/reset-password-request", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Check response
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Contains(t, response, "message")
	assert.Contains(t, response, "token")
	assert.NotEmpty(t, response["token"])
}

func TestRequestPasswordResetHandler_UserNotFound(t *testing.T) {
	router := setupTestRouter()

	// Create request payload for nonexistent user
	requestBody := map[string]string{
		"username": "nonexistent",
	}
	jsonData, _ := json.Marshal(requestBody)

	// Create test request
	req, _ := http.NewRequest("POST", "/auth/reset-password-request", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Check response
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestResetPasswordHandler(t *testing.T) {
	router := setupTestRouter()

	// Create test user and reset token
	user := models.User{
		Username: "handlerresetuser",
		Password: "oldpassword",
	}
	database.DB.Create(&user)

	resetToken := models.PasswordResetToken{
		UserID:    user.ID,
		Token:     "handlertoken123",
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}
	database.DB.Create(&resetToken)

	// Create request payload
	requestBody := map[string]string{
		"token":        "handlertoken123",
		"new_password": "newsecurepassword",
	}
	jsonData, _ := json.Marshal(requestBody)

	// Create test request
	req, _ := http.NewRequest("POST", "/auth/reset-password", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Check response
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Password has been reset successfully", response["message"])
}

func TestResetPasswordHandler_InvalidToken(t *testing.T) {
	router := setupTestRouter()

	// Create request payload with invalid token
	requestBody := map[string]string{
		"token":        "invalidtoken",
		"new_password": "newpassword",
	}
	jsonData, _ := json.Marshal(requestBody)

	// Create test request
	req, _ := http.NewRequest("POST", "/auth/reset-password", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Check response - should be error
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
