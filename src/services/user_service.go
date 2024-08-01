package services

import (
	"errors"

	"github.com/im-mk/user-service/src/models"
	"github.com/im-mk/user-service/src/repositories"
	"github.com/im-mk/user-service/src/utils"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	UserRepo *repositories.UserRepository
	JwtKey   []byte
}

func NewUserService(userRepo *repositories.UserRepository, jwtKey []byte) *UserService {
	return &UserService{UserRepo: userRepo, JwtKey: jwtKey}
}

func (s *UserService) Login(creds models.LoginRequest) (string, error) {
	user, err := s.UserRepo.GetUserByUsername(creds.Username)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := utils.GenerateJWT(user.Username, s.JwtKey)
	if err != nil {
		return "", err
	}

	return token, nil
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
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	return s.UserRepo.CreateUser(user)
}
