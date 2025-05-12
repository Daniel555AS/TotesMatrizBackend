package controllers

import (
	"net/http"
	"totesbackend/config"
	"totesbackend/controllers/utilities"
	"totesbackend/services"

	"github.com/gin-gonic/gin"
)

type OrderStateTypeController struct {
	Service *services.OrderStateTypeService
	Auth    *utilities.AuthorizationUtil
	Log     *utilities.LogUtil
}

func NewOrderStateTypeController(service *services.OrderStateTypeService, auth *utilities.AuthorizationUtil, log *utilities.LogUtil) *OrderStateTypeController {
	return &OrderStateTypeController{Service: service, Auth: auth, Log: log}
}

// GetOrderStateTypeByID godoc
// @Summary      Get order state type by ID
// @Description  Retrieves a specific order state type using its ID.
// @Tags         order-state-types
// @Produce      json
// @Param        id   path      int                         true  "Order State Type ID"
// @Success      200  {object}  models.OrderStateType       "Order state type retrieved successfully"
// @Failure      403  {object}  models.ErrorResponse        "Access denied"
// @Failure      404  {object}  models.ErrorResponse        "Order state type not found"
// @Failure      500  {object}  models.ErrorResponse        "Internal server error"
// @Security     ApiKeyAuth
// @Router       /order-state-types/{id} [get]
func (ostc *OrderStateTypeController) GetOrderStateTypeByID(c *gin.Context) {
	permissionId := config.PERMISSION_GET_ORDER_STATE_TYPE_BY_ID

	if ostc.Log.RegisterLog(c, "Attempting to get order state type by ID") != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
		return
	}

	if !ostc.Auth.CheckPermission(c, permissionId) {
		_ = ostc.Log.RegisterLog(c, "Access denied for GetOrderStateTypeByID")
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	id := c.Param("id")

	orderStateType, err := ostc.Service.GetOrderStateTypeByID(id)
	if err != nil {
		_ = ostc.Log.RegisterLog(c, "Order state type not found with ID: "+id)
		c.JSON(http.StatusNotFound, gin.H{"error": "Order State Type not found"})
		return
	}

	_ = ostc.Log.RegisterLog(c, "Successfully retrieved order state type with ID: "+id)
	c.JSON(http.StatusOK, orderStateType)
}

// GetAllOrderStateTypes godoc
// @Summary      Get all order state types
// @Description  Retrieves a list of all available order state types.
// @Tags         order-state-types
// @Produce      json
// @Success      200  {array}   models.OrderStateType       "List of order state types"
// @Failure      403  {object}  models.ErrorResponse        "Access denied"
// @Failure      500  {object}  models.ErrorResponse        "Internal server error"
// @Security     ApiKeyAuth
// @Router       /order-state-types [get]
func (ostc *OrderStateTypeController) GetAllOrderStateTypes(c *gin.Context) {
	permissionId := config.PERMISSION_GET_ALL_ORDER_STATE_TYPES

	if ostc.Log.RegisterLog(c, "Attempting to get all order state types") != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
		return
	}

	if !ostc.Auth.CheckPermission(c, permissionId) {
		_ = ostc.Log.RegisterLog(c, "Access denied for GetAllOrderStateTypes")
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	orderStateTypes, err := ostc.Service.GetAllOrderStateTypes()
	if err != nil {
		_ = ostc.Log.RegisterLog(c, "Error retrieving order state types: "+err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving Order State Types"})
		return
	}

	_ = ostc.Log.RegisterLog(c, "Successfully retrieved all order state types")
	c.JSON(http.StatusOK, orderStateTypes)
}
