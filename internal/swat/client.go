package swat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type SWATClient struct {
	BaseURL    string
	httpClient *http.Client
	AuthToken  string
	TestUsers  []string // usernames of users created for testing
}

type AuthResponse struct {
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
}

type GenericResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// NewSWATClient creates a new SWAT client for testing the API
func NewSWATClient(baseURL string) *SWATClient {
	return &SWATClient{
		BaseURL:    baseURL,
		httpClient: &http.Client{Timeout: 30 * time.Second},
		TestUsers:  []string{},
	}
}

// MakeRequest performs an HTTP request and returns the response
func (c *SWATClient) MakeRequest(method, endpoint string, body interface{}, headers map[string]string) (*http.Response, []byte, error) {
	url := c.BaseURL + endpoint

	var jsonBody []byte
	var err error
	if body != nil {
		jsonBody, err = json.Marshal(body)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to marshal request body: %v", err)
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Default headers
	req.Header.Set("Content-Type", "application/json")

	// Add custom headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Add auth token if available
	if c.AuthToken != "" && headers["Authorization"] == "" {
		req.Header.Set("Authorization", "Bearer "+c.AuthToken)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to make request: %v", err)
	}

	responseBody, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	return resp, responseBody, err
}

// RegisterTestUser registers a test user and stores the username for cleanup
func (c *SWATClient) RegisterTestUser(username, password string) error {
	userData := map[string]string{
		"username": username,
		"password": password,
	}

	resp, body, err := c.MakeRequest("POST", "/auth/register", userData, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		var errorResp GenericResponse
		json.Unmarshal(body, &errorResp)
		return fmt.Errorf("registration failed: %s", errorResp.Error)
	}

	c.TestUsers = append(c.TestUsers, username)
	return nil
}

// LoginUser logs in a user and sets the auth token
func (c *SWATClient) LoginUser(username, password string) error {
	authData := map[string]string{
		"username": username,
		"password": password,
	}

	resp, body, err := c.MakeRequest("POST", "/auth/login", authData, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		var errorResp GenericResponse
		json.Unmarshal(body, &errorResp)
		return fmt.Errorf("login failed: %s", errorResp.Error)
	}

	var authResp AuthResponse
	if err := json.Unmarshal(body, &authResp); err != nil {
		return fmt.Errorf("failed to parse login response: %v", err)
	}

	c.AuthToken = authResp.Token
	return nil
}

// CreateTestUser creates a unique test user, registers them, and logs them in
func (c *SWATClient) CreateTestUser(baseUsername, password string) (string, error) {
	timestamp := time.Now().Unix()
	username := fmt.Sprintf("%s_%d", baseUsername, timestamp)

	err := c.RegisterTestUser(username, password)
	if err != nil {
		return "", err
	}

	err = c.LoginUser(username, password)
	if err != nil {
		return "", err
	}

	return username, nil
}

// ClearToken clears the current authentication token
func (c *SWATClient) ClearToken() {
	c.AuthToken = ""
}

// GetTestEndpoint generates a unique endpoint for testing (to avoid conflicts)
func (c *SWATClient) GetTestEndpoint(baseEndpoint string) string {
	if !strings.Contains(baseEndpoint, "?") {
		return baseEndpoint + "?_test=1"
	}
	return baseEndpoint + "&_test=1"
}

// Cleanup removes test data created during testing
func (c *SWATClient) Cleanup() error {
	// Note: Actual cleanup would require API endpoints to delete test users
	// For now, we just log the users that were created
	fmt.Printf("Test cleanup: Created %d test users: %v\n", len(c.TestUsers), c.TestUsers)
	return nil
}

// TestResult represents the result of a test
type TestResult struct {
	Name     string        `json:"name"`
	Status   string        `json:"status"`
	Duration time.Duration `json:"duration"`
	Error    string        `json:"error,omitempty"`
}

// Helper function to create test results
func NewTestResult(name string, duration time.Duration, err error) TestResult {
	result := TestResult{
		Name:     name,
		Duration: duration,
	}

	if err != nil {
		result.Status = "FAIL"
		result.Error = err.Error()
	} else {
		result.Status = "PASS"
	}

	return result
}

// StartTime returns the current time for measuring test duration
func StartTime() time.Time {
	return time.Now()
}

// Elapsed returns the duration since the start time
func Elapsed(start time.Time) time.Duration {
	return time.Since(start)
}
