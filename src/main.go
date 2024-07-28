package main

import (
	_ "github.com/im-mk/user-service/src/docs"
)

// @title           user-service
// @version         1.0
// @description     service to manager users
// @contact.name   im-mk
// @contact.url    http://github.com/im-mk
// @host      localhost:8080
// @BasePath  /
// @securityDefinitions.basic  BasicAuth
// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	initDB()
	registerRoutes()
}
