package controllers

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/im-mk/user-service/src/models"
	"github.com/im-mk/user-service/src/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) UserExists(username, email string) (bool, error) {
	args := m.Called(username, email)
	return args.Bool(0), args.Error(1)
}
func (m *MockUserRepo) GetUserByUsername(username string) (*models.User, error) {
	args := m.Called(username)
	if u := args.Get(0); u != nil {
		return u.(*models.User), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockUserRepo) GetUserByID(id int) (*models.User, error) {
	args := m.Called(id)
	if u := args.Get(0); u != nil {
		return u.(*models.User), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockUserRepo) CreateUser(user models.User) error {
	args := m.Called(user)
	return args.Error(0)
}
func (m *MockUserRepo) AnyUserExists() (bool, error) {
	args := m.Called()
	return args.Bool(0), args.Error(1)
}

type MockRefreshRepo struct {
	mock.Mock
}

func (m *MockRefreshRepo) SaveRefreshToken(hash, userID string, expiresAt time.Time) error {
	args := m.Called(hash, userID, expiresAt)
	return args.Error(0)
}
func (m *MockRefreshRepo) GetRefreshToken(hash string) (string, error) {
	args := m.Called(hash)
	return args.String(0), args.Error(1)
}
func (m *MockRefreshRepo) DeleteRefreshToken(hash string) error {
	args := m.Called(hash)
	return args.Error(0)
}
func (m *MockRefreshRepo) DeleteExpiredTokens() error {
	args := m.Called()
	return args.Error(0)
}

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

func TestAuthController_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)
	config := models.AuthConfig{
		Issuer:             "test",
		Audience:           "test",
		TokenExpirySeconds: 3600,
	}

	t.Run("successful login", func(t *testing.T) {
		mockUserRepo := new(MockUserRepo)
		mockRefreshRepo := new(MockRefreshRepo)
		mockTokenProv := new(MockTokenProvider)

		hashedPw, _ := bcrypt.GenerateFromPassword([]byte("pass123"), bcrypt.DefaultCost)
		mockUserRepo.On("GetUserByUsername", "alice").Return(&models.User{
			ID:         1,
			Username:   "alice",
			Password:   string(hashedPw),
			IsActive:   true,
			IsVerified: true,
		}, nil)
		mockTokenProv.On("GenerateAccessToken", "1", "alice").Return("access_token", nil)
		mockTokenProv.On("GenerateRefreshToken").Return("refresh_token", nil)
		mockRefreshRepo.On("SaveRefreshToken", mock.Anything, "1", mock.Anything).Return(nil)

		svc := services.NewAuthService(mockUserRepo, mockRefreshRepo, mockTokenProv, config)
		ctrl := NewAuthController(svc)
		r := gin.New()
		r.POST("/login", ctrl.Login)

		body := models.LoginRequest{
			Username: "alice",
			Password: "pass123",
		}
		jsonBody, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]string
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "access_token", resp["token"])
		assert.Equal(t, "refresh_token", resp["refresh_token"])
	})

	t.Run("invalid credentials", func(t *testing.T) {
		mockUserRepo := new(MockUserRepo)
		mockRefreshRepo := new(MockRefreshRepo)
		mockTokenProv := new(MockTokenProvider)

		mockUserRepo.On("GetUserByUsername", "alice").Return(nil, errors.New("not found"))

		svc := services.NewAuthService(mockUserRepo, mockRefreshRepo, mockTokenProv, config)
		ctrl := NewAuthController(svc)
		r := gin.New()
		r.POST("/login", ctrl.Login)

		body := models.LoginRequest{
			Username: "alice",
			Password: "wrongpass",
		}
		jsonBody, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}

func TestAuthController_Refresh(t *testing.T) {
	gin.SetMode(gin.TestMode)
	config := models.AuthConfig{
		Issuer:             "test",
		Audience:           "test",
		TokenExpirySeconds: 3600,
	}

	t.Run("successful refresh", func(t *testing.T) {
		mockUserRepo := new(MockUserRepo)
		mockRefreshRepo := new(MockRefreshRepo)
		mockTokenProv := new(MockTokenProvider)

		oldRefreshToken := "old_refresh_token"
		tokenHash := hashToken(oldRefreshToken)

		mockRefreshRepo.On("GetRefreshToken", tokenHash).Return("1", nil)
		mockUserRepo.On("GetUserByID", 1).Return(&models.User{
			ID:       1,
			Username: "alice",
		}, nil)
		mockTokenProv.On("GenerateAccessToken", "1", "alice").Return("new_access", nil)
		mockTokenProv.On("GenerateRefreshToken").Return("new_refresh", nil)
		mockRefreshRepo.On("DeleteRefreshToken", tokenHash).Return(nil)
		newRefreshHash := hashToken("new_refresh")
		mockRefreshRepo.On("SaveRefreshToken", newRefreshHash, "1", mock.Anything).Return(nil)

		svc := services.NewAuthService(mockUserRepo, mockRefreshRepo, mockTokenProv, config)
		ctrl := NewAuthController(svc)
		r := gin.New()
		r.POST("/refresh", ctrl.Refresh)

		body := models.RefreshRequest{
			RefreshToken: oldRefreshToken,
		}
		jsonBody, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/refresh", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]string
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "new_access", resp["token"])
		assert.Equal(t, "new_refresh", resp["refresh_token"])
	})
}
