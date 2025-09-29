package repositories

import (
	"fmt"
	"time"

	db "github.com/CodeAndCraft-Online/cortex-api/internal/database"
	"github.com/CodeAndCraft-Online/cortex-api/internal/models"
)

func GetSubs(username string) (*[]models.Sub, error) {
	var subs []models.Sub

	if username == "" {
		// ✅ User is not authenticated → Return only public subs
		if err := db.DB.Where("private = ?", false).Order("created_at DESC").Find(&subs).Error; err != nil {
			return nil, fmt.Errorf("failed to fetch subs")
		}
	} else {
		// ✅ User is authenticated → Fetch user data
		var user models.User
		if err := db.DB.Where("username = ?", username).First(&user).Error; err != nil {
			return nil, fmt.Errorf("user not found")
		}

		// ✅ Fetch public subs + private subs where user is the owner or a member
		if err := db.DB.Raw(`
			SELECT DISTINCT subs.* 
			FROM subs
			LEFT JOIN sub_memberships ON subs.id = sub_memberships.sub_id
			WHERE subs.private = false
			OR subs.owner_id = ?
			OR (sub_memberships.user_id = ? AND sub_memberships.sub_id IS NOT NULL)
			ORDER BY subs.created_at DESC
		`, user.ID, user.ID).Scan(&subs).Error; err != nil {
			return nil, fmt.Errorf("failed to fetch subs")
		}
	}

	return &subs, nil
}

func CreateSub(username string, subRequest models.SubRequest) (*models.Sub, error) {
	// Fetch user ID from the database
	var user models.User
	if err := db.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Ensure the sub name is unique
	var existingSub models.Sub
	if err := db.DB.Where("name = ?", subRequest.Name).First(&existingSub).Error; err == nil {
		return nil, fmt.Errorf("sub name already taken")
	}

	// Create the sub
	newSub := models.Sub{
		Name:        subRequest.Name,
		Description: subRequest.Description,
		OwnerID:     user.ID,
		Private:     subRequest.Private,
		CreatedAt:   time.Now(),
	}

	db.DB.Create(&newSub)

	return &newSub, nil
}

func JoinSub(username, subID string) (*models.SubMembership, error) {
	// Fetch user ID
	var user models.User
	if err := db.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Check if the sub exists
	var sub models.Sub
	if err := db.DB.First(&sub, subID).Error; err != nil {
		return nil, fmt.Errorf("sub not found")
	}

	// ✅ If the sub is private, check for an invitation
	if sub.Private {
		var invitation models.SubInvitation
		if err := db.DB.Where("sub_id = ? AND invitee_id = ? AND status = ?", sub.ID, user.ID, "pending").First(&invitation).Error; err != nil {
			return nil, fmt.Errorf("you need an invitation to join this private sub")
		}

		// ✅ Mark the invitation as accepted
		invitation.Status = "accepted"
		db.DB.Save(&invitation)
	}

	// ✅ Add user to sub_memberships
	newMembership := models.SubMembership{
		SubID:    sub.ID,
		UserID:   user.ID,
		JoinedAt: time.Now(),
	}
	db.DB.Create(&newMembership)

	return &newMembership, nil
}

func InviteUser(subID, username string, inviteRequest models.InviteRequest) error {

	// Fetch user ID (inviter)
	var inviter models.User
	if err := db.DB.Where("username = ?", username).First(&inviter).Error; err != nil {
		return fmt.Errorf("user not found")
	}

	// Check if the sub exists
	var sub models.Sub
	if err := db.DB.First(&sub, subID).Error; err != nil {
		return fmt.Errorf("sub not found")
	}

	// ✅ Ensure invitations are only for private subs
	if !sub.Private {
		return fmt.Errorf("can only invite users to private subs")
	}

	// ✅ Ensure the inviter is the sub owner
	if sub.OwnerID != inviter.ID {
		return fmt.Errorf("only the owner can invite users")
	}

	// Fetch invitee user
	var invitee models.User
	if err := db.DB.Where("username = ?", inviteRequest.InviteeUsername).First(&invitee).Error; err != nil {
		return fmt.Errorf("invitee user not found")
	}

	// ✅ Check if an invitation already exists
	var existingInvite models.SubInvitation
	if err := db.DB.Where("sub_id = ? AND invitee_id = ?", sub.ID, invitee.ID).First(&existingInvite).Error; err == nil {
		return fmt.Errorf("user is already invited")
	}

	// ✅ Create invitation
	newInvite := models.SubInvitation{
		SubID:     sub.ID,
		InviterID: inviter.ID,
		InviteeID: invitee.ID,
		Status:    "pending",
		CreatedAt: time.Now(),
	}
	db.DB.Create(&newInvite)

	return nil
}

