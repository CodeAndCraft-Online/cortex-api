package services

import (
	"errors"

	"github.com/CodeAndCraft-Online/cortex-api/internal/models"
	"github.com/CodeAndCraft-Online/cortex-api/internal/repositories"
)

// GetPostByID fetches a single post by ID
func GetPostByID(postID string) (*models.PostResponse, error) {
	postResponse, err := repositories.GetPostByID(postID)
	if err != nil {
		return nil, errors.New("post not found")
	}
	return postResponse, nil
}

// GetAllPosts fetches all posts
func GetAllPosts() ([]models.Post, error) {
	return repositories.FindAllPosts()
}

func CreatPost(username string, post models.Post) (*models.Post, error) {
	newPost, err := repositories.CreatPost(username, post)
	if err != nil {
		return nil, err
	}

	return newPost, nil
}

func GetPosts() (*[]models.PostResponse, error) {
	commentResponse, err := repositories.GetPosts()
	if err != nil {
		return nil, err
	}

	return commentResponse, nil
}
