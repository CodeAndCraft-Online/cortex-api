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

// TestMain is defined in auth_handlers_test.go for the handlers package

// Helper function to create test router for votes
func setupVotesTestRouter(username string, userID uint) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Apply auth middleware to protected routes
	router.Use(func(c *gin.Context) {
		c.Set("username", username)
		c.Set("user_id", userID)
		c.Next()
	})

	router.POST("/posts/upvote", UpvotePost)
	router.POST("/posts/downvote", DownvotePost)

	return router
}

func TestUpvotePostHandler(t *testing.T) {
	if database.DB == nil {
		t.Skip("Database not available, skipping integration test")
		return
	}

	// Setup test data
	user := models.User{Username: "voter", Password: "password"}
	database.DB.Create(&user)

	sub := models.Sub{Name: "voteSub", Description: "Sub for voting", OwnerID: user.ID, Private: false}
	database.DB.Create(&sub)

	post := models.Post{
		Title:   "Post for Upvote Test",
		Content: "Content for upvote testing",
		UserID:  user.ID,
		SubID:   sub.ID,
	}
	database.DB.Create(&post)

	router := setupVotesTestRouter("voter", user.ID)

	t.Run("successful upvote on new post", func(t *testing.T) {
		voteRequest := map[string]interface{}{"postID": post.ID}
		jsonData, _ := json.Marshal(voteRequest)

		req, _ := http.NewRequest("POST", "/posts/upvote", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "Vote recorded", response["message"])
	})

	t.Run("upvote existing post with downvote", func(t *testing.T) {
		// Clear any existing votes first
		database.DB.Exec("DELETE FROM votes WHERE user_id = ? AND post_id = ?", user.ID, post.ID)

		// Create existing downvote
		existingVote := models.Vote{
			UserID: user.ID,
			PostID: post.ID,
			Vote:   -1,
		}
		database.DB.Create(&existingVote)

		voteRequest := map[string]interface{}{"postID": post.ID}
		jsonData, _ := json.Marshal(voteRequest)

		req, _ := http.NewRequest("POST", "/posts/upvote", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "Vote updated", response["message"])
	})

	t.Run("remove existing upvote", func(t *testing.T) {
		// Clear any existing votes first
		database.DB.Exec("DELETE FROM votes WHERE user_id = ? AND post_id = ?", user.ID, post.ID)

		// Create existing upvote
		existingVote := models.Vote{
			UserID: user.ID,
			PostID: post.ID,
			Vote:   1,
		}
		database.DB.Create(&existingVote)

		voteRequest := map[string]interface{}{"postID": post.ID}
		jsonData, _ := json.Marshal(voteRequest)

		req, _ := http.NewRequest("POST", "/posts/upvote", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "Vote removed", response["message"])
	})
}

func TestDownvotePostHandler(t *testing.T) {
	if database.DB == nil {
		t.Skip("Database not available, skipping integration test")
		return
	}

	// Setup test data
	user := models.User{Username: "downvoter", Password: "password"}
	database.DB.Create(&user)

	sub := models.Sub{Name: "downvoteSub", Description: "Sub for downvoting", OwnerID: user.ID, Private: false}
	database.DB.Create(&sub)

	post := models.Post{
		Title:   "Post for Downvote Test",
		Content: "Content for downvote testing",
		UserID:  user.ID,
		SubID:   sub.ID,
	}
	database.DB.Create(&post)

	router := setupVotesTestRouter("downvoter", user.ID)

	t.Run("successful downvote on new post", func(t *testing.T) {
		voteRequest := map[string]interface{}{"postID": post.ID}
		jsonData, _ := json.Marshal(voteRequest)

		req, _ := http.NewRequest("POST", "/posts/downvote", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "Vote recorded", response["message"])
	})

	t.Run("remove existing downvote", func(t *testing.T) {
		// Create existing downvote
		existingVote := models.Vote{
			UserID: user.ID,
			PostID: post.ID,
			Vote:   -1,
		}
		database.DB.Create(&existingVote)

		voteRequest := map[string]interface{}{"postID": post.ID}
		jsonData, _ := json.Marshal(voteRequest)

		req, _ := http.NewRequest("POST", "/posts/downvote", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "Vote removed", response["message"])
	})
}

func TestVoteHandlerErrors(t *testing.T) {
	if database.DB == nil {
		t.Skip("Database not available, skipping integration test")
		return
	}

	router := setupVotesTestRouter("testuser", 1)

	t.Run("unauthorized user - no auth middleware", func(t *testing.T) {
		r := gin.Default() // No auth middleware
		r.POST("/posts/upvote", UpvotePost)

		voteRequest := map[string]interface{}{"postID": 999}
		jsonData, _ := json.Marshal(voteRequest)

		req, _ := http.NewRequest("POST", "/posts/upvote", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("invalid JSON request", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/posts/upvote", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("post not found", func(t *testing.T) {
		voteRequest := map[string]interface{}{"postID": 99999}
		jsonData, _ := json.Marshal(voteRequest)

		req, _ := http.NewRequest("POST", "/posts/upvote", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "Post not found", response["error"])
	})

	t.Run("user not found", func(t *testing.T) {
		// Create a post
		sub := models.Sub{Name: "errorSub", Description: "Sub for error testing", OwnerID: 1, Private: false}
		database.DB.Create(&sub)

		post := models.Post{
			Title:   "Post for Error Test",
			Content: "Content for error testing",
			UserID:  1,
			SubID:   sub.ID,
		}
		database.DB.Create(&post)

		// Delete the user to simulate user not found
		database.DB.Delete(&models.User{}, 1)

		r := gin.Default()
		r.Use(func(c *gin.Context) {
			c.Set("username", "nonexistent")
			c.Next()
		})
		r.POST("/posts/upvote", UpvotePost)

		voteRequest := map[string]interface{}{"postID": post.ID}
		jsonData, _ := json.Marshal(voteRequest)

		req, _ := http.NewRequest("POST", "/posts/upvote", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "User not found", response["error"])
	})
}
