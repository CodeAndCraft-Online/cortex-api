package handlers

import (
	"net/http"

	"github.com/CodeAndCraft-Online/cortex-api/internal/services"
	"github.com/gin-gonic/gin"
)

// GetSubs returns public subs + private subs for authorized users
func GetSubs(c *gin.Context) {

	subs, err := services.GetSubs(c.Param("username"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, subs)
}

// CreateSub allows users to create a new subreddit
func CreateSub(c *gin.Context) {
	var subRequest struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Private     bool   `json:"private"`
	}

	newSub, err := services.CreateSub(c.Param("username"), subRequest)
	if err != nil {
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

	sub, err := services.JoinSub(c.Param("username"), c.Param("subID"))
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, gin.H{"joined": sub.ID})
}

// InviteUser allows sub owners to invite users to a private sub
func InviteUser(c *gin.Context) {
	var inviteRequest struct {
		InviteeUsername string `json:"invitee"`
	}

	err := services.InviteUser(c.Param("subID"), c.Param("username"), inviteRequest)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Invitation sent to " + inviteRequest.InviteeUsername})
}

// AcceptInvite allows users to accept an invitation
func AcceptInvite(c *gin.Context) {

	err := services.AcceptInvite(c.Param("inviteID"), c.Param("username"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "You have joined the sub."})
}

// ListSubPosts fetches all posts for a specific sub
func ListSubPosts(c *gin.Context) {

	formattedPosts, err := services.ListSubPosts(c.Param("subID"), c.Param("username"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, formattedPosts)
}

// LeaveSub allows users to leave a subreddit
func LeaveSub(c *gin.Context) {

	sub, err := services.LeaveSub(c.Param("subID"), c.Param("username"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Left " + sub.Name})
}

// Get a count of how many posts are in a sub
func GetPostCountPerSub(c *gin.Context) {

	postCount, err := services.GetPostCountPerSub(c.Param("subID"), c.Param("username"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, postCount)
}
