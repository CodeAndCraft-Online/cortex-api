package repositories

import (
	"testing"

	"github.com/CodeAndCraft-Online/cortex-api/internal/database"
	"github.com/CodeAndCraft-Online/cortex-api/internal/models"
	"github.com/stretchr/testify/assert"
)

// Note: TestMain is defined in auth_repository_test.go for the repositories package

func TestGetPostByID(t *testing.T) {
	// Create test user
	user := models.User{
		Username: "testuser",
		Password: "password",
	}
	database.DB.Create(&user)

	// Create test sub
	sub := models.Sub{
		Name:    "testsub",
		OwnerID: user.ID,
	}
	database.DB.Create(&sub)

	// Create test post
	post := models.Post{
		Title:   "Test Post",
		Content: "Test content",
		UserID:  user.ID,
		SubID:   sub.ID,
	}
	database.DB.Create(&post)

	t.Run("valid post ID", func(t *testing.T) {
		response, err := GetPostByID("1") // Assuming ID is 1

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, "Test Post", response.Title)
		assert.Equal(t, "Test content", response.Content)
		assert.Equal(t, "testuser", response.Username)
	})

	t.Run("invalid post ID", func(t *testing.T) {
		response, err := GetPostByID("999")

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "post not found")
	})
}

func TestFindAllPosts(t *testing.T) {
	// Create test user
	user := models.User{
		Username: "allpostsuser",
		Password: "password",
	}
	database.DB.Create(&user)

	// Create test sub
	sub := models.Sub{
		Name:    "allpostsub",
		OwnerID: user.ID,
	}
	database.DB.Create(&sub)

	// Create test posts
	posts := []models.Post{
		{Title: "Post 1", Content: "Content 1", UserID: user.ID, SubID: sub.ID},
		{Title: "Post 2", Content: "Content 2", UserID: user.ID, SubID: sub.ID},
	}
	for _, post := range posts {
		database.DB.Create(&post)
	}

	result, err := FindAllPosts()

	assert.NoError(t, err)
	assert.True(t, len(result) >= 2)
}

func TestCreatePost(t *testing.T) {
	// Create test user
	user := models.User{
		Username: "createuser",
		Password: "password",
	}
	database.DB.Create(&user)

	// Create test sub
	sub := models.Sub{
		Name:    "createsub",
		OwnerID: user.ID,
	}
	database.DB.Create(&sub)

	newPost := models.Post{
		Title:   "New Post",
		Content: "New content",
		UserID:  user.ID,
		SubID:   sub.ID,
	}

	createdPost, err := CreatPost("createuser", newPost)

	assert.NoError(t, err)
	assert.NotNil(t, createdPost)
	assert.Equal(t, "New Post", createdPost.Title)
	assert.NotZero(t, createdPost.ID)
}
