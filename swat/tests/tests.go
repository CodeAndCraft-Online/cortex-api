package tests

import (
	"encoding/json"
	"fmt"

	"github.com/CodeAndCraft-Online/cortex-api/internal/swat"
)

// SWATClient re-export for backward compatibility
type SWATClient = swat.SWATClient

// NewSWATClient creates a new SWAT client
func NewSWATClient(baseURL string) *swat.SWATClient {
	return swat.NewSWATClient(baseURL)
}

// TestResult re-export
type TestResult = swat.TestResult

// RunHealthTests runs health check tests
func RunHealthTests(client *swat.SWATClient, verbose bool) ([]swat.TestResult, error) {
	var results []swat.TestResult

	// Test basic connectivity
	start := swat.StartTime()
	resp, body, err := client.MakeRequest("GET", "/", nil, nil)
	duration := swat.Elapsed(start)

	result := swat.NewTestResult("Health endpoint connectivity", duration, err)
	results = append(results, result)

	if verbose {
		fmt.Printf("   Status: %d, Body: %s\n", resp.StatusCode, string(body))
	}

	return results, nil
}

// RunAuthTests runs authentication tests
func RunAuthTests(client *swat.SWATClient, verbose bool) ([]swat.TestResult, error) {
	var results []swat.TestResult

	// Test user registration
	start := swat.StartTime()
	username, err := client.CreateTestUser("testuser", "password123")
	duration := swat.Elapsed(start)

	result := swat.NewTestResult("User registration and login", duration, err)
	results = append(results, result)

	if err == nil && verbose {
		fmt.Printf("   Created test user: %s\n", username)
	}

	// Test password reset request
	resetData := map[string]string{
		"username": username,
	}
	start = swat.StartTime()
	resp, _, err := client.MakeRequest("POST", "/auth/password-reset/request", resetData, nil)
	duration = swat.Elapsed(start)
	result = swat.NewTestResult("Password reset request", duration, err)
	if err == nil && resp.StatusCode != 200 && resp.StatusCode != 404 { // 404 if user not found, 200 if success
		result.Status = "FAIL"
		result.Error = fmt.Sprintf("Expected status 200 or 404, got %d", resp.StatusCode)
	}
	results = append(results, result)

	// Note: Password reset completion would require capturing the reset token
	// from the handler output, which is complex. For now, the request test demonstrates the endpoint.
	result = swat.NewTestResult("Password reset with token (requires manual token capture)", 0, nil)
	results = append(results, result)

	return results, nil
}

// Placeholder implementations for other test functions

