package repositories

import (
	"database/sql"

	"github.com/im-mk/user-service/src/models"
)

type UserRepositoryInterface interface {
	UserExists(username, email string) (bool, error)
	GetUserByUsername(username string) (*models.User, error)
	CreateUser(user models.User) error
	AnyUserExists() (bool, error)
}

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

func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.DB.QueryRow(`SELECT id, username, password FROM users WHERE username = $1`, username).
		Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) CreateUser(user models.User) error {
	_, err := r.DB.Exec(`INSERT INTO users (username, email, password) VALUES ($1, $2, $3)`,
		user.Username, user.Email, user.Password)
	return err
}

func (r *UserRepository) AnyUserExists() (bool, error) {
	var exists bool
	err := r.DB.QueryRow(`SELECT EXISTS (SELECT 1 FROM users)`).Scan(&exists)
	return exists, err
}