func AcceptInvite(inviteID, username string) error {
	// Fetch user ID
	var user models.User
	if err := db.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return fmt.Errorf("user not found")
	}

	// Find invitation
	var invitation models.SubInvitation
	if err := db.DB.First(&invitation, inviteID).Error; err != nil {
		return fmt.Errorf("invitation not found")
	}

	// ✅ Ensure the user is the invitee
	if invitation.InviteeID != user.ID {
		return fmt.Errorf("you are not the invitee for this invitation")
	}

	// ✅ Accept invitation
	invitation.Status = "accepted"
	db.DB.Save(&invitation)

	// ✅ Add user to sub_memberships
	newMembership := models.SubMembership{
		SubID:    invitation.SubID,
		UserID:   user.ID,
		JoinedAt: time.Now(),
	}
	db.DB.Create(&newMembership)

	return nil
}

func ListSubPosts(subID, username string) (*[]models.PostResponse, error) {
	// Check if the sub exists
	var sub models.Sub
	if err := db.DB.First(&sub, subID).Error; err != nil {
		return nil, fmt.Errorf("sub not found")
	}

	// ✅ If the sub is private, check if the user is a member
	if sub.Private {
		if username == "" {
			return nil, fmt.Errorf("this is a private sub. Please log in to view posts")
		}

		var user models.User
		if err := db.DB.Where("username = ?", username).First(&user).Error; err != nil {
			return nil, fmt.Errorf("user not found")
		}

		var count int64
		db.DB.Model(&models.SubMembership{}).Where("sub_id = ? AND user_id = ?", sub.ID, user.ID).Count(&count)
		if count == 0 {
			return nil, fmt.Errorf("you must be a member to view this sub")
		}
	}

	// Fetch posts from the sub
	var posts []models.Post
	if err := db.DB.Preload("User").Where("sub_id = ?", subID).Order("created_at DESC").Find(&posts).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch posts")
	}

	// Convert posts to response format
	var formattedPosts []models.PostResponse
	for _, post := range posts {
		var upvotes, downvotes int64

		// Count upvotes and downvotes
		db.DB.Model(&models.Vote{}).Where("post_id = ? AND vote = 1", post.ID).Count(&upvotes)
		db.DB.Model(&models.Vote{}).Where("post_id = ? AND vote = -1", post.ID).Count(&downvotes)

		formattedPosts = append(formattedPosts, models.PostResponse{
			ID:        post.ID,
			Title:     post.Title,
			Content:   post.Content,
			Username:  post.User.Username,
			Upvotes:   int(upvotes),
			Downvotes: int(downvotes),
			SubID:     post.SubID,
			CreatedAt: post.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &formattedPosts, nil
}

func LeaveSub(subID, username string) (*models.Sub, error) {
	// Fetch user ID
	var user models.User
	if err := db.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Check if the sub exists
	var sub models.Sub
	if err := db.DB.First(&sub, subID).Error; err != nil {
		return nil, fmt.Errorf("sub not found")
	}

	// Check if the user is actually a member
	var count int64
	db.DB.Model(&models.SubMembership{}).Where("sub_id = ? AND user_id = ?", sub.ID, user.ID).Count(&count)
	if count == 0 {
		return nil, fmt.Errorf("you are not a member of this sub")
	}

	// Remove the membership
	if err := db.DB.Where("sub_id = ? AND user_id = ?", sub.ID, user.ID).Delete(&models.SubMembership{}).Error; err != nil {
		return nil, fmt.Errorf("failed to leave sub")
	}

	return &sub, nil
}

func GetPostCountPerSub(subID, username string) (int, error) {
	// Check if the sub exists
	var sub models.Sub
	if err := db.DB.First(&sub, subID).Error; err != nil {
		return -1, fmt.Errorf("sub not found")
	}

	// ✅ If the sub is private, check if the user is a member
	if sub.Private {
		if username == "" {
			return -1, fmt.Errorf("this is a private sub. Please log in to view post count")
		}

		var user models.User
		if err := db.DB.Where("username = ?", username).First(&user).Error; err != nil {
			return -1, fmt.Errorf("user not found")
		}

		var count int64
		db.DB.Model(&models.SubMembership{}).Where("sub_id = ? AND user_id = ?", sub.ID, user.ID).Count(&count)
		if count == 0 {
			return -1, fmt.Errorf("you must be a member to view this sub")
		}
	}

	// Fetch posts from the sub
	var posts []models.Post
	if err := db.DB.Preload("User").Where("sub_id = ?", subID).Order("created_at DESC").Find(&posts).Error; err != nil {
		return -1, fmt.Errorf("failed to fetch posts")
	}

	postCount := len(posts)

	return postCount, nil
}

func UpdateSub(subID, ownerID uint, updateRequest models.SubRequest) (*models.Sub, error) {
	// Check if the sub exists and verify ownership
	var sub models.Sub
	if err := db.DB.Preload("Owner").First(&sub, subID).Error; err != nil {
		return nil, fmt.Errorf("sub not found")
	}

	// Verify the user is the owner
	if sub.OwnerID != ownerID {
		return nil, fmt.Errorf("only the sub owner can update the sub")
	}

	// Update allowed fields (only description and private flag)
	sub.Description = updateRequest.Description
	sub.Private = updateRequest.Private

	// Save the updated sub
	if err := db.DB.Save(&sub).Error; err != nil {
		return nil, fmt.Errorf("failed to update sub")
	}

	return &sub, nil
}

func DeleteSub(subID, ownerID uint) error {
	// Check if the sub exists and verify ownership
	var sub models.Sub
	if err := db.DB.First(&sub, subID).Error; err != nil {
		return fmt.Errorf("sub not found")
	}

	// Verify the user is the owner
	if sub.OwnerID != ownerID {
		return fmt.Errorf("only the sub owner can delete the sub")
	}

	// Delete the sub (cascade delete will handle related records)
	if err := db.DB.Delete(&sub).Error; err != nil {
		return fmt.Errorf("failed to delete sub")
	}

	return nil
}

func GetSubMembers(subID string, userID uint, isOwner bool) ([]models.SubMemberResponse, error) {
	// Check if the sub exists
	var sub models.Sub
	if err := db.DB.First(&sub, subID).Error; err != nil {
		return nil, fmt.Errorf("sub not found")
	}

	// Access control: private subs can only be viewed by members/owners
	if sub.Private {
		isMember := false
		if isOwner && sub.OwnerID == userID {
			isMember = true
		} else {
			var count int64
			db.DB.Model(&models.SubMembership{}).Where("sub_id = ? AND user_id = ?", sub.ID, userID).Count(&count)
			if count > 0 {
				isMember = true
			}
		}
		if !isMember {
			return nil, fmt.Errorf("you must be a member to view this sub's members")
		}
	}

	// Get all members of the sub
	var memberships []models.SubMembership
	if err := db.DB.Preload("User").Where("sub_id = ?", subID).Order("joined_at ASC").Find(&memberships).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch sub members")
	}

	// Convert to response format
	var memberResponses []models.SubMemberResponse
	for _, membership := range memberships {
		memberResponses = append(memberResponses, models.SubMemberResponse{
			Username: membership.User.Username,
			JoinedAt: membership.JoinedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return memberResponses, nil
}

func GetPendingInvites(subID string, ownerID uint) ([]models.InviteResponse, error) {
	// Check if the sub exists and verify ownership
	var sub models.Sub
	if err := db.DB.First(&sub, subID).Error; err != nil {
		return nil, fmt.Errorf("sub not found")
	}

	// Only sub owners can view pending invites
	if sub.OwnerID != ownerID {
		return nil, fmt.Errorf("only the sub owner can view pending invites")
	}

	// Get all pending invites for this sub
	var invites []models.SubInvitation
	if err := db.DB.Preload("Invitee").Where("sub_id = ? AND status = ?", subID, "pending").Order("created_at ASC").Find(&invites).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch pending invites")
	}

	// Convert to response format
	var inviteResponses []models.InviteResponse
	for _, invite := range invites {
		inviteResponses = append(inviteResponses, models.InviteResponse{
			InviteeUsername: invite.Invitee.Username,
			CreatedAt:       invite.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return inviteResponses, nil
}
