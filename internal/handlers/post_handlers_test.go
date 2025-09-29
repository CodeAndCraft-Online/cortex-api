package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CodeAndCraft-Online/cortex-api/internal/database"
	"github.com/CodeAndCraft-Online/cortex-api/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestMain is defined in auth_handlers_test.go for the handlers package

// Helper function to create test router
func setupPostTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Apply auth middleware to protected routes
	router.Use(func(c *gin.Context) {
		// Mock auth middleware - set a test username
		c.Set("username", "testuser")
		c.Next()
	})

	router.GET("/posts/:postID", GetPostByID)
	router.POST("/posts", CreatePost)
	router.GET("/posts", GetPosts)
	router.POST("/posts/:postID/comments", CreateComment)
	router.GET("/posts/:postID/comments", GetCommentsByPostID)

	return router
}

func TestGetPostByIDHandler(t *testing.T) {
	if database.DB == nil {
		t.Skip("Database not available, skipping integration test")
		return
	}

	router := setupPostTestRouter()

	// Create test data
	user := models.User{Username: "testuser", Password: "password"}
	database.DB.Create(&user)

	sub := models.Sub{Name: "postsub", OwnerID: user.ID}
	database.DB.Create(&sub)

	post := models.Post{
		Title:   "Test Post for Handler",
		Content: "Handler test content",
		UserID:  user.ID,
		SubID:   sub.ID,
	}
	database.DB.Create(&post)

	// Test GET request
	req, _ := http.NewRequest("GET", "/posts/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Test Post for Handler", response["title"])
	assert.Equal(t, "testuser", response["username"])
}

func TestGetPostByIDHandler_PostNotFound(t *testing.T) {
	if database.DB == nil {
		t.Skip("Database not available, skipping integration test")
		return
	}

	router := setupPostTestRouter()

	// Test with non-existent post
	req, _ := http.NewRequest("GET", "/posts/999", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Post not found", response["error"])
}

func TestCreatePostHandler(t *testing.T) {
	if database.DB == nil {
		t.Skip("Database not available, skipping integration test")
		return
	}

	router := setupPostTestRouter()

	// Find or create test user
	user := models.User{Username: "testuser"}
	database.DB.Where("username = ?", "testuser").FirstOrCreate(&user)

	sub := models.Sub{Name: "createposthandlersub", OwnerID: user.ID}
	database.DB.Create(&sub)

	// Create request payload
	postData := map[string]interface{}{
		"title":   "New Post via Handler",
		"content": "Handler test content",
		"subID":   sub.ID,
	}
	jsonData, _ := json.Marshal(postData)

	// Create test request
	req, _ := http.NewRequest("POST", "/posts", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "New Post via Handler", response["Title"])
	assert.Equal(t, float64(user.ID), response["UserID"])
}

func TestGetPostsHandler(t *testing.T) {
	if database.DB == nil {
		t.Skip("Database not available, skipping integration test")
		return
	}

	router := setupPostTestRouter()

	// Create test data
	user := models.User{Username: "getpostsuser", Password: "password"}
	database.DB.Create(&user)

	sub := models.Sub{Name: "getpostssub", OwnerID: user.ID}
	database.DB.Create(&sub)

	post := models.Post{
		Title:   "Post for GetPosts Test",
		Content: "Test content for get posts",
		UserID:  user.ID,
		SubID:   sub.ID,
	}
	database.DB.Create(&post)

	// Test GET all posts
	req, _ := http.NewRequest("GET", "/posts", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, len(response) >= 1)
}

func TestCreateCommentHandler(t *testing.T) {
	if database.DB == nil {
		t.Skip("Database not available, skipping integration test")
		return
	}

	// Clear comments to avoid ID conflicts
	database.DB.Exec("TRUNCATE TABLE comments RESTART IDENTITY CASCADE")

	router := setupPostTestRouter()

	// Find or create test user
	user := models.User{Username: "testuser"}
	database.DB.Where("username = ?", "testuser").FirstOrCreate(&user)

	sub := models.Sub{Name: "commenthandlersub", OwnerID: user.ID}
	database.DB.Create(&sub)

	post := models.Post{
		Title:   "Post for Comment Test",
		Content: "Content for comment testing",
		UserID:  user.ID,
		SubID:   sub.ID,
	}
	database.DB.Create(&post)

	// Create comment payload
	commentData := map[string]interface{}{
		"postID":  post.ID,
		"content": "This is a test comment from handler",
	}
	jsonData, _ := json.Marshal(commentData)

	// Create test request
	req, _ := http.NewRequest("POST", fmt.Sprintf("/posts/%d/comments", post.ID), bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "This is a test comment from handler", response["content"])
	assert.Equal(t, "testuser", response["username"])
}

func TestGetCommentsByPostIDHandler(t *testing.T) {
	if database.DB == nil {
		t.Skip("Database not available, skipping integration test")
		return
	}

	router := setupPostTestRouter()

	// Create test data
	user := models.User{Username: "getcommentsuser", Password: "password"}
	database.DB.Create(&user)

	sub := models.Sub{Name: "getcommentssub", OwnerID: user.ID}
	database.DB.Create(&sub)

	post := models.Post{
		Title:   "Post with Comments",
		Content: "Content for comments test",
		UserID:  user.ID,
		SubID:   sub.ID,
	}
	database.DB.Create(&post)

	comment := models.Comment{
		PostID:  post.ID,
		UserID:  user.ID,
		Content: "Test comment for handler",
	}
	database.DB.Create(&comment)

	// Test GET comments
	req, _ := http.NewRequest("GET", fmt.Sprintf("/posts/%d/comments", post.ID), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, len(response) >= 1)
}
