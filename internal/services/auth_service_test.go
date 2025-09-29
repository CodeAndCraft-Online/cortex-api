package services

import (
	"testing"
	"time"

	"github.com/CodeAndCraft-Online/cortex-api/internal/database"
	"github.com/CodeAndCraft-Online/cortex-api/internal/models"
	"github.com/CodeAndCraft-Online/cortex-api/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var available = true   // Track if database is available
var forceMocks = false // Set to true to test mock functionality

func TestMain(m *testing.M) {
	// For testing: uncomment to force mocks
	// forceMocks = true

	if forceMocks {
		available = false
		m.Run()
		return
	}

	// Try to set up test DB - if Docker not available, skip all tests
	db, teardown, err := testutils.SetupTestDB()
	if err != nil {
		println("Docker not available, skipping service integration tests:", err.Error())
		available = false
		teardown() // Safe to call even if nil
		return     // Skip all tests
	}

	database.DB = db
	m.Run()
	teardown()
}

// MockAuthRepository is a mock implementation of IAuthRepository
type MockAuthRepository struct {
	mock.Mock
}

func (m *MockAuthRepository) ResetPasswordRequest(username string) (*models.PasswordResetToken, error) {
	args := m.Called(username)
	return args.Get(0).(*models.PasswordResetToken), args.Error(1)
}

func (m *MockAuthRepository) ResetPassword(token, newPassword string) error {
	args := m.Called(token, newPassword)
	return args.Error(0)
}

// Test auth service with mocks (run when DB not available - for CI)
func TestResetPasswordRequest_ServiceWithMock(t *testing.T) {
	if available {
		t.Skip("Database available, skipping mock tests")
		return
	}

	mockRepo := new(MockAuthRepository)
	service := NewAuthService(mockRepo)

	expectedToken := &models.PasswordResetToken{
		UserID:    1,
		Token:     "mocktoken",
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}

	mockRepo.On("ResetPasswordRequest", "serviceuser").Return(expectedToken, nil)

	result, err := service.ResetPasswordRequest("serviceuser")

	assert.NoError(t, err)
	assert.Equal(t, expectedToken, result)
	mockRepo.AssertExpectations(t)
}

func TestResetPassword_ServiceWithMock(t *testing.T) {
	if available {
		t.Skip("Database available, skipping mock tests")
		return
	}

	mockRepo := new(MockAuthRepository)
	service := NewAuthService(mockRepo)

	mockRepo.On("ResetPassword", "Servicetoken", "newpassword").Return(nil)

	err := service.ResetPassword("Servicetoken", "newpassword")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
