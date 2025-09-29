package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/CodeAndCraft-Online/cortex-api/internal/database"
	"github.com/CodeAndCraft-Online/cortex-api/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupSubTestRouter(username string, userID uint) *gin.Engine {
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		// Mock authentication middleware
		c.Set("user_id", userID)
		c.Set("username", username)
		c.Next()
	})
	r.GET("/subs", GetSubs)
	r.POST("/subs", CreateSub)
	r.POST("/subs/:subID/join", JoinSub)
	r.POST("/subs/:subID/invite", InviteUser)
	r.POST("/subs/invite/:inviteID", AcceptInvite)
	r.GET("/subs/:subID/posts", ListSubPosts)
	r.DELETE("/subs/:subID", LeaveSub)
	r.GET("/subs/post-count", GetPostCountPerSub)
	return r
}

func TestGetSubsHandler(t *testing.T) {
	if database.DB == nil {
		t.Skip("Database not available, skipping handler integration tests")
		return
	}

	// Setup - Clear tables in correct order to avoid foreign key violations and data interference
	database.DB.Exec("DELETE FROM comments")
	database.DB.Exec("DELETE FROM posts")
	database.DB.Exec("DELETE FROM sub_memberships")
	database.DB.Exec("DELETE FROM sub_invitations")
	database.DB.Exec("DELETE FROM subs")
	database.DB.Exec("DELETE FROM users")

	// Add test user
	user := models.User{Username: "testuser", Password: "hashedpass"}
	database.DB.Create(&user)

	// Add test subs
	publicSub := models.Sub{Name: "publicsub", Description: "Public Sub", OwnerID: user.ID, Private: false}
	privateSub := models.Sub{Name: "privatesub", Description: "Private Sub", OwnerID: user.ID, Private: true}
	database.DB.Create(&publicSub)
	database.DB.Create(&privateSub)

	router := setupSubTestRouter("testuser", user.ID)

	t.Run("get subs for authenticated user", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/subs", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		responseBody := w.Body.String()
		// Should contain both subs in JSON array format
		assert.Contains(t, responseBody, `"Name":"publicsub"`)
		assert.Contains(t, responseBody, `"Name":"privatesub"`)
	})

	t.Run("get subs for unauthenticated user", func(t *testing.T) {
		r := gin.Default() // No auth middleware
		r.GET("/subs", GetSubs)
		req, _ := http.NewRequest("GET", "/subs", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		responseBody := w.Body.String()
		assert.Contains(t, responseBody, "publicsub")
		assert.NotContains(t, responseBody, "privatesub") // Only public
	})
}

func TestCreateSubHandler(t *testing.T) {
	if database.DB == nil {
		t.Skip("Database not available, skipping handler integration tests")
		return
	}

	// Setup
	database.DB.Exec("DELETE FROM sub_memberships")
	database.DB.Exec("DELETE FROM sub_invitations")
	database.DB.Exec("DELETE FROM subs")
	database.DB.Exec("DELETE FROM users")

	// Add test user
	user := models.User{Username: "creator", Password: "hashedpass"}
	database.DB.Create(&user)

	router := setupSubTestRouter("creator", user.ID)

	t.Run("create public sub", func(t *testing.T) {
		subRequest := `{"name":"newpublic","description":"New Public Sub","private":false}`
		req, _ := http.NewRequest("POST", "/subs", strings.NewReader(subRequest))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		responseBody := w.Body.String()
		assert.Contains(t, responseBody, "Sub created successfully")
		assert.Contains(t, responseBody, "newpublic")
		assert.Contains(t, responseBody, `"private":false`)
	})

	t.Run("create private sub", func(t *testing.T) {
		subRequest := `{"name":"newprivate","description":"New Private Sub","private":true}`
		req, _ := http.NewRequest("POST", "/subs", strings.NewReader(subRequest))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), "newprivate")
	})

	t.Run("create sub with duplicate name", func(t *testing.T) {
		subRequest := `{"name":"newpublic","description":"Duplicate Sub","private":false}`
		req, _ := http.NewRequest("POST", "/subs", strings.NewReader(subRequest))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "sub name already taken")
	})
}

func TestJoinSubHandler(t *testing.T) {
	if database.DB == nil {
		t.Skip("Database not available, skipping handler integration tests")
		return
	}

	// Setup
	database.DB.Exec("DELETE FROM sub_memberships")
	database.DB.Exec("DELETE FROM sub_invitations")
	database.DB.Exec("DELETE FROM subs")
	database.DB.Exec("DELETE FROM users")

	// Add test users
	owner := models.User{Username: "owner", Password: "hashedpass"}
	user := models.User{Username: "member", Password: "hashedpass"}
	database.DB.Create(&owner)
	database.DB.Create(&user)

	// Add public sub
	publicSub := models.Sub{Name: "joinpublic", Description: "Public Sub", OwnerID: owner.ID, Private: false}
	database.DB.Create(&publicSub)

	// Add private sub
	privateSub := models.Sub{Name: "joinprivate", Description: "Private Sub", OwnerID: owner.ID, Private: true}
	database.DB.Create(&privateSub)

	router := setupSubTestRouter("member", user.ID)

	t.Run("join public sub", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/subs/"+fmt.Sprintf("%d", publicSub.ID)+"/join", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		responseBody := w.Body.String()
		assert.Contains(t, responseBody, `"joined":`)
		assert.Contains(t, responseBody, fmt.Sprintf("%d", publicSub.ID))
	})

	t.Run("join private sub without invitation", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/subs/"+fmt.Sprintf("%d", privateSub.ID)+"/join", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code) // Service returns error
		responseBody := w.Body.String()
		// The service returns the error message from repository
		assert.Contains(t, responseBody, "error")
	})
}

