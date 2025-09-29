package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CodeAndCraft-Online/cortex-api/internal/database"
	"github.com/CodeAndCraft-Online/cortex-api/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestGetCommentByID_Handler tests the GetCommentByID handler
func TestGetCommentByID_Handler(t *testing.T) {
	if database.DB == nil {
		t.Skip("Database not available, skipping handler integration test")
		return
	}

	// Create test data
	user := models.User{
		Username: "getcommenthandleruser",
		Password: "password",
	}
	database.DB.Where(&user).FirstOrCreate(&user)

	sub := models.Sub{Name: "getcommenthandlersub", OwnerID: user.ID}
	database.DB.Where(&sub).FirstOrCreate(&sub)

	post := models.Post{
		Title:   "Post for Handler Test",
		Content: "Post content",
		UserID:  user.ID,
		SubID:   sub.ID,
	}
	database.DB.Where(&post).FirstOrCreate(&post)

	comment := models.Comment{
		Content:  "Test comment for handler",
		PostID:   post.ID,
		UserID:   user.ID,
		ParentID: nil,
	}
	database.DB.Create(&comment)

	// Create test router
	router := gin.Default()
	router.GET("/comments/:id", GetCommentByID)

	t.Run("get comment successfully", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/comments/1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.CommentResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "getcommenthandleruser", response.Username)
		assert.Equal(t, "Test comment for handler", response.Content)
	})

	t.Run("get non-existent comment", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/comments/999", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var errorResponse gin.H
		err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
		assert.NoError(t, err)
		assert.Contains(t, errorResponse["error"], "not found")
	})
}

// TestUpdateComment_Handler tests the UpdateComment handler
func TestUpdateComment_Handler(t *testing.T) {
	if database.DB == nil {
		t.Skip("Database not available, skipping handler integration test")
		return
	}

	// Create test data
	user := models.User{
		Username: "updatecommenthandleruser",
		Password: "password",
	}
	database.DB.Where(&user).FirstOrCreate(&user)

	sub := models.Sub{Name: "updatecommenthandlersub", OwnerID: user.ID}
	database.DB.Where(&sub).FirstOrCreate(&sub)

	post := models.Post{
		Title:   "Post for Update Handler Test",
		Content: "Post content",
		UserID:  user.ID,
		SubID:   sub.ID,
	}
	database.DB.Where(&post).FirstOrCreate(&post)

	comment := models.Comment{
		Content: "Original comment content",
		PostID:  post.ID,
		UserID:  user.ID,
	}
	database.DB.Create(&comment)

	// Create test router with mock authentication
	router := gin.Default()
	router.Use(func(c *gin.Context) {
		// Mock authentication middleware
		c.Set("username", "updatecommenthandleruser")
		c.Next()
	})
	router.PUT("/comments/:id", UpdateComment)

	t.Run("update comment successfully", func(t *testing.T) {
		updateData := models.CommentUpdateRequest{
			Content: "Updated comment content",
		}
		jsonData, _ := json.Marshal(updateData)

		req, _ := http.NewRequest("PUT", "/comments/1", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var updatedComment models.Comment
		json.Unmarshal(w.Body.Bytes(), &updatedComment)
		assert.Contains(t, string(w.Body.Bytes()), "Updated comment content")
	})

	t.Run("update comment with invalid data", func(t *testing.T) {
		updateData := map[string]interface{}{
			"content": "",
		}
		jsonData, _ := json.Marshal(updateData)

		req, _ := http.NewRequest("PUT", "/comments/1", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var errorResponse gin.H
		json.Unmarshal(w.Body.Bytes(), &errorResponse)
		assert.Contains(t, errorResponse["error"], "required")
	})
}

// TestDeleteComment_Handler tests the DeleteComment handler
func TestDeleteComment_Handler(t *testing.T) {
	if database.DB == nil {
		t.Skip("Database not available, skipping handler integration test")
		return
	}

	// Create test data
	user := models.User{
		Username: "deletecommenthandleruser",
		Password: "password",
	}
	database.DB.Where(&user).FirstOrCreate(&user)

	sub := models.Sub{Name: "deletecommenthandlersub", OwnerID: user.ID}
	database.DB.Where(&sub).FirstOrCreate(&sub)

	post := models.Post{
		Title:   "Post for Delete Handler Test",
		Content: "Post content",
		UserID:  user.ID,
		SubID:   sub.ID,
	}
	database.DB.Where(&post).FirstOrCreate(&post)

	comment := models.Comment{
		Content: "Comment to be deleted by handler",
		PostID:  post.ID,
		UserID:  user.ID,
	}
	database.DB.Create(&comment)

	// Create test router with mock authentication
	router := gin.Default()
	router.Use(func(c *gin.Context) {
		// Mock authentication middleware
		c.Set("username", "deletecommenthandleruser")
		c.Next()
	})
	router.DELETE("/comments/:id", DeleteComment)

	t.Run("delete comment successfully", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/comments/1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var successResponse gin.H
		json.Unmarshal(w.Body.Bytes(), &successResponse)
		assert.Contains(t, successResponse["message"], "deleted successfully")
	})
}

// Test comment handlers without authentication
func TestComment_Handlers_NoAuth(t *testing.T) {
	// Create test router
	router := gin.Default()
	router.PUT("/comments/:id", UpdateComment)
	router.DELETE("/comments/:id", DeleteComment)

	t.Run("update comment without authentication", func(t *testing.T) {
		updateData := models.CommentUpdateRequest{
			Content: "Unauthorized update",
		}
		jsonData, _ := json.Marshal(updateData)

		req, _ := http.NewRequest("PUT", "/comments/1", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var errorResponse gin.H
		json.Unmarshal(w.Body.Bytes(), &errorResponse)
		assert.Contains(t, errorResponse["error"], "Unauthorized")
	})

	t.Run("delete comment without authentication", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/comments/1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var errorResponse gin.H
		json.Unmarshal(w.Body.Bytes(), &errorResponse)
		assert.Contains(t, errorResponse["error"], "Unauthorized")
	})
}
