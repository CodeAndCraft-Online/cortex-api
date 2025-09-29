package services

import (
	"fmt"
	"testing"

	"github.com/CodeAndCraft-Online/cortex-api/internal/database"
	"github.com/CodeAndCraft-Online/cortex-api/internal/models"
	"github.com/stretchr/testify/assert"
)

// TestMain is defined in auth_service_test.go

func TestGetSubs_Service(t *testing.T) {
	if database.DB == nil {
		t.Skip("Database not available, skipping integration test")
		return
	}

	// Clear tables to avoid conflicts
	database.DB.Exec("DELETE FROM comments")
	database.DB.Exec("DELETE FROM sub_memberships")
	database.DB.Exec("DELETE FROM sub_invitations")
	database.DB.Exec("DELETE FROM posts")
	database.DB.Exec("DELETE FROM subs")
	database.DB.Exec("DELETE FROM users")

	// Create test user
	user := models.User{Username: "subsuser", Password: "password"}
	database.DB.Create(&user)

	// Create test subs
	sub1 := models.Sub{Name: "sub1", OwnerID: user.ID}
	sub2 := models.Sub{Name: "sub2", OwnerID: user.ID}
	database.DB.Create(&sub1)
	database.DB.Create(&sub2)

	// Memberships for user
	membership1 := models.SubMembership{UserID: user.ID, SubID: sub1.ID}
	membership2 := models.SubMembership{UserID: user.ID, SubID: sub2.ID}
	database.DB.Create(&membership1)
	database.DB.Create(&membership2)

	subs, err := GetSubs("subsuser")

	assert.NoError(t, err)
	assert.NotNil(t, subs)
	assert.GreaterOrEqual(t, len(*subs), 2)
}

func TestCreateSub_Service(t *testing.T) {
	if database.DB == nil {
		t.Skip("Database not available, skipping integration test")
		return
	}

	// Clear tables to avoid conflicts
	database.DB.Exec("DELETE FROM subs")
	database.DB.Exec("DELETE FROM users")

	// Create test user
	user := models.User{Username: "createsubuser", Password: "password"}
	database.DB.Create(&user)

	subRequest := models.SubRequest{
		Name:        "newsub",
		Description: "A new sub",
	}

	createdSub, err := CreateSub("createsubuser", subRequest)

	assert.NoError(t, err)
	assert.NotNil(t, createdSub)
	assert.Equal(t, "newsub", createdSub.Name)
	assert.Equal(t, user.ID, createdSub.OwnerID)
}

func TestJoinSub_Service(t *testing.T) {
	if database.DB == nil {
		t.Skip("Database not available, skipping integration test")
		return
	}

	// Clear tables to avoid conflicts
	database.DB.Exec("DELETE FROM sub_memberships")
	database.DB.Exec("DELETE FROM subs")
	database.DB.Exec("DELETE FROM users")

	// Create test user and sub
	user := models.User{Username: "joinsubuser", Password: "password"}
	database.DB.Create(&user)

	sub := models.Sub{Name: "joinsub", OwnerID: user.ID}
	database.DB.Create(&sub)

	membership, err := JoinSub("joinsubuser", fmt.Sprintf("%d", sub.ID))

	assert.NoError(t, err)
	assert.NotNil(t, membership)
	assert.Equal(t, user.ID, membership.UserID)
}

func TestInviteUser_Service(t *testing.T) {
	if database.DB == nil {
		t.Skip("Database not available, skipping integration test")
		return
	}

	// Clear tables to avoid conflicts
	database.DB.Exec("DELETE FROM sub_invitations")
	database.DB.Exec("DELETE FROM subs")
	database.DB.Exec("DELETE FROM users")

	// Create owner and invitee
	owner := models.User{Username: "owner", Password: "password"}
	database.DB.Create(&owner)

	invitee := models.User{Username: "invitee", Password: "password"}
	database.DB.Create(&invitee)

	sub := models.Sub{Name: "invitesub", OwnerID: owner.ID, Private: true}
	database.DB.Create(&sub)

	inviteRequest := models.InviteRequest{
		InviteeUsername: "invitee",
	}

	err := InviteUser(fmt.Sprintf("%d", sub.ID), "owner", inviteRequest)

	assert.NoError(t, err)
}

