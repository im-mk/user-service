package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/im-mk/user-service/src/controllers"
	_ "github.com/im-mk/user-service/src/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func registerRoutes(userController *controllers.UserController, port string) {
	router := gin.Default()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.POST("/login", userController.Login)

	auth := router.Group("/")
	auth.Use(authMiddleware())
	{
		auth.POST("/users", userController.CreateUser)
	}
	router.Run(fmt.Sprintf("127.0.0.1:%s", port))
}
