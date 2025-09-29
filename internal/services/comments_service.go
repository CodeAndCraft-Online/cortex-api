package services

import (
	"github.com/CodeAndCraft-Online/cortex-api/internal/models"
	"github.com/CodeAndCraft-Online/cortex-api/internal/repositories"
)

// CommentsService handles comment business logic
type CommentsService struct {
	commentRepo repositories.ICommentRepository
}

// NewCommentsService creates a new comments service with dependency injection
func NewCommentsService(commentRepo repositories.ICommentRepository) *CommentsService {
	return &CommentsService{
		commentRepo: commentRepo,
	}
}

func (s *CommentsService) GetCommentsByPostID(postID string) (*[]models.CommentResponse, error) {
	commentResponse, err := s.commentRepo.GetCommentsByPostID(postID)
	if err != nil {
		return nil, err
	}

	return commentResponse, nil
}

func (s *CommentsService) CreateComment(username string, commentReq models.CommentRequest, post models.Post) (*models.Comment, error) {
	comment, err := s.commentRepo.CreateComment(username, commentReq, post)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

// Legacy global functions for backward compatibility
func GetCommentsByPostID(postID string) (*[]models.CommentResponse, error) {
	service := NewCommentsService(repositories.NewCommentRepository())
	return service.GetCommentsByPostID(postID)
}

func CreateComment(username string, commentReq models.CommentRequest, post models.Post) (*models.Comment, error) {
	service := NewCommentsService(repositories.NewCommentRepository())
	return service.CreateComment(username, commentReq, post)
}
