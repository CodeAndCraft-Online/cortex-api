package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CodeAndCraft-Online/cortex-api/internal/database"
	"github.com/CodeAndCraft-Online/cortex-api/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

// Helper function to create test router
func setupUserTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Apply auth middleware to protected routes
	router.Use(func(c *gin.Context) {
		// Mock auth middleware - set a test username
		c.Set("username", "testuser")
		c.Next()
	})

	router.GET("/user/{username}", GetUserProfile)
	router.GET("/user/profile", GetCurrentUserProfile)
	router.PUT("/user/profile", UpdateUserProfile)
	router.DELETE("/user/profile", DeleteUserAccount)

	return router
}

func TestGetUserProfileHandler(t *testing.T) {
	if database.DB == nil {
		t.Skip("Database not available, skipping integration test")
		return
	}

	router := setupUserTestRouter()

	// Create test user
	user := models.User{
		Username:    "testuser",
		Password:    "password",
		DisplayName: "Test User",
		Bio:         "Test bio",
		Email:       stringPtr("test@example.com"),
		AvatarURL:   stringPtr("http://example.com/avatar.jpg"),
		IsPrivate:   false,
	}
	database.DB.Create(&user)

	// Test GET request
	req, _ := http.NewRequest("GET", "/user/testuser", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.UserResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, user.Username, response.Username)
	assert.Equal(t, user.DisplayName, response.DisplayName)
	assert.Equal(t, user.Bio, response.Bio)
	assert.Equal(t, *user.AvatarURL, *response.AvatarURL)
	assert.Equal(t, user.IsPrivate, response.IsPrivate)
}

func TestGetCurrentUserProfileHandler(t *testing.T) {
	if database.DB == nil {
		t.Skip("Database not available, skipping integration test")
		return
	}

	router := setupUserTestRouter()

	// Ensure test user exists
	user := models.User{
		Username:    "testuser",
		Password:    "password",
		DisplayName: "Test User",
		Bio:         "Test bio",
		Email:       stringPtr("test@example.com"),
		IsPrivate:   true,
	}
	database.DB.Where("username = ?", "testuser").FirstOrCreate(&user)

	// Test GET request
	req, _ := http.NewRequest("GET", "/user/profile", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.UserProfileResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, user.Username, response.Username)
	assert.Equal(t, *user.Email, *response.Email)
	assert.Equal(t, user.DisplayName, response.DisplayName)
	assert.Equal(t, user.Bio, response.Bio)
	assert.Equal(t, user.IsPrivate, response.IsPrivate)
	assert.NotEmpty(t, response.UpdatedAt)
}

func TestUpdateUserProfileHandler(t *testing.T) {
	if database.DB == nil {
		t.Skip("Database not available, skipping integration test")
		return
	}

	router := setupUserTestRouter()

	// Create test user
	user := models.User{
		Username:    "testuser",
		Password:    "password",
		DisplayName: "Original Display Name",
		Bio:         "Original bio",
		Email:       stringPtr("original@example.com"),
		IsPrivate:   false,
	}
	database.DB.Where("username = ?", "testuser").FirstOrCreate(&user)

	// Create request payload
	updateData := map[string]interface{}{
		"email":        "updated@example.com",
		"display_name": "Updated Display Name",
		"bio":          "Updated bio",
		"avatar_url":   "http://example.com/newavatar.jpg",
		"is_private":   true,
	}
	jsonData, _ := json.Marshal(updateData)

	// Test PUT request
	req, _ := http.NewRequest("PUT", "/user/profile", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.UserProfileResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "updated@example.com", *response.Email)
	assert.Equal(t, "Updated Display Name", response.DisplayName)
	assert.Equal(t, "Updated bio", response.Bio)
	assert.Equal(t, "http://example.com/newavatar.jpg", *response.AvatarURL)
	assert.Equal(t, true, response.IsPrivate)
}

func TestUpdateUserProfileHandler_InvalidEmail(t *testing.T) {
	if database.DB == nil {
		t.Skip("Database not available, skipping integration test")
		return
	}

	router := setupUserTestRouter()

	// Create request payload with invalid email
	updateData := map[string]interface{}{
		"email": "invalid-email",
	}
	jsonData, _ := json.Marshal(updateData)

	req, _ := http.NewRequest("PUT", "/user/profile", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response["error"], "invalid email format")
}

func TestDeleteUserAccountHandler(t *testing.T) {
	if database.DB == nil {
		t.Skip("Database not available, skipping integration test")
		return
	}

	router := setupUserTestRouter()

	// Create test user with hashed password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
	user := models.User{
		Username:    "deletetestuser",
		Password:    string(hashedPassword),
		DisplayName: "Delete Test User",
		Bio:         "Delete bio",
	}
	database.DB.Create(&user)

	// Temporarily modify router to use this specific username
	router.Use(func(c *gin.Context) {
		c.Set("username", "deletetestuser")
		c.Next()
	})

	// Create request payload with password
	deleteData := map[string]interface{}{
		"password": "correctpassword",
	}
	jsonData, _ := json.Marshal(deleteData)

	req, _ := http.NewRequest("DELETE", "/user/profile", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Account deleted successfully", response["message"])
}

func TestDeleteUserAccountHandler_InvalidPassword(t *testing.T) {
	if database.DB == nil {
		t.Skip("Database not available, skipping integration test")
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
	user := models.User{
		Username:    "wrongpasstestuser",
		Password:    string(hashedPassword),
		DisplayName: "Wrong Password User",
	}
	database.DB.Create(&user)

	// Create router with this user
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("username", "wrongpasstestuser")
		c.Next()
	})
	r.DELETE("/user/profile", DeleteUserAccount)

	deleteData := map[string]interface{}{
		"password": "wrongpassword",
	}
	jsonData, _ := json.Marshal(deleteData)

	req, _ := http.NewRequest("DELETE", "/user/profile", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Invalid password", response["error"])
}

func TestGetUserProfileHandler_UserNotFound(t *testing.T) {
	if database.DB == nil {
		t.Skip("Database not available, skipping integration test")
		return
	}

	router := setupUserTestRouter()

	req, _ := http.NewRequest("GET", "/user/nonexistentuser", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "User not found", response["error"])
}

func TestUpdateUserProfileHandler_DuplicateEmail(t *testing.T) {
	if database.DB == nil {
		t.Skip("Database not available, skipping integration test")
		return
	}

	// Create first user
	user1 := models.User{
		Username:    "user1",
		Password:    "password",
		Email:       stringPtr("taken@example.com"),
		DisplayName: "User One",
	}
	database.DB.Create(&user1)

	// Create second user
	user2 := models.User{
		Username:    "user2",
		Password:    "password",
		Email:       stringPtr("other@example.com"),
		DisplayName: "User Two",
	}
	database.DB.Create(&user2)

	// Try to update user2 to use the same email as user1
	updateData := map[string]interface{}{
		"email": "taken@example.com",
	}
	jsonData, _ := json.Marshal(updateData)

	req, _ := http.NewRequest("PUT", "/user/profile", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	setupUserTestRouter().ServeHTTP(w, req) // Create router inline

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response["error"], "email already taken")
}

// Helper functions for pointers
func stringPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}
