package models

type Vote struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint
	PostID    uint
	CommentID *uint `gorm:"default:null"` // Nullable, used for comment votes
	Vote      int   // 1 for upvote, -1 for downvote
}
