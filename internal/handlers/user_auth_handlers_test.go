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

// TestMain is defined in auth_handlers_test.go for the handlers package

// Helper function to create test router for user auth
func setupUserAuthTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.POST("/auth/register", Register)
	router.POST("/auth/login", Login)

	return router
}

func TestRegisterHandler(t *testing.T) {
	if database.DB == nil {
		t.Skip("Database not available, skipping integration test")
		return
	}

	// Setup - Clear users table
	database.DB.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE")

	router := setupUserAuthTestRouter()

	t.Run("successful registration", func(t *testing.T) {
		userData := map[string]string{
			"username": "newuser",
			"password": "securepassword",
		}
		jsonData, _ := json.Marshal(userData)

		req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "User registered successfully", response["message"])

		// Verify user was created in database
		var user models.User
		err := database.DB.Where("username = ?", "newuser").First(&user).Error
		assert.NoError(t, err)
		assert.Equal(t, "newuser", user.Username)

		// Verify password was hashed
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("securepassword"))
		assert.NoError(t, err)
	})

	t.Run("registration with missing username", func(t *testing.T) {
		userData := map[string]interface{}{
			"password": "securepassword",
		}
		jsonData, _ := json.Marshal(userData)

		req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("registration with missing password", func(t *testing.T) {
		userData := map[string]interface{}{
			"username": "anotheruser",
		}
		jsonData, _ := json.Marshal(userData)

		req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("registration with invalid JSON", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("registration with duplicate username", func(t *testing.T) {
		// First registration
		userData := map[string]string{
			"username": "duplicateuser",
			"password": "password",
		}
		jsonData1, _ := json.Marshal(userData)

		req1, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonData1))
		req1.Header.Set("Content-Type", "application/json")

		w1 := httptest.NewRecorder()
		router.ServeHTTP(w1, req1)

		assert.Equal(t, http.StatusOK, w1.Code)

		// Second registration with same username should fail due to unique constraint
		jsonData2, _ := json.Marshal(userData)
		req2, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonData2))
		req2.Header.Set("Content-Type", "application/json")

		w2 := httptest.NewRecorder()
		router.ServeHTTP(w2, req2)

		// This might return different status codes based on DB constraints, but should not be 200
		assert.NotEqual(t, http.StatusOK, w2.Code)
	})
}

func TestLoginHandler(t *testing.T) {
	if database.DB == nil {
		t.Skip("Database not available, skipping integration test")
		return
	}

	// Setup - Clear users table and create test user
	database.DB.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE")

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("testpass"), bcrypt.DefaultCost)
	testUser := models.User{
		Username: "testuser",
		Password: string(hashedPassword),
	}
	database.DB.Create(&testUser)

	router := setupUserAuthTestRouter()

	t.Run("successful login", func(t *testing.T) {
		loginData := map[string]string{
			"username": "testuser",
			"password": "testpass",
		}
		jsonData, _ := json.Marshal(loginData)

		req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Contains(t, response, "token")
		assert.NotEmpty(t, response["token"])
	})

	t.Run("login with non-existent user", func(t *testing.T) {
		loginData := map[string]string{
			"username": "nonexistent",
			"password": "password",
		}
		jsonData, _ := json.Marshal(loginData)

		req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "Invalid credentials", response["error"])
	})

	t.Run("login with wrong password", func(t *testing.T) {
		loginData := map[string]string{
			"username": "testuser",
			"password": "wrongpass",
		}
		jsonData, _ := json.Marshal(loginData)

		req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "Invalid credentials", response["error"])
	})

	t.Run("login with missing username", func(t *testing.T) {
		loginData := map[string]interface{}{
			"password": "testpass",
		}
		jsonData, _ := json.Marshal(loginData)

		req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("login with missing password", func(t *testing.T) {
		loginData := map[string]interface{}{
			"username": "testuser",
		}
		jsonData, _ := json.Marshal(loginData)

		req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("login with invalid JSON", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
