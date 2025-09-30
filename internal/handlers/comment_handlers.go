package handlers

import (
	"net/http"
	"strconv"

	"github.com/CodeAndCraft-Online/cortex-api/internal/models"
	"github.com/CodeAndCraft-Online/cortex-api/internal/repositories"
	"github.com/CodeAndCraft-Online/cortex-api/internal/services"
	"github.com/gin-gonic/gin"
)

// GetCommentByID retrieves a single comment by ID
// @Summary Get comment by ID
// @Description Get a single comment with its details
// @Tags comments
// @Accept  json
// @Produce  json
// @Param id path int true "Comment ID"
// @Success 200 {object} models.CommentResponse
// @Failure 401 {object} gin.H
// @Failure 404 {object} gin.H
// @Router /comments/{id} [get]
func GetCommentByID(c *gin.Context) {
	// Get comment ID from URL
	commentIDStr := c.Param("id")
	commentID, err := strconv.ParseUint(commentIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	// Get comment from service
	service := services.NewCommentsService(repositories.NewCommentRepository())
	comment, err := service.GetCommentByID(uint(commentID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, comment)
}

// UpdateComment updates an existing comment
// @Summary Update comment
// @Description Update a comment's content (only by comment author)
// @Tags comments
// @Accept  json
// @Produce  json
// @Param id path int true "Comment ID"
// @Param comment body models.CommentUpdateRequest true "Comment update data"
// @Success 200 {object} models.Comment
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Failure 403 {object} gin.H
// @Failure 404 {object} gin.H
// @Router /comments/{id} [put]
func UpdateComment(c *gin.Context) {
	// Get username from JWT token
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get comment ID from URL
	commentIDStr := c.Param("id")
	commentID, err := strconv.ParseUint(commentIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	// Parse request body
	var commentReq models.CommentUpdateRequest
	if err := c.ShouldBindJSON(&commentReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate input
	if commentReq.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Content is required"})
		return
	}

	// Update comment via service
	service := services.NewCommentsService(repositories.NewCommentRepository())
	updatedComment, err := service.UpdateComment(uint(commentID), username.(string), commentReq)
	if err != nil {
		if err.Error() == "unauthorized: can only edit own comments" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedComment)
}

// DeleteComment deletes a comment
// @Summary Delete comment
// @Description Delete a comment (only by comment author)
// @Tags comments
// @Accept  json
// @Produce  json
// @Param id path int true "Comment ID"
// @Success 200 {object} gin.H
// @Failure 401 {object} gin.H
// @Failure 403 {object} gin.H
// @Failure 404 {object} gin.H
// @Router /comments/{id} [delete]
func DeleteComment(c *gin.Context) {
	// Get username from JWT token
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get comment ID from URL
	commentIDStr := c.Param("id")
	commentID, err := strconv.ParseUint(commentIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	// Delete comment via service
	service := services.NewCommentsService(repositories.NewCommentRepository())
	err = service.DeleteComment(uint(commentID), username.(string))
	if err != nil {
		if err.Error() == "unauthorized: can only delete own comments" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}
