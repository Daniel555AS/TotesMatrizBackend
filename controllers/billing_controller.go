package controllers

import (
	"net/http"
	"strconv"
	"totesbackend/config"
	"totesbackend/controllers/utilities"
	"totesbackend/dtos"
	"totesbackend/services"

	"github.com/gin-gonic/gin"
)

type BillingController struct {
	Service *services.BillingService
	Auth    *utilities.AuthorizationUtil
}

func NewBillingController(service *services.BillingService, auth *utilities.AuthorizationUtil) *BillingController {
	return &BillingController{Service: service, Auth: auth}
}

type SubtotalResponse struct {
	Subtotal float64 `json:"subtotal"`
}

// CalculateSubtotal godoc
// @Summary      Calculate subtotal
// @Description  Calculates the subtotal based on a list of billing items. Requires permission.
// @Tags         billing
// @Accept       json
// @Produce      json
// @Param        items  body      []dtos.BillingItemDTO  true  "List of billing items"
// @Success      200    {object}  SubtotalResponse       "Calculated subtotal"
// @Failure      400    {object}  models.ErrorResponse    "Invalid request data"
// @Failure      401    {object}  models.ErrorResponse    "Unauthorized or permission denied"
// @Failure      404    {object}  models.ErrorResponse    "Calculation error (e.g., related data not found)"
// @Security     ApiKeyAuth
// @Router       /billing/subtotal [post]
func (bc *BillingController) CalculateSubtotal(c *gin.Context) {
	permissionId := config.PERMISSION_CALCULATE_SUBTOTAL

	if !bc.Auth.CheckPermission(c, permissionId) {
		return
	}

	var itemsDTO []dtos.BillingItemDTO
	if err := c.ShouldBindJSON(&itemsDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	subtotal, err := bc.Service.CalculateSubtotal(itemsDTO)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"subtotal": subtotal})
}

type TotalResponse struct {
	Total float64 `json:"total"`
}

// CalculateTotal godoc
// @Summary      Calculate total
// @Description  Calculates the total amount based on billing items, discounts, and tax types. Requires permission.
// @Tags         billing
// @Accept       json
// @Produce      json
// @Param        body  body  dtos.CalculateTotalRequestDTO  true  "Billing total calculation input"
// @Success      200   {object}  TotalResponse         "Calculated total"
// @Failure      400   {object}  models.ErrorResponse       "Invalid request data"
// @Failure      401   {object}  models.ErrorResponse       "Unauthorized or permission denied"
// @Failure      404   {object}  models.ErrorResponse       "Calculation error (e.g., related data not found)"
// @Security     ApiKeyAuth
// @Router       /billing/total [post]
func (bc *BillingController) CalculateTotal(c *gin.Context) {
	permissionId := config.PERMISSION_CALCULATE_TOTAL

	if !bc.Auth.CheckPermission(c, permissionId) {
		return
	}
	var request dtos.CalculateTotalRequestDTO

	// Estructura del request con arrays de enteros
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	discountTypesIdsStr := make([]string, len(request.DiscountTypesIds))
	for i, id := range request.DiscountTypesIds {
		discountTypesIdsStr[i] = strconv.Itoa(id)
	}

	taxTypesIdsStr := make([]string, len(request.TaxTypesIds))
	for i, id := range request.TaxTypesIds {
		taxTypesIdsStr[i] = strconv.Itoa(id)
	}

	total, err := bc.Service.CalculateTotal(discountTypesIdsStr, taxTypesIdsStr, request.ItemsDTO)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"total": total})
}
