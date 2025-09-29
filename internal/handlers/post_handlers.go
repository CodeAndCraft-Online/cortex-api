package handlers

import (
	"net/http"

	db "github.com/CodeAndCraft-Online/cortex-api/internal/database"
	models "github.com/CodeAndCraft-Online/cortex-api/internal/models"
	"github.com/CodeAndCraft-Online/cortex-api/internal/services"
	"github.com/gin-gonic/gin"
)

func GetPostByID(c *gin.Context) {
	postID := c.Param("postID") // Get postID from URL parameter

	// ✅ Fetch the post and preload user details
	postResponse, err := services.GetPostByID(postID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	c.JSON(http.StatusOK, postResponse)
}

// CreatePost handles creating a new post
func CreatePost(c *gin.Context) {

	// Get username from the JWT token stored in the Gin context
	_, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "you must login to post"})
		return
	}

	var post models.Post
	// Bind request JSON to post model
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	postResponse, err := services.CreatPost(c.Param("username"), post)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
	}

	c.JSON(http.StatusCreated, postResponse)
}

// GetCommentsByPostID retrieves all comments for a specific post
func GetCommentsByPostID(c *gin.Context) {
	postID := c.Param("postID") // Get postID from URL parameter

	comments, err := services.GetCommentsByPostID(postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
	}

	c.JSON(http.StatusOK, comments)
}

// GetPosts retrieves all posts with user and associated comments (including upvote/downvote counts)
func GetPosts(c *gin.Context) {

	postResponse, err := services.GetPosts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
	}

	c.JSON(http.StatusOK, postResponse)
}

// CreateComment handles adding a comment to a post
func CreateComment(c *gin.Context) {
	var commentReq models.CommentRequest

	// Get username from the JWT token stored in the Gin context
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Fetch user ID from the database based on username
	var user models.User
	if err := db.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// Bind JSON input to CommentRequest struct
	if err := c.ShouldBindJSON(&commentReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ensure the post exists before adding a comment
	var post models.Post
	if err := db.DB.First(&post, commentReq.PostID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	comment, err := services.CreateComment(c.Param("username"), commentReq, post)
	if err != nil {
		return
	}

	// Save the comment to the database
	if err := db.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":        comment.ID,
		"content":   comment.Content,
		"postID":    comment.PostID,
		"username":  user.Username,
		"imageURL":  comment.ImageURL,
		"createdAt": comment.CreatedAt.Format("2006-01-02 15:04:05"),
	})
}
