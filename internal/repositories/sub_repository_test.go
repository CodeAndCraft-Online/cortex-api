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

func TestUpdateSub(t *testing.T) {
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
	owner := models.User{Username: "subowner", Password: "password"}
	nonOwner := models.User{Username: "notowner", Password: "password"}
	database.DB.Create(&owner)
	database.DB.Create(&nonOwner)

	// Create test sub
	sub := models.Sub{Name: "updateablesub", Description: "Original description", OwnerID: owner.ID, Private: false}
	database.DB.Create(&sub)

	t.Run("update sub successfully", func(t *testing.T) {
		updateRequest := models.SubRequest{
			Description: "Updated description",
			Private:     true,
		}

		updatedSub, err := UpdateSub(sub.ID, owner.ID, updateRequest)

		assert.NoError(t, err)
		assert.NotNil(t, updatedSub)
		assert.Equal(t, "Updated description", updatedSub.Description)
		assert.True(t, updatedSub.Private)
		assert.Equal(t, sub.Name, updatedSub.Name)    // Name should not change
		assert.Equal(t, owner.ID, updatedSub.OwnerID) // Owner should not change
	})

	t.Run("non-owner cannot update sub", func(t *testing.T) {
		updateRequest := models.SubRequest{
			Description: "Unauthorized update",
			Private:     false,
		}

		updatedSub, err := UpdateSub(sub.ID, nonOwner.ID, updateRequest)

		assert.Error(t, err)
		assert.Nil(t, updatedSub)
		assert.Contains(t, err.Error(), "only the sub owner can update")
	})

	t.Run("update non-existent sub", func(t *testing.T) {
		updateRequest := models.SubRequest{
			Description: "Update non-existent",
			Private:     false,
		}

		updatedSub, err := UpdateSub(999, owner.ID, updateRequest)

		assert.Error(t, err)
		assert.Nil(t, updatedSub)
		assert.Contains(t, err.Error(), "sub not found")
	})
}

func TestDeleteSub(t *testing.T) {
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
	owner := models.User{Username: "deleteowner", Password: "password"}
	nonOwner := models.User{Username: "deletenonowner", Password: "password"}
	database.DB.Create(&owner)
	database.DB.Create(&nonOwner)

	// Create test sub
	sub := models.Sub{Name: "deletablesub", Description: "To be deleted", OwnerID: owner.ID, Private: false}
	database.DB.Create(&sub)

	t.Run("delete sub successfully", func(t *testing.T) {
		err := DeleteSub(sub.ID, owner.ID)

		assert.NoError(t, err)

		// Verify sub is deleted
		var deletedSub models.Sub
		err = database.DB.First(&deletedSub, sub.ID).Error
		assert.Error(t, err) // Should return an error since sub is deleted
	})

	t.Run("non-owner cannot delete sub", func(t *testing.T) {
		// Create another sub for testing
		sub2 := models.Sub{Name: "nondeletablesub", Description: "Cannot be deleted", OwnerID: owner.ID, Private: false}
		database.DB.Create(&sub2)

		err := DeleteSub(sub2.ID, nonOwner.ID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "only the sub owner can delete")

		// Verify sub still exists
		var stillExists models.Sub
		err = database.DB.First(&stillExists, sub2.ID).Error
		assert.NoError(t, err)
	})

	t.Run("delete non-existent sub", func(t *testing.T) {
		err := DeleteSub(999, owner.ID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "sub not found")
	})
}

