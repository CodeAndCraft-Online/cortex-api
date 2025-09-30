package services

import (
	"testing"

	"github.com/CodeAndCraft-Online/cortex-api/internal/database"
	"github.com/CodeAndCraft-Online/cortex-api/internal/models"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestGetUserProfile(t *testing.T) {
	if database.DB == nil {
		t.Skip("Database not available, skipping integration test")
		return
	}

	// Clear users table
	database.DB.Exec("DELETE FROM users")

	// Create test user
	user := models.User{
		Username:    "profileuser",
		Password:    "password",
		DisplayName: "Profile User",
		Bio:         "Profile bio",
		Email:       stringPtr("profile@example.com"),
		AvatarURL:   stringPtr("http://example.com/avatar.jpg"),
		IsPrivate:   false,
	}
	database.DB.Create(&user)

	t.Run("successful profile retrieval", func(t *testing.T) {
		profile, err := GetUserProfile(user.Username, nil)

		assert.NoError(t, err)
		assert.Equal(t, user.Username, profile.Username)
		assert.Equal(t, user.DisplayName, profile.DisplayName)
		assert.Equal(t, user.Bio, profile.Bio)
		assert.Equal(t, *user.AvatarURL, *profile.AvatarURL)
		assert.Equal(t, user.IsPrivate, profile.IsPrivate)
	})

	t.Run("user not found", func(t *testing.T) {
		profile, err := GetUserProfile("nonexistentuser", nil)

		assert.Error(t, err)
		assert.Nil(t, profile)
	})
}

func TestGetUserProfileInternal(t *testing.T) {
	if database.DB == nil {
		t.Skip("Database not available, skipping integration test")
		return
	}

	// Clear users table
	database.DB.Exec("DELETE FROM users")

	// Create test user
	user := models.User{
		Username:    "internaluser",
		Password:    "password",
		DisplayName: "Internal User",
		Bio:         "Internal bio",
		Email:       stringPtr("internal@example.com"),
		AvatarURL:   stringPtr("http://example.com/avatar.jpg"),
		IsPrivate:   true,
	}
	database.DB.Create(&user)

	t.Run("successful internal profile retrieval", func(t *testing.T) {
		profile, err := GetUserProfileInternal(user.ID)

		assert.NoError(t, err)
		assert.Equal(t, user.Username, profile.Username)
		assert.Equal(t, *user.Email, *profile.Email)
		assert.Equal(t, user.DisplayName, profile.DisplayName)
		assert.Equal(t, user.Bio, profile.Bio)
		assert.Equal(t, *user.AvatarURL, *profile.AvatarURL)
		assert.Equal(t, user.IsPrivate, profile.IsPrivate)
		assert.NotEmpty(t, profile.UpdatedAt)
	})
}

func TestUpdateUserProfile(t *testing.T) {
	if database.DB == nil {
		t.Skip("Database not available, skipping integration test")
		return
	}

	// Clear users table
	database.DB.Exec("DELETE FROM users")

	// Create test user
	user := models.User{
		Username:    "updateprofileuser",
		Password:    "password",
		DisplayName: "Update Profile User",
		Bio:         "Original bio",
		Email:       stringPtr("original@example.com"),
		IsPrivate:   false,
	}
	database.DB.Create(&user)

	t.Run("successful profile update", func(t *testing.T) {
		updates := models.UserUpdateRequest{
			Email:       stringPtr("updated@example.com"),
			DisplayName: stringPtr("Updated Display Name"),
			Bio:         stringPtr("Updated bio"),
			AvatarURL:   stringPtr("http://example.com/newavatar.jpg"),
			IsPrivate:   boolPtr(true),
		}

		result, err := UpdateUserProfile(user.ID, updates)

		assert.NoError(t, err)
		assert.Equal(t, *updates.Email, *result.Email)
		assert.Equal(t, *updates.DisplayName, result.DisplayName)
		assert.Equal(t, *updates.Bio, result.Bio)
		assert.Equal(t, *updates.AvatarURL, *result.AvatarURL)
		assert.Equal(t, *updates.IsPrivate, result.IsPrivate)
	})

	t.Run("invalid email format", func(t *testing.T) {
		updates := models.UserUpdateRequest{
			Email: stringPtr("invalid-email"),
		}

		result, err := UpdateUserProfile(user.ID, updates)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid email format")
	})

	t.Run("duplicate email", func(t *testing.T) {
		// Create another user with different email first
		user2 := models.User{
			Username: "duplicateuser",
			Password: "password",
			Email:    stringPtr("different@example.com"),
		}
		database.DB.Create(&user2)

		updates := models.UserUpdateRequest{
			Email: stringPtr("different@example.com"), // Same as user2
		}

		result, err := UpdateUserProfile(user.ID, updates)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "email already taken")
	})
}

func TestDeleteUserAccount(t *testing.T) {
	if database.DB == nil {
		t.Skip("Database not available, skipping integration test")
		return
	}

	// Clear users table
	database.DB.Exec("DELETE FROM users")

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("originalpassword"), bcrypt.DefaultCost)

	// Create test user
	user := models.User{
		Username:    "deleteaccountuser",
		Password:    string(hashedPassword),
		DisplayName: "Delete Account User",
		Bio:         "Delete bio",
	}
	database.DB.Create(&user)

	t.Run("successful account deletion", func(t *testing.T) {
		// Verify user exists
		profile, err := GetUserProfileInternal(user.ID)
		assert.NoError(t, err)
		assert.NotNil(t, profile)

		// Delete account
		err = DeleteUserAccount(user.ID, "originalpassword")
		assert.NoError(t, err)

		// Verify user is gone
		profile, err = GetUserProfileInternal(user.ID)
		assert.Error(t, err)
		assert.Nil(t, profile)
	})

	t.Run("invalid password", func(t *testing.T) {
		// Create another user for this test
		hashedPassword2, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
		user2 := models.User{
			Username: "invalidpassuser",
			Password: string(hashedPassword2),
		}
		database.DB.Create(&user2)
		defer database.DB.Unscoped().Delete(&user2)

		// Try to delete with wrong password
		err := DeleteUserAccount(user2.ID, "wrongpassword")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid password")
	})
}

// Helper functions for pointers
func stringPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}
