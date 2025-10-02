package models

import "time"

// Response struct to format the output
type PostResponse struct {
	ID        uint              `json:"id"`
	Title     string            `json:"title"`
	Content   string            `json:"content"`
	ImageURL  *string           `json:"imageURL,omitempty"` // Link to an image
	Username  string            `json:"username"`
	Upvotes   int               `json:"upvotes"`
	Downvotes int               `json:"downvotes"`
	CreatedAt string            `json:"created_at"`
	SubID     uint              `json:"sub_id"`
	Comments  []CommentResponse `json:"comments"`
}

type Post struct {
	ID        uint    `gorm:"primaryKey"`
	Title     string  `json:"title" binding:"required"`
	SubID     uint    `json:"sub_id" binding:"required"`
	Content   string  `json:"content" binding:"required"`
	Upvotes   int     `json:"upvotes"`
	Downvotes int     `json:"downvotes"`
	ImageURL  *string `json:"imageURL,omitempty"` // Link to an image
	UserID    uint
	User      User
	CreatedAt time.Time
}
