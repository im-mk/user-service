package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/im-mk/user-service/src/models"
	"github.com/im-mk/user-service/src/services"
	_ "github.com/lib/pq"
)

type AuthController struct {
	AuthService *services.AuthService
}

func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{AuthService: authService}
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
func (ctrl *AuthController) Login(c *gin.Context) {
	var creds models.LoginRequest
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	access, refresh, err := ctrl.AuthService.Login(creds)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": access, "refresh_token": refresh})
}

// Refresh exchanges a refresh token for a new access token
func (ctrl *AuthController) Refresh(c *gin.Context) {
	var req models.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	access, refresh, err := ctrl.AuthService.Refresh(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": access, "refresh_token": refresh})
}
