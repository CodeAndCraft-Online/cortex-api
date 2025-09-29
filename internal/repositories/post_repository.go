package repositories

import (
	"errors"
	"fmt"

	db "github.com/CodeAndCraft-Online/cortex-api/internal/database"
	"github.com/CodeAndCraft-Online/cortex-api/internal/models"
)

// FindPostByID retrieves a single post from the database by ID
func GetPostByID(postID string) (*models.PostResponse, error) {
	var upvotes, downvotes int64

	var post models.Post
	if err := db.DB.Preload("User").Where("id = ?", postID).First(&post).Error; err != nil {
		return nil, errors.New("post not found")
	}

	// Count upvotes (vote = 1)
	db.DB.Model(&models.Vote{}).Where("post_id = ? AND vote = 1", post.ID).Count(&upvotes)
	// Count downvotes (vote = -1)
	db.DB.Model(&models.Vote{}).Where("post_id = ? AND vote = -1", post.ID).Count(&downvotes)

	// ✅ Return a properly formatted PostResponse
	postResponse := models.PostResponse{
		ID:        post.ID,
		Title:     post.Title,
		Content:   post.Content,
		ImageURL:  post.ImageURL,
		Upvotes:   int(upvotes),
		Downvotes: int(downvotes),
		Username:  post.User.Username, // ✅ Ensure "User" is preloaded
		CreatedAt: post.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	return &postResponse, nil
}

// FindAllPosts retrieves all posts
func FindAllPosts() ([]models.Post, error) {
	var posts []models.Post
	if err := db.DB.Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

func CreatePost(username string, post models.Post) (*models.Post, error) {

	var user models.User
	if err := db.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, fmt.Errorf("user not found")
	}
	post.UserID = user.ID

	// Validate that the sub exists
	var sub models.Sub
	if err := db.DB.Where("id = ?", post.SubID).First(&sub).Error; err != nil {
		return nil, fmt.Errorf("sub not found")
	}

	// Save post to the database
	if err := db.DB.Create(&post).Error; err != nil {
		return nil, fmt.Errorf("failed to create post")
	}

	return &post, nil
}

func GetPosts() (*[]models.PostResponse, error) {
	var posts []models.Post

	// Fetch posts and preload user details
	if err := db.DB.Preload("User").Order("created_at DESC").Find(&posts).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch posts")
	}

	// Format response to include votes and comments for each post
	var formattedPosts []models.PostResponse
	for _, post := range posts {
		var upvotes, downvotes int64

		// Count upvotes (vote = 1)
		db.DB.Model(&models.Vote{}).Where("post_id = ? AND vote = 1", post.ID).Count(&upvotes)

		// Count downvotes (vote = -1)
		db.DB.Model(&models.Vote{}).Where("post_id = ? AND vote = -1", post.ID).Count(&downvotes)

		// Fetch comments for the post
		var comments []models.Comment
		db.DB.Preload("User").Where("post_id = ?", post.ID).Order("created_at ASC").Find(&comments)

		// Convert comments into formatted response
		var formattedComments []models.CommentResponse
		for _, comment := range comments {
			formattedComments = append(formattedComments, models.CommentResponse{
				ID:        comment.ID,
				Content:   comment.Content,
				ImageURL:  comment.ImageURL,
				Username:  comment.User.Username,
				CreatedAt: comment.CreatedAt.Format("2006-01-02 15:04:05"),
			})
		}

		formattedPosts = append(formattedPosts, models.PostResponse{
			ID:        post.ID,
			Title:     post.Title,
			Content:   post.Content,
			ImageURL:  post.ImageURL,
			Username:  post.User.Username,
			Upvotes:   int(upvotes),
			Downvotes: int(downvotes),
			Comments:  formattedComments,
			CreatedAt: post.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &formattedPosts, nil
}
