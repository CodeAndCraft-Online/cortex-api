// Package repositories handles user data operations
package repositories

import (
	"fmt"
	"time"

	db "github.com/CodeAndCraft-Online/cortex-api/internal/database"
	"github.com/CodeAndCraft-Online/cortex-api/internal/models"
)

// IUserRepository defines methods for user repository operations
type IUserRepository interface {
	GetUserByID(id uint) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	UpdateUser(id uint, updates models.UserUpdateRequest) (*models.User, error)
	DeleteUser(id uint) error
}

// UserRepository implements IUserRepository
type UserRepository struct{}

// NewUserRepository creates a new user repository instance
func NewUserRepository() IUserRepository {
	return &UserRepository{}
}

// Legacy global functions for backward compatibility

// GetUserByID retrieves a user by their ID
func GetUserByID(id uint) (*models.User, error) {
	repo := NewUserRepository()
	return repo.GetUserByID(id)
}

// GetUserByUsername retrieves a user by their username
func GetUserByUsername(username string) (*models.User, error) {
	repo := NewUserRepository()
	return repo.GetUserByUsername(username)
}

// UpdateUser updates user profile information
func UpdateUser(id uint, updates models.UserUpdateRequest) (*models.User, error) {
	repo := NewUserRepository()
	return repo.UpdateUser(id, updates)
}

// DeleteUser permanently deletes a user account
func DeleteUser(id uint) error {
	repo := NewUserRepository()
	return repo.DeleteUser(id)
}

// GetUserByID implementation
func (r *UserRepository) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := db.DB.First(&user, id).Error; err != nil {
		return nil, fmt.Errorf("user not found")
	}
	return &user, nil
}

// GetUserByUsername implementation
func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	if err := db.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, fmt.Errorf("user not found")
	}
	return &user, nil
}

// UpdateUser implementation
func (r *UserRepository) UpdateUser(id uint, updates models.UserUpdateRequest) (*models.User, error) {
	// First, get the current user
	user, err := r.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	// Apply updates only for non-nil fields
	if updates.Email != nil {
		user.Email = updates.Email
	}
	if updates.DisplayName != nil {
		user.DisplayName = *updates.DisplayName
	}
	if updates.Bio != nil {
		user.Bio = *updates.Bio
	}
	if updates.AvatarURL != nil {
		user.AvatarURL = updates.AvatarURL
	}
	if updates.IsPrivate != nil {
		user.IsPrivate = *updates.IsPrivate
	}

	// Update the timestamp
	user.UpdatedAt = time.Now()

	// Save the user
	if err := db.DB.Save(user).Error; err != nil {
		return nil, fmt.Errorf("failed to update user: %v", err)
	}

	return user, nil
}

// DeleteUser implementation
func (r *UserRepository) DeleteUser(id uint) error {
	// Verify user exists first
	_, err := r.GetUserByID(id)
	if err != nil {
		return err
	}

	// Delete the user (cascade delete will handle related records)
	if err := db.DB.Delete(&models.User{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}

	return nil
}
