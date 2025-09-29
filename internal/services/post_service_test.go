package services

import (
	"testing"

	"github.com/CodeAndCraft-Online/cortex-api/internal/database"
	"github.com/CodeAndCraft-Online/cortex-api/internal/models"
	"github.com/stretchr/testify/assert"
)

// Note: Service tests would require Docker for database access
// TestMain is defined in auth_service_test.go
func TestGetPostByID_Service(t *testing.T) {
	// Create test data
	user := models.User{Username: "postuser", Password: "password"}
	database.DB.Create(&user)

	sub := models.Sub{Name: "postsub", OwnerID: user.ID}
	database.DB.Create(&sub)

	post := models.Post{
		Title:   "Service Test Post",
		Content: "Service test content",
		UserID:  user.ID,
		SubID:   sub.ID,
	}
	database.DB.Create(&post)

	// Test service
	response, err := GetPostByID("1") // Assuming ID is 1

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "Service Test Post", response.Title)
	assert.Equal(t, "postuser", response.Username)
}

func TestGetPostByID_Service_PostNotFound(t *testing.T) {
	response, err := GetPostByID("999")

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "post not found")
}

func TestGetAllPosts_Service(t *testing.T) {
	// Clear any existing posts first (this is a simplified approach)
	database.DB.Exec("DELETE FROM posts")

	// Create test data
	user := models.User{Username: "allpostsuser", Password: "password"}
	database.DB.Create(&user)

	sub := models.Sub{Name: "allpostssub", OwnerID: user.ID}
	database.DB.Create(&sub)

	// Create test posts
	post1 := models.Post{Title: "Post 1", Content: "Content 1", UserID: user.ID, SubID: sub.ID}
	post2 := models.Post{Title: "Post 2", Content: "Content 2", UserID: user.ID, SubID: sub.ID}
	database.DB.Create(&post1)
	database.DB.Create(&post2)

	posts, err := GetAllPosts()

	assert.NoError(t, err)
	assert.True(t, len(posts) >= 2)
}

func TestCreatPost_Service(t *testing.T) {
	user := models.User{Username: "createpostuser", Password: "password"}
	database.DB.Create(&user)

	sub := models.Sub{Name: "createpostsub", OwnerID: user.ID}
	database.DB.Create(&sub)

	newPost := models.Post{
		Title:   "New Service Post",
		Content: "New service content",
		SubID:   sub.ID,
	}

	createdPost, err := CreatPost("createpostuser", newPost)

	assert.NoError(t, err)
	assert.NotNil(t, createdPost)
	assert.Equal(t, "New Service Post", createdPost.Title)
	assert.Equal(t, user.ID, createdPost.UserID)
	assert.NotZero(t, createdPost.ID)
}
