package middleware

import (
	"crypto/rand"
	"crypto/rsa"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/im-mk/user-service/src/models"
	"github.com/stretchr/testify/assert"
)

func makeTestToken(t *testing.T, priv *rsa.PrivateKey, cfg models.AuthConfig, sub, username string) string {
	claims := jwt.MapClaims{
		"sub":      sub,
		"username": username,
		"iss":      cfg.Issuer,
		"aud":      []string{cfg.Audience},
		"exp":      time.Now().Add(time.Duration(cfg.TokenExpirySeconds) * time.Second).Unix(),
		"iat":      time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = cfg.Kid
	sign, err := token.SignedString(priv)
	assert.NoError(t, err)
	return sign
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)
	pub := &priv.PublicKey

	cfg := models.AuthConfig{
		Issuer:             "test-issuer",
		Audience:           "test-aud",
		TokenExpirySeconds: 3600,
		Kid:                "kid1",
	}

	token := makeTestToken(t, priv, cfg, "42", "bob")

	r := gin.New()
	r.Use(AuthMiddleware(pub, cfg))
	r.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthMiddleware_InvalidAudience(t *testing.T) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)
	pub := &priv.PublicKey

	cfg := models.AuthConfig{
		Issuer:             "issuer",
		Audience:           "aud1",
		TokenExpirySeconds: 3600,
		Kid:                "kid",
	}

	// generate token with wrong audience
	badCfg := cfg
	badCfg.Audience = "wrong"
	token := makeTestToken(t, priv, badCfg, "1", "name")

	r := gin.New()
	r.Use(AuthMiddleware(pub, cfg))
	r.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
