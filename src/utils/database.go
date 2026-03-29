package utils

import (
	"fmt"
	"log"

	_ "github.com/im-mk/user-service/src/docs"
	"github.com/im-mk/user-service/src/models"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func InitDB(dbConnection models.DBConfig) *sqlx.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbConnection.Host,
		dbConnection.User,
		dbConnection.Password,
		dbConnection.DBName,
		dbConnection.Port,
	)
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Printf("Failed to ping database: %v", err)
	}

	return db
}
