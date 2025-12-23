package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/im-mk/user-service/src/models"
)

type JWTGenerator func(userId string, username string, jwtKey []byte) (string, error)

func GenerateJWT(userID string, username string, jwtKey []byte) (string, error) {
	claims := models.Claims{
		UserID:   userID,
		Username: username,
		Scope:    "orders.read orders.write",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "https://userservice.internal",
			Audience:  []string{"orders-api"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}
