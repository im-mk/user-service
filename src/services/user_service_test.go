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
			Username:   "newuser",
			Email:      "newuser@example.com",
			Password:   "password123",
			FirstName:  "New",
			MiddleName: "",
			LastName:   "User",
		}

		mockRepo.On("UserExists", req.Username, req.Email).Return(false, nil).Once()
		// validate that the passed user has the expected fields set
		mockRepo.On("CreateUser", mock.MatchedBy(func(u models.User) bool {
			return u.Username == req.Username &&
				u.Email == req.Email &&
				u.FirstName == req.FirstName &&
				u.MiddleName == req.MiddleName &&
				u.LastName == req.LastName &&
				u.IsActive == true &&
				u.IsVerified == false
		})).Return(nil).Once()

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

	// add a test for GetUser mapping

	t.Run("bootstrap user is verified", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		userService := NewUserService(mockRepo)

		req := models.CreateUserRequest{
			Username: "first",
			Email:    "first@example.com",
			Password: "pass",
		}

		// no existing users
		mockRepo.On("AnyUserExists").Return(false, nil).Once()
		// when CreateUser is called we expect IsVerified true
		mockRepo.On("CreateUser", mock.MatchedBy(func(u models.User) bool {
			return u.IsVerified == true && u.Username == req.Username
		})).Return(nil).Once()

		err := userService.Bootstrap(req)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("bootstrap fails when users exist", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		userService := NewUserService(mockRepo)

		mockRepo.On("AnyUserExists").Return(true, nil).Once()

		err := userService.Bootstrap(models.CreateUserRequest{})
		assert.Error(t, err)
	})

	t.Run("get user details mapping", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		userService := NewUserService(mockRepo)

		domainUser := &models.User{
			ID:         42,
			Username:   "bob",
			Email:      "bob@example.com",
			FirstName:  "Bob",
			MiddleName: "Q",
			LastName:   "Builder",
			IsActive:   true,
			IsVerified: true,
		}

		mockRepo.On("GetUserByID", 42).Return(domainUser, nil).Once()

		details, err := userService.GetUser(42)
		assert.NoError(t, err)
		assert.Equal(t, 42, details.ID)
		assert.Equal(t, "bob", details.Username)
		assert.Equal(t, "bob@example.com", details.Email)
		assert.Equal(t, "Bob", details.FirstName)
		assert.Equal(t, "Q", details.MiddleName)
		assert.Equal(t, "Builder", details.LastName)
		assert.True(t, details.IsActive)
		assert.True(t, details.IsVerified)
		mockRepo.AssertExpectations(t)
	})
}
