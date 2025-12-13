package main

import (
	"github.com/im-mk/user-service/src/controllers"
	_ "github.com/im-mk/user-service/src/docs"
	"github.com/im-mk/user-service/src/repositories"
	"github.com/im-mk/user-service/src/services"
	"github.com/im-mk/user-service/src/utils"
)

// @title						user-service
// @version					1.0
// @description				service to manage users
// @contact.name				im-mk
// @contact.url				http://github.com/im-mk
// @host						localhost:8080
// @BasePath					/
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
func main() {

	appConfig := GetConfig()
	db := initDB(appConfig.DB)
	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo, []byte(appConfig.JWTKey), utils.GenerateJWT)
	userController := controllers.NewUserController(userService)

	registerRoutes(userController, appConfig.App, []byte(appConfig.JWTKey))
}
