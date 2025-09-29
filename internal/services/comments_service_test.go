package services

import (
	"testing"

	"github.com/CodeAndCraft-Online/cortex-api/internal/database"
	"github.com/CodeAndCraft-Online/cortex-api/internal/models"
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

func TestGetCommentsByPostID_Service(t *testing.T) {
	if database.DB == nil {
		t.Skip("Database not available, skipping integration test")
		return
	}

	// Clear tables to avoid conflicts
	database.DB.Exec("DELETE FROM comments")
	database.DB.Exec("DELETE FROM posts")
	database.DB.Exec("DELETE FROM subs")
	database.DB.Exec("DELETE FROM users")

	// Create test data
	refreshToken1 := "refresh_commentuser"
	user := models.User{Username: "commentuser", Password: "password", RefreshToken: &refreshToken1}
	database.DB.Create(&user)

	sub := models.Sub{Name: "commentsub", OwnerID: user.ID}
	database.DB.Create(&sub)

	post := models.Post{
		Title:   "Test Post for Comments",
		Content: "Test content",
		UserID:  user.ID,
		SubID:   sub.ID,
	}
	database.DB.Create(&post)

	// Create a comment
	comment := models.Comment{
		Content: "Test comment",
		PostID:  post.ID,
		UserID:  user.ID,
	}
	database.DB.Create(&comment)

	// Test service
	comments, err := GetCommentsByPostID("1")

	assert.NoError(t, err)
	assert.NotNil(t, comments)
	assert.GreaterOrEqual(t, len(*comments), 1)
}

func TestCreateComment_Service(t *testing.T) {
	if database.DB == nil {
		t.Skip("Database not available, skipping integration test")
		return
	}

	// Create test data (do not clear tables, use existing or create)
	refreshToken2 := "refresh_createcommentuser"
	user := models.User{Username: "createcommentuser", Password: "password", RefreshToken: &refreshToken2}
	err := database.DB.Create(&user).Error
	if err != nil {
		// User may already exist, try to find
		database.DB.Where("username = ?", "createcommentuser").First(&user)
	}

	sub := models.Sub{Name: "createcommentsub", OwnerID: user.ID}
	err = database.DB.Create(&sub).Error
	if err != nil {
		// Sub may already exist, try to find
		database.DB.Where("name = ?", "createcommentsub").First(&sub)
	}

	post := models.Post{
		Title:   "Post for New Comment",
		Content: "Content",
		UserID:  user.ID,
		SubID:   sub.ID,
	}
	err = database.DB.Create(&post).Error
	if err != nil {
		// Post may already exist, try to find
		database.DB.Where("title = ? AND content = ?", "Post for New Comment", "Content").First(&post)
	}

	commentReq := models.CommentRequest{
		Content: "New comment content",
		PostID:  post.ID,
	}

	createdComment, err := CreateComment("createcommentuser", commentReq, post)

	assert.NoError(t, err)
	assert.NotNil(t, createdComment)
	assert.Equal(t, "New comment content", createdComment.Content)
	assert.Equal(t, post.ID, createdComment.PostID)
	assert.Equal(t, user.ID, createdComment.UserID)
}
