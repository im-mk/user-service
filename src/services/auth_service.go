package services

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strconv"
	"time"

	"github.com/im-mk/user-service/src/models"
	"github.com/im-mk/user-service/src/repositories"
	"golang.org/x/crypto/bcrypt"
)

type TokenProvider interface {
	GenerateAccessToken(userID, username string) (string, error)
	GenerateRefreshToken() (string, error)
}

type AuthService struct {
	UserRepo         repositories.UserRepositoryInterface
	RefreshTokenRepo repositories.RefreshTokenRepositoryInterface
	TokenProvider    TokenProvider
	AuthConfig       models.AuthConfig
}

func NewAuthService(
	userRepo repositories.UserRepositoryInterface,
	refreshRepo repositories.RefreshTokenRepositoryInterface,
	tokenProvider TokenProvider,
	authConfig models.AuthConfig,
) *AuthService {

	return &AuthService{
		UserRepo:         userRepo,
		RefreshTokenRepo: refreshRepo,
		TokenProvider:    tokenProvider,
		AuthConfig:       authConfig,
	}
}

func (s *AuthService) Login(creds models.LoginRequest) (string, string, error) {

	user, err := s.UserRepo.GetUserByUsername(creds.Username)
	if err != nil {
		return "", "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
		return "", "", errors.New("invalid credentials")
	}

	if !user.IsActive {
		return "", "", errors.New("account inactive")
	}

	if !user.IsVerified {
		return "", "", errors.New("account unverified")
	}

	accessToken, err := s.TokenProvider.GenerateAccessToken(
		strconv.Itoa(user.ID),
		user.Username,
	)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := s.TokenProvider.GenerateRefreshToken()
	if err != nil {
		return "", "", err
	}

	hash := hashToken(refreshToken)

	err = s.RefreshTokenRepo.SaveRefreshToken(
		hash,
		strconv.Itoa(user.ID),
		time.Now().Add(7*24*time.Hour),
	)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *AuthService) Refresh(oldRefreshToken string) (string, string, error) {
	hash := hashToken(oldRefreshToken)

	userID, err := s.RefreshTokenRepo.GetRefreshToken(hash)
	if err != nil {
		return "", "", errors.New("invalid refresh token")
	}

	_ = s.RefreshTokenRepo.DeleteRefreshToken(hash)

	newRefreshToken, err := s.TokenProvider.GenerateRefreshToken()
	if err != nil {
		return "", "", err
	}

	newHash := hashToken(newRefreshToken)

	err = s.RefreshTokenRepo.SaveRefreshToken(
		newHash,
		userID,
		time.Now().Add(time.Duration(s.AuthConfig.RefreshTokenExpirySeconds)*time.Second),
	)

	if err != nil {
		return "", "", err
	}

	uid, err := strconv.Atoi(userID)
	if err != nil {
		return "", "", errors.New("invalid user ID")
	}

	user, err := s.UserRepo.GetUserByID(uid)
	if err != nil {
		return "", "", err
	}

	accessToken, err := s.TokenProvider.GenerateAccessToken(
		userID,
		user.Username,
	)
	if err != nil {
		return "", "", err
	}

	return accessToken, newRefreshToken, nil
}

func (s *AuthService) Logout(refreshToken string) error {
	hash := hashToken(refreshToken)
	return s.RefreshTokenRepo.DeleteRefreshToken(hash)
}

func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}
