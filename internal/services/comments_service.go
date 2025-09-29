package services

import (
	"errors"

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

func (s *CommentsService) GetCommentByID(commentID uint) (*models.CommentResponse, error) {
	comment, err := s.commentRepo.GetCommentByID(commentID)
	if err != nil {
		return nil, err
	}
	return comment, nil
}

func (s *CommentsService) UpdateComment(commentID uint, username string, commentReq models.CommentUpdateRequest) (*models.Comment, error) {
	// First get the comment to check ownership
	currentComment, err := s.commentRepo.GetCommentByID(commentID)
	if err != nil {
		return nil, err
	}

	// Check if user owns the comment
	if currentComment.Username != username {
		return nil, errors.New("unauthorized: can only edit own comments")
	}

	updatedComment, err := s.commentRepo.UpdateComment(commentID, commentReq)
	if err != nil {
		return nil, err
	}

	return updatedComment, nil
}

func (s *CommentsService) DeleteComment(commentID uint, username string) error {
	// First get the comment to check ownership
	currentComment, err := s.commentRepo.GetCommentByID(commentID)
	if err != nil {
		return err
	}

	// Check if user owns the comment
	if currentComment.Username != username {
		return errors.New("unauthorized: can only delete own comments")
	}

	return s.commentRepo.DeleteComment(commentID)
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
