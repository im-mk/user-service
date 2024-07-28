package main

import (
	"database/sql"
	"log"

	_ "github.com/im-mk/user-service/src/docs"
)

var db *sql.DB

func initDB() {
	var err error
	dsn := "host=localhost user=postgres password=postgres dbname=user-service port=5432 sslmode=disable"
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Ensure the database is available
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
}
