package repositories

import (
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
)

type RefreshTokenRepositoryInterface interface {
	SaveRefreshToken(tokenHash, userID string, expiresAt time.Time) error
	GetRefreshToken(tokenHash string) (string, error)
	DeleteRefreshToken(tokenHash string) error
	DeleteExpiredTokens() error
}

type RefreshTokenRepository struct {
	DB *sqlx.DB
}

func NewRefreshTokenRepository(db *sqlx.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{DB: db}
}

func (r *RefreshTokenRepository) SaveRefreshToken(
	tokenHash string,
	userID string,
	expiresAt time.Time,
) error {
	_, err := r.DB.Exec(`
		INSERT INTO refresh_tokens (token_hash, user_id, expires_at)
		VALUES ($1, $2, $3)
	`, tokenHash, userID, expiresAt)

	return err
}

func (r *RefreshTokenRepository) GetRefreshToken(tokenHash string) (string, error) {
	var userID string
	err := r.DB.Get(&userID, `
		SELECT user_id
		FROM refresh_tokens
		WHERE token_hash = $1
		  AND expires_at > NOW()
	`, tokenHash)

	if err == sql.ErrNoRows {
		return "", errors.New("refresh token not found or expired")
	}
	if err != nil {
		return "", err
	}

	return userID, nil
}

func (r *RefreshTokenRepository) DeleteRefreshToken(tokenHash string) error {
	_, err := r.DB.Exec(`
		DELETE FROM refresh_tokens
		WHERE token_hash = $1
	`, tokenHash)

	return err
}

func (r *RefreshTokenRepository) DeleteExpiredTokens() error {
	_, err := r.DB.Exec(`
		DELETE FROM refresh_tokens
		WHERE expires_at <= NOW()
	`)
	return err
}
