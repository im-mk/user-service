package controllers

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestJwksController_JwksHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("returns valid jwks", func(t *testing.T) {
		priv, err := rsa.GenerateKey(rand.Reader, 2048)
		assert.NoError(t, err)
		pub := &priv.PublicKey

		ctrl := NewJwksController(pub)
		r := gin.New()
		r.GET("/.well-known/jwks.json", ctrl.JwksHandler)

		req := httptest.NewRequest(http.MethodGet, "/.well-known/jwks.json", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)

		// check structure
		assert.Contains(t, resp, "keys")
		keys := resp["keys"].([]interface{})
		assert.Len(t, keys, 1)

		key := keys[0].(map[string]interface{})
		assert.Equal(t, "RSA", key["kty"])
		assert.Equal(t, "sig", key["use"])
		assert.Equal(t, "auth-key-v1", key["kid"])
		assert.Equal(t, "RS256", key["alg"])
		assert.Contains(t, key, "n")
		assert.Contains(t, key, "e")
	})

	t.Run("sets cache header", func(t *testing.T) {
		priv, err := rsa.GenerateKey(rand.Reader, 2048)
		assert.NoError(t, err)
		pub := &priv.PublicKey

		ctrl := NewJwksController(pub)
		r := gin.New()
		r.GET("/.well-known/jwks.json", ctrl.JwksHandler)

		req := httptest.NewRequest(http.MethodGet, "/.well-known/jwks.json", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		cacheControl := w.Header().Get("Cache-Control")
		assert.Equal(t, "public, max-age=3600", cacheControl)
	})
}
