package handlers

import (
	"net/http"

	"github.com/CodeAndCraft-Online/cortex-api/internal/models"
	"github.com/CodeAndCraft-Online/cortex-api/internal/services"
	"github.com/gin-gonic/gin"
)

// @Summary Get all subs
// @Description Returns all public subs and private subs the user is authorized to access
// @Tags Subs
// @Produce json
// @Success 200 {array} interface{} "Array of available subs"
// @Failure 500 {object} map[string]string "error: Internal server error"
// @Security BearerAuth
// @Router /subs/ [get]
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

// @Summary Create a new sub
// @Description Creates a new subreddit (community)
// @Tags Subs
// @Accept json
// @Produce json
// @Param sub body models.SubRequest true "Sub creation data"
// @Success 201 {object} interface{} "Created sub with ID and details"
// @Failure 400 {object} map[string]string "error: Bad request or sub name already taken"
// @Failure 401 {object} map[string]string "error: Unauthorized"
// @Failure 500 {object} map[string]string "error: Internal server error"
// @Security BearerAuth
// @Router /subs/ [post]
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

// @Summary Join a sub
// @Description Allows users to join a subreddit (public or invited private subs)
// @Tags Subs
// @Produce json
// @Param subID path string true "Sub ID"
// @Success 200 {object} interface{} "membership.SubID: ID of joined sub"
// @Failure 401 {object} map[string]string "error: Unauthorized"
// @Failure 500 {object} map[string]string "error: Internal server error"
// @Security BearerAuth
// @Router /sub/sub/{subID}/join [post]
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

// @Summary Invite user to private sub
// @Description Allows sub owners and moderators to invite users to private subs
// @Tags Subs
// @Accept json
// @Produce json
// @Param subID path string true "Sub ID"
// @Param invite body models.InviteRequest true "User to invite"
// @Success 200 {object} map[string]string "message: Invitation sent"
// @Failure 401 {object} map[string]string "error: Unauthorized"
// @Failure 404 {object} map[string]string "error: Sub not found or permission denied"
// @Security BearerAuth
// @Router /sub/sub/{subID}/invite [post]
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

// @Summary Accept sub invitation
// @Description Allows users to accept invitations to join private subs
// @Tags Subs
// @Produce json
// @Param inviteID path string true "Invite ID"
// @Success 200 {object} map[string]string "message: You have joined the sub"
// @Failure 401 {object} map[string]string "error: Unauthorized"
// @Failure 404 {object} map[string]string "error: Invitation not found or expired"
// @Failure 500 {object} map[string]string "error: Internal server error"
// @Security BearerAuth
// @Router /user/invite/{inviteID}/accept [post]
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

// @Summary Get posts by sub ID
// @Description Fetches all posts for a specific subreddit
// @Tags Subs
// @Produce json
// @Param subID path string true "Sub ID"
// @Success 200 {array} interface{} "Array of posts in the sub"
// @Failure 401 {object} map[string]string "error: Unauthorized access to private sub"
// @Failure 500 {object} map[string]string "error: Internal server error"
// @Security BearerAuth
// @Router /sub/sub/{subID}/posts [get]
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

// @Summary Leave a sub
// @Description Allows users to leave a subreddit they are currently a member of
// @Tags Subs
// @Produce json
// @Param subID path string true "Sub ID"
// @Success 200 {object} map[string]string "message: Left [sub name]"
// @Failure 401 {object} map[string]string "error: Unauthorized"
// @Failure 404 {object} map[string]string "error: Sub not found or not a member"
// @Failure 500 {object} map[string]string "error: Internal server error"
// @Security BearerAuth
// @Router /sub/sub/{subID}/leave [post]
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

// @Summary Get post count for sub
// @Description Returns the total number of posts in a specified sub
// @Tags Subs
// @Produce json
// @Param subID query string true "Sub ID"
// @Success 200 {integer} int "Post count in the sub"
// @Failure 401 {object} map[string]string "error: Unauthorized access to private sub"
// @Failure 404 {object} map[string]string "error: Sub not found"
// @Failure 500 {object} map[string]string "error: Internal server error"
// @Security BearerAuth
// @Router /sub/sub/{subID}/postCount [get]
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
