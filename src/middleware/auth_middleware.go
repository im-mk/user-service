package middleware

import (
	"crypto/rsa"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/im-mk/user-service/src/models"
)

func AuthMiddleware(publicKey *rsa.PublicKey, cfg models.AuthConfig) gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		opts := []jwt.ParserOption{
			jwt.WithValidMethods([]string{"RS256"}),
		}
		if cfg.Issuer != "" {
			opts = append(opts, jwt.WithIssuer(cfg.Issuer))
		}
		if cfg.Audience != "" {
			opts = append(opts, jwt.WithAudience(cfg.Audience))
		}
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return publicKey, nil
		}, opts...)

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token claims",
			})
			return
		}

		c.Set("userID", claims["sub"])
		c.Set("username", claims["username"])
		c.Set("scope", claims["scope"])

		c.Next()
	}
}
