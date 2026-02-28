package services

import (
	"errors"
	"testing"

	"github.com/im-mk/user-service/src/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthService_Login(t *testing.T) {
	authConfig := models.AuthConfig{
		Issuer:                    "test-issuer",
		Audience:                  "test-audience",
		TokenExpirySeconds:        3600,
		RefreshTokenExpirySeconds: 7200,
		Kid:                       "test-kid",
	}
	mockRepo := new(MockUserRepository)
	mockRefreshRepo := new(MockRefreshTokenRepository)
	mockTokenProv := new(MockTokenProvider)
	authService := NewAuthService(mockRepo, mockRefreshRepo, mockTokenProv, authConfig)

	t.Run("successful login", func(t *testing.T) {
		password := "password123"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		mockUser := &models.User{
			Username:   "testuser",
			Password:   string(hashedPassword),
			IsActive:   true,
			IsVerified: true,
		}

		mockRepo.On("GetUserByUsername", "testuser").Return(mockUser, nil)
		mockTokenProv.On("GenerateAccessToken", "0", "testuser").Return("tokA", nil)
		mockTokenProv.On("GenerateRefreshToken").Return("tokR", nil)
		mockRefreshRepo.On("SaveRefreshToken", mock.Anything, mock.Anything, mock.Anything).Return(nil)

		access, refresh, err := authService.Login(models.LoginRequest{
			Username: "testuser",
			Password: password,
		})

		assert.NoError(t, err)
		assert.NotEmpty(t, access)
		assert.NotEmpty(t, refresh)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid credentials - wrong password", func(t *testing.T) {
		password := "password123"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		mockUser := &models.User{
			Username:   "testuser",
			Password:   string(hashedPassword),
			IsActive:   true,
			IsVerified: true,
		}

		mockRepo.On("GetUserByUsername", "testuser").Return(mockUser, nil)

		access, refresh, err := authService.Login(models.LoginRequest{
			Username: "testuser",
			Password: "wrongpassword",
		})

		assert.Error(t, err)
		assert.Equal(t, "", access)
		assert.Equal(t, "", refresh)
		assert.EqualError(t, err, "invalid credentials")
		mockRepo.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		mockRepo.On("GetUserByUsername", "unknownuser").Return(&models.User{}, errors.New("user not found"))
		// provider not invoked

		access, refresh, err := authService.Login(models.LoginRequest{
			Username: "unknownuser",
			Password: "password123",
		})

		assert.Error(t, err)
		assert.Equal(t, "", access)
		assert.Equal(t, "", refresh)
		assert.EqualError(t, err, "invalid credentials")
		mockRepo.AssertExpectations(t)
	})

	t.Run("inactive account", func(t *testing.T) {
		// Fresh mocks for this subtest
		inactiveMockRepo := new(MockUserRepository)
		inactiveMockRefreshRepo := new(MockRefreshTokenRepository)
		inactiveTokenProv := new(MockTokenProvider)
		inactiveAuthService := NewAuthService(inactiveMockRepo, inactiveMockRefreshRepo, inactiveTokenProv, authConfig)

		password := "password123"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		mockUser := &models.User{
			Username:   "testuser",
			Password:   string(hashedPassword),
			IsActive:   false,
			IsVerified: true,
		}

		inactiveMockRepo.On("GetUserByUsername", "testuser").Return(mockUser, nil)

		access, refresh, err := inactiveAuthService.Login(models.LoginRequest{
			Username: "testuser",
			Password: password,
		})

		assert.Error(t, err)
		assert.EqualError(t, err, "account inactive")
		assert.Equal(t, "", access)
		assert.Equal(t, "", refresh)
		inactiveMockRepo.AssertExpectations(t)
	})

	t.Run("unverified account", func(t *testing.T) {
		// Fresh mocks for this subtest
		unverifiedMockRepo := new(MockUserRepository)
		unverifiedMockRefreshRepo := new(MockRefreshTokenRepository)
		unverifiedTokenProv := new(MockTokenProvider)
		unverifiedAuthService := NewAuthService(unverifiedMockRepo, unverifiedMockRefreshRepo, unverifiedTokenProv, authConfig)

		password := "password123"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		mockUser := &models.User{
			Username:   "testuser",
			Password:   string(hashedPassword),
			IsActive:   true,
			IsVerified: false,
		}

		unverifiedMockRepo.On("GetUserByUsername", "testuser").Return(mockUser, nil)

		access, refresh, err := unverifiedAuthService.Login(models.LoginRequest{
			Username: "testuser",
			Password: password,
		})

		assert.Error(t, err)
		assert.EqualError(t, err, "account unverified")
		assert.Equal(t, "", access)
		assert.Equal(t, "", refresh)
		unverifiedMockRepo.AssertExpectations(t)
	})
}
