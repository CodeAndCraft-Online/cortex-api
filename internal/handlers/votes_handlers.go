package handlers

import (
	"net/http"

	db "github.com/CodeAndCraft-Online/cortex-api/internal/database"
	models "github.com/CodeAndCraft-Online/cortex-api/internal/models"
	"github.com/gin-gonic/gin"
)

// DownvotePost allows a user to downvote a post
func DownvotePost(c *gin.Context) {
	handleVote(c, -1) // -1 = Downvote
}

func UpvotePost(c *gin.Context) {
	handleVote(c, 1) // 1 = Upvote
}

// handleVote processes upvotes and downvotes
func handleVote(c *gin.Context, voteValue int) {
	var voteRequest struct {
		PostID uint `json:"postID"`
	}

	// Extract username from JWT
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Bind JSON input
	if err := c.ShouldBindJSON(&voteRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Fetch user ID
	var user models.User
	if err := db.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// Check if the post exists
	var post models.Post
	if err := db.DB.First(&post, voteRequest.PostID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// ✅ Check if the user has already voted
	var existingVote models.Vote
	err := db.DB.Where("user_id = ? AND post_id = ?", user.ID, voteRequest.PostID).First(&existingVote).Error

	if err != nil && err.Error() != "record not found" {
		// Only return an error if it's a real DB error, not a "record not found"
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if err == nil {
		// ✅ If the user has already voted and selects the same vote, remove it
		if existingVote.Vote == voteValue {
			db.DB.Delete(&existingVote)
			c.JSON(http.StatusOK, gin.H{"message": "Vote removed"})
			return
		}

		// ✅ If the user voted differently before, update it
		existingVote.Vote = voteValue
		db.DB.Save(&existingVote)
		c.JSON(http.StatusOK, gin.H{"message": "Vote updated"})
		return
	}

	// ✅ If no previous vote exists, create a new one
	newVote := models.Vote{
		UserID: user.ID,
		PostID: voteRequest.PostID,
		Vote:   voteValue,
	}
	db.DB.Create(&newVote)

	c.JSON(http.StatusCreated, gin.H{"message": "Vote recorded"})
}
