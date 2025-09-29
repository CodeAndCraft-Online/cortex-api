package services

import (
	"github.com/CodeAndCraft-Online/cortex-api/internal/models"
	"github.com/CodeAndCraft-Online/cortex-api/internal/repositories"
)

func GetCommentsByPostID(postID string) (*[]models.CommentResponse, error) {
	commentResponse, err := repositories.GetCommentsByPostID(postID)
	if err != nil {
		return nil, err
	}

	return commentResponse, nil
}

func CreateComment(username string, commentReq models.CommentRequest, post models.Post) (*models.Comment, error) {
	comment, err := repositories.CreateComment(username, commentReq, post)
	if err != nil {
		return nil, err
	}

	return comment, nil
}
