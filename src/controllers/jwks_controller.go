package controllers

import (
	"crypto/rsa"
	"encoding/base64"
	"math/big"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type JwksController struct {
	SigningKey *rsa.PublicKey
}

func NewJwksController(signingKey *rsa.PublicKey) *JwksController {
	return &JwksController{SigningKey: signingKey}
}

// @Summary Return JWKS
// @Description Returns the JSON Web Key Set
// @Tags auth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /.well-known/jwks.json [get]
func (ctrl *JwksController) JwksHandler(c *gin.Context) {
	var keyID = "auth-key-v1"
	publicKey := ctrl.SigningKey
	n := base64.RawURLEncoding.EncodeToString(publicKey.N.Bytes())
	e := base64.RawURLEncoding.EncodeToString(big.NewInt(int64(publicKey.E)).Bytes())

	jwks := map[string]interface{}{
		"keys": []interface{}{
			map[string]interface{}{
				"kty": "RSA",
				"use": "sig",
				"kid": keyID,
				"alg": "RS256",
				"n":   n,
				"e":   e,
			},
		},
	}

	c.Header("Cache-Control", "public, max-age=3600")
	c.JSON(http.StatusOK, jwks)
}
