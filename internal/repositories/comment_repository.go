package repositories

import (
	"fmt"
	"time"

	db "github.com/CodeAndCraft-Online/cortex-api/internal/database"
	models "github.com/CodeAndCraft-Online/cortex-api/internal/models"
)

// ICommentRepository defines methods for comment repository
type ICommentRepository interface {
	GetCommentsByPostID(postID string) (*[]models.CommentResponse, error)
	GetCommentByID(commentID uint) (*models.CommentResponse, error)
	CreateComment(username string, commentReq models.CommentRequest, post models.Post) (*models.Comment, error)
	UpdateComment(commentID uint, commentReq models.CommentUpdateRequest) (*models.Comment, error)
	DeleteComment(commentID uint) error
}

// CommentRepository implements ICommentRepository
type CommentRepository struct{}

// NewCommentRepository creates a new comment repository
func NewCommentRepository() ICommentRepository {
	return &CommentRepository{}
}

func (r *CommentRepository) GetCommentsByPostID(postID string) (*[]models.CommentResponse, error) {
	var comments []models.Comment

	// Fetch comments and preload user details
	if err := db.DB.Preload("User").Where("post_id = ?", postID).Order("created_at DESC").Find(&comments).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch comments")
	}

	// Format response to exclude sensitive data
	var formattedComments []models.CommentResponse
	for _, comment := range comments {
		response := models.CommentResponse{
			ID:        comment.ID,
			Content:   comment.Content,
			ImageURL:  comment.ImageURL,
			ParentID:  comment.ParentID,
			Username:  comment.User.Username, // Include only username, not full User object
			CreatedAt: comment.CreatedAt.Format("2006-01-02 15:04:05"),
		}

		if comment.UpdatedAt != nil {
			updatedAtStr := comment.UpdatedAt.Format("2006-01-02 15:04:05")
			response.UpdatedAt = updatedAtStr
		}

		formattedComments = append(formattedComments, response)
	}

	return &formattedComments, nil
}

func (r *CommentRepository) GetCommentByID(commentID uint) (*models.CommentResponse, error) {
	var comment models.Comment

	// Fetch comment and preload user details
	if err := db.DB.Preload("User").First(&comment, commentID).Error; err != nil {
		return nil, fmt.Errorf("comment not found")
	}

	// Format response
	response := models.CommentResponse{
		ID:        comment.ID,
		Content:   comment.Content,
		ImageURL:  comment.ImageURL,
		ParentID:  comment.ParentID,
		Username:  comment.User.Username,
		CreatedAt: comment.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	if comment.UpdatedAt != nil {
		updatedAtStr := comment.UpdatedAt.Format("2006-01-02 15:04:05")
		response.UpdatedAt = updatedAtStr
	}

	return &response, nil
}

func (r *CommentRepository) CreateComment(username string, commentReq models.CommentRequest, post models.Post) (*models.Comment, error) {

	// Fetch user ID from the database based on username
	var user models.User
	if err := db.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}

	// Create the new comment with the correct user ID
	comment := models.Comment{
		PostID:   post.ID,
		ParentID: commentReq.ParentID,
		Content:  commentReq.Content,
		ImageURL: commentReq.ImageURL,
		UserID:   user.ID, // Assign the authenticated user's ID
	}

	// Save the comment to the database
	if err := db.DB.Create(&comment).Error; err != nil {
		return nil, err
	}

	return &comment, nil
}

func (r *CommentRepository) UpdateComment(commentID uint, commentReq models.CommentUpdateRequest) (*models.Comment, error) {
	var comment models.Comment

	// Find the comment
	if err := db.DB.First(&comment, commentID).Error; err != nil {
		return nil, fmt.Errorf("comment not found")
	}

	// Update fields
	now := time.Now()
	comment.Content = commentReq.Content
	comment.ImageURL = commentReq.ImageURL
	comment.UpdatedAt = &now

	// Save updates
	if err := db.DB.Save(&comment).Error; err != nil {
		return nil, err
	}

	return &comment, nil
}

func (r *CommentRepository) DeleteComment(commentID uint) error {
	// Delete the comment
	if err := db.DB.Delete(&models.Comment{}, commentID).Error; err != nil {
		return fmt.Errorf("failed to delete comment")
	}

	return nil
}

func GetCommentsByPostID(postID string) (*[]models.CommentResponse, error) {
	var comments []models.Comment

	// Fetch comments and preload user details
	if err := db.DB.Preload("User").Where("post_id = ?", postID).Order("created_at DESC").Find(&comments).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch comments")
	}

	// Format response to exclude sensitive data
	var formattedComments []models.CommentResponse
	for _, comment := range comments {
		formattedComments = append(formattedComments, models.CommentResponse{
			ID:        comment.ID,
			Content:   comment.Content,
			ImageURL:  comment.ImageURL,
			Username:  comment.User.Username, // Include only username, not full User object
			CreatedAt: comment.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &formattedComments, nil
}

func CreateComment(username string, commentReq models.CommentRequest, post models.Post) (*models.Comment, error) {

	// Fetch user ID from the database based on username
	var user models.User
	if err := db.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}

	// Create the new comment with the correct user ID
	comment := models.Comment{
		PostID:   post.ID,
		Content:  commentReq.Content,
		ImageURL: commentReq.ImageURL,
		UserID:   user.ID, // Assign the authenticated user's ID
	}

	// Save the comment to the database
	if err := db.DB.Create(&comment).Error; err != nil {
		return nil, err
	}

	return &comment, nil
}
