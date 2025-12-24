package services

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"

	"github.com/im-mk/user-service/src/models"
)

func TestAuthService_Login(t *testing.T) {
	mockRepo := new(MockUserRepository)
	jwtKey := []byte("my_secret_key")
	mockRefreshRepo := new(MockRefreshTokenRepository)
	authService := NewAuthService(mockRepo, mockRefreshRepo, jwtKey)

	t.Run("successful login", func(t *testing.T) {
		password := "password123"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		mockUser := &models.User{
			Username: "testuser",
			Password: string(hashedPassword),
		}

		mockRepo.On("GetUserByUsername", "testuser").Return(mockUser, nil)
		// save refresh token expected
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
			Username: "testuser",
			Password: string(hashedPassword),
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
}
