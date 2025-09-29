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

// TestGetCommentByID_Service tests the GetCommentByID service method
func TestGetCommentByID_Service(t *testing.T) {
	if database.DB == nil {
		t.Skip("Database not available, skipping service integration tests")
		return
	}

	// Create test user
	refreshToken3 := "refresh_getcommentuser"
	user := models.User{Username: "getcommentuser", Password: "password", RefreshToken: &refreshToken3}
	err := database.DB.Create(&user).Error
	if err != nil {
		database.DB.Where("username = ?", "getcommentuser").First(&user)
	}

	// Create test sub and post
	sub := models.Sub{Name: "getcommentsub", OwnerID: user.ID}
	err = database.DB.Create(&sub).Error
	post := models.Post{
		Title:   "Post for GetCommentByID Test",
		Content: "Post content",
		UserID:  user.ID,
		SubID:   sub.ID,
	}
	err = database.DB.Create(&post).Error

	// Create test comment
	comment := models.Comment{
		Content:  "Test comment for service",
		PostID:   post.ID,
		UserID:   user.ID,
		ParentID: nil,
	}
	err = database.DB.Create(&comment).Error

	t.Run("get comment by ID", func(t *testing.T) {
		service := NewCommentsService(NewCommentRepository())
		retrievedComment, err := service.GetCommentByID(comment.ID)

		assert.NoError(t, err)
		assert.NotNil(t, retrievedComment)
		assert.Equal(t, comment.ID, retrievedComment.ID)
		assert.Equal(t, comment.Content, retrievedComment.Content)
		assert.Equal(t, "getcommentuser", retrievedComment.Username)
	})
}

// TestUpdateComment_Service tests the UpdateComment service method
func TestUpdateComment_Service(t *testing.T) {
	if database.DB == nil {
		t.Skip("Database not available, skipping service integration tests")
		return
	}

	// Create test user
	refreshToken4 := "refresh_updatecommentuser"
	user := models.User{Username: "updatecommentuser", Password: "password", RefreshToken: &refreshToken4}
	err := database.DB.Create(&user).Error

	// Create test sub and post
	sub := models.Sub{Name: "updatecommentsub", OwnerID: user.ID}
	err = database.DB.Create(&sub).Error

	post := models.Post{
		Title:   "Post for UpdateComment Test",
		Content: "Post content",
		UserID:  user.ID,
		SubID:   sub.ID,
	}
	err = database.DB.Create(&post).Error

	// Create test comment
	comment := models.Comment{
		Content: "Original content",
		PostID:  post.ID,
		UserID:  user.ID,
	}
	err = database.DB.Create(&comment).Error

	t.Run("update comment owned by user", func(t *testing.T) {
		service := NewCommentsService(NewCommentRepository())
		updateReq := models.CommentUpdateRequest{
			Content: "Updated content",
		}

		updatedComment, err := service.UpdateComment(comment.ID, "updatecommentuser", updateReq)

		assert.NoError(t, err)
		assert.NotNil(t, updatedComment)
		assert.Equal(t, "Updated content", updatedComment.Content)
		assert.NotNil(t, updatedComment.UpdatedAt)
	})

	t.Run("update comment not owned by user", func(t *testing.T) {
		service := NewCommentsService(NewCommentRepository())
		updateReq := models.CommentUpdateRequest{
			Content: "Unauthorized update",
		}

		updatedComment, err := service.UpdateComment(comment.ID, "differentuser", updateReq)

		assert.Error(t, err)
		assert.Nil(t, updatedComment)
		assert.Contains(t, err.Error(), "unauthorized")
	})
}

// TestDeleteComment_Service tests the DeleteComment service method
func TestDeleteComment_Service(t *testing.T) {
	if database.DB == nil {
		t.Skip("Database not available, skipping service integration tests")
		return
	}

	// Create test user
	refreshToken5 := "refresh_deletecommentuser"
	user := models.User{Username: "deletecommentuser", Password: "password", RefreshToken: &refreshToken5}
	err := database.DB.Create(&user).Error

	// Create test sub and post
	sub := models.Sub{Name: "deletecommentsub", OwnerID: user.ID}
	err = database.DB.Create(&sub).Error

	post := models.Post{
		Title:   "Post for DeleteComment Test",
		Content: "Post content",
		UserID:  user.ID,
		SubID:   sub.ID,
	}
	err = database.DB.Create(&post).Error

	// Create test comment
	comment := models.Comment{
		Content: "Comment to be deleted",
		PostID:  post.ID,
		UserID:  user.ID,
	}
	err = database.DB.Create(&comment).Error

	t.Run("delete comment owned by user", func(t *testing.T) {
		service := NewCommentsService(NewCommentRepository())
		err := service.DeleteComment(comment.ID, "deletecommentuser")

		assert.NoError(t, err)

		// Verify comment is deleted
		var deletedComment models.Comment
		err = database.DB.First(&deletedComment, comment.ID).Error
		assert.Error(t, err) // Should fail to find deleted record
	})

	t.Run("delete comment not owned by user", func(t *testing.T) {
		// Create another comment for this test
		comment2 := models.Comment{
			Content: "Another comment to be deleted",
			PostID:  post.ID,
			UserID:  user.ID,
		}
		err = database.DB.Create(&comment2).Error

		service := NewCommentsService(NewCommentRepository())
		err = service.DeleteComment(comment2.ID, "unauthorizeduser")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unauthorized")
	})
}
