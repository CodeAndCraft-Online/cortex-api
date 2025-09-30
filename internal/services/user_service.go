// Package services handles user business logic
package services

import (
	"errors"

	"github.com/CodeAndCraft-Online/cortex-api/internal/models"
	"github.com/CodeAndCraft-Online/cortex-api/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

// UserService handles user business logic
type UserService struct {
	userRepo repositories.IUserRepository
}

// NewUserService creates a new user service with dependency injection
func NewUserService(userRepo repositories.IUserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// GetUserProfile returns the publicly visible profile for a user
func (s *UserService) GetUserProfile(username string, requestingUserID *uint) (*models.UserResponse, error) {
	user, err := s.userRepo.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	response := &models.UserResponse{
		ID:          user.ID,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Bio:         user.Bio,
		AvatarURL:   user.AvatarURL,
		IsPrivate:   user.IsPrivate,
		CreatedAt:   user.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	return response, nil
}

// GetUserProfileInternal returns the full profile (including private data) for the authenticated user
func (s *UserService) GetUserProfileInternal(userID uint) (*models.UserProfileResponse, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	response := &models.UserProfileResponse{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		DisplayName: user.DisplayName,
		Bio:         user.Bio,
		AvatarURL:   user.AvatarURL,
		IsPrivate:   user.IsPrivate,
		CreatedAt:   user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   user.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return response, nil
}

// UpdateUserProfile updates the authenticated user's profile
func (s *UserService) UpdateUserProfile(userID uint, updates models.UserUpdateRequest) (*models.UserProfileResponse, error) {
	// Validate updates
	if err := s.validateUserUpdates(updates); err != nil {
		return nil, err
	}

	// Update the user
	_, err := s.userRepo.UpdateUser(userID, updates)
	if err != nil {
		return nil, err
	}

	// Return the updated profile
	return s.GetUserProfileInternal(userID)
}

// DeleteUserAccount deletes the authenticated user's account
func (s *UserService) DeleteUserAccount(userID uint, password string) error {
	// First, verify the password for security
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return err
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return errors.New("invalid password")
	}

	// Delete the user
	return s.userRepo.DeleteUser(userID)
}

// ChangePassword allows users to change their password - TODO: Implement password change functionality
func (s *UserService) ChangePassword(userID uint, currentPassword, newPassword string) error {
	// TODO: Implement password change functionality
	// This would require:
	// 1. Adding a method to IUserRepository interface
	// 2. Implementing password hashing and validation
	// 3. Adding proper audit trail
	return errors.New("password change functionality not yet implemented")
}

// validateUserUpdates performs basic validation on profile updates
func (s *UserService) validateUserUpdates(updates models.UserUpdateRequest) error {
	if updates.DisplayName != nil && len(*updates.DisplayName) > 100 {
		return errors.New("display name must be 100 characters or less")
	}

	if updates.Bio != nil && len(*updates.Bio) > 500 {
		return errors.New("bio must be 500 characters or less")
	}

	if updates.Email != nil {
		// Basic email validation - could be enhanced with regex
		if len(*updates.Email) == 0 {
			return errors.New("email cannot be empty if provided")
		}
		// Check if email is already taken by another user would go here
	}

	return nil
}

// Legacy global functions for backward compatibility
func GetUserProfile(username string, requestingUserID *uint) (*models.UserResponse, error) {
	service := NewUserService(repositories.NewUserRepository())
	return service.GetUserProfile(username, requestingUserID)
}

func GetUserProfileInternal(userID uint) (*models.UserProfileResponse, error) {
	service := NewUserService(repositories.NewUserRepository())
	return service.GetUserProfileInternal(userID)
}

func UpdateUserProfile(userID uint, updates models.UserUpdateRequest) (*models.UserProfileResponse, error) {
	service := NewUserService(repositories.NewUserRepository())
	return service.UpdateUserProfile(userID, updates)
}

func DeleteUserAccount(userID uint, password string) error {
	service := NewUserService(repositories.NewUserRepository())
	return service.DeleteUserAccount(userID, password)
}

func ChangePassword(userID uint, currentPassword, newPassword string) error {
	service := NewUserService(repositories.NewUserRepository())
	return service.ChangePassword(userID, currentPassword, newPassword)
}
