package repositories

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/CodeAndCraft-Online/cortex-api/internal/database"
	"github.com/CodeAndCraft-Online/cortex-api/internal/models"
	"github.com/CodeAndCraft-Online/cortex-api/internal/testutils"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func TestMain(m *testing.M) {
	// Try to set up test DB - if Docker not available, skip all tests in this package
	db, teardown, err := testutils.SetupTestDB()
	if err != nil {
		log.Printf("Docker not available, skipping all repository tests: %v", err)
		return // Skip all tests in this package
	}
	defer teardown()

	database.DB = db
	os.Exit(m.Run())
}

func TestGenerateToken(t *testing.T) {
	token, err := GenerateToken()

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.Len(t, token, 64) // 32 bytes = 64 hex characters
}

func TestHashPassword(t *testing.T) {
	password := "testpassword123"
	hash, err := HashPassword(password)

	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, password, hash)

	// Verify the hash can be checked
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	assert.NoError(t, err)
}

func TestResetPasswordIntegration(t *testing.T) {
	// Create test user
	user := models.User{
		Username: "testuser",
		Password: "oldhash",
	}
	err := database.DB.Create(&user).Error
	assert.NoError(t, err)

	// Test reset password request
	token, err := ResetPasswordRequest("testuser")
	assert.NoError(t, err)
	assert.NotNil(t, token)

	// Verify token was created
	var dbToken models.PasswordResetToken
	err = database.DB.Where("user_id = ?", user.ID).First(&dbToken).Error
	assert.NoError(t, err)
	assert.Equal(t, token.Token, dbToken.Token)

	// Test password reset
	err = ResetPassword(dbToken.Token, "newpassword123")
	assert.NoError(t, err)

	// Verify password was updated
	var updatedUser models.User
	database.DB.First(&updatedUser, user.ID)
	err = bcrypt.CompareHashAndPassword([]byte(updatedUser.Password), []byte("newpassword123"))
	assert.NoError(t, err)

	// Verify token was deleted
	err = database.DB.First(&models.PasswordResetToken{}, dbToken.ID).Error
	assert.Error(t, err) // Should not find the token
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestResetPasswordRequest_UserNotFound(t *testing.T) {
	token, err := ResetPasswordRequest("nonexistentuser")

	assert.Error(t, err)
	assert.Nil(t, token)
	assert.Contains(t, err.Error(), "user not found")
}

func TestResetPassword_InvalidToken(t *testing.T) {
	err := ResetPassword("invalidtoken", "newpassword")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid or expired token")
}

func TestResetPassword_ExpiredToken(t *testing.T) {
	// Create test user
	user := models.User{
		Username: "expireduser",
		Password: "oldhash",
	}
	database.DB.Create(&user)

	// Create expired token
	expiredToken := models.PasswordResetToken{
		UserID:    user.ID,
		Token:     "expiredtoken",
		ExpiresAt: time.Now().Add(-1 * time.Hour), // Expired
	}
	database.DB.Create(&expiredToken)

	// Test reset with expired token
	err := ResetPassword("expiredtoken", "newpassword")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "reset token has expired")

	// Clean up
	database.DB.Unscoped().Delete(&user)
	database.DB.Unscoped().Delete(&expiredToken)
}
