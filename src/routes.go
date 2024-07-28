package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/im-mk/user-service/src/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func registerRoutes() {
	router := gin.Default()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.POST("/login", loginHandler)
	router.POST("/users", createUserHandler)
	router.Run(":8080")
}
