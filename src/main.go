package main

import (
	"github.com/im-mk/user-service/src/controllers"
	_ "github.com/im-mk/user-service/src/docs"
	"github.com/im-mk/user-service/src/repositories"
	"github.com/im-mk/user-service/src/services"
)

// @title						user-service
// @version					1.0
// @description				service to manage users
// @contact.name				im-mk
// @contact.url				http://github.com/im-mk
// @host						localhost:8040
// @BasePath					/
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
func main() {

	appConfig := GetConfig()
	db := initDB(appConfig.DB)
	userRepo := repositories.NewUserRepository(db)
	refreshTokenRepo := repositories.NewRefreshTokenRepository(db)
	userService := services.NewUserService(userRepo)
	authService := services.NewAuthService(userRepo, refreshTokenRepo, []byte(appConfig.JWTKey))

	userController := controllers.NewUserController(userService)
	authController := controllers.NewAuthController(authService)

	registerRoutes(userController, authController, appConfig.App, []byte(appConfig.JWTKey))
}
