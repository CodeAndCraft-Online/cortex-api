package models

import "time"

type User struct {
	ID           uint   `gorm:"primaryKey"`
	Username     string `gorm:"unique"`
	Password     string
	CreatedAt    time.Time
	RefreshToken string `gorm:"unique"`
	TokenExpires time.Time
}