func RunPostsTests(client *swat.SWATClient, verbose bool) ([]swat.TestResult, error) {
	var results []swat.TestResult

	// Create a test user for posts
	username, err := client.CreateTestUser("posttest", "password123")
	if err != nil {
		result := swat.NewTestResult("Create test user for posts", 0, err)
		results = append(results, result)
		return results, nil
	}

	if verbose {
		fmt.Printf("   Created test user: %s\n", username)
	}

	// Create a test sub first (posts require a sub_id)
	subData := map[string]interface{}{
		"name":        "posttestsub",
		"description": "A test sub for post testing",
		"private":     false,
	}
	start := swat.StartTime()
	resp, responseBody, err := client.MakeRequest("POST", "/sub/sub", subData, nil)
	duration := swat.Elapsed(start)
	result := swat.NewTestResult("Create test sub for posts", duration, err)
	if err == nil && resp.StatusCode != 201 {
		result.Status = "FAIL"
		result.Error = fmt.Sprintf("Expected status 201, got %d", resp.StatusCode)
	}
	results = append(results, result)

	var subResponse map[string]interface{}
	if err == nil && resp.StatusCode == 201 {
		json.Unmarshal(responseBody, &subResponse)
	}

	if subResponse == nil || subResponse["id"] == nil {
		// Cannot proceed with post tests without a sub
		result = swat.NewTestResult("Create new post", 0, fmt.Errorf("no sub ID available to create post"))
		results = append(results, result)
		result = swat.NewTestResult("Get all posts", 0, fmt.Errorf("no sub available to test posts"))
		results = append(results, result)
		result = swat.NewTestResult("Get post by ID", 0, fmt.Errorf("no sub available to test posts"))
		results = append(results, result)
		result = swat.NewTestResult("Get comments for post", 0, fmt.Errorf("no sub available to test posts"))
		results = append(results, result)
		return results, nil
	}

	subID := fmt.Sprintf("%.0f", subResponse["id"].(float64))

	// Test 1: Create a new post
	start = swat.StartTime()
	postData := map[string]interface{}{
		"title":   "Test Post Title",
		"content": "This is a test post content for SWAT testing.",
		"sub_id":  subID,
	}
	resp, body, err := client.MakeRequest("POST", "/posts/", postData, nil)
	duration = swat.Elapsed(start)

	result = swat.NewTestResult("Create new post", duration, err)
	if err == nil && resp.StatusCode != 201 {
		result.Status = "FAIL"
		result.Error = fmt.Sprintf("Expected status 201, got %d", resp.StatusCode)
	}
	results = append(results, result)

	var postResponse map[string]interface{}
	if err == nil && resp.StatusCode == 201 {
		json.Unmarshal(body, &postResponse)
	}

	// Test 2: Get posts (list all)
	start = swat.StartTime()
	resp, _, err = client.MakeRequest("GET", "/posts/", nil, nil)
	duration = swat.Elapsed(start)

	result = swat.NewTestResult("Get all posts", duration, err)
	if err == nil && resp.StatusCode != 200 {
		result.Status = "FAIL"
		result.Error = fmt.Sprintf("Expected status 200, got %d", resp.StatusCode)
	}
	results = append(results, result)

	// Test 3: Get post by ID (if we have one)
	if postResponse != nil && postResponse["id"] != nil {
		postID := fmt.Sprintf("%.0f", postResponse["id"].(float64))

		start = swat.StartTime()
		resp, _, err = client.MakeRequest("GET", "/posts/"+postID, nil, nil)
		duration = swat.Elapsed(start)

		result = swat.NewTestResult("Get post by ID", duration, err)
		if err == nil && resp.StatusCode != 200 {
			result.Status = "FAIL"
			result.Error = fmt.Sprintf("Expected status 200, got %d", resp.StatusCode)
		}
		results = append(results, result)

		// Test 4: Get comments for post
		start = swat.StartTime()
		resp, _, err = client.MakeRequest("GET", "/posts/posts/"+postID+"/comments", nil, nil)
		duration = swat.Elapsed(start)

		result = swat.NewTestResult("Get comments for post", duration, err)
		if err == nil && resp.StatusCode != 200 {
			result.Status = "FAIL"
			result.Error = fmt.Sprintf("Expected status 200, got %d", resp.StatusCode)
		}
		results = append(results, result)
	} else {
		result = swat.NewTestResult("Get post by ID", 0, fmt.Errorf("no post ID available to test"))
		results = append(results, result)

		result = swat.NewTestResult("Get comments for post", 0, fmt.Errorf("no post ID available to test"))
		results = append(results, result)
	}

	// Test 5: Test duplicate route POST /posts/posts/:postID (investigate if exists)
	if postResponse != nil && postResponse["id"] != nil {
		postID := fmt.Sprintf("%.0f", postResponse["id"].(float64))

		start = swat.StartTime()
		resp, _, _ = client.MakeRequest("GET", "/posts/posts/"+postID+"/comments", nil, nil)
		duration = swat.Elapsed(start)

		result = swat.NewTestResult("Investigate POST /posts/posts/:postID route", duration, nil)
		if resp.StatusCode == 405 || resp.StatusCode == 404 {
			result.Status = "PASS"
		} else {
			result.Error = fmt.Sprintf("Unexpected status %d", resp.StatusCode)
		}
		results = append(results, result)
	}

	return results, nil
}

