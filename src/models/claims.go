package models

import "github.com/golang-jwt/jwt/v4"

type Claims struct {
	UserID   string `json:"sub"`
	Username string `json:"username"`
	Scope    string `json:"scope"`
	jwt.RegisteredClaims
}
