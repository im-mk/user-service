package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/im-mk/user-service/src/models"
	"github.com/im-mk/user-service/src/services"
	_ "github.com/lib/pq"
)

type UserController struct {
	UserService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{UserService: userService}
}

// @Summary		Logs in a user
// @Description	Logs in a user and returns a JWT token
// @Tags			auth
// @Accept			json
// @Produce		json
// @Param			credentials	body		models.LoginRequest	true	"User credentials"
// @Success		200			{string}	string				"token"
// @Failure		400			{object}	gin.H				"Invalid request"
// @Failure		500			{object}	gin.H				"Could not create token"
// @Router			/login [post]
func (ctrl *UserController) Login(c *gin.Context) {
	var creds models.LoginRequest
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := ctrl.UserService.Login(creds)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// @Summary		Create a new user
// @Description	Create a new user with the input payload
// @Tags			users
// @Accept			json
// @Produce		json
// @Param			user	body		models.CreateUserRequest	true	"Create User Request"
// @Success		200		{object}	models.User
// @Failure		400		{object}	gin.H
// @Failure		500		{object}	gin.H
// @Router			/users [post]
func (ctrl *UserController) CreateUser(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.UserService.CreateUser(req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

// @Summary    Bootstrap first user
// @Description Create the first user if no users exist yet
// @Tags       bootstrap
// @Accept     json
// @Produce    json
// @Param      user  body      models.CreateUserRequest  true  "Create User Request"
// @Success    200   {object}  gin.H
// @Failure    400   {object}  gin.H
// @Failure    500   {object}  gin.H
// @Router     /bootstrap [post]
func (ctrl *UserController) Bootstrap(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.UserService.Bootstrap(req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bootstrap user created"})
}