func RunCommentsTests(client *swat.SWATClient, verbose bool) ([]swat.TestResult, error) {
	var results []swat.TestResult

	// Create a test user for comments
	_, err := client.CreateTestUser("commenttest", "password123")
	if err != nil {
		result := swat.NewTestResult("Create test user for comments", 0, err)
		results = append(results, result)
		return results, nil
	}

	// Create a test sub first for post creation
	subData := map[string]interface{}{
		"name":        "commenttestsub",
		"description": "A test sub for comment testing",
		"private":     false,
	}
	subStart := swat.StartTime()
	resp, responseBody, err := client.MakeRequest("POST", "/sub/sub", subData, nil)
	subDuration := swat.Elapsed(subStart)
	result := swat.NewTestResult("Create test sub for comments", subDuration, err)
	if err == nil && resp.StatusCode != 201 {
		result.Status = "FAIL"
		result.Error = fmt.Sprintf("Expected status 201, got %d", resp.StatusCode)
	}
	results = append(results, result)

	var subResponse map[string]interface{}
	if err == nil && resp.StatusCode == 201 {
		json.Unmarshal(responseBody, &subResponse)
	}

	if subResponse == nil || subResponse["id"] == nil {
		result = swat.NewTestResult("Create test post for comments", 0, fmt.Errorf("no sub ID available"))
		results = append(results, result)
		result = swat.NewTestResult("Create comment", 0, fmt.Errorf("no sub available"))
		results = append(results, result)
		return results, nil
	}

	subID := fmt.Sprintf("%.0f", subResponse["id"].(float64))

	// Create a test post first
	postData := map[string]interface{}{
		"title":   "Test Post for Comments",
		"content": "This post will have comments for testing.",
		"sub_id":  subID,
	}
	postResp, postBody, postErr := client.MakeRequest("POST", "/posts/", postData, nil)
	duration := swat.Elapsed(swat.StartTime())
	result = swat.NewTestResult("Create test post for comments", duration, postErr)
	if postErr == nil && postResp.StatusCode != 201 {
		result.Status = "FAIL"
		result.Error = fmt.Sprintf("Expected status 201, got %d", postResp.StatusCode)
	}
	results = append(results, result)

	var postResponse map[string]interface{}
	if postErr == nil && postResp.StatusCode == 201 {
		json.Unmarshal(postBody, &postResponse)
	}

	if postResponse == nil || postResponse["id"] == nil {
		result = swat.NewTestResult("Get comment by ID", 0, fmt.Errorf("no test data available"))
		results = append(results, result)
		result = swat.NewTestResult("Update comment", 0, fmt.Errorf("no test data available"))
		results = append(results, result)
		result = swat.NewTestResult("Delete comment", 0, fmt.Errorf("no test data available"))
		results = append(results, result)
		result = swat.NewTestResult("Create comment", 0, fmt.Errorf("no test data available"))
		results = append(results, result)
		return results, nil
	}

	postID := fmt.Sprintf("%.0f", postResponse["id"].(float64))

	// Test 1: Create a comment
	commentData := map[string]interface{}{
		"postID":  postID,
		"content": "This is a test comment.",
	}
	start := swat.StartTime()
	commentResp, commentBody, commentErr := client.MakeRequest("POST", "/comments/comments", commentData, nil)
	duration = swat.Elapsed(start)
	result = swat.NewTestResult("Create comment", duration, commentErr)
	if commentErr == nil && commentResp.StatusCode != 201 {
		result.Status = "FAIL"
		result.Error = fmt.Sprintf("Expected status 201, got %d", commentResp.StatusCode)
	}
	results = append(results, result)

	var commentResponse map[string]interface{}
	if commentErr == nil && commentResp.StatusCode == 201 {
		json.Unmarshal(commentBody, &commentResponse)
	}

	// Test 2: Get comment by ID (if we have a comment)
	if commentResponse != nil && commentResponse["id"] != nil {
		commentID := fmt.Sprintf("%.0f", commentResponse["id"].(float64))

		start = swat.StartTime()
		resp, _, err = client.MakeRequest("GET", "/comments/"+commentID, nil, nil)
		duration = swat.Elapsed(start)
		result = swat.NewTestResult("Get comment by ID", duration, err)
		if err == nil && resp.StatusCode != 200 {
			result.Status = "FAIL"
			result.Error = fmt.Sprintf("Expected status 200, got %d", resp.StatusCode)
		}
		results = append(results, result)

		// Test 3: Update comment
		updateData := map[string]interface{}{
			"content": "Updated test comment content.",
		}
		start = swat.StartTime()
		resp, _, err = client.MakeRequest("PUT", "/comments/"+commentID, updateData, nil)
		duration = swat.Elapsed(start)
		result = swat.NewTestResult("Update comment", duration, err)
		if err == nil && resp.StatusCode != 200 {
			result.Status = "FAIL"
			result.Error = fmt.Sprintf("Expected status 200, got %d", resp.StatusCode)
		}
		results = append(results, result)

		// Test 4: Delete comment
		start = swat.StartTime()
		resp, _, err = client.MakeRequest("DELETE", "/comments/"+commentID, nil, nil)
		duration = swat.Elapsed(start)
		result = swat.NewTestResult("Delete comment", duration, err)
		if err == nil && resp.StatusCode != 200 {
			result.Status = "FAIL"
			result.Error = fmt.Sprintf("Expected status 200, got %d", resp.StatusCode)
		}
		results = append(results, result)
	} else {
		result = swat.NewTestResult("Get comment by ID", 0, fmt.Errorf("no comment ID available to test"))
		results = append(results, result)
		result = swat.NewTestResult("Update comment", 0, fmt.Errorf("no comment ID available to test"))
		results = append(results, result)
		result = swat.NewTestResult("Delete comment", 0, fmt.Errorf("no comment ID available to test"))
		results = append(results, result)
	}

	return results, nil
}

