package services

import (
	"errors"

	"github.com/CodeAndCraft-Online/cortex-api/internal/models"
	"github.com/CodeAndCraft-Online/cortex-api/internal/repositories"
)

// GetPostByID fetches a single post by ID
func ResetPasswordRequest(username string) (*models.PasswordResetToken, error) {
	user, err := repositories.ResetPasswordRequest(username)
	if err != nil {
		return nil, errors.New("post not found")
	}
	return user, nil
}

func ResetPassword(token, newPassword string) error {
	err := repositories.ResetPassword(token, newPassword)
	if err != nil {
		return err
	}

	return nil
}
