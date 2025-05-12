package controllers

import (
	"net/http"

	"totesbackend/config"
	"totesbackend/controllers/utilities"
	"totesbackend/models"
	"totesbackend/services"

	"github.com/gin-gonic/gin"
)

type DiscountTypeController struct {
	Service *services.DiscountTypeService
	Auth    *utilities.AuthorizationUtil
	Log     *utilities.LogUtil //
}

func NewDiscountTypeController(service *services.DiscountTypeService, auth *utilities.AuthorizationUtil, log *utilities.LogUtil) *DiscountTypeController {
	return &DiscountTypeController{Service: service, Auth: auth, Log: log}
}

// GetDiscountTypeByID godoc
// @Summary      Get a discount type by ID
// @Description  Retrieves a specific discount type based on its ID. Requires appropriate permission.
// @Tags         discount-types
// @Accept       json
// @Produce      json
// @Param        id  path     string              true  "Discount Type ID"
// @Success      200 {object} models.DiscountType "The requested discount type"
// @Failure      400 {object} models.ErrorResponse "Invalid ID format"
// @Failure      401 {object} models.ErrorResponse "Unauthorized or permission denied"
// @Failure      404 {object} models.ErrorResponse "Discount type not found"
// @Failure      500 {object} models.ErrorResponse "Internal server error or failure in retrieving discount type"
// @Security     ApiKeyAuth
// @Router       /discount-types/{id} [get]
func (dtc *DiscountTypeController) GetDiscountTypeByID(c *gin.Context) {
	id := c.Param("id")

	if dtc.Log.RegisterLog(c, "Attempting to retrieve discount type with ID: "+id) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
		return
	}

	permissionId := config.PERMISSION_GET_DISCOUNT_TYPE_BY_ID
	if !dtc.Auth.CheckPermission(c, permissionId) {
		_ = dtc.Log.RegisterLog(c, "Access denied for GetDiscountTypeByID")
		return
	}

	discountType, err := dtc.Service.GetDiscountTypeByID(id)
	if err != nil {
		_ = dtc.Log.RegisterLog(c, "Discount Type with ID "+id+" not found: "+err.Error())
		c.JSON(http.StatusNotFound, gin.H{"error": "Discount Type not found"})
		return
	}

	_ = dtc.Log.RegisterLog(c, "Successfully retrieved discount type with ID: "+id)
	c.JSON(http.StatusOK, discountType)
}

// GetAllDiscountTypes godoc
// @Summary      Get all discount types
// @Description  Retrieves all available discount types. Requires appropriate permission.
// @Tags         discount-types
// @Accept       json
// @Produce      json
// @Success      200 {array} models.DiscountType "List of all discount types"
// @Failure      401 {object} models.ErrorResponse "Unauthorized or permission denied"
// @Failure      500 {object} models.ErrorResponse "Internal server error or failure in retrieving discount types"
// @Security     ApiKeyAuth
// @Router       /discount-types [get]
func (dtc *DiscountTypeController) GetAllDiscountTypes(c *gin.Context) {
	if dtc.Log.RegisterLog(c, "Attempting to retrieve all discount types") != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
		return
	}

	permissionId := config.PERMISSION_GET_ALL_DISCOUNT_TYPES
	if !dtc.Auth.CheckPermission(c, permissionId) {
		_ = dtc.Log.RegisterLog(c, "Access denied for GetAllDiscountTypes")
		return
	}

	discountTypes, err := dtc.Service.GetAllDiscountTypes()
	if err != nil {
		_ = dtc.Log.RegisterLog(c, "Error retrieving discount types: "+err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving Discount Types"})
		return
	}

	_ = dtc.Log.RegisterLog(c, "Successfully retrieved all discount types")
	c.JSON(http.StatusOK, discountTypes)
}

// CreateDiscountType godoc
// @Summary      Create a new discount type
// @Description  Allows the creation of a new discount type in the system. Requires appropriate permissions.
// @Tags         discount-types
// @Accept       json
// @Produce      json
// @Param        discountType body models.DiscountType true "Discount type details"
// @Success      201 {object} models.DiscountType "Successfully created discount type"
// @Failure      400 {object} models.ErrorResponse "Invalid input data"
// @Failure      401 {object} models.ErrorResponse "Unauthorized or permission denied"
// @Failure      500 {object} models.ErrorResponse "Internal server error or failure in creating the discount type"
// @Security     ApiKeyAuth
// @Router       /discount-types [post]
func (dtc *DiscountTypeController) CreateDiscountType(c *gin.Context) {
	if dtc.Log.RegisterLog(c, "Attempting to create a new discount type") != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
		return
	}

	permissionId := config.PERMISSION_CREATE_DISCOUNT_TYPE
	if !dtc.Auth.CheckPermission(c, permissionId) {
		_ = dtc.Log.RegisterLog(c, "Access denied for CreateDiscountType")
		return
	}

	var discount models.DiscountType
	if err := c.ShouldBindJSON(&discount); err != nil {
		_ = dtc.Log.RegisterLog(c, "Invalid input for discount creation: "+err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err := dtc.Service.CreateDiscountType(&discount)
	if err != nil {
		_ = dtc.Log.RegisterLog(c, "Failed to create discount type: "+err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create discount type"})
		return
	}

	_ = dtc.Log.RegisterLog(c, "Successfully created new discount type")
	c.JSON(http.StatusCreated, discount)
}
