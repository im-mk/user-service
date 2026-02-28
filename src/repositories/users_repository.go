package repositories

import (
	"github.com/jmoiron/sqlx"

	"github.com/im-mk/user-service/src/models"
)

type UserRepositoryInterface interface {
	UserExists(username, email string) (bool, error)
	GetUserByUsername(username string) (*models.User, error)
	GetUserByID(userID int) (*models.User, error)
	CreateUser(user models.User) error
	AnyUserExists() (bool, error)
}

type UserRepository struct {
	DB *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) UserExists(username, email string) (bool, error) {
	var exists bool
	err := r.DB.Get(&exists, `SELECT EXISTS (
        SELECT 1 FROM users WHERE username = $1 OR email = $2
    )`, username, email)
	return exists, err
}

func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.DB.Get(&user, `
		SELECT id, username, email, password,
		       first_name, middle_name, last_name,
		       is_active, is_verified
		FROM users
		WHERE username = $1
	`, username)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUserByID(userID int) (*models.User, error) {
	var user models.User
	err := r.DB.Get(&user, `
		SELECT id, username, email, password,
		       first_name, middle_name, last_name,
		       is_active, is_verified
		FROM users
		WHERE id = $1
	`, userID)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) CreateUser(user models.User) error {
	_, err := r.DB.Exec(`INSERT INTO users (username, email, password, first_name, middle_name, last_name, is_active, is_verified) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		user.Username, user.Email, user.Password,
		user.FirstName, user.MiddleName, user.LastName, user.IsActive, user.IsVerified)
	return err
}

func (r *UserRepository) AnyUserExists() (bool, error) {
	var exists bool
	err := r.DB.Get(&exists, `SELECT EXISTS (SELECT 1 FROM users)`)
	return exists, err
}
