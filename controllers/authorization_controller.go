package controllers

import (
	"net/http"
	"strconv"

	"totesbackend/controllers/utilities"
	"totesbackend/services"

	"github.com/gin-gonic/gin"
)

type AuthorizationController struct {
	Service *services.AuthorizationService
	Log     *utilities.LogUtil
}

func NewAuthorizationController(service *services.AuthorizationService, log *utilities.LogUtil) *AuthorizationController {
	return &AuthorizationController{Service: service, Log: log}
}

// CheckUserPermission godoc
// @Summary      Check if a user has a specific permission
// @Description  Verifies if the user with the provided email has the specified permission ID
// @Tags         authorization
// @Accept       json
// @Produce      json
// @Param        email       query     string  true  "User's email address"
// @Param        permission_id  query  string  true  "Permission ID to check"
// @Success      200        {object}  models.MessageResponse   "Response with the permission status"
// @Failure      400        {object}  models.ErrorResponse   "Invalid or missing parameters"
// @Failure      500        {object}  models.ErrorResponse   "Error checking permission"
// @Router       /auth/check-permission [get]
func (ac *AuthorizationController) CheckUserPermission(c *gin.Context) {
	email := c.Query("email")
	permissionID := c.Query("permission_id")
	permissionStr, err := strconv.Atoi(permissionID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Permission ID"})
		return
	}

	if email == "" || permissionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email and permission_id are required"})
		return
	}

	hasPermission, err := ac.Service.UserHasPermission(email, permissionStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking permission"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"has_permission": hasPermission})
}
