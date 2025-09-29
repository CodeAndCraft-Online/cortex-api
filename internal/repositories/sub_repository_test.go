package repositories

import (
	"fmt"
	"testing"

	"github.com/CodeAndCraft-Online/cortex-api/internal/database"
	"github.com/CodeAndCraft-Online/cortex-api/internal/models"
	"github.com/stretchr/testify/assert"
)

// TestMain is defined in auth_repository_test.go for the repositories package

func TestGetSubs(t *testing.T) {
	if database.DB == nil {
		t.Skip("Database not available, skipping repository integration tests")
		return
	}

	// Clear tables to avoid conflicts
	database.DB.Exec("DELETE FROM sub_memberships")
	database.DB.Exec("DELETE FROM sub_invitations")
	database.DB.Exec("DELETE FROM subs")
	database.DB.Exec("DELETE FROM users")

	// Create test users
	user1 := models.User{Username: "user1", Password: "password"}
	user2 := models.User{Username: "user2", Password: "password"}
	database.DB.Create(&user1)
	database.DB.Create(&user2)

	// Create test subs
	publicSub1 := models.Sub{Name: "publicsub1", Description: "Public Sub 1", OwnerID: user1.ID, Private: false}
	publicSub2 := models.Sub{Name: "publicsub2", Description: "Public Sub 2", OwnerID: user1.ID, Private: false}
	privateSub := models.Sub{Name: "privatesub", Description: "Private Sub", OwnerID: user1.ID, Private: true}
	database.DB.Create(&publicSub1)
	database.DB.Create(&publicSub2)
	database.DB.Create(&privateSub)

	t.Run("get subs for unauthenticated user", func(t *testing.T) {
		subs, err := GetSubs("")

		assert.NoError(t, err)
		assert.NotNil(t, subs)
		assert.True(t, len(*subs) >= 2)

		// Should only include public subs
		for _, sub := range *subs {
			assert.False(t, sub.Private)
		}
	})

	t.Run("get subs for authenticated user", func(t *testing.T) {
		subs, err := GetSubs("user1")

		assert.NoError(t, err)
		assert.NotNil(t, subs)
		assert.True(t, len(*subs) >= 3) // Public subs + private sub where user is owner

		// Should include public subs and owned private subs
		subNames := make(map[string]bool)
		for _, sub := range *subs {
			subNames[sub.Name] = true
		}
		assert.True(t, subNames["publicsub1"])
		assert.True(t, subNames["publicsub2"])
		assert.True(t, subNames["privatesub"]) // User is owner
	})

	t.Run("get subs for non-existent user", func(t *testing.T) {
		subs, err := GetSubs("nonexistent")

		assert.Error(t, err)
		assert.Nil(t, subs)
		assert.Contains(t, err.Error(), "user not found")
	})
}

func TestCreateSub(t *testing.T) {
	if database.DB == nil {
		t.Skip("Database not available, skipping repository integration tests")
		return
	}

	// Clear tables to avoid conflicts
	database.DB.Exec("DELETE FROM sub_memberships")
	database.DB.Exec("DELETE FROM sub_invitations")
	database.DB.Exec("DELETE FROM subs")
	database.DB.Exec("DELETE FROM users")

	// Create test user
	user := models.User{Username: "createsubuser", Password: "password"}
	database.DB.Create(&user)

	t.Run("create public sub successfully", func(t *testing.T) {
		subRequest := models.SubRequest{
			Name:        "newpublicsub",
			Description: "New public sub",
			Private:     false,
		}

		sub, err := CreateSub("createsubuser", subRequest)

		assert.NoError(t, err)
		assert.NotNil(t, sub)
		assert.Equal(t, subRequest.Name, sub.Name)
		assert.Equal(t, subRequest.Description, sub.Description)
		assert.Equal(t, user.ID, sub.OwnerID)
		assert.False(t, sub.Private)
	})

	t.Run("create private sub successfully", func(t *testing.T) {
		subRequest := models.SubRequest{
			Name:        "newprivatesub",
			Description: "New private sub",
			Private:     true,
		}

		sub, err := CreateSub("createsubuser", subRequest)

		assert.NoError(t, err)
		assert.NotNil(t, sub)
		assert.Equal(t, subRequest.Name, sub.Name)
		assert.True(t, sub.Private)
	})

	t.Run("create sub with duplicate name", func(t *testing.T) {
		subRequest := models.SubRequest{
			Name:        "newpublicsub", // Already exists
			Description: "Duplicate name",
			Private:     false,
		}

		sub, err := CreateSub("createsubuser", subRequest)

		assert.Error(t, err)
		assert.Nil(t, sub)
		assert.Contains(t, err.Error(), "sub name already taken")
	})
}

