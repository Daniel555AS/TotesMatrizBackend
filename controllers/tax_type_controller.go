package controllers

import (
	"net/http"

	"totesbackend/config"
	"totesbackend/controllers/utilities"
	"totesbackend/models"
	"totesbackend/services"

	"github.com/gin-gonic/gin"
)

type TaxTypeController struct {
	Service *services.TaxTypeService
	Auth    *utilities.AuthorizationUtil
	Log     *utilities.LogUtil
}

func NewTaxTypeController(service *services.TaxTypeService,
	auth *utilities.AuthorizationUtil, log *utilities.LogUtil) *TaxTypeController {
	return &TaxTypeController{Service: service, Auth: auth, Log: log}
}

// GetTaxTypeByID godoc
// @Summary      Retrieve a tax type by its ID
// @Description  Fetches the details of a tax type specified by the provided ID.
// @Tags         tax-types
// @Produce      json
// @Param        id  path  string  true  "Tax Type ID"
// @Success      200  {object}  models.TaxType  "Details of the tax type"
// @Failure      400  {object}  models.ErrorResponse  "Invalid ID format"
// @Failure      403  {object}  models.ErrorResponse  "Permission denied"
// @Failure      404  {object}  models.ErrorResponse  "Tax type not found"
// @Security     ApiKeyAuth
// @Router       /tax-types/{id} [get]
func (ttc *TaxTypeController) GetTaxTypeByID(c *gin.Context) {
	permissionId := config.PERMISSION_GET_TAX_TYPE_BY_ID

	if !ttc.Auth.CheckPermission(c, permissionId) {
		return
	}

	id := c.Param("id")
	taxType, err := ttc.Service.GetTaxTypeByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tax Type not found"})
		return
	}

	c.JSON(http.StatusOK, taxType)
}

// GetAllTaxTypes godoc
// @Summary      Retrieve all tax types
// @Description  Fetches the list of all available tax types.
// @Tags         tax-types
// @Produce      json
// @Success      200  {array}   models.TaxType  "List of tax types"
// @Failure      403  {object}  models.ErrorResponse  "Permission denied"
// @Failure      500  {object}  models.ErrorResponse  "Error retrieving tax types"
// @Security     ApiKeyAuth
// @Router       /tax-types [get]
func (ttc *TaxTypeController) GetAllTaxTypes(c *gin.Context) {
	permissionId := config.PERMISSION_GET_ALL_TAX_TYPES

	if !ttc.Auth.CheckPermission(c, permissionId) {
		return
	}

	taxTypes, err := ttc.Service.GetAllTaxTypes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving Tax Types"})
		return
	}
	c.JSON(http.StatusOK, taxTypes)
}

// CreateTaxType godoc
// @Summary      Create a new tax type
// @Description  Creates a new tax type by providing its details (name, percentage, etc).
// @Tags         tax-types
// @Accept       json
// @Produce      json
// @Param        tax  body     models.TaxType  true  "Tax Type Details"
// @Success      201  {object}  models.TaxType  "Successfully created tax type"
// @Failure      400  {object}  models.ErrorResponse  "Invalid input data"
// @Failure      403  {object}  models.ErrorResponse  "Permission denied"
// @Failure      500  {object}  models.ErrorResponse  "Error creating tax type"
// @Security     ApiKeyAuth
// @Router       /tax-types [post]
func (ttc *TaxTypeController) CreateTaxType(c *gin.Context) {
	if ttc.Log.RegisterLog(c, "Attempting to create a new tax type") != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
		return
	}

	permissionId := config.PERMISSION_CREATE_TAX_TYPE
	if !ttc.Auth.CheckPermission(c, permissionId) {
		_ = ttc.Log.RegisterLog(c, "Access denied for CreateTaxType")
		return
	}

	var tax models.TaxType
	if err := c.ShouldBindJSON(&tax); err != nil {
		_ = ttc.Log.RegisterLog(c, "Invalid input for tax type creation: "+err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos del impuesto"})
		return
	}

	err := ttc.Service.CreateTaxType(&tax)
	if err != nil {
		_ = ttc.Log.RegisterLog(c, "Failed to create tax type: "+err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear el impuesto"})
		return
	}

	_ = ttc.Log.RegisterLog(c, "Successfully created new tax type")
	c.JSON(http.StatusCreated, tax)
}
