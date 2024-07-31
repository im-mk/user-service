package repositories

import (
	"database/sql"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) UserExists(username, email string) (bool, error) {
	var exists bool
	err := r.DB.QueryRow(`SELECT EXISTS (
        SELECT 1 FROM users WHERE username = $1 OR email = $2
    )`, username, email).Scan(&exists)
	return exists, err
}

func (r *UserRepository) CreateUser(username, email, password string) error {
	_, err := r.DB.Exec(`INSERT INTO users (username, email, password) VALUES ($1, $2, $3)`,
		username, email, password)
	return err
}
