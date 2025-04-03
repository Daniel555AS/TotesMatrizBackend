package controllers

import (
	"net/http"
	"totesbackend/config"
	"totesbackend/controllers/utilities"
	"totesbackend/services"

	"github.com/gin-gonic/gin"
)

type HistoricalItemPriceController struct {
	Service *services.HistoricalItemPriceService
	Auth    *utilities.AuthorizationUtil
	Log     *utilities.LogUtil
}

func NewHistoricalItemPriceController(service *services.HistoricalItemPriceService, auth *utilities.AuthorizationUtil, log *utilities.LogUtil) *HistoricalItemPriceController {
	return &HistoricalItemPriceController{
		Service: service,
		Auth:    auth,
		Log:     log,
	}
}

func (c *HistoricalItemPriceController) GetHistoricalItemPrice(ctx *gin.Context) {
	itemID := ctx.Param("id")

	if c.Log.RegisterLog(ctx, "Attempting to retrieve historical item price for item ID: "+itemID) != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
		return
	}

	permissionId := config.PERMISSION_GET_HISTORICAL_ITEM_PRICE
	if !c.Auth.CheckPermission(ctx, permissionId) {
		_ = c.Log.RegisterLog(ctx, "Access denied for GetHistoricalItemPrice")
		return
	}

	historicalPrices, err := c.Service.GetHistoricalItemPrice(itemID)
	if err != nil {
		_ = c.Log.RegisterLog(ctx, "Error retrieving historical prices for item ID "+itemID+": "+err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve historical prices"})
		return
	}

	if len(historicalPrices) == 0 {
		_ = c.Log.RegisterLog(ctx, "No historical prices found for item ID "+itemID)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "No historical prices found"})
		return
	}

	_ = c.Log.RegisterLog(ctx, "Successfully retrieved historical prices for item ID: "+itemID)
	ctx.JSON(http.StatusOK, historicalPrices)
}
