package services

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/im-mk/user-service/src/models"
)

type DefaultTokenProvider struct {
	PrivateKey *rsa.PrivateKey
	AuthConfig models.AuthConfig
}

func NewDefaultTokenProvider(
	privateKey *rsa.PrivateKey,
	authConfig models.AuthConfig,
) *DefaultTokenProvider {

	return &DefaultTokenProvider{
		PrivateKey: privateKey,
		AuthConfig: authConfig,
	}
}

func (p *DefaultTokenProvider) GenerateAccessToken(userID, username string) (string, error) {

	claims := jwt.MapClaims{
		"sub":      userID,
		"username": username,
		"iss":      p.AuthConfig.Issuer,
		"aud":      []string{p.AuthConfig.Audience},
		"exp":      time.Now().Add(time.Duration(p.AuthConfig.TokenExpirySeconds) * time.Second).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = p.AuthConfig.Kid

	return token.SignedString(p.PrivateKey)
}

func (p *DefaultTokenProvider) GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
