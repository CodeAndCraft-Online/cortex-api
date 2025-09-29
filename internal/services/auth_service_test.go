package services

import (
	"testing"

	"github.com/CodeAndCraft-Online/cortex-api/internal/database"
	"github.com/CodeAndCraft-Online/cortex-api/internal/models"
	"github.com/CodeAndCraft-Online/cortex-api/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// Try to set up test DB - if Docker not available, skip all tests
	db, teardown, err := testutils.SetupTestDB()
	if err != nil {
		println("Docker not available, skipping service integration tests:", err.Error())
		return // Skip all tests
	}

	database.DB = db
	m.Run()
	teardown()
}

// Note: Service tests would require Docker for database access

func TestResetPasswordRequest_Service(t *testing.T) {
	// Create test user
	user := models.User{
		Username: "serviceuser",
		Password: "password",
	}
	database.DB.Create(&user)

	// Test service layer
	token, err := ResetPasswordRequest("serviceuser")

	assert.NoError(t, err)
	assert.NotNil(t, token)
	assert.Equal(t, uint(1), token.UserID) // Assuming ID is 1
	assert.NotEmpty(t, token.Token)
}

func TestResetPassword_Service(t *testing.T) {
	// Create test user and token
	user := models.User{
		Username: "resetuser",
		Password: "oldhash",
	}
	database.DB.Create(&user)

	token := models.PasswordResetToken{
		UserID: 1,
		Token:  "Servicetoken",
	}
	database.DB.Create(&token)

	// Test service layer
	err := ResetPassword("Servicetoken", "newpassword")

	assert.NoError(t, err)

	// Verify password was updated
	var updatedUser models.User
	database.DB.First(&updatedUser, user.ID)
	assert.NotEqual(t, "oldhash", updatedUser.Password) // Password should be hashed
}