func TestJoinSub(t *testing.T) {
	if database.DB == nil {
		t.Skip("Database not available, skipping repository integration tests")
		return
	}

	// Clear tables to avoid conflicts
	database.DB.Exec("DELETE FROM sub_memberships")
	database.DB.Exec("DELETE FROM sub_invitations")
	database.DB.Exec("DELETE FROM subs")
	database.DB.Exec("DELETE FROM users")

	// Create test users
	user := models.User{Username: "joinuser", Password: "password"}
	otherUser := models.User{Username: "otheruser", Password: "password"}
	database.DB.Create(&user)
	database.DB.Create(&otherUser)

	// Create test subs
	publicSub := models.Sub{Name: "publicjoin", OwnerID: otherUser.ID, Private: false}
	privateSub := models.Sub{Name: "privatejoin", OwnerID: otherUser.ID, Private: true}
	database.DB.Create(&publicSub)
	database.DB.Create(&privateSub)

	t.Run("join public sub", func(t *testing.T) {
		membership, err := JoinSub("joinuser", fmt.Sprintf("%d", publicSub.ID))

		assert.NoError(t, err)
		assert.NotNil(t, membership)
		assert.Equal(t, publicSub.ID, membership.SubID)
		assert.Equal(t, user.ID, membership.UserID)
	})

	t.Run("join private sub without invitation", func(t *testing.T) {
		membership, err := JoinSub("joinuser", fmt.Sprintf("%d", privateSub.ID))

		assert.Error(t, err)
		assert.Nil(t, membership)
		assert.Contains(t, err.Error(), "you need an invitation")
	})

	t.Run("join non-existent sub", func(t *testing.T) {
		membership, err := JoinSub("joinuser", "999")

		assert.Error(t, err)
		assert.Nil(t, membership)
		assert.Contains(t, err.Error(), "sub not found")
	})
}

func TestListSubPosts(t *testing.T) {
	if database.DB == nil {
		t.Skip("Database not available, skipping repository integration tests")
		return
	}

	// Clear tables to avoid conflicts
	database.DB.Exec("DELETE FROM posts")
	database.DB.Exec("DELETE FROM votes")
	database.DB.Exec("DELETE FROM sub_memberships")
	database.DB.Exec("DELETE FROM sub_invitations")
	database.DB.Exec("DELETE FROM subs")
	database.DB.Exec("DELETE FROM users")

	// Create test users
	user := models.User{Username: "postuser", Password: "password"}
	memberUser := models.User{Username: "memberuser", Password: "password"}
	database.DB.Create(&user)
	database.DB.Create(&memberUser)

	// Create test subs
	publicSub := models.Sub{Name: "postpublic", OwnerID: user.ID, Private: false}
	privateSub := models.Sub{Name: "postprivate", OwnerID: user.ID, Private: true}
	database.DB.Create(&publicSub)
	database.DB.Create(&privateSub)

	// Create test post in public sub
	publicPost := models.Post{
		Title:   "Public Post",
		Content: "Public content",
		UserID:  user.ID,
		SubID:   publicSub.ID,
	}
	database.DB.Create(&publicPost)

	t.Run("list posts from public sub", func(t *testing.T) {
		posts, err := ListSubPosts(fmt.Sprintf("%d", publicSub.ID), "")

		assert.NoError(t, err)
		assert.NotNil(t, posts)
		assert.True(t, len(*posts) >= 1)

		firstPost := (*posts)[0]
		assert.Equal(t, "Public Post", firstPost.Title)
		assert.Equal(t, "Public content", firstPost.Content)
	})

	t.Run("list posts from private sub without membership", func(t *testing.T) {
		posts, err := ListSubPosts(fmt.Sprintf("%d", privateSub.ID), "")

		assert.Error(t, err)
		assert.Nil(t, posts)
		assert.Contains(t, err.Error(), "Please log in to view posts")
	})

	t.Run("list posts from non-existent sub", func(t *testing.T) {
		posts, err := ListSubPosts("999", "")

		assert.Error(t, err)
		assert.Nil(t, posts)
		assert.Contains(t, err.Error(), "sub not found")
	})
}

func TestLeaveSub(t *testing.T) {
	if database.DB == nil {
		t.Skip("Database not available, skipping repository integration tests")
		return
	}

	// Clear tables to avoid conflicts
	database.DB.Exec("DELETE FROM sub_memberships")
	database.DB.Exec("DELETE FROM sub_invitations")
	database.DB.Exec("DELETE FROM subs")
	database.DB.Exec("DELETE FROM users")

	// Create test users
	user := models.User{Username: "leaveuser", Password: "password"}
	database.DB.Create(&user)

	// Create test sub
	sub := models.Sub{Name: "leavesub", OwnerID: user.ID, Private: false}
	database.DB.Create(&sub)

	// Create membership
	membership := models.SubMembership{SubID: sub.ID, UserID: user.ID}
	database.DB.Create(&membership)

	t.Run("leave sub successfully", func(t *testing.T) {
		result, err := LeaveSub(fmt.Sprintf("%d", sub.ID), "leaveuser")

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, sub.Name, result.Name)
	})

	t.Run("leave non-existent sub", func(t *testing.T) {
		result, err := LeaveSub("999", "leaveuser")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "sub not found")
	})
}
