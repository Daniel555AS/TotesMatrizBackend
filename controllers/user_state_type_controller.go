package controllers

import (
	"net/http"
	"totesbackend/config"
	"totesbackend/controllers/utilities"
	"totesbackend/services"

	"github.com/gin-gonic/gin"
)

type UserStateTypeController struct {
	Service *services.UserStateTypeService
	Auth    *utilities.AuthorizationUtil
	Log     *utilities.LogUtil
}

func NewUserStateTypeController(service *services.UserStateTypeService, auth *utilities.AuthorizationUtil, log *utilities.LogUtil) *UserStateTypeController {
	return &UserStateTypeController{Service: service, Auth: auth, Log: log}
}

// GetUserStateTypeByID godoc
// @Summary      Get user state type by ID
// @Description  Retrieves a user state type by its unique ID.
// @Tags         user_state_types
// @Accept       json
// @Produce      json
// @Param        id      path     string  true  "User State Type ID"
// @Success      200     {object}  models.UserStateType  "User State Type details"
// @Failure      403     {object}  models.ErrorResponse  "Permission denied"
// @Failure      404     {object}  models.ErrorResponse  "User State Type not found"
// @Failure      500     {object}  models.ErrorResponse  "Internal server error"
// @Security     ApiKeyAuth
// @Router       /user-state-types/{id} [get]
func (ustc *UserStateTypeController) GetUserStateTypeByID(c *gin.Context) {
	permissionId := config.PERMISSION_GET_USER_STATE_TYPE_BY_ID

	if !ustc.Auth.CheckPermission(c, permissionId) {
		_ = ustc.Log.RegisterLog(c, "Access denied for GetUserStateTypeByID")
		return
	}

	id := c.Param("id")

	if ustc.Log.RegisterLog(c, "Attempting to retrieve user state type with ID: "+id) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
		return
	}

	userStateType, err := ustc.Service.GetUserStateTypeByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User State Type not found"})
		return
	}

	_ = ustc.Log.RegisterLog(c, "Successfully retrieved user state type with ID: "+id)

	c.JSON(http.StatusOK, userStateType)
}

// GetAllUserStateTypes godoc
// @Summary      Get all user state types
// @Description  Retrieves a list of all user state types available in the system.
// @Tags         user_state_types
// @Accept       json
// @Produce      json
// @Success      200     {array}   models.UserStateType  "List of User State Types"
// @Failure      403     {object}  models.ErrorResponse  "Permission denied"
// @Failure      500     {object}  models.ErrorResponse  "Internal server error"
// @Security     ApiKeyAuth
// @Router       /user-state-types [get]
func (ustc *UserStateTypeController) GetAllUserStateTypes(c *gin.Context) {
	permissionId := config.PERMISSION_GET_ALL_USER_STATE_TYPES

	if !ustc.Auth.CheckPermission(c, permissionId) {
		_ = ustc.Log.RegisterLog(c, "Access denied for GetAllUserStateTypes")
		return
	}

	if ustc.Log.RegisterLog(c, "Attempting to retrieve all user state types") != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
		return
	}

	userStateTypes, err := ustc.Service.GetAllUserStateTypes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving User State Types"})
		return
	}

	_ = ustc.Log.RegisterLog(c, "Successfully retrieved all user state types")

	c.JSON(http.StatusOK, userStateTypes)
}
