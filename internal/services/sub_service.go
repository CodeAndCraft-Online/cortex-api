package services

import (
	"fmt"
	"strconv"

	db "github.com/CodeAndCraft-Online/cortex-api/internal/database"
	"github.com/CodeAndCraft-Online/cortex-api/internal/models"
	"github.com/CodeAndCraft-Online/cortex-api/internal/repositories"
)

func GetSubs(username string) (*[]models.Sub, error) {
	subs, err := repositories.GetSubs(username)
	if err != nil {
		return nil, err
	}

	return subs, nil
}

func CreateSub(username string, subRequest models.SubRequest) (*models.Sub, error) {
	newSub, err := repositories.CreateSub(username, subRequest)
	if err != nil {
		return nil, err
	}

	return newSub, nil
}

func JoinSub(username, subID string) (*models.SubMembership, error) {
	sub, err := repositories.JoinSub(username, subID)
	if err != nil {
		return nil, err
	}

	return sub, nil
}

func InviteUser(subID, username string, inviteRequest models.InviteRequest) error {
	err := repositories.InviteUser(subID, username, inviteRequest)
	if err != nil {
		return err
	}

	return nil
}

func AcceptInvite(inviteID, username string) error {
	err := repositories.AcceptInvite(inviteID, username)
	if err != nil {
		return err
	}

	return nil
}

func ListSubPosts(subID, username string) (*[]models.PostResponse, error) {
	subPosts, err := repositories.ListSubPosts(subID, username)
	if err != nil {
		return nil, err
	}

	return subPosts, nil
}

func LeaveSub(subID, username string) (*models.Sub, error) {
	sub, err := repositories.LeaveSub(subID, username)
	if err != nil {
		return nil, err
	}

	return sub, nil
}

func GetPostCountPerSub(subID, username string) (int, error) {
	postCount, err := repositories.GetPostCountPerSub(subID, username)
	if err != nil {
		return -1, err
	}

	return postCount, nil
}

func UpdateSub(subID, username string, updateRequest models.SubRequest) (*models.Sub, error) {
	// Get user ID
	var user models.User
	if err := db.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Convert subID string to uint
	subIDUint, err := strconv.ParseUint(subID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid sub ID")
	}

	sub, err := repositories.UpdateSub(uint(subIDUint), user.ID, updateRequest)
	if err != nil {
		return nil, err
	}

	return sub, nil
}

func DeleteSub(subID, username string) error {
	// Get user ID
	var user models.User
	if err := db.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return fmt.Errorf("user not found")
	}

	// Convert subID string to uint
	subIDUint, err := strconv.ParseUint(subID, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid sub ID")
	}

	err = repositories.DeleteSub(uint(subIDUint), user.ID)
	if err != nil {
		return err
	}

	return nil
}

func GetSubMembers(subID, username string) ([]models.SubMemberResponse, error) {
	// Get user information for access control
	var user models.User
	var isOwner bool

	if username != "" {
		if err := db.DB.Where("username = ?", username).First(&user).Error; err != nil {
			return nil, fmt.Errorf("user not found")
		}

		// Check if user is the owner (only needed for private subs, but calculate once)
		var sub models.Sub
		if err := db.DB.Where("name = ?", subID).First(&sub).Error; err != nil {
			return nil, fmt.Errorf("sub not found")
		}
		isOwner = sub.OwnerID == user.ID
	}

	members, err := repositories.GetSubMembers(subID, user.ID, isOwner)
	if err != nil {
		return nil, err
	}

	return members, nil
}

func GetPendingInvites(subID, username string) ([]models.InviteResponse, error) {
	// Get user ID
	var user models.User
	if err := db.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, fmt.Errorf("user not found")
	}

	invites, err := repositories.GetPendingInvites(subID, user.ID)
	if err != nil {
		return nil, err
	}

	return invites, nil
}