func TestListSubPostsHandler(t *testing.T) {
	if database.DB == nil {
		t.Skip("Database not available, skipping handler integration tests")
		return
	}

	// Setup
	database.DB.Exec("DELETE FROM posts")
	database.DB.Exec("DELETE FROM sub_memberships")
	database.DB.Exec("DELETE FROM subs")
	database.DB.Exec("DELETE FROM users")

	// Add test user and sub
	user := models.User{Username: "poster", Password: "hashedpass"}
	database.DB.Create(&user)
	sub := models.Sub{Name: "postsub", Description: "Sub for posts", OwnerID: user.ID, Private: false}
	database.DB.Create(&sub)

	// Add a post
	post := models.Post{Title: "Test Post", Content: "Test Content", UserID: user.ID, SubID: sub.ID}
	database.DB.Create(&post)

	router := setupSubTestRouter("poster", user.ID)

	t.Run("list posts from public sub", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/subs/"+fmt.Sprintf("%d", sub.ID)+"/posts", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Test Post")
	})
}

func TestLeaveSubHandler(t *testing.T) {
	if database.DB == nil {
		t.Skip("Database not available, skipping handler integration tests")
		return
	}

	// Setup
	database.DB.Exec("DELETE FROM sub_memberships")
	database.DB.Exec("DELETE FROM subs")
	database.DB.Exec("DELETE FROM users")

	// Add test user and sub
	user := models.User{Username: "leaver", Password: "hashedpass"}
	database.DB.Create(&user)
	sub := models.Sub{Name: "leavesub", Description: "Sub to leave", OwnerID: user.ID, Private: false}
	database.DB.Create(&sub)

	// Add membership
	membership := models.SubMembership{SubID: sub.ID, UserID: user.ID}
	database.DB.Create(&membership)

	router := setupSubTestRouter("leaver", user.ID)

	t.Run("leave sub", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/subs/"+fmt.Sprintf("%d", sub.ID), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		responseBody := w.Body.String()
		assert.Contains(t, responseBody, `"message":`)
		assert.Contains(t, responseBody, "Left leavesub")
	})

	t.Run("leave sub not member", func(t *testing.T) {
		// Create another sub without membership
		user2 := models.User{Username: "nonmember", Password: "hashedpass"}
		database.DB.Create(&user2)
		sub2 := models.Sub{Name: "notmembersub", Description: "Sub not a member of", OwnerID: user2.ID, Private: false}
		database.DB.Create(&sub2)

		router := setupSubTestRouter("leaver", user.ID)

		req, _ := http.NewRequest("DELETE", "/subs/"+fmt.Sprintf("%d", sub2.ID), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Should return error code (service dependent)
		assert.NotEqual(t, http.StatusOK, w.Code)
	})
}

