package controllers

import (
	"fmt"
	"net/http"
	"totesbackend/config"
	"totesbackend/controllers/utilities"
	"totesbackend/dtos"
	"totesbackend/services"

	"github.com/gin-gonic/gin"
)

type UserTypeController struct {
	Service *services.UserTypeService
	Auth    *utilities.AuthorizationUtil
	Log     *utilities.LogUtil
}

func NewUserTypeController(service *services.UserTypeService, auth *utilities.AuthorizationUtil, log *utilities.LogUtil) *UserTypeController {
	return &UserTypeController{Service: service, Auth: auth, Log: log}
}

// GetUserTypeByID godoc
// @Summary      Get user type by ID
// @Description  Retrieves a user type based on the provided user type ID.
// @Tags         user_types
// @Accept       json
// @Produce      json
// @Param        id      path     int     true  "User Type ID"
// @Success      200     {object}  dtos.UserTypeDTO  "User type information"
// @Failure      400     {object}  models.ErrorResponse  "Invalid user type ID format"
// @Failure      403     {object}  models.ErrorResponse  "Permission denied"
// @Failure      404     {object}  models.ErrorResponse  "User type not found"
// @Failure      500     {object}  models.ErrorResponse  "Internal server error"
// @Security     ApiKeyAuth
// @Router       /user-types/{id} [get]
func (utc *UserTypeController) GetUserTypeByID(c *gin.Context) {
	permissionId := config.PERMISSION_GET_USER_TYPE_BY_ID

	if !utc.Auth.CheckPermission(c, permissionId) {
		return
	}

	idParam := c.Param("id")

	if utc.Log.RegisterLog(c, "Attempting to retrieve user type with ID: "+idParam) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
		return
	}

	var id uint
	if _, err := fmt.Sscanf(idParam, "%d", &id); err != nil {
		_ = utc.Log.RegisterLog(c, "Invalid user type ID format: "+idParam)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user type ID"})
		return
	}

	userType, err := utc.Service.GetUserTypeByID(id)
	if err != nil {
		_ = utc.Log.RegisterLog(c, "User type not found with ID: "+idParam)
		c.JSON(http.StatusNotFound, gin.H{"error": "User type not found"})
		return
	}

	roleIDs, err := utc.Service.GetRolesForUserType(id)
	if err != nil {
		_ = utc.Log.RegisterLog(c, "Error retrieving roles for user type with ID: "+idParam)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving roles for user type"})
		return
	}

	userTypeDTO := dtos.UserTypeDTO{
		ID:          userType.ID,
		Name:        userType.Name,
		Description: userType.Description,
		Roles:       make([]string, len(roleIDs)),
	}

	for i, roleID := range roleIDs {
		userTypeDTO.Roles[i] = fmt.Sprintf("%d", roleID)
	}

	_ = utc.Log.RegisterLog(c, "Successfully retrieved user type with ID: "+idParam)
	c.JSON(http.StatusOK, userTypeDTO)
}

// GetAllUserTypes godoc
// @Summary      Get all user types
// @Description  Retrieves a list of all user types along with their associated roles.
// @Tags         user_types
// @Accept       json
// @Produce      json
// @Success      200     {array}   dtos.UserTypeDTO  "List of user types"
// @Failure      403     {object}  models.ErrorResponse  "Permission denied"
// @Failure      500     {object}  models.ErrorResponse  "Error retrieving user types"
// @Security     ApiKeyAuth
// @Router       /user-types [get]
func (utc *UserTypeController) GetAllUserTypes(c *gin.Context) {
	permissionId := config.PERMISSION_GET_ALL_USER_TYPES

	if !utc.Auth.CheckPermission(c, permissionId) {
		return
	}

	if utc.Log.RegisterLog(c, "Attempting to retrieve all user types") != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
		return
	}

	userTypes, err := utc.Service.ObtainAllUserTypes()
	if err != nil {
		_ = utc.Log.RegisterLog(c, "Error retrieving all user types")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving user types"})
		return
	}

	var userTypesDTO []dtos.UserTypeDTO
	for _, userType := range userTypes {
		roleIDs, err := utc.Service.GetRolesForUserType(userType.ID)
		if err != nil {
			_ = utc.Log.RegisterLog(c, fmt.Sprintf("Error retrieving roles for user type ID: %d", userType.ID))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving roles for user type"})
			return
		}

		userTypeDTO := dtos.UserTypeDTO{
			ID:          userType.ID,
			Name:        userType.Name,
			Description: userType.Description,
			Roles:       make([]string, len(roleIDs)),
		}

		for i, roleID := range roleIDs {
			userTypeDTO.Roles[i] = fmt.Sprintf("%d", roleID)
		}

		userTypesDTO = append(userTypesDTO, userTypeDTO)
	}

	_ = utc.Log.RegisterLog(c, "Successfully retrieved all user types")
	c.JSON(http.StatusOK, userTypesDTO)
}

