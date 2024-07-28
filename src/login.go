package main

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	_ "github.com/im-mk/user-service/src/docs"
	"github.com/im-mk/user-service/src/models"
)

var jwtKey = []byte("my_secret_key")

// @Summary Logs in a user
// @Description Logs in a user and returns a JWT token
// @Tags auth
// @Accept  json
// @Produce  json
// @Param   credentials body models.Credentials true "User credentials"
// @Success 200 {string} string "token"
// @Failure 400 {object} gin.H "Invalid request"
// @Failure 500 {object} gin.H "Could not create token"
// @Router /login [post]
func loginHandler(c *gin.Context) {
	var creds models.Credentials
	if err := c.BindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// temp logic
	if creds.Username != "user1" || creds.Password != "password1" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &models.Claims{
		Username: creds.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
