package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/im-mk/user-service/src/models"
)

type JWTGenerator func(username string, jwtKey []byte) (string, error)

func GenerateJWT(username string, jwtKey []byte) (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &models.Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