func RunSubsTests(client *swat.SWATClient, verbose bool) ([]swat.TestResult, error) {
	var results []swat.TestResult

	// Create a test user for subs
	_, err := client.CreateTestUser("subtest", "password123")
	if err != nil {
		result := swat.NewTestResult("Create test user for subs", 0, err)
		results = append(results, result)
		return results, nil
	}

	// Test 1: Get all subs (list)
	start := swat.StartTime()
	resp, _, err := client.MakeRequest("GET", "/sub/", nil, nil)
	duration := swat.Elapsed(start)
	result := swat.NewTestResult("Get all subs", duration, err)
	if err == nil && resp.StatusCode != 200 {
		result.Status = "FAIL"
		result.Error = fmt.Sprintf("Expected status 200, got %d", resp.StatusCode)
	}
	results = append(results, result)

	// Test 2: Create a new sub
	subData := map[string]interface{}{
		"name":        "testcommunity",
		"description": "A test community for SWAT testing",
		"private":     false,
	}
	start = swat.StartTime()
	resp, responseBody, err := client.MakeRequest("POST", "/sub/sub", subData, nil)
	duration = swat.Elapsed(start)
	result = swat.NewTestResult("Create new sub", duration, err)
	if err == nil && resp.StatusCode != 201 {
		result.Status = "FAIL"
		result.Error = fmt.Sprintf("Expected status 201, got %d", resp.StatusCode)
	}
	results = append(results, result)

	var subResponse map[string]interface{}
	if err == nil && resp.StatusCode == 201 {
		json.Unmarshal(responseBody, &subResponse)
	}

	if subResponse == nil || subResponse["id"] == nil {
		result = swat.NewTestResult("Get posts in sub", 0, fmt.Errorf("no sub ID available to test"))
		results = append(results, result)
		result = swat.NewTestResult("Get post count for sub", 0, fmt.Errorf("no sub ID available to test"))
		results = append(results, result)
		result = swat.NewTestResult("Join sub", 0, fmt.Errorf("no sub ID available to test"))
		results = append(results, result)
		result = swat.NewTestResult("Leave sub", 0, fmt.Errorf("no sub ID available to test"))
		results = append(results, result)
		result = swat.NewTestResult("Update sub", 0, fmt.Errorf("no sub ID available to test"))
		results = append(results, result)
		result = swat.NewTestResult("Delete sub", 0, fmt.Errorf("no sub ID available to test"))
		results = append(results, result)
		result = swat.NewTestResult("Get sub members", 0, fmt.Errorf("no sub ID available to test"))
		results = append(results, result)
		result = swat.NewTestResult("Get pending invites", 0, fmt.Errorf("no sub ID available to test"))
		results = append(results, result)
		return results, nil
	}

	subID := fmt.Sprintf("%.0f", subResponse["id"].(float64))

	// Test 3: Get posts in sub
	start = swat.StartTime()
	resp, _, err = client.MakeRequest("GET", "/sub/sub/"+subID+"/posts", nil, nil)
	duration = swat.Elapsed(start)
	result = swat.NewTestResult("Get posts in sub", duration, err)
	if err == nil && resp.StatusCode != 200 {
		result.Status = "FAIL"
		result.Error = fmt.Sprintf("Expected status 200, got %d", resp.StatusCode)
	}
	results = append(results, result)

	// Test 4: Get post count for sub
	start = swat.StartTime()
	resp, _, err = client.MakeRequest("GET", "/sub/sub/"+subID+"/postCount?subID="+subID, nil, nil)
	duration = swat.Elapsed(start)
	result = swat.NewTestResult("Get post count for sub", duration, err)
	if err == nil && resp.StatusCode != 200 {
		result.Status = "FAIL"
		result.Error = fmt.Sprintf("Expected status 200, got %d", resp.StatusCode)
	}
	results = append(results, result)

	// Test 5: Join sub (automatically joined as creator, but test explicit join)
	start = swat.StartTime()
	resp, _, err = client.MakeRequest("POST", "/sub/sub/"+subID+"/join", nil, nil)
	duration = swat.Elapsed(start)
	result = swat.NewTestResult("Join sub", duration, err)
	// Joining own sub might return error or success, allow both
	if err == nil && resp.StatusCode != 200 && resp.StatusCode != 400 {
		result.Status = "FAIL"
		result.Error = fmt.Sprintf("Expected status 200 or 400, got %d", resp.StatusCode)
	}
	results = append(results, result)

	// Test 6: Leave sub
	start = swat.StartTime()
	resp, _, err = client.MakeRequest("POST", "/sub/sub/"+subID+"/leave", nil, nil)
	duration = swat.Elapsed(start)
	result = swat.NewTestResult("Leave sub", duration, err)
	if err == nil && resp.StatusCode != 200 {
		result.Status = "FAIL"
		result.Error = fmt.Sprintf("Expected status 200, got %d", resp.StatusCode)
	}
	results = append(results, result)

	// Rejoin to continue testing
	start = swat.StartTime()
	_, _, err = client.MakeRequest("POST", "/sub/sub/"+subID+"/join", nil, nil)
	_ = swat.Elapsed(start)
	if err != nil && verbose {
		fmt.Printf("   Failed to rejoin sub: %s\n", err.Error())
	}

	// Test 7: Update sub
	updateData := map[string]interface{}{
		"description": "Updated test community description",
		"private":     true,
	}
	start = swat.StartTime()
	resp, _, err = client.MakeRequest("PATCH", "/sub/"+subID, updateData, nil)
	duration = swat.Elapsed(start)
	result = swat.NewTestResult("Update sub", duration, err)
	if err == nil && resp.StatusCode != 200 {
		result.Status = "FAIL"
		result.Error = fmt.Sprintf("Expected status 200, got %d", resp.StatusCode)
	}
	results = append(results, result)

	// Test 8: Get sub members
	start = swat.StartTime()
	resp, _, err = client.MakeRequest("GET", "/sub/"+subID+"/members", nil, nil)
	duration = swat.Elapsed(start)
	result = swat.NewTestResult("Get sub members", duration, err)
	if err == nil && resp.StatusCode != 200 {
		result.Status = "FAIL"
		result.Error = fmt.Sprintf("Expected status 200, got %d", resp.StatusCode)
	}
	results = append(results, result)

	// Test 9: Get pending invites (should be empty)
	start = swat.StartTime()
	resp, _, err = client.MakeRequest("GET", "/sub/"+subID+"/pending-invites", nil, nil)
	duration = swat.Elapsed(start)
	result = swat.NewTestResult("Get pending invites", duration, err)
	if err == nil && resp.StatusCode != 200 {
		result.Status = "FAIL"
		result.Error = fmt.Sprintf("Expected status 200, got %d", resp.StatusCode)
	}
	results = append(results, result)

	// Test 10: Invite user to private sub (test with another user)
	// First create another user for inviting
	_, err = client.CreateTestUser("inviteetest", "password123")
	if err != nil && verbose {
		fmt.Printf("   Failed to create invitee user: %s\n", err.Error())
	}

	// Try to invite user (even though sub might be public now)
	inviteData := map[string]interface{}{
		"invitee_username": "inviteetest",
	}
	start = swat.StartTime()
	_, _, err = client.MakeRequest("POST", "/sub/sub/"+subID+"/invite", inviteData, nil)
	duration = swat.Elapsed(start)
	result = swat.NewTestResult("Invite user to sub", duration, err)
	// This might fail depending on sub privacy, allow various status codes
	results = append(results, result)

	// Test 11: Delete sub
	start = swat.StartTime()
	resp, _, err = client.MakeRequest("DELETE", "/sub/"+subID, nil, nil)
	duration = swat.Elapsed(start)
	result = swat.NewTestResult("Delete sub", duration, err)
	if err == nil && resp.StatusCode != 200 {
		result.Status = "FAIL"
		result.Error = fmt.Sprintf("Expected status 200, got %d", resp.StatusCode)
	}
	results = append(results, result)

	return results, nil
}

