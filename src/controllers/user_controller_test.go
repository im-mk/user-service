package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/im-mk/user-service/src/models"
	"github.com/im-mk/user-service/src/services"
)

// fakeRepo is a simple in-test implementation of UserRepositoryInterface
type fakeRepo struct {
	f func(int) (*models.User, error)
}

func (r *fakeRepo) UserExists(username, email string) (bool, error)         { return false, nil }
func (r *fakeRepo) GetUserByUsername(username string) (*models.User, error) { return nil, nil }
func (r *fakeRepo) GetUserByID(id int) (*models.User, error)                { return r.f(id) }
func (r *fakeRepo) CreateUser(user models.User) error                       { return nil }
func (r *fakeRepo) AnyUserExists() (bool, error)                            { return false, nil }

func TestGetUserHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		funcRepo := &fakeRepo{f: func(id int) (*models.User, error) {
			return &models.User{
				ID:         id,
				Username:   "alice",
				Email:      "alice@example.com",
				Password:   "x",
				FirstName:  "Alice",
				LastName:   "Liddell",
				IsActive:   true,
				IsVerified: true,
			}, nil
		}}

		userSvc := services.NewUserService(funcRepo)
		ctrl := NewUserController(userSvc)

		r := gin.New()
		r.GET("/users/:id", ctrl.GetUser)

		req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var got models.UserDetails
		err := json.Unmarshal(w.Body.Bytes(), &got)
		assert.NoError(t, err)
		assert.Equal(t, 1, got.ID)
		assert.Equal(t, "alice", got.Username)
		assert.Equal(t, "alice@example.com", got.Email)
		assert.Equal(t, "Alice", got.FirstName)
		assert.Equal(t, "Liddell", got.LastName)
		assert.True(t, got.IsActive)
		assert.True(t, got.IsVerified)
	})

	t.Run("not found", func(t *testing.T) {
		funcRepo := &fakeRepo{f: func(id int) (*models.User, error) { return nil, errors.New("not found") }}

		userSvc := services.NewUserService(funcRepo)
		ctrl := NewUserController(userSvc)

		r := gin.New()
		r.GET("/users/:id", ctrl.GetUser)

		req := httptest.NewRequest(http.MethodGet, "/users/2", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
