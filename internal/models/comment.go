package models

import "time"

// CommentResponse struct to format comment output
type CommentResponse struct {
	ID        uint    `json:"id"`
	Content   string  `json:"content"`
	Username  string  `json:"username"`
	ImageURL  *string `json:"imageURL,omitempty"` // Link to an image
	Upvotes   int     `json:"upvotes"`
	Downvotes int     `json:"downvotes"`
	CreatedAt string  `json:"created_at"`
}

// CommentRequest struct for incoming JSON data
type CommentRequest struct {
	PostID   uint    `json:"postID"`
	Content  string  `json:"content"`
	ImageURL *string `json:"imageURL,omitempty"` // Link to an image
}

type Comment struct {
	ID        uint `gorm:"primaryKey"`
	Content   string
	UserID    uint
	PostID    uint
	ImageURL  *string `json:"imageURL,omitempty"` // Link to an image
	User      User
	Post      Post
	CreatedAt time.Time
}