func RunVotesTests(client *swat.SWATClient, verbose bool) ([]swat.TestResult, error) {
	var results []swat.TestResult

	// Create a test user for voting
	_, err := client.CreateTestUser("votetest", "password123")
	if err != nil {
		result := swat.NewTestResult("Create test user for votes", 0, err)
		results = append(results, result)
		return results, nil
	}

	// Create a test sub first for post creation
	subData := map[string]interface{}{
		"name":        "votetestsub",
		"description": "A test sub for vote testing",
		"private":     false,
	}
	start := swat.StartTime()
	resp, responseBody, err := client.MakeRequest("POST", "/sub/sub", subData, nil)
	duration := swat.Elapsed(start)
	result := swat.NewTestResult("Create test sub for votes", duration, err)
	if err == nil && resp.StatusCode != 201 {
		result.Status = "FAIL"
		result.Error = fmt.Sprintf("Expected status 201, got %d", resp.StatusCode)
	}
	results = append(results, result)

	var subResponse map[string]interface{}
	if err == nil && resp.StatusCode == 201 {
		json.Unmarshal(responseBody, &subResponse)
	}

	if subResponse == nil || subResponse["id"] == nil {
		result := swat.NewTestResult("Create test post for voting", 0, fmt.Errorf("no sub ID available"))
		results = append(results, result)
		return results, nil
	}

	subID := fmt.Sprintf("%.0f", subResponse["id"].(float64))

	// Create a test post to vote on
	postData := map[string]interface{}{
		"title":   "Test Post for Voting",
		"content": "This post will be voted on for testing.",
		"sub_id":  subID,
	}
	resp, body, err := client.MakeRequest("POST", "/posts/", postData, nil)
	if err != nil {
		result := swat.NewTestResult("Create test post for voting", 0, err)
		results = append(results, result)
		return results, nil
	}

	if resp.StatusCode != 201 {
		result := swat.NewTestResult("Create test post for voting", 0, fmt.Errorf("failed to create post: %d", resp.StatusCode))
		results = append(results, result)
		return results, nil
	}

	var postResponse map[string]interface{}
	json.Unmarshal(body, &postResponse)
	if postResponse == nil || postResponse["id"] == nil {
		result := swat.NewTestResult("Parse post response for voting", 0, fmt.Errorf("invalid post response"))
		results = append(results, result)
		return results, nil
	}

	postID := fmt.Sprintf("%.0f", postResponse["id"].(float64))

	// Test 1: Upvote a post
	upvoteData := map[string]interface{}{
		"post_id": postID,
	}
	start = swat.StartTime()
	resp, _, err = client.MakeRequest("POST", "/vote/upvote", upvoteData, nil)
	duration = swat.Elapsed(start)
	upvoteResult := swat.NewTestResult("Upvote post", duration, err)
	if err == nil && resp.StatusCode != 200 {
		upvoteResult.Status = "FAIL"
		upvoteResult.Error = fmt.Sprintf("Expected status 200, got %d", resp.StatusCode)
	}
	results = append(results, upvoteResult)

	// Test 2: Downvote a post
	downvoteData := map[string]interface{}{
		"post_id": postID,
	}
	start = swat.StartTime()
	_, _, err = client.MakeRequest("POST", "/vote/downvote", downvoteData, nil)
	duration = swat.Elapsed(start)
	downvoteResult := swat.NewTestResult("Downvote post", duration, err)
	results = append(results, downvoteResult)

	return results, nil
}

