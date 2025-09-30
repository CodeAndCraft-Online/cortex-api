package models

import "time"

// CommentResponse struct to format comment output
type CommentResponse struct {
	ID        uint    `json:"id"`
	Content   string  `json:"content"`
	Username  string  `json:"username"`
	ImageURL  *string `json:"imageURL,omitempty"` // Link to an image
	ParentID  *uint   `json:"parentID,omitempty"` // For comment threading
	Upvotes   int     `json:"upvotes"`
	Downvotes int     `json:"downvotes"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at,omitempty"`
}

// CommentRequest struct for incoming JSON data
type CommentRequest struct {
	PostID   uint    `json:"postID"`
	ParentID *uint   `json:"parentID,omitempty"` // For replying to comments
	Content  string  `json:"content"`
	ImageURL *string `json:"imageURL,omitempty"` // Link to an image
}

// CommentUpdateRequest struct for updating comments
type CommentUpdateRequest struct {
	Content  string  `json:"content"`
	ImageURL *string `json:"imageURL,omitempty"`
}

type Comment struct {
	ID        uint `gorm:"primaryKey"`
	Content   string
	UserID    uint
	PostID    uint
	ParentID  *uint   `json:",omitempty"`         // For comment threading
	ImageURL  *string `json:"imageURL,omitempty"` // Link to an image
	User      User
	Post      Post
	CreatedAt time.Time
	UpdatedAt *time.Time `json:",omitempty"`
}
