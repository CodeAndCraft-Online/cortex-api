package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUser_JSONSerialization(t *testing.T) {
	// Create a test user
	refreshToken := "refreshtoken123"
	user := User{
		ID:           1,
		Username:     "testuser",
		Password:     "hashedpassword",
		CreatedAt:    time.Now(),
		RefreshToken: &refreshToken,
		TokenExpires: time.Now().Add(24 * time.Hour),
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(user)
	assert.NoError(t, err)
	assert.NotEmpty(t, jsonData)

	// Verify JSON contains expected fields
	var jsonMap map[string]interface{}
	err = json.Unmarshal(jsonData, &jsonMap)
	assert.NoError(t, err)

	assert.Equal(t, float64(1), jsonMap["ID"])
	assert.Equal(t, "testuser", jsonMap["Username"])
	assert.Equal(t, "hashedpassword", jsonMap["Password"])
	assert.Contains(t, jsonMap, "CreatedAt")
	assert.Equal(t, "refreshtoken123", jsonMap["RefreshToken"])
	assert.Contains(t, jsonMap, "TokenExpires")
}

func TestUser_JSONUnmarshaling(t *testing.T) {
	jsonStr := `{
		"ID": 1,
		"Username": "testuser",
		"Password": "hashedpassword",
		"CreatedAt": "2023-01-01T00:00:00Z",
		"RefreshToken": "refreshtoken123",
		"TokenExpires": "2023-01-02T00:00:00Z"
	}`

	var user User
	err := json.Unmarshal([]byte(jsonStr), &user)
	assert.NoError(t, err)

	assert.Equal(t, uint(1), user.ID)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "hashedpassword", user.Password)
	assert.Equal(t, "refreshtoken123", *user.RefreshToken)

	// Verify timestamps were parsed
	expectedCreatedAt, _ := time.Parse(time.RFC3339, "2023-01-01T00:00:00Z")
	expectedTokenExpires, _ := time.Parse(time.RFC3339, "2023-01-02T00:00:00Z")
	assert.True(t, expectedCreatedAt.Equal(user.CreatedAt))
	assert.True(t, expectedTokenExpires.Equal(user.TokenExpires))
}

func TestUser_EmptyStruct(t *testing.T) {
	var user User

	// Test zero values
	assert.Equal(t, uint(0), user.ID)
	assert.Equal(t, "", user.Username)
	assert.Equal(t, "", user.Password)
	assert.True(t, user.CreatedAt.IsZero())
	assert.Nil(t, user.RefreshToken)
	assert.True(t, user.TokenExpires.IsZero())
}

func TestUser_StructFields(t *testing.T) {
	token := "token"
	user := User{
		ID:           1,
		Username:     "testuser",
		Password:     "password",
		RefreshToken: &token,
	}

	// Test field assignments
	assert.Equal(t, uint(1), user.ID)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "password", user.Password)
	assert.Equal(t, "token", *user.RefreshToken)
}
