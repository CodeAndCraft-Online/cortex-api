package services

import (
	"errors"

	"github.com/CodeAndCraft-Online/cortex-api/internal/models"
	"github.com/CodeAndCraft-Online/cortex-api/internal/repositories"
)

// AuthService handles authentication business logic
type AuthService struct {
	authRepo repositories.IAuthRepository
}

// NewAuthService creates a new auth service with dependency injection
func NewAuthService(authRepo repositories.IAuthRepository) *AuthService {
	return &AuthService{
		authRepo: authRepo,
	}
}

// GetPostByID fetches a single post by ID
func (s *AuthService) ResetPasswordRequest(username string) (*models.PasswordResetToken, error) {
	user, err := s.authRepo.ResetPasswordRequest(username)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (s *AuthService) ResetPassword(token, newPassword string) error {
	err := s.authRepo.ResetPassword(token, newPassword)
	if err != nil {
		return err
	}

	return nil
}

// Legacy global functions for backward compatibility
func ResetPasswordRequest(username string) (*models.PasswordResetToken, error) {
	service := NewAuthService(repositories.NewAuthRepository())
	return service.ResetPasswordRequest(username)
}

func ResetPassword(token, newPassword string) error {
	service := NewAuthService(repositories.NewAuthRepository())
	return service.ResetPassword(token, newPassword)
}
