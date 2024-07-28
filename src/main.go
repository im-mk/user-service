package main

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	_ "github.com/im-mk/user-service/src/docs"
	"github.com/im-mk/user-service/src/models"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	// temp logic
	if creds.Username != "user1" || creds.Password != "password1" {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
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
		c.JSON(500, gin.H{"error": "Could not create token"})
		return
	}

	c.JSON(200, gin.H{"token": tokenString})
}

// @title           user-service
// @version         1.0
// @description     service to manager users
// @termsOfService  http://swagger.io/terms/

// @contact.name   im-mk
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.POST("/login", loginHandler)

	router.Run(":8080")
}
