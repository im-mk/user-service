package services

import (
	"time"

	"github.com/stretchr/testify/mock"
)

type MockRefreshTokenRepository struct {
	mock.Mock
}

func (m *MockRefreshTokenRepository) SaveRefreshToken(tokenHash, userID string, expiresAt time.Time) error {
	args := m.Called(tokenHash, userID, expiresAt)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) GetRefreshToken(tokenHash string) (string, error) {
	args := m.Called(tokenHash)
	return args.String(0), args.Error(1)
}

func (m *MockRefreshTokenRepository) DeleteRefreshToken(tokenHash string) error {
	args := m.Called(tokenHash)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) DeleteExpiredTokens() error {
	args := m.Called()
	return args.Error(0)
}
