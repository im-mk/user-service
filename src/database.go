package main

import (
	"fmt"
	"log"

	_ "github.com/im-mk/user-service/src/docs"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // postgres driver
)

// initDB establishes a connection using sqlx and returns the wrapper DB.
func initDB(dbConnection DBConfig) *sqlx.DB {
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

	// Ensure the database is available
	if err := db.Ping(); err != nil {
		log.Printf("Failed to ping database: %v", err)
	}

	return db
}
