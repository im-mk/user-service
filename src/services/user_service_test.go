package services

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/im-mk/user-service/src/models"
)

func TestUserService_CreateUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := NewUserService(mockRepo)

	t.Run("successful user creation", func(t *testing.T) {
		req := models.CreateUserRequest{
			Username: "newuser",
			Email:    "newuser@example.com",
			Password: "password123",
		}

		mockRepo.On("UserExists", req.Username, req.Email).Return(false, nil).Once()
		mockRepo.On("CreateUser", mock.Anything).Return(nil).Once()

		err := userService.CreateUser(req)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("user already exists", func(t *testing.T) {
		req := models.CreateUserRequest{
			Username: "existinguser",
			Email:    "existinguser@example.com",
			Password: "password123",
		}

		mockRepo.On("UserExists", req.Username, req.Email).Return(true, nil).Once()

		err := userService.CreateUser(req)

		assert.Error(t, err)
		assert.EqualError(t, err, "username or email already exists")
		mockRepo.AssertExpectations(t)
	})

	t.Run("error checking user existence", func(t *testing.T) {
		req := models.CreateUserRequest{
			Username: "newuser",
			Email:    "newuser@example.com",
			Password: "password123",
		}

		mockRepo.On("UserExists", req.Username, req.Email).Return(false, errors.New("db error")).Once()

		err := userService.CreateUser(req)

		assert.Error(t, err)
		assert.EqualError(t, err, "failed to check for existing user")
		mockRepo.AssertExpectations(t)
	})
}
