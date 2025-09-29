package repositories

import (
	"testing"

	"github.com/CodeAndCraft-Online/cortex-api/internal/database"
	"github.com/CodeAndCraft-Online/cortex-api/internal/models"
	"github.com/stretchr/testify/assert"
)

// TestMain is defined in auth_repository_test.go for the repositories package

func TestGetCommentsByPostID(t *testing.T) {
	if !dbAvailable {
		t.Skip("Database not available, skipping repository integration tests")
		return
	}

	// Create test user
	user := models.User{
		Username: "getcommentsuser",
		Password: "password",
	}
	database.DB.Create(&user)

	// Create test sub
	sub := models.Sub{
		Name:    "commentsub",
		OwnerID: user.ID,
	}
	database.DB.Create(&sub)

	// Create test post
	post := models.Post{
		Title:   "Test Post for Comments",
		Content: "Post content",
		UserID:  user.ID,
		SubID:   sub.ID,
	}
	database.DB.Create(&post)

	// Create test comments
	comment1 := models.Comment{
		PostID:  post.ID,
		UserID:  user.ID,
		Content: "First comment",
	}
	comment2 := models.Comment{
		PostID:  post.ID,
		UserID:  user.ID,
		Content: "Second comment",
	}
	database.DB.Create(&comment1)
	database.DB.Create(&comment2)

	t.Run("get comments for existing post", func(t *testing.T) {
		comments, err := GetCommentsByPostID("1") // Assuming ID is 1

		assert.NoError(t, err)
		assert.NotNil(t, comments)
		assert.True(t, len(*comments) >= 2)

		// Check first comment
		firstComment := (*comments)[0]
		assert.Equal(t, "getcommentsuser", firstComment.Username)
		assert.Contains(t, []string{"First comment", "Second comment"}, firstComment.Content)
	})

	t.Run("get comments for post with no comments", func(t *testing.T) {
		// Create another post with no comments
		emptyPost := models.Post{
			Title:   "Empty Post",
			Content: "No comments here",
			UserID:  user.ID,
			SubID:   sub.ID,
		}
		database.DB.Create(&emptyPost)

		comments, err := GetCommentsByPostID("2") // Assuming ID is 2

		assert.NoError(t, err)
		assert.NotNil(t, comments)
		// Should return empty slice, not nil
		assert.Equal(t, 0, len(*comments))
	})

	t.Run("get comments for non-existent post", func(t *testing.T) {
		comments, err := GetCommentsByPostID("999")

		assert.NoError(t, err) // This should not error, just return empty
		assert.NotNil(t, comments)
		assert.Equal(t, 0, len(*comments))
	})
}

func TestCreateComment(t *testing.T) {
	if !dbAvailable {
		t.Skip("Database not available, skipping repository integration tests")
		return
	}

	// Create test user
	user := models.User{
		Username: "createcommentuser",
		Password: "password",
	}
	database.DB.Create(&user)

	// Create test sub
	sub := models.Sub{
		Name:    "createcommentsub",
		OwnerID: user.ID,
	}
	database.DB.Create(&sub)

	// Create test post
	post := models.Post{
		Title:   "Post for New Comment",
		Content: "Post content",
		UserID:  user.ID,
		SubID:   sub.ID,
	}
	database.DB.Create(&post)

	t.Run("create comment successfully", func(t *testing.T) {
		commentReq := models.CommentRequest{
			PostID:  post.ID,
			Content: "This is a new comment",
		}

		// Create the post object as expected by CreateComment
		postObj := models.Post{ID: post.ID}

		createdComment, err := CreateComment("createcommentuser", commentReq, postObj)

		assert.NoError(t, err)
		assert.NotNil(t, createdComment)
		assert.Equal(t, commentReq.Content, createdComment.Content)
		assert.Equal(t, user.ID, createdComment.UserID)
		assert.Equal(t, post.ID, createdComment.PostID)
		assert.NotZero(t, createdComment.ID)
	})

	t.Run("create comment with invalid user", func(t *testing.T) {
		commentReq := models.CommentRequest{
			PostID:  post.ID,
			Content: "Comment from invalid user",
		}

		postObj := models.Post{ID: post.ID}

		createdComment, err := CreateComment("nonexistentuser", commentReq, postObj)

		assert.Error(t, err)
		assert.Nil(t, createdComment)
	})

	t.Run("create comment with invalid post", func(t *testing.T) {
		commentReq := models.CommentRequest{
			PostID:  999, // Invalid post ID
			Content: "Comment on invalid post",
		}

		postObj := models.Post{ID: 999}

		createdComment, err := CreateComment("createcommentuser", commentReq, postObj)

		assert.Error(t, err)
		assert.Nil(t, createdComment)
	})
}
