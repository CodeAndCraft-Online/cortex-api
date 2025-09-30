package models

import (
	"time"
)

// User represents a registered user account
type User struct {
	ID           uint      `gorm:"primaryKey"`
	Username     string    `gorm:"unique;not null" json:"username"`
	Password     string    `gorm:"not null" json:"-"`
	Email        *string   `gorm:"unique" json:"email,omitempty"`
	DisplayName  string    `gorm:"default:''" json:"display_name"`
	Bio          string    `gorm:"type:text" json:"bio"`
	AvatarURL    *string   `json:"avatar_url,omitempty"`
	IsPrivate    bool      `gorm:"default:false" json:"is_private"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	RefreshToken *string   `gorm:"unique" json:"-"`
	TokenExpires time.Time `json:"-"`
}

// UserResponse represents user data in API responses (excludes sensitive fields)
type UserResponse struct {
	ID          uint    `json:"id"`
	Username    string  `json:"username"`
	DisplayName string  `json:"display_name"`
	Bio         string  `json:"bio"`
	AvatarURL   *string `json:"avatar_url,omitempty"`
	IsPrivate   bool    `json:"is_private"`
	CreatedAt   string  `json:"created_at"`
}

// UserUpdateRequest represents data that can be updated by the user
type UserUpdateRequest struct {
	Email       *string `json:"email,omitempty"`
	DisplayName *string `json:"display_name,omitempty"`
	Bio         *string `json:"bio,omitempty"`
	AvatarURL   *string `json:"avatar_url,omitempty"`
	IsPrivate   *bool   `json:"is_private,omitempty"`
}

// UserProfileResponse represents the full profile view (includes private data for owner)
type UserProfileResponse struct {
	ID          uint    `json:"id"`
	Username    string  `json:"username"`
	Email       *string `json:"email,omitempty"`
	DisplayName string  `json:"display_name"`
	Bio         string  `json:"bio"`
	AvatarURL   *string `json:"avatar_url,omitempty"`
	IsPrivate   bool    `json:"is_private"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}
