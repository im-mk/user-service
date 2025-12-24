package services

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strconv"
	"time"

	"github.com/im-mk/user-service/src/models"
	"github.com/im-mk/user-service/src/repositories"
	"github.com/im-mk/user-service/src/utils"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepo         repositories.UserRepositoryInterface
	RefreshTokenRepo repositories.RefreshTokenRepositoryInterface
	AccessKey        []byte
}

func NewAuthService(
	userRepo repositories.UserRepositoryInterface,
	refreshRepo repositories.RefreshTokenRepositoryInterface,
	accessKey []byte,
) *AuthService {
	return &AuthService{
		UserRepo:         userRepo,
		RefreshTokenRepo: refreshRepo,
		AccessKey:        accessKey,
	}
}

// Login issues access + refresh tokens
func (s *AuthService) Login(creds models.LoginRequest) (string, string, error) {
	user, err := s.UserRepo.GetUserByUsername(creds.Username)
	if err != nil {
		return "", "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
		return "", "", errors.New("invalid credentials")
	}

	accessToken, err := utils.GenerateAccessToken(
		strconv.Itoa(user.ID),
		user.Username,
		s.AccessKey,
	)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := utils.GenerateRefreshToken()
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

// Refresh rotates refresh token and issues new access token
func (s *AuthService) Refresh(oldRefreshToken string) (string, string, error) {
	hash := hashToken(oldRefreshToken)

	userID, err := s.RefreshTokenRepo.GetRefreshToken(hash)
	if err != nil {
		return "", "", errors.New("invalid refresh token")
	}

	// rotate refresh token
	_ = s.RefreshTokenRepo.DeleteRefreshToken(hash)

	newRefreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		return "", "", err
	}

	newHash := hashToken(newRefreshToken)

	err = s.RefreshTokenRepo.SaveRefreshToken(
		newHash,
		userID,
		time.Now().Add(24*time.Hour),
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

	accessToken, err := utils.GenerateAccessToken(
		userID,
		user.Username,
		s.AccessKey,
	)
	if err != nil {
		return "", "", err
	}

	return accessToken, newRefreshToken, nil
}

// Logout revokes refresh token
func (s *AuthService) Logout(refreshToken string) error {
	hash := hashToken(refreshToken)
	return s.RefreshTokenRepo.DeleteRefreshToken(hash)
}

// hashToken hashes refresh tokens before storage
func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}
