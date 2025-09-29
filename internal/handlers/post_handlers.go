package handlers

import (
	"net/http"

	db "github.com/CodeAndCraft-Online/cortex-api/internal/database"
	models "github.com/CodeAndCraft-Online/cortex-api/internal/models"
	"github.com/CodeAndCraft-Online/cortex-api/internal/services"
	"github.com/gin-gonic/gin"
)

// @Summary Get post by ID
// @Description Retrieves a specific post with user details and comment count
// @Tags Posts
// @Produce json
// @Param postID path string true "Post ID"
// @Success 200 {object} interface{} "Post with user and comment details"
// @Failure 404 {object} map[string]string "error: Post not found"
// @Failure 500 {object} map[string]string "error: Internal server error"
// @Security BearerAuth
// @Router /posts/{postID} [get]
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

// @Summary Create a new post
// @Description Creates a new post for the authenticated user
// @Tags Posts
// @Accept json
// @Produce json
// @Param post body models.Post true "Post data"
// @Success 201 {object} interface{} "Created post with details"
// @Failure 400 {object} map[string]string "error: Bad request or validation error"
// @Failure 401 {object} map[string]string "error: must login to post"
// @Failure 500 {object} map[string]string "error: Internal server error"
// @Security BearerAuth
// @Router /posts/ [post]
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

	username, _ := c.Get("username")
	postResponse, err := services.CreatePost(username.(string), post)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusCreated, postResponse)
}

// @Summary Get comments by post ID
// @Description Retrieves all comments for a specific post
// @Tags Posts
// @Produce json
// @Param postID path string true "Post ID"
// @Success 200 {array} interface{} "Array of comments for the post"
// @Failure 500 {object} map[string]string "error: Internal server error"
// @Security BearerAuth
// @Router /posts/posts/{postID}/comments [get]
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

// @Summary Get all posts
// @Description Retrieves all posts with user details and comment/vote counts
// @Tags Posts
// @Produce json
// @Success 200 {array} interface{} "Array of posts with user and vote details"
// @Failure 500 {object} map[string]string "error: Internal server error"
// @Security BearerAuth
// @Router /posts/ [get]
func GetPosts(c *gin.Context) {

	postResponse, err := services.GetPosts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
	}

	c.JSON(http.StatusOK, postResponse)
}

// @Summary Create a comment
// @Description Creates a new comment on a specified post
// @Tags Comments
// @Accept json
// @Produce json
// @Param comment body models.CommentRequest true "Comment data with postID"
// @Success 201 {object} interface{} "Created comment with details"
// @Failure 400 {object} map[string]string "error: Bad request or validation error"
// @Failure 401 {object} map[string]string "error: Unauthorized or user not found"
// @Failure 404 {object} map[string]string "error: Post not found"
// @Failure 500 {object} map[string]string "error: Internal server error"
// @Security BearerAuth
// @Router /posts/comments [post]
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

	comment, err := services.CreateComment(username.(string), commentReq, post)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
