package utils

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/im-mk/user-service/src/models"
)

type AccessTokenGenerator func(userID, username string, key []byte) (string, error)

// GenerateAccessToken creates a short-lived JWT access token
func GenerateAccessToken(userID, username string, key []byte) (string, error) {
	claims := models.Claims{
		UserID:   userID,
		Username: username,
		Scope:    "orders.read orders.write",
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(),
			Issuer:    "https://userservice.internal",
			Audience:  []string{"orders-api"},
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(10 * time.Minute)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(key)
}

// GenerateRefreshToken creates a cryptographically secure opaque token
func GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
