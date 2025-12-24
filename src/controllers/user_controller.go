package controllers

import (
	"net/http"
	"strconv"

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

// @Summary		Return user details
// @Description	Return user details with the given ID
// @Tags			users
// @Accept			json
// @Produce		json
// @Param			user	body		models.CreateUserRequest	true	"Create User Request"
// @Success		200		{object}	models.User
// @Failure		400		{object}	gin.H
// @Failure		500		{object}	gin.H
// @Router			/users [post]
func (ctrl *UserController) GetUser(c *gin.Context) {
	userId := c.Param("id")
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	user, err := ctrl.UserService.GetUser(userIdInt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, user)
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
