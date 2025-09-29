package handlers

import (
	"net/http"

	"github.com/CodeAndCraft-Online/cortex-api/internal/models"
	"github.com/CodeAndCraft-Online/cortex-api/internal/services"
	"github.com/gin-gonic/gin"
)

// GetSubs returns public subs + private subs for authorized users
func GetSubs(c *gin.Context) {
	username, exists := c.Get("username")
	user := ""
	if exists {
		user = username.(string)
	}

	subs, err := services.GetSubs(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, subs)
}

// CreateSub allows users to create a new subreddit
func CreateSub(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var subRequest models.SubRequest
	if err := c.ShouldBindJSON(&subRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newSub, err := services.CreateSub(username.(string), subRequest)
	if err != nil {
		// Check for specific error types
		if err.Error() == "sub name already taken" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Sub created successfully",
		"id":      newSub.ID,
		"name":    newSub.Name,
		"private": newSub.Private,
	})
}

// JoinSub allows users to join a subreddit (only public or invited private subs)
func JoinSub(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	membership, err := services.JoinSub(username.(string), c.Param("subID"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"joined": membership.SubID})
}

// InviteUser allows sub owners to invite users to a private sub
func InviteUser(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var inviteRequest models.InviteRequest
	if err := c.ShouldBindJSON(&inviteRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := services.InviteUser(c.Param("subID"), username.(string), inviteRequest)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Invitation sent to " + inviteRequest.InviteeUsername})
}

// AcceptInvite allows users to accept an invitation
func AcceptInvite(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	err := services.AcceptInvite(c.Param("inviteID"), username.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "You have joined the sub."})
}

// ListSubPosts fetches all posts for a specific sub
func ListSubPosts(c *gin.Context) {
	username, exists := c.Get("username")
	user := ""
	if exists {
		user = username.(string)
	}

	formattedPosts, err := services.ListSubPosts(c.Param("subID"), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, formattedPosts)
}

// LeaveSub allows users to leave a subreddit
func LeaveSub(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	sub, err := services.LeaveSub(c.Param("subID"), username.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Left " + sub.Name})
}

// Get a count of how many posts are in a sub
func GetPostCountPerSub(c *gin.Context) {
	username, exists := c.Get("username")
	user := ""
	if exists {
		user = username.(string)
	}

	postCount, err := services.GetPostCountPerSub(c.Query("subID"), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, postCount)
}
