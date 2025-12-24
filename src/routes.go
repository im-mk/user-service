package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/im-mk/user-service/src/controllers"
	_ "github.com/im-mk/user-service/src/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func registerRoutes(userController *controllers.UserController, authController *controllers.AuthController, appConfig AppConfig, jwtKey []byte) {
	router := gin.Default()
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, "healthy")
	})
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.POST("/login", authController.Login)
	router.POST("/refresh", authController.Refresh)
	router.POST("/bootstrap", userController.Bootstrap)

	auth := router.Group("/")
	auth.Use(authMiddleware(jwtKey))
	{
		auth.POST("/users", userController.CreateUser)
		auth.GET("/users/:id", userController.GetUser)
	}
	router.Run(fmt.Sprintf("%s:%s", appConfig.Host, appConfig.Port))
}
