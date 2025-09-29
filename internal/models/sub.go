package models

import "time"

// Sub represents community
type Sub struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"unique;not null"`
	Description string
	Private     bool `json:"private" gorm:"default:false"`
	OwnerID     uint `gorm:"not null"`
	Owner       User `gorm:"foreignKey:OwnerID"` // ✅ Define the relationship
	CreatedAt   time.Time
}

// SubInvitation represents an invitation to join a private sub
type SubInvitation struct {
	ID        uint   `gorm:"primaryKey"`
	SubID     uint   `gorm:"not null"`
	InviterID uint   `gorm:"not null"`
	InviteeID uint   `gorm:"not null"`
	Status    string `gorm:"default:'pending'"` // pending, accepted, rejected
	CreatedAt time.Time
	Invitee   User `gorm:"foreignKey:InviteeID"` // For preloading invitee data
}

// SubMembership tracks users who join subs
type SubMembership struct {
	ID       uint `gorm:"primaryKey"`
	SubID    uint `gorm:"not null"`
	UserID   uint `gorm:"not null"`
	JoinedAt time.Time
	User     User `gorm:"foreignKey:UserID"` // For preloading user data
}

// SubResponse struct for formatted output
type SubResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Owner       string `json:"owner"`
	Private     bool   `json:"private"`
	CreatedAt   string `json:"created_at"`
}

type SubRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Private     bool   `json:"private"`
}

type InviteRequest struct {
	InviteeUsername string `json:"invitee"`
}

// SubMemberResponse represents a sub member in API responses
type SubMemberResponse struct {
	Username string `json:"username"`
	JoinedAt string `json:"joined_at"`
}

// InviteResponse represents a pending invitation in API responses
type InviteResponse struct {
	InviteeUsername string `json:"invitee_username"`
	CreatedAt       string `json:"created_at"`
}
