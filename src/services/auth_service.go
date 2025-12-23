package services

import (
	"errors"
	"strconv"

	"github.com/im-mk/user-service/src/models"
	"github.com/im-mk/user-service/src/repositories"
	"github.com/im-mk/user-service/src/utils"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepo    repositories.UserRepositoryInterface
	JwtKey      []byte
	GenerateJWT utils.JWTGenerator
}

func NewAuthService(userRepo repositories.UserRepositoryInterface, jwtKey []byte, generateJWT utils.JWTGenerator) *AuthService {
	if generateJWT == nil {
		generateJWT = utils.GenerateJWT
	}

	return &AuthService{UserRepo: userRepo, JwtKey: jwtKey, GenerateJWT: generateJWT}
}

func (s *AuthService) Login(creds models.LoginRequest) (string, error) {
	user, err := s.UserRepo.GetUserByUsername(creds.Username)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := s.GenerateJWT(strconv.Itoa(user.ID), user.Username, s.JwtKey)
	if err != nil {
		return "", err
	}

	return token, nil
}
