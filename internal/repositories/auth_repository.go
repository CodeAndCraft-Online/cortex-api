package repositories

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	db "github.com/CodeAndCraft-Online/cortex-api/internal/database"
	models "github.com/CodeAndCraft-Online/cortex-api/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// IAuthRepository defines methods for authentication repository
type IAuthRepository interface {
	ResetPasswordRequest(username string) (*models.PasswordResetToken, error)
	ResetPassword(token, newPassword string) error
}

// AuthRepository implements IAuthRepository
type AuthRepository struct{}

// NewAuthRepository creates a new auth repository
func NewAuthRepository() IAuthRepository {
	return &AuthRepository{}
}

// Legacy global functions for backward compatibility
func ResetPasswordRequest(username string) (*models.PasswordResetToken, error) {
	repo := NewAuthRepository()
	return repo.ResetPasswordRequest(username)
}

func ResetPassword(token, newPassword string) error {
	repo := NewAuthRepository()
	return repo.ResetPassword(token, newPassword)
}

func (r *AuthRepository) ResetPasswordRequest(username string) (*models.PasswordResetToken, error) {

	// Find user by username
	var user models.User
	if err := db.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Generate a reset token
	token, err := GenerateToken()
	if err != nil {
		return nil, fmt.Errorf("could not generate reset token")
	}

	// Save token in the database
	resetToken := models.PasswordResetToken{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(15 * time.Minute), // Expires in 15 minutes
	}
	db.DB.Create(&resetToken)

	return &resetToken, nil
}

func (r *AuthRepository) ResetPassword(token, newPassword string) error {
	// Find the token
	var resetToken models.PasswordResetToken
	if err := db.DB.Where("token = ?", token).First(&resetToken).Error; err != nil {
		return fmt.Errorf("invalid or expired token")
	}

	// Check expiration
	if time.Now().After(resetToken.ExpiresAt) {
		return fmt.Errorf("reset token has expired")
	}

	// Find the user
	var user models.User
	if err := db.DB.First(&user, resetToken.UserID).Error; err != nil {
		return fmt.Errorf("user not found")
	}

	// Hash new password
	hashedPassword, err := HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password")
	}

	// Update password
	user.Password = hashedPassword
	db.DB.Save(&user)

	// Delete used reset token
	db.DB.Delete(&resetToken)
	return nil
}

// GenerateToken creates a random token
func GenerateToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
