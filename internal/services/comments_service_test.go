package services

import (
	"testing"

	"github.com/CodeAndCraft-Online/cortex-api/internal/database"
	"github.com/CodeAndCraft-Online/cortex-api/internal/models"
	"github.com/CodeAndCraft-Online/cortex-api/internal/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCommentRepository is a mock implementation of ICommentRepository
type MockCommentRepository struct {
	mock.Mock
}

func (m *MockCommentRepository) GetCommentsByPostID(postID string) (*[]models.CommentResponse, error) {
	args := m.Called(postID)
	return args.Get(0).(*[]models.CommentResponse), args.Error(1)
}

func (m *MockCommentRepository) CreateComment(username string, commentReq models.CommentRequest, post models.Post) (*models.Comment, error) {
	args := m.Called(username, commentReq, post)
	return args.Get(0).(*models.Comment), args.Error(1)
}

func (m *MockCommentRepository) GetCommentByID(commentID uint) (*models.CommentResponse, error) {
	args := m.Called(commentID)
	return args.Get(0).(*models.CommentResponse), args.Error(1)
}

func (m *MockCommentRepository) UpdateComment(commentID uint, commentReq models.CommentUpdateRequest) (*models.Comment, error) {
	args := m.Called(commentID, commentReq)
	return args.Get(0).(*models.Comment), args.Error(1)
}

func (m *MockCommentRepository) DeleteComment(commentID uint) error {
	args := m.Called(commentID)
	return args.Error(0)
}

// Service layer tests - basic functionality
func TestGetCommentsByPostID_Service(t *testing.T) {
	if database.DB == nil {
		t.Skip("Database not available, skipping integration test")
		return
	}

	// Test with real database
	service := NewCommentsService(repositories.NewCommentRepository())
	comments, err := service.GetCommentsByPostID("1") // Test post ID

	if err != nil {
		// No data, but service should not error
		assert.NoError(t, err)
		return
	}

	assert.NotNil(t, comments)
	// Test pagination/filtering if comments exist
}

// Basic service integration test
func TestCommentsService_BasicFunctionality(t *testing.T) {
	if database.DB == nil {
		t.Skip("Database not available, skipping integration test")
		return
	}

	// Just test that the service can be instantiated
	service := NewCommentsService(repositories.NewCommentRepository())
	assert.NotNil(t, service)

	// Test GetCommentByID with non-existent ID (should handle gracefully)
	_, err := service.GetCommentByID(999)
	// Error is expected for non-existent comment
	assert.Error(t, err)
}
