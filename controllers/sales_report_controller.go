package controllers

import (
	"net/http"
	"time"
	"totesbackend/config"
	"totesbackend/controllers/utilities"
	"totesbackend/dtos"
	"totesbackend/models"
	"totesbackend/services"

	"github.com/gin-gonic/gin"
)

type SalesReportController struct {
	Service *services.SalesReportService
	Auth    *utilities.AuthorizationUtil
	Log     *utilities.LogUtil
}

func NewSalesReportController(service *services.SalesReportService, auth *utilities.AuthorizationUtil, log *utilities.LogUtil) *SalesReportController {
	return &SalesReportController{Service: service, Auth: auth, Log: log}
}

// GetInvoicesBetweenDates godoc
// @Summary      Fetch invoices between specified dates
// @Description  Retrieve invoices that were generated between the given start and end dates.
// @Tags         sales-report
// @Produce      json
// @Param        startDate  query  string  true  "Start Date (RFC3339 format)"
// @Param        endDate    query  string  true  "End Date (RFC3339 format)"
// @Success      200  {array}  dtos.SalesReportInvoiceDTO  "List of invoices between the given dates"
// @Failure      400  {object}  models.ErrorResponse  "Invalid date format"
// @Failure      403  {object}  models.ErrorResponse  "Permission denied"
// @Failure      500  {object}  models.ErrorResponse  "Error fetching invoices"
// @Security     ApiKeyAuth
// @Router       /sales-report/invoices [get]
func (src *SalesReportController) GetInvoicesBetweenDates(c *gin.Context) {
	startDateStr := c.Query("startDate")
	endDateStr := c.Query("endDate")

	if src.Log.RegisterLog(c, "Request to fetch invoices between "+startDateStr+" and "+endDateStr) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
		return
	}

	permissionId := config.PERMISSION_VIEW_SALES_REPORT
	if !src.Auth.CheckPermission(c, permissionId) {
		_ = src.Log.RegisterLog(c, "Access denied for GetInvoicesBetweenDates")
		return
	}

	startDate, err := time.Parse(time.RFC3339, startDateStr)
	if err != nil {
		_ = src.Log.RegisterLog(c, "Invalid startDate: "+startDateStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid startDate format. Use RFC3339 format: yyyy-mm-ddTHH:MM:SSZ"})
		return
	}

	endDate, err := time.Parse(time.RFC3339, endDateStr)
	if err != nil {
		_ = src.Log.RegisterLog(c, "Invalid endDate: "+endDateStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid endDate format. Use RFC3339 format: yyyy-mm-ddTHH:MM:SSZ"})
		return
	}

	invoices, err := src.Service.GetInvoicesBetweenDates(startDate, endDate)
	if err != nil {
		_ = src.Log.RegisterLog(c, "Error fetching invoices: "+err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching invoices"})
		return
	}

	// Mapeo de []models.Invoice a []dtos.SalesReportInvoiceDTO
	var invoiceDTOs []dtos.SalesReportInvoiceDTO
	for _, invoice := range invoices {
		invoiceDTO := mapInvoiceToSalesReportDTO(invoice)
		invoiceDTOs = append(invoiceDTOs, invoiceDTO)
	}

	_ = src.Log.RegisterLog(c, "Successfully fetched invoices between "+startDateStr+" and "+endDateStr)
	c.JSON(http.StatusOK, invoiceDTOs)
}

// Función para mapear un Invoice a SalesReportInvoiceDTO
func mapInvoiceToSalesReportDTO(invoice models.Invoice) dtos.SalesReportInvoiceDTO {
	// Convertir los items
	var billingItems []dtos.BillingItemDTO
	for _, item := range invoice.Items {
		billingItems = append(billingItems, dtos.BillingItemDTO{
			ID:    item.ItemID,
			Stock: item.Amount,
			// Puedes agregar más campos si tu BillingItemDTO tiene más propiedades
		})
	}

	return dtos.SalesReportInvoiceDTO{
		ID:        invoice.ID,
		DateTime:  invoice.DateTime,
		Total:     invoice.Total,
		Subtotal:  invoice.Subtotal,
		Items:     billingItems,
		Discounts: invoice.Discounts,
		Taxes:     invoice.Taxes,
	}
}
