package controllers

import (
	"net/http"

	"totesbackend/controllers/utilities"
	"totesbackend/services"

	"github.com/gin-gonic/gin"
)

type UserCredentialValidationController struct {
	Service *services.UserCredentialValidationService
	Auth    *utilities.AuthorizationUtil
	Log     *utilities.LogUtil
}

func NewUserCredentialValidationController(service *services.UserCredentialValidationService, auth *utilities.AuthorizationUtil, log *utilities.LogUtil) *UserCredentialValidationController {
	return &UserCredentialValidationController{Service: service, Auth: auth, Log: log}
}

// LoginData defines the structure for user login request
type LoginData struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// ValidateUserCredentials godoc
// @Summary      Validate user credentials
// @Description  Validates the user's credentials (email and password) for login.
// @Tags         authentication
// @Accept       json
// @Produce      json
// @Param        body    body     LoginData  true  "User credentials to validate"
// @Success      200     {object}   models.MessageResponse "Login successful message"
// @Failure      400     {object}  models.ErrorResponse  "Invalid request body"
// @Failure      403     {object}  models.ErrorResponse  "User account is not active"
// @Failure      401     {object}  models.ErrorResponse  "Invalid email or password"
// @Failure      500     {object}  models.ErrorResponse  "Error validating credentials"
// @Security     ApiKeyAuth
// @Router       /login [post]
func (ucvc *UserCredentialValidationController) ValidateUserCredentials(c *gin.Context) {
	if ucvc.Log.RegisterLog(c, "Attempting user login") != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
		return
	}

	var loginData LoginData

	if err := c.ShouldBindJSON(&loginData); err != nil {
		_ = ucvc.Log.RegisterLog(c, "Invalid request body for login")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err := ucvc.Service.ValidateUserCredentials(loginData.Email, loginData.Password)
	if err != nil {
		if err.Error() == "user is not active" {
			_ = ucvc.Log.RegisterLog(c, "Login attempt for inactive user: "+loginData.Email)
			c.JSON(http.StatusForbidden, gin.H{"error": "User account is not active"})
			return
		}

		_ = ucvc.Log.RegisterLog(c, "Login failed for user: "+loginData.Email)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	_ = ucvc.Log.RegisterLog(c, "Login successful for user: "+loginData.Email)

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
	})

}