func TestGetSubMembers(t *testing.T) {
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
	owner := models.User{Username: "memowner", Password: "password"}
	member1 := models.User{Username: "memmember1", Password: "password"}
	member2 := models.User{Username: "memmember2", Password: "password"}
	nonMember := models.User{Username: "memnonmember", Password: "password"}
	database.DB.Create(&owner)
	database.DB.Create(&member1)
	database.DB.Create(&member2)
	database.DB.Create(&nonMember)

	// Create test subs
	publicSub := models.Sub{Name: "publicwithmembers", OwnerID: owner.ID, Private: false}
	privateSub := models.Sub{Name: "privatewithmembers", OwnerID: owner.ID, Private: true}
	database.DB.Create(&publicSub)
	database.DB.Create(&privateSub)

	// Create memberships
	membership1 := models.SubMembership{SubID: publicSub.ID, UserID: member1.ID}
	membership2 := models.SubMembership{SubID: publicSub.ID, UserID: member2.ID}
	membership3 := models.SubMembership{SubID: privateSub.ID, UserID: member1.ID}
	database.DB.Create(&membership1)
	database.DB.Create(&membership2)
	database.DB.Create(&membership3)

	t.Run("get members of public sub", func(t *testing.T) {
		members, err := GetSubMembers(fmt.Sprintf("%d", publicSub.ID), nonMember.ID, false)

		assert.NoError(t, err)
		assert.NotNil(t, members)
		assert.Len(t, members, 2)

		memberUsernames := make(map[string]bool)
		for _, member := range members {
			memberUsernames[member.Username] = true
		}
		assert.True(t, memberUsernames["memmember1"])
		assert.True(t, memberUsernames["memmember2"])
	})

	t.Run("get members of private sub as non-member", func(t *testing.T) {
		members, err := GetSubMembers(fmt.Sprintf("%d", privateSub.ID), nonMember.ID, false)

		assert.Error(t, err)
		assert.Nil(t, members)
		assert.Contains(t, err.Error(), "you must be a member to view")
	})

	t.Run("get members of private sub as owner", func(t *testing.T) {
		members, err := GetSubMembers(fmt.Sprintf("%d", privateSub.ID), owner.ID, true)

		assert.NoError(t, err)
		assert.NotNil(t, members)
		assert.Len(t, members, 1)
		assert.Equal(t, "memmember1", members[0].Username)
	})

	t.Run("get members of private sub as member", func(t *testing.T) {
		members, err := GetSubMembers(fmt.Sprintf("%d", privateSub.ID), member1.ID, false)

		assert.NoError(t, err)
		assert.NotNil(t, members)
		assert.Len(t, members, 1)
		assert.Equal(t, "memmember1", members[0].Username)
	})

	t.Run("get members of non-existent sub", func(t *testing.T) {
		members, err := GetSubMembers("999", owner.ID, true)

		assert.Error(t, err)
		assert.Nil(t, members)
		assert.Contains(t, err.Error(), "sub not found")
	})
}

func TestGetPendingInvites(t *testing.T) {
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
	owner := models.User{Username: "inviteowner", Password: "password"}
	nonOwner := models.User{Username: "invitenonowner", Password: "password"}
	invitee1 := models.User{Username: "invitee1", Password: "password"}
	invitee2 := models.User{Username: "invitee2", Password: "password"}
	database.DB.Create(&owner)
	database.DB.Create(&nonOwner)
	database.DB.Create(&invitee1)
	database.DB.Create(&invitee2)

	// Create private sub
	sub := models.Sub{Name: "inviteprivatesub", OwnerID: owner.ID, Private: true}
	database.DB.Create(&sub)

	// Create pending invites
	invite1 := models.SubInvitation{SubID: sub.ID, InviterID: owner.ID, InviteeID: invitee1.ID, Status: "pending"}
	invite2 := models.SubInvitation{SubID: sub.ID, InviterID: owner.ID, InviteeID: invitee2.ID, Status: "pending"}
	database.DB.Create(&invite1)
	database.DB.Create(&invite2)

	t.Run("get pending invites as owner", func(t *testing.T) {
		invites, err := GetPendingInvites(fmt.Sprintf("%d", sub.ID), owner.ID)

		assert.NoError(t, err)
		assert.NotNil(t, invites)
		assert.Len(t, invites, 2)

		inviteeUsernames := make(map[string]bool)
		for _, invite := range invites {
			inviteeUsernames[invite.InviteeUsername] = true
		}
		assert.True(t, inviteeUsernames["invitee1"])
		assert.True(t, inviteeUsernames["invitee2"])
	})

	t.Run("non-owner cannot view pending invites", func(t *testing.T) {
		invites, err := GetPendingInvites(fmt.Sprintf("%d", sub.ID), nonOwner.ID)

		assert.Error(t, err)
		assert.Nil(t, invites)
		assert.Contains(t, err.Error(), "only the sub owner can view")
	})

	t.Run("get pending invites for non-existent sub", func(t *testing.T) {
		invites, err := GetPendingInvites("999", owner.ID)

		assert.Error(t, err)
		assert.Nil(t, invites)
		assert.Contains(t, err.Error(), "sub not found")
	})
}
