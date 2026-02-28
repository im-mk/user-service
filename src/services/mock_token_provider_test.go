package services

import (
	"github.com/stretchr/testify/mock"
)

// MockTokenProvider is a testify mock for TokenProvider interface used in tests.
// It allows setting expectations on access/refresh token generation.

type MockTokenProvider struct {
	mock.Mock
}

func (m *MockTokenProvider) GenerateAccessToken(userID, username string) (string, error) {
	args := m.Called(userID, username)
	return args.String(0), args.Error(1)
}

func (m *MockTokenProvider) GenerateRefreshToken() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}