// ExistsUserType godoc
// @Summary      Check if a user type exists
// @Description  Checks if a user type exists based on the provided user type ID.
// @Tags         user_types
// @Accept       json
// @Produce      json
// @Param        id      path     string  true  "User Type ID"
// @Success      200     {object}  models.MessageResponse "Existence status of the user type"
// @Failure      400     {object}  models.ErrorResponse  "Invalid user type ID"
// @Failure      403     {object}  models.ErrorResponse  "Permission denied"
// @Failure      500     {object}  models.ErrorResponse  "Error checking user type existence"
// @Security     ApiKeyAuth
// @Router       /user-types/{id}/exists [get]
func (utc *UserTypeController) ExistsUserType(c *gin.Context) {
	permissionId := config.PERMISSION_EXIST_USER_TYPE

	if !utc.Auth.CheckPermission(c, permissionId) {
		return
	}

	idParam := c.Param("id")

	if utc.Log.RegisterLog(c, "Attempting to check existence of user type with ID: "+idParam) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
		return
	}

	var id uint
	if _, err := fmt.Sscanf(idParam, "%d", &id); err != nil {
		_ = utc.Log.RegisterLog(c, "Invalid user type ID format: "+idParam)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user type ID"})
		return
	}

	exists, err := utc.Service.Exists(id)
	if err != nil {
		_ = utc.Log.RegisterLog(c, "Error checking existence for user type ID: "+idParam)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking user type existence"})
		return
	}

	_ = utc.Log.RegisterLog(c, "Checked existence of user type ID: "+idParam+", exists: "+fmt.Sprintf("%v", exists))
	c.JSON(http.StatusOK, gin.H{"exists": exists})
}

// SearchUserTypesByID godoc
// @Summary      Search user types by ID
// @Description  Searches for user types based on the provided user type ID query.
// @Tags         user_types
// @Accept       json
// @Produce      json
// @Param        id      query    string  true  "User Type ID Query"
// @Success      200     {array}   dtos.UserTypeDTO  "List of user types matching the ID query"
// @Failure      400     {object}  models.ErrorResponse  "Invalid query parameter"
// @Failure      403     {object}  models.ErrorResponse  "Permission denied"
// @Failure      500     {object}  models.ErrorResponse  "Error searching user types"
// @Security     ApiKeyAuth
// @Router       /user-types/searchByID [get]
func (utc *UserTypeController) SearchUserTypesByID(c *gin.Context) {
	permissionId := config.PERMISSION_SEARCH_USER_TYPES_BY_ID

	if !utc.Auth.CheckPermission(c, permissionId) {
		return
	}

	query := c.Query("id")

	if utc.Log.RegisterLog(c, "Attempting to search user types by ID: "+query) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
		return
	}

	userTypes, err := utc.Service.SearchUserTypesByID(query)
	if err != nil {
		_ = utc.Log.RegisterLog(c, "Error retrieving user types by ID query: "+query)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving user types"})
		return
	}

	var userTypesDTO []dtos.UserTypeDTO
	for _, userType := range userTypes {
		roleIDs, _ := utc.Service.GetRolesForUserType(userType.ID)

		userTypeDTO := dtos.UserTypeDTO{
			ID:          userType.ID,
			Name:        userType.Name,
			Description: userType.Description,
			Roles:       make([]string, len(roleIDs)),
		}

		for i, roleID := range roleIDs {
			userTypeDTO.Roles[i] = fmt.Sprintf("%d", roleID)
		}

		userTypesDTO = append(userTypesDTO, userTypeDTO)
	}

	_ = utc.Log.RegisterLog(c, "Successfully searched user types by ID query: "+query)
	c.JSON(http.StatusOK, userTypesDTO)
}

// SearchUserTypesByName godoc
// @Summary      Search user types by name
// @Description  Searches for user types based on the provided user type name query.
// @Tags         user_types
// @Accept       json
// @Produce      json
// @Param        name    query    string  true  "User Type Name Query"
// @Success      200     {array}   dtos.UserTypeDTO  "List of user types matching the name query"
// @Failure      400     {object}  models.ErrorResponse  "Invalid query parameter"
// @Failure      403     {object}  models.ErrorResponse  "Permission denied"
// @Failure      500     {object}  models.ErrorResponse  "Error searching user types"
// @Security     ApiKeyAuth
// @Router       /user-types/searchByName [get]
func (utc *UserTypeController) SearchUserTypesByName(c *gin.Context) {
	permissionId := config.PERMISSION_SEARCH_USER_TYPES_BY_NAME

	if !utc.Auth.CheckPermission(c, permissionId) {
		return
	}

	query := c.Query("name")

	if utc.Log.RegisterLog(c, "Attempting to search user types by name: "+query) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
		return
	}

	userTypes, err := utc.Service.SearchUserTypesByName(query)
	if err != nil {
		_ = utc.Log.RegisterLog(c, "Error retrieving user types by name query: "+query)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving user types"})
		return
	}

	var userTypesDTO []dtos.UserTypeDTO
	for _, userType := range userTypes {
		roleIDs, _ := utc.Service.GetRolesForUserType(userType.ID)

		userTypeDTO := dtos.UserTypeDTO{
			ID:          userType.ID,
			Name:        userType.Name,
			Description: userType.Description,
			Roles:       make([]string, len(roleIDs)),
		}

		for i, roleID := range roleIDs {
			userTypeDTO.Roles[i] = fmt.Sprintf("%d", roleID)
		}

		userTypesDTO = append(userTypesDTO, userTypeDTO)
	}

	_ = utc.Log.RegisterLog(c, "Successfully searched user types by name query: "+query)
	c.JSON(http.StatusOK, userTypesDTO)
}
