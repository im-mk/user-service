package services

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/im-mk/user-service/src/models"
	"github.com/stretchr/testify/assert"
)

func TestDefaultTokenProvider_GenerateAccessToken(t *testing.T) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)

	cfg := models.AuthConfig{
		Issuer:             "issuer",
		Audience:           "audience",
		TokenExpirySeconds: 3600,
		Kid:                "kid",
	}

	provider := NewDefaultTokenProvider(priv, cfg)

	tokStr, err := provider.GenerateAccessToken("user123", "alice")
	assert.NoError(t, err)
	assert.NotEmpty(t, tokStr)

	// parse token using public key to verify claims
	pub := &priv.PublicKey
	parsed, err := jwt.Parse(tokStr, func(token *jwt.Token) (interface{}, error) {
		return pub, nil
	}, jwt.WithIssuer(cfg.Issuer), jwt.WithAudience(cfg.Audience), jwt.WithValidMethods([]string{"RS256"}))
	assert.NoError(t, err)
	assert.True(t, parsed.Valid)

	claims := parsed.Claims.(jwt.MapClaims)
	assert.Equal(t, "user123", claims["sub"])
	assert.Equal(t, "alice", claims["username"])
	// exp should be roughly now + cfg.TokenExpirySeconds
	exp := int64(claims["exp"].(float64))
	assert.WithinDuration(t, time.Now().Add(time.Duration(cfg.TokenExpirySeconds)*time.Second), time.Unix(exp, 0), 5*time.Second)
}

func TestDefaultTokenProvider_GenerateRefreshToken(t *testing.T) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)

	provider := NewDefaultTokenProvider(priv, models.AuthConfig{})

	rt1, err := provider.GenerateRefreshToken()
	assert.NoError(t, err)
	assert.NotEmpty(t, rt1)

	rt2, err := provider.GenerateRefreshToken()
	assert.NoError(t, err)
	assert.NotEmpty(t, rt2)

	assert.NotEqual(t, rt1, rt2)
}
