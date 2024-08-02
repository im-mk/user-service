package main

import (
	"database/sql"

	"github.com/im-mk/user-service/src/config"
	"github.com/im-mk/user-service/src/controllers"
	_ "github.com/im-mk/user-service/src/docs"
	"github.com/im-mk/user-service/src/repositories"
	"github.com/im-mk/user-service/src/services"
)

var (
	db     *sql.DB
	jwtKey = []byte("my_secret_key")
)

// @title						user-service
// @version					1.0
// @description				service to manager users
// @contact.name				im-mk
// @contact.url				http://github.com/im-mk
// @host						localhost:8080
// @BasePath					/
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
func main() {

	appConfig := config.GetConfig()
	initDB()
	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo, jwtKey)
	userController := controllers.NewUserController(userService)

	registerRoutes(userController, appConfig.Port)
}
