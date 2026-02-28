package services

import (
	"errors"

	"github.com/im-mk/user-service/src/models"
	"github.com/im-mk/user-service/src/repositories"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	UserRepo repositories.UserRepositoryInterface
}

func NewUserService(userRepo repositories.UserRepositoryInterface) *UserService {
	return &UserService{UserRepo: userRepo}
}

func (s *UserService) CreateUser(req models.CreateUserRequest) error {

	exists, err := s.UserRepo.UserExists(req.Username, req.Email)
	if err != nil {
		return errors.New("failed to check for existing user")
	}

	if exists {
		return errors.New("username or email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := models.User{
		Username:   req.Username,
		Email:      req.Email,
		Password:   string(hashedPassword),
		FirstName:  req.FirstName,
		MiddleName: req.MiddleName,
		LastName:   req.LastName,
		IsActive:   true,
		IsVerified: false,
	}

	return s.UserRepo.CreateUser(user)
}

func (s *UserService) GetUser(userId int) (*models.UserDetails, error) {

	user, err := s.UserRepo.GetUserByID(userId)
	if err != nil {
		return nil, err
	}

	userDetails := &models.UserDetails{
		ID:         user.ID,
		Username:   user.Username,
		Email:      user.Email,
		FirstName:  user.FirstName,
		MiddleName: user.MiddleName,
		LastName:   user.LastName,
		IsActive:   user.IsActive,
		IsVerified: user.IsVerified,
	}

	return userDetails, nil
}

// Bootstrap creates the first user if no users exist in the system.
func (s *UserService) Bootstrap(req models.CreateUserRequest) error {

	any, err := s.UserRepo.AnyUserExists()
	if err != nil || any {
		return errors.New("an error occurred whilst performing bootstrap")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := models.User{
		Username:   req.Username,
		Email:      req.Email,
		Password:   string(hashedPassword),
		FirstName:  req.FirstName,
		MiddleName: req.MiddleName,
		LastName:   req.LastName,
		IsActive:   true,
		IsVerified: true, // verify bootstrap user by default
	}

	return s.UserRepo.CreateUser(user)
}
