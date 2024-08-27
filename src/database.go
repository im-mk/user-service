package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/im-mk/user-service/src/docs"
)

func initDB(dbConnection DBConfig) *sql.DB {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbConnection.Host,
		dbConnection.User,
		dbConnection.Password,
		dbConnection.DBName,
		dbConnection.Port,
	)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Ensure the database is available
	// if err := db.Ping(); err != nil {
	// 	log.Fatalf("Failed to ping database: %v", err)
	// }

	return db
}
