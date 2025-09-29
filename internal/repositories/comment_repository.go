package repositories

import (
	"fmt"

	db "github.com/CodeAndCraft-Online/cortex-api/internal/database"
	models "github.com/CodeAndCraft-Online/cortex-api/internal/models"
)

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

	if err := db.DB.First(&post, commentReq.PostID).Error; err != nil {
		return nil, err
	}

	// Create the new comment with the correct user ID
	comment := models.Comment{
		PostID:   commentReq.PostID,
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