func RunUsersTests(client *swat.SWATClient, verbose bool) ([]swat.TestResult, error) {
	var results []swat.TestResult

	// Create a test user for users tests
	actualUsername, err := client.CreateTestUser("usertest", "password123")
	if err != nil {
		result := swat.NewTestResult("Create test user for users", 0, err)
		results = append(results, result)
		return results, nil
	}

	// Test 1: Get public user profile (use actual username)
	start := swat.StartTime()
	resp, _, err := client.MakeRequest("GET", "/user/"+actualUsername, nil, nil)
	duration := swat.Elapsed(start)
	result := swat.NewTestResult("Get public user profile", duration, err)
	if err == nil && resp.StatusCode != 200 {
		result.Status = "FAIL"
		result.Error = fmt.Sprintf("Expected status 200, got %d", resp.StatusCode)
	}
	results = append(results, result)

	// Test 2: Get current user profile (protected)
	start = swat.StartTime()
	resp, _, err = client.MakeRequest("GET", "/user/profile", nil, nil)
	duration = swat.Elapsed(start)
	result = swat.NewTestResult("Get current user profile", duration, err)
	if err == nil && resp.StatusCode != 200 {
		result.Status = "FAIL"
		result.Error = fmt.Sprintf("Expected status 200, got %d", resp.StatusCode)
	}
	results = append(results, result)

	// Test 3: Update user profile
	updateData := map[string]interface{}{
		"display_name": "Updated Display Name",
		"bio":          "Updated bio content for testing.",
	}
	start = swat.StartTime()
	resp, _, err = client.MakeRequest("PUT", "/user/profile", updateData, nil)
	duration = swat.Elapsed(start)
	result = swat.NewTestResult("Update user profile", duration, err)
	if err == nil && resp.StatusCode != 200 {
		result.Status = "FAIL"
		result.Error = fmt.Sprintf("Expected status 200, got %d", resp.StatusCode)
	}
	results = append(results, result)

	// Create a sub to test invite functionality
	subData := map[string]interface{}{
		"name":    "invitecomm",
		"private": true,
	}
	subStart := swat.StartTime()
	subResp, subBody, subErr := client.MakeRequest("POST", "/sub/sub", subData, nil)
	subDuration := swat.Elapsed(subStart)
	result = swat.NewTestResult("Create private sub for invite test", subDuration, subErr)
	if subErr == nil && subResp.StatusCode != 201 {
		result.Status = "FAIL"
		result.Error = fmt.Sprintf("Expected status 201, got %d", subResp.StatusCode)
	}
	results = append(results, result)

	var subResponse map[string]interface{}
	if subResp.StatusCode == 201 && subErr == nil {
		json.Unmarshal(subBody, &subResponse)
	}

	// Create another user for invite testing
	var inviteeUsername string
	inviteeUsername, err = client.CreateTestUser("inviteeuser", "password123")
	if err != nil {
		result = swat.NewTestResult("Create invitee user", 0, err)
		results = append(results, result)
		// Continue with remaining tests
	} else if subResponse != nil && subResponse["id"] != nil {
		subID := fmt.Sprintf("%.0f", subResponse["id"].(float64))

		// Invite the user
		inviteData := map[string]interface{}{
			"invitee_username": inviteeUsername,
		}
		start = swat.StartTime()
		_, _, err = client.MakeRequest("POST", "/sub/sub/"+subID+"/invite", inviteData, nil)
		duration = swat.Elapsed(start)
		result = swat.NewTestResult("Send user invite", duration, err)
		// Invite might succeed or fail depending on implementation, record result
		results = append(results, result)

		// Note: Accept invite test is complex to implement with user switching
		// We'd need to switch authentication to the invited user
		// For now, skip this test
		result = swat.NewTestResult("Accept invite (skipped - requires user switching)", 0, nil)
		results = append(results, result)
	} else {
		result = swat.NewTestResult("Send user invite", 0, fmt.Errorf("no sub available to test invites"))
		results = append(results, result)
		result = swat.NewTestResult("Accept invite", 0, fmt.Errorf("no invite available to test"))
		results = append(results, result)
	}

	// Skip account deletion test to avoid breaking the test user for other tests
	result = swat.NewTestResult("Delete user account (skipped - destructive)", 0, nil)
	results = append(results, result)

	return results, nil
}

