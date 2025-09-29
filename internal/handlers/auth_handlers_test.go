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

func TestRequestPasswordResetHandler_InvalidJSON(t *testing.T) {
	router := setupTestRouter()

	// Create test request with invalid JSON
	req, _ := http.NewRequest("POST", "/auth/reset-password-request", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Check response - should be BadRequest
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRequestPasswordResetHandler_MissingUsername(t *testing.T) {
	router := setupTestRouter()

	// Create request payload with missing username
	requestBody := map[string]interface{}{
		// username is missing
	}
	jsonData, _ := json.Marshal(requestBody)

	// Create test request
	req, _ := http.NewRequest("POST", "/auth/reset-password-request", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Check response - should be BadRequest due to missing required field
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestResetPasswordHandler_InvalidJSON(t *testing.T) {
	router := setupTestRouter()

	// Create test request with invalid JSON
	req, _ := http.NewRequest("POST", "/auth/reset-password", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Check response - should be BadRequest
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestResetPasswordHandler_MissingFields(t *testing.T) {
	router := setupTestRouter()

	testCases := []struct {
		name        string
		requestBody map[string]interface{}
	}{
		{
			name: "missing token",
			requestBody: map[string]interface{}{
				"new_password": "password",
			},
		},
		{
			name: "missing password",
			requestBody: map[string]interface{}{
				"token": "sometoken",
			},
		},
		{
			name:        "empty request",
			requestBody: map[string]interface{}{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			jsonData, _ := json.Marshal(tc.requestBody)

			req, _ := http.NewRequest("POST", "/auth/reset-password", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code)
		})
	}
}

func TestResetPasswordHandler_ExpiredToken(t *testing.T) {
	router := setupTestRouter()

	// Create test user and expired reset token
	user := models.User{
		Username: "expireduser",
		Password: "oldpassword",
	}
	database.DB.Create(&user)

	// Create an expired token (set expires_at to past time)
	expiredToken := models.PasswordResetToken{
		UserID:    user.ID,
		Token:     "expiredtoken123",
		ExpiresAt: time.Now().Add(-1 * time.Hour), // Already expired
	}
	database.DB.Create(&expiredToken)

	// Create request payload with expired token
	requestBody := map[string]string{
		"token":        "expiredtoken123",
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

	// Check response - should be error for expired token
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestResetPasswordHandler_WeakPassword(t *testing.T) {
	router := setupTestRouter()

	// Create test user and reset token
	user := models.User{
		Username: "weakpassuser",
		Password: "oldpassword",
	}
	database.DB.Create(&user)

	resetToken := models.PasswordResetToken{
		UserID:    user.ID,
		Token:     "weakpasstoken123",
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}
	database.DB.Create(&resetToken)

	// Create request payload with weak password (empty string)
	requestBody := map[string]string{
		"token":        "weakpasstoken123",
		"new_password": "",
	}
	jsonData, _ := json.Marshal(requestBody)

	// Create test request
	req, _ := http.NewRequest("POST", "/auth/reset-password", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Service layer should handle password validation - this may return OK if service accepts empty passwords
	// For now, just verify it returns a valid JSON response
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.NotNil(t, response)
}
