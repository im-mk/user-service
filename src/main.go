package main

import (
	"database/sql"

	_ "github.com/im-mk/user-service/src/docs"
	"github.com/im-mk/user-service/src/repositories"
)

var (
	db       *sql.DB
	jwtKey   = []byte("my_secret_key")
	userRepo *repositories.UserRepository
)

// @title           user-service
// @version         1.0
// @description     service to manager users
// @contact.name   im-mk
// @contact.url    http://github.com/im-mk
// @host      localhost:8080
// @BasePath  /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	initDB()
	userRepo = repositories.NewUserRepository(db)
	registerRoutes()
}