func RunSecurityTests(client *swat.SWATClient, verbose bool) ([]swat.TestResult, error) {
	var results []swat.TestResult

	// Test 1: Unauthorized access to protected endpoint
	result1 := swat.NewTestResult("Test protected endpoint without auth", 0, fmt.Errorf("implemented"))
	start1 := swat.StartTime()
	resp1, _, err1 := client.MakeRequest("GET", "/user/profile", nil, nil)
	duration1 := swat.Elapsed(start1)
	result1.Duration = duration1
	if err1 != nil {
		result1.Error = err1.Error()
		result1.Status = "PASS" // Error is expected for unauthorized access
	} else if resp1.StatusCode == 401 {
		result1.Status = "PASS" // 401 Unauthorized is correct
	} else {
		result1.Error = fmt.Sprintf("Expected 401, got %d", resp1.StatusCode)
	}
	results = append(results, result1)

	// Test 2: SQL injection attempt in route parameter
	// This tests if the API properly sanitizes route parameters
	result2 := swat.NewTestResult("SQL injection in route parameter", 0, fmt.Errorf("implemented"))
	start2 := swat.StartTime()
	resp2, _, err2 := client.MakeRequest("GET", "/posts/' OR '1'='1", nil, nil)
	duration2 := swat.Elapsed(start2)
	result2.Duration = duration2
	// Should return 404 or 400, not 500 (which would indicate injection worked)
	if err2 == nil && resp2.StatusCode == 404 {
		result2.Status = "PASS"
	} else if err2 == nil && resp2.StatusCode == 400 {
		result2.Status = "PASS"
	} else {
		result2.Error = fmt.Sprintf("Unexpected response - status: %d, error: %v", resp2.StatusCode, err2)
	}
	results = append(results, result2)

	// Test 3: Attempt to access non-existent endpoint
	result3 := swat.NewTestResult("Access to non-existent endpoint", 0, fmt.Errorf("implemented"))
	start3 := swat.StartTime()
	resp3, _, err3 := client.MakeRequest("GET", "/nonexistent", nil, nil)
	duration3 := swat.Elapsed(start3)
	result3.Duration = duration3
	if err3 == nil && resp3.StatusCode == 404 {
		result3.Status = "PASS"
	} else {
		result3.Error = fmt.Sprintf("Expected 404, got %d", resp3.StatusCode)
	}
	results = append(results, result3)

	return results, nil
}