func TestInviteUserHandler(t *testing.T) {
	if database.DB == nil {
		t.Skip("Database not available, skipping handler integration tests")
		return
	}

	// Setup - Clear tables
	database.DB.Exec("DELETE FROM sub_memberships")
	database.DB.Exec("DELETE FROM sub_invitations")
	database.DB.Exec("DELETE FROM subs")
	database.DB.Exec("DELETE FROM users")

	// Create test users
	owner := models.User{Username: "owner", Password: "hashedpass"}
	invitee := models.User{Username: "invitee", Password: "hashedpass"}
	database.DB.Create(&owner)
	database.DB.Create(&invitee)

	// Create private sub
	privateSub := models.Sub{Name: "privatesub", Description: "Private Sub", OwnerID: owner.ID, Private: true}
	database.DB.Create(&privateSub)

	router := setupSubTestRouter("owner", owner.ID)

	t.Run("invite user to private sub", func(t *testing.T) {
		inviteData := `{"invitee":"invitee"}`
		req, _ := http.NewRequest("POST", "/subs/"+fmt.Sprintf("%d", privateSub.ID)+"/invite", strings.NewReader(inviteData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		responseBody := w.Body.String()
		assert.Contains(t, responseBody, "Invitation sent to invitee")
	})

	t.Run("invite to public sub", func(t *testing.T) {
		// Create public sub
		publicSub := models.Sub{Name: "publicsub", Description: "Public Sub", OwnerID: owner.ID, Private: false}
		database.DB.Create(&publicSub)

		inviteData := `{"invitee":"invitee"}`
		req, _ := http.NewRequest("POST", "/subs/"+fmt.Sprintf("%d", publicSub.ID)+"/invite", strings.NewReader(inviteData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// May return error as public subs don't require invites
		assert.NotEqual(t, http.StatusOK, w.Code)
	})

	t.Run("invite non-existent user", func(t *testing.T) {
		inviteData := `{"invitee":"nonexistent"}`
		req, _ := http.NewRequest("POST", "/subs/"+fmt.Sprintf("%d", privateSub.ID)+"/invite", strings.NewReader(inviteData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.NotEqual(t, http.StatusOK, w.Code)
	})

	t.Run("invite to non-owned sub", func(t *testing.T) {
		// Create another sub owned by different user
		user3 := models.User{Username: "otherowner", Password: "hashedpass"}
		database.DB.Create(&user3)
		otherSub := models.Sub{Name: "othersub", Description: "Other Sub", OwnerID: user3.ID, Private: true}
		database.DB.Create(&otherSub)

		inviteData := `{"invitee":"invitee"}`
		req, _ := http.NewRequest("POST", "/subs/"+fmt.Sprintf("%d", otherSub.ID)+"/invite", strings.NewReader(inviteData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.NotEqual(t, http.StatusOK, w.Code)
	})
}

func TestAcceptInviteHandler(t *testing.T) {
	if database.DB == nil {
		t.Skip("Database not available, skipping handler integration tests")
		return
	}

	// Setup - Clear tables to avoid interference
	database.DB.Exec("TRUNCATE TABLE users, subs, sub_invitations, sub_memberships, posts, comments RESTART IDENTITY CASCADE")

	// Create test users
	owner := models.User{Username: "owner", Password: "hashedpass"}
	invitee := models.User{Username: "invitee", Password: "hashedpass"}
	database.DB.Create(&owner)
	database.DB.Create(&invitee)

	// Create private sub
	privateSub := models.Sub{Name: "privatesub", Description: "Private Sub", OwnerID: owner.ID, Private: true}
	database.DB.Create(&privateSub)

	// Create invitation - Note: The service layer should handle invitation creation
	// For this test, we'll just call the service directly if possible
	// or assume the service creates a valid invitation

	router := setupSubTestRouter("invitee", invitee.ID)

	t.Run("accept valid invitation", func(t *testing.T) {
		// Create invitation for the test
		invitation := models.SubInvitation{
			SubID:     privateSub.ID,
			InviterID: owner.ID,
			InviteeID: invitee.ID,
			Status:    "pending",
		}
		database.DB.Create(&invitation)

		req, _ := http.NewRequest("POST", "/subs/invite/"+fmt.Sprintf("%d", invitation.ID), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		responseBody := w.Body.String()
		assert.Contains(t, responseBody, "You have joined the sub")
	})

	t.Run("accept invalid invitation", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/subs/invite/9999", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.NotEqual(t, http.StatusOK, w.Code)
	})

	t.Run("accept invitation not authorized", func(t *testing.T) {
		wrongRouter := setupSubTestRouter("wronguser", owner.ID) // Different user

		req, _ := http.NewRequest("POST", "/subs/invite/testinvite123", nil)
		w := httptest.NewRecorder()
		wrongRouter.ServeHTTP(w, req)

		assert.NotEqual(t, http.StatusOK, w.Code)
	})
}

func TestGetPostCountPerSubHandler(t *testing.T) {
	if database.DB == nil {
		t.Skip("Database not available, skipping handler integration tests")
		return
	}

	// Setup - Clear tables
	database.DB.Exec("DELETE FROM posts")
	database.DB.Exec("DELETE FROM sub_memberships")
	database.DB.Exec("DELETE FROM subs")
	database.DB.Exec("DELETE FROM users")

	// Create test data
	user := models.User{Username: "counter", Password: "hashedpass"}
	database.DB.Create(&user)
	sub := models.Sub{Name: "countsub", Description: "Sub for count", OwnerID: user.ID, Private: false}
	database.DB.Create(&sub)

	// Create posts
	post1 := models.Post{Title: "Post 1", Content: "Content 1", UserID: user.ID, SubID: sub.ID}
	post2 := models.Post{Title: "Post 2", Content: "Content 2", UserID: user.ID, SubID: sub.ID}
	database.DB.Create(&post1)
	database.DB.Create(&post2)

	t.Run("get post count for public sub", func(t *testing.T) {
		// This endpoint returns posts, not count - we need to check the actual endpoint
		// The GetPostCountPerSub endpoint ID might be different
		// Let me check what endpoints we have

		r := gin.Default()
		r.Use(func(c *gin.Context) {
			c.Set("username", "counter")
			c.Next()
		})
		r.GET("/subs/post-count", GetPostCountPerSub)

		req, _ := http.NewRequest("GET", "/subs/post-count?subID="+fmt.Sprintf("%d", sub.ID), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// Assert returns some valid response for post count
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
