package repositories

import (
	"testing"

	"github.com/CodeAndCraft-Online/cortex-api/internal/database"
	"github.com/CodeAndCraft-Online/cortex-api/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestGetUserByID(t *testing.T) {
	if !dbAvailable {
		t.Skip("Database not available, skipping repository integration tests")
		return
	}

	// Clear users table
	database.DB.Exec("DELETE FROM users")

	// Create test user
	user := models.User{
		Username:    "getuserbyid",
		Password:    "password",
		DisplayName: "Get User By ID",
		Bio:         "Test bio",
	}
	database.DB.Create(&user)

	t.Run("valid user ID", func(t *testing.T) {
		repo := NewUserRepository()
		result, err := repo.GetUserByID(user.ID)

		assert.NoError(t, err)
		assert.Equal(t, user.ID, result.ID)
		assert.Equal(t, user.Username, result.Username)
		assert.Equal(t, user.DisplayName, result.DisplayName)
		assert.Equal(t, user.Bio, result.Bio)
	})

	t.Run("invalid user ID", func(t *testing.T) {
		repo := NewUserRepository()
		result, err := repo.GetUserByID(9999)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "user not found")
	})
}

func TestGetUserByUsername(t *testing.T) {
	if !dbAvailable {
		t.Skip("Database not available, skipping repository integration tests")
		return
	}

	// Clear users table
	database.DB.Exec("DELETE FROM users")

	// Create test user
	user := models.User{
		Username:    "getuserbyusername",
		Password:    "password",
		DisplayName: "Get User By Username",
		Bio:         "Test bio",
	}
	database.DB.Create(&user)

	t.Run("valid username", func(t *testing.T) {
		repo := NewUserRepository()
		result, err := repo.GetUserByUsername(user.Username)

		assert.NoError(t, err)
		assert.Equal(t, user.Username, result.Username)
		assert.Equal(t, user.DisplayName, result.DisplayName)
		assert.Equal(t, user.Bio, result.Bio)
	})

	t.Run("invalid username", func(t *testing.T) {
		repo := NewUserRepository()
		result, err := repo.GetUserByUsername("nonexistentuser")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "user not found")
	})
}

func TestUpdateUser(t *testing.T) {
	if !dbAvailable {
		t.Skip("Database not available, skipping repository integration tests")
		return
	}

	// Clear users table
	database.DB.Exec("DELETE FROM users")

	// Create test user
	user := models.User{
		Username:    "updateuser",
		Password:    "password",
		DisplayName: "Update User",
		Bio:         "Original bio",
	}
	database.DB.Create(&user)

	t.Run("successful update", func(t *testing.T) {
		repo := NewUserRepository()
		updates := models.UserUpdateRequest{
			Email:       stringPtr("updated@example.com"),
			DisplayName: stringPtr("Updated User Name"),
			Bio:         stringPtr("Updated bio"),
			AvatarURL:   stringPtr("http://example.com/avatar.jpg"),
			IsPrivate:   boolPtr(true),
		}

		result, err := repo.UpdateUser(user.ID, updates)

		assert.NoError(t, err)
		assert.Equal(t, *updates.Email, *result.Email)
		assert.Equal(t, *updates.DisplayName, result.DisplayName)
		assert.Equal(t, *updates.Bio, result.Bio)
		assert.Equal(t, *updates.AvatarURL, *result.AvatarURL)
		assert.Equal(t, *updates.IsPrivate, result.IsPrivate)
		assert.NotEqual(t, user.UpdatedAt, result.UpdatedAt)
	})
}

func TestDeleteUser(t *testing.T) {
	if !dbAvailable {
		t.Skip("Database not available, skipping repository integration tests")
		return
	}

	// Clear users table
	database.DB.Exec("DELETE FROM users")

	// Create test user
	user := models.User{
		Username:    "deleteuser",
		Password:    "password",
		DisplayName: "Delete User",
		Bio:         "Test bio",
	}
	database.DB.Create(&user)

	t.Run("successful delete", func(t *testing.T) {
		repo := NewUserRepository()

		// Verify user exists
		_, err := repo.GetUserByID(user.ID)
		assert.NoError(t, err)

		// Delete user
		err = repo.DeleteUser(user.ID)
		assert.NoError(t, err)

		// Verify user is gone
		_, err = repo.GetUserByID(user.ID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user not found")
	})
}

// Helper functions for pointers
func stringPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}
