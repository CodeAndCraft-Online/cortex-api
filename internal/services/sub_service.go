package services

import (
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
