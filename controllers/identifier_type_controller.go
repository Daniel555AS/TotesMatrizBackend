package controllers

import (
	"net/http"
	"totesbackend/config"
	"totesbackend/controllers/utilities"
	"totesbackend/services"

	"github.com/gin-gonic/gin"
)

type IdentifierTypeController struct {
	Service *services.IdentifierTypeService
	Auth    *utilities.AuthorizationUtil
	Log     *utilities.LogUtil
}

func NewIdentifierTypeController(service *services.IdentifierTypeService, auth *utilities.AuthorizationUtil, log *utilities.LogUtil) *IdentifierTypeController {
	return &IdentifierTypeController{Service: service, Auth: auth, Log: log}
}

// GetAllIdentifierTypes godoc
// @Summary      Get all identifier types
// @Description  Retrieves a list of all available identifier types.
// @Tags         identifier-types
// @Accept       json
// @Produce      json
// @Success      200 {array} models.IdentifierType "Successfully retrieved identifier types"
// @Failure      500 {object} models.ErrorResponse "Error retrieving identifier types"
// @Failure      403 {object} models.ErrorResponse "Access denied"
// @Security     ApiKeyAuth
// @Router       /identifier-types [get]
func (itc *IdentifierTypeController) GetAllIdentifierTypes(c *gin.Context) {
	permissionId := config.PERMISSION_GET_ALL_IDENTIFIER_TYPES

	if itc.Log.RegisterLog(c, "Attempting to get all identifier types") != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
		return
	}

	if !itc.Auth.CheckPermission(c, permissionId) {
		_ = itc.Log.RegisterLog(c, "Access denied for GetAllIdentifierTypes")
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	identifierTypes, err := itc.Service.GetAllIdentifierTypes()
	if err != nil {
		_ = itc.Log.RegisterLog(c, "Error retrieving identifier types: "+err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving Identifier Types"})
		return
	}

	_ = itc.Log.RegisterLog(c, "Successfully retrieved all identifier types")
	c.JSON(http.StatusOK, identifierTypes)
}

// GetIdentifierTypeByID godoc
// @Summary      Get identifier type by ID
// @Description  Retrieves an identifier type by its ID.
// @Tags         identifier-types
// @Accept       json
// @Produce      json
// @Param        id  path      string  true  "Identifier Type ID"
// @Success      200 {object} models.IdentifierType "Successfully retrieved identifier type"
// @Failure      404 {object} models.ErrorResponse "Identifier Type not found"
// @Failure      403 {object} models.ErrorResponse "Access denied"
// @Security     ApiKeyAuth
// @Router       /identifier-types/{id} [get]
func (itc *IdentifierTypeController) GetIdentifierTypeByID(c *gin.Context) {
	permissionId := config.PERMISSION_GET_IDENTIFIER_TYPE_BY_ID

	if itc.Log.RegisterLog(c, "Attempting to get identifier type by ID") != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
		return
	}

	if !itc.Auth.CheckPermission(c, permissionId) {
		_ = itc.Log.RegisterLog(c, "Access denied for GetIdentifierTypeByID")
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	id := c.Param("id")

	identifierType, err := itc.Service.GetIdentifierTypeByID(id)
	if err != nil {
		_ = itc.Log.RegisterLog(c, "Identifier type not found with ID: "+id)
		c.JSON(http.StatusNotFound, gin.H{"error": "Identifier Type not found"})
		return
	}

	_ = itc.Log.RegisterLog(c, "Successfully retrieved identifier type with ID: "+id)
	c.JSON(http.StatusOK, identifierType)
}