func TestAcceptInvite_Service(t *testing.T) {
	if database.DB == nil {
		t.Skip("Database not available, skipping integration test")
		return
	}

	// Clear tables to avoid conflicts
	database.DB.Exec("DELETE FROM sub_memberships")
	database.DB.Exec("DELETE FROM sub_invitations")
	database.DB.Exec("DELETE FROM subs")
	database.DB.Exec("DELETE FROM users")

	// Create invite
	owner := models.User{Username: "acceptowner", Password: "password"}
	database.DB.Create(&owner)

	invitee := models.User{Username: "acceptinvitee", Password: "password"}
	database.DB.Create(&invitee)

	sub := models.Sub{Name: "acceptsub", OwnerID: owner.ID}
	database.DB.Create(&sub)

	invite := models.SubInvitation{SubID: sub.ID, InviterID: owner.ID, InviteeID: invitee.ID}
	database.DB.Create(&invite)

	err := AcceptInvite(fmt.Sprintf("%d", invite.ID), "acceptinvitee")

	assert.NoError(t, err)
}

func TestListSubPosts_Service(t *testing.T) {
	if database.DB == nil {
		t.Skip("Database not available, skipping integration test")
		return
	}

	// Clear tables to avoid conflicts
	database.DB.Exec("DELETE FROM posts")
	database.DB.Exec("DELETE FROM subs")
	database.DB.Exec("DELETE FROM users")

	// Create test data
	user := models.User{Username: "listsubuser", Password: "password"}
	database.DB.Create(&user)

	sub := models.Sub{Name: "listsub", OwnerID: user.ID}
	database.DB.Create(&sub)

	post := models.Post{
		Title:   "Sub Post",
		Content: "Content",
		UserID:  user.ID,
		SubID:   sub.ID,
	}
	database.DB.Create(&post)

	posts, err := ListSubPosts(fmt.Sprintf("%d", sub.ID), "listsubuser")

	assert.NoError(t, err)
	assert.NotNil(t, posts)
	assert.GreaterOrEqual(t, len(*posts), 1)
}

func TestLeaveSub_Service(t *testing.T) {
	if database.DB == nil {
		t.Skip("Database not available, skipping integration test")
		return
	}

	// Clear tables to avoid conflicts
	database.DB.Exec("DELETE FROM sub_memberships")
	database.DB.Exec("DELETE FROM subs")
	database.DB.Exec("DELETE FROM users")

	// Create user and sub
	user := models.User{Username: "leavesubuser", Password: "password"}
	database.DB.Create(&user)

	sub := models.Sub{Name: "leavesub", OwnerID: user.ID}
	database.DB.Create(&sub)

	membership := models.SubMembership{UserID: user.ID, SubID: sub.ID}
	database.DB.Create(&membership)

	leftSub, err := LeaveSub(fmt.Sprintf("%d", sub.ID), "leavesubuser")

	assert.NoError(t, err)
	assert.NotNil(t, leftSub)
	assert.Equal(t, "leavesub", leftSub.Name)
}

func TestGetPostCountPerSub_Service(t *testing.T) {
	if database.DB == nil {
		t.Skip("Database not available, skipping integration test")
		return
	}

	// Clear tables to avoid conflicts
	database.DB.Exec("DELETE FROM posts")
	database.DB.Exec("DELETE FROM subs")
	database.DB.Exec("DELETE FROM users")

	// Create test data
	user := models.User{Username: "countuser", Password: "password"}
	database.DB.Create(&user)

	sub := models.Sub{Name: "countsub", OwnerID: user.ID}
	database.DB.Create(&sub)

	post := models.Post{
		Title:   "Count Post",
		Content: "Content",
		UserID:  user.ID,
		SubID:   sub.ID,
	}
	database.DB.Create(&post)

	count, err := GetPostCountPerSub(fmt.Sprintf("%d", sub.ID), "countuser")

	assert.NoError(t, err)
	assert.Equal(t, 1, count)
}
