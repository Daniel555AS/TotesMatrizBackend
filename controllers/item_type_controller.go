package controllers

import (
	"net/http"

	"totesbackend/config"
	"totesbackend/controllers/utilities"
	"totesbackend/services"

	"github.com/gin-gonic/gin"
)

type ItemTypeController struct {
	Service *services.ItemTypeService
	Auth    *utilities.AuthorizationUtil
	Log     *utilities.LogUtil
}

func NewItemTypeController(service *services.ItemTypeService, auth *utilities.AuthorizationUtil, log *utilities.LogUtil) *ItemTypeController {
	return &ItemTypeController{Service: service, Auth: auth, Log: log}
}

// GetItemTypeByID godoc
// @Summary      Get item type by ID
// @Description  Retrieves the details of a specific item type by its ID.
// @Tags         item-types
// @Produce      json
// @Param        id   path      string                 true  "Item Type ID"
// @Success      200  {object}  models.ItemType        "Item Type retrieved successfully"
// @Failure      404  {object}  models.ErrorResponse   "Item Type not found"
// @Failure      500  {object}  models.ErrorResponse   "Error registering log"
// @Security     ApiKeyAuth
// @Router       /item-types/{id} [get]
func (itc *ItemTypeController) GetItemTypeByID(c *gin.Context) {
	id := c.Param("id")

	if itc.Log.RegisterLog(c, "Attempting to retrieve ItemType with ID: "+id) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
		return
	}

	permissionId := config.PERMISSION_GET_ITEM_BY_ID
	if !itc.Auth.CheckPermission(c, permissionId) {
		_ = itc.Log.RegisterLog(c, "Access denied for GetItemTypeByID")
		return
	}

	itemType, err := itc.Service.GetItemTypeByID(id)
	if err != nil {
		_ = itc.Log.RegisterLog(c, "Error retrieving ItemType with ID "+id+": "+err.Error())
		c.JSON(http.StatusNotFound, gin.H{"error": "Item Type not found"})
		return
	}

	_ = itc.Log.RegisterLog(c, "Successfully retrieved ItemType with ID: "+id)
	c.JSON(http.StatusOK, itemType)
}

// GetItemTypes godoc
// @Summary      Get all item types
// @Description  Retrieves a list of all item types.
// @Tags         item-types
// @Produce      json
// @Success      200  {array}   models.ItemType         "List of item types retrieved successfully"
// @Failure      500  {object}  models.ErrorResponse    "Error retrieving item types or registering log"
// @Security     ApiKeyAuth
// @Router       /item-types [get]
func (itc *ItemTypeController) GetItemTypes(c *gin.Context) {
	if itc.Log.RegisterLog(c, "Attempting to retrieve all ItemTypes") != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
		return
	}

	permissionId := config.PERMISSION_GET_ITEM_TYPES
	if !itc.Auth.CheckPermission(c, permissionId) {
		_ = itc.Log.RegisterLog(c, "Access denied for GetItemTypes")
		return
	}

	itemTypes, err := itc.Service.GetAllItemTypes()
	if err != nil {
		_ = itc.Log.RegisterLog(c, "Error retrieving ItemTypes: "+err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving Item Types"})
		return
	}

	_ = itc.Log.RegisterLog(c, "Successfully retrieved all ItemTypes")
	c.JSON(http.StatusOK, itemTypes)
}
