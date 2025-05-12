package controllers

import (
	"net/http"
	"strconv"
	"totesbackend/config"
	"totesbackend/controllers/utilities"
	"totesbackend/dtos"
	"totesbackend/models"
	"totesbackend/services"

	"github.com/gin-gonic/gin"
)

type InvoiceController struct {
	Service *services.InvoiceService
	Auth    *utilities.AuthorizationUtil
	Log     *utilities.LogUtil //
}

func NewInvoiceController(
	service *services.InvoiceService, auth *utilities.AuthorizationUtil, log *utilities.LogUtil) *InvoiceController {
	return &InvoiceController{
		Service: service,
		Auth:    auth,
		Log:     log,
	}
}

// GetAllInvoices godoc
// @Summary      Get all invoices
// @Description  Retrieves all invoices from the system.
// @Tags         invoices
// @Accept       json
// @Produce      json
// @Success      200 {array} dtos.GetInvoiceDTO "List of all invoices"
// @Failure      404 {object} models.ErrorResponse "No invoices found"
// @Failure      403 {object} models.ErrorResponse "Access denied"
// @Security     ApiKeyAuth
// @Router       /invoices [get]
func (ic *InvoiceController) GetAllInvoices(c *gin.Context) {
	if ic.Log.RegisterLog(c, "Attempting to retrieve all invoices") != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
		return
	}

	permissionId := config.PERMISSION_GET_ALL_INVOICES
	if !ic.Auth.CheckPermission(c, permissionId) {
		_ = ic.Log.RegisterLog(c, "Access denied for GetAllInvoices")
		return
	}

	invoices, err := ic.Service.GetAllInvoices()
	if err != nil {
		_ = ic.Log.RegisterLog(c, "Error retrieving invoices: "+err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve invoices"})
		return
	}

	if len(invoices) == 0 {
		_ = ic.Log.RegisterLog(c, "No invoices found")
		c.JSON(http.StatusNotFound, gin.H{"error": "No invoices found"})
		return
	}

	var invoiceDTOs []dtos.GetInvoiceDTO
	for _, invoice := range invoices {
		invoiceDTOs = append(invoiceDTOs, dtos.GetInvoiceDTO{
			ID:             invoice.ID,
			EnterpriseData: invoice.EnterpriseData,
			DateTime:       invoice.DateTime,
			CustomerID:     invoice.CustomerID,
			Subtotal:       invoice.Subtotal,
			Total:          invoice.Total,
			Items:          extractInvoiceBillingItems(invoice.Items),
			Discounts:      extractDiscountIds(invoice.Discounts),
			Taxes:          extractTaxIds(invoice.Taxes),
		})
	}

	_ = ic.Log.RegisterLog(c, "Successfully retrieved all invoices")
	c.JSON(http.StatusOK, invoiceDTOs)
}

// GetInvoiceByID godoc
// @Summary      Get invoice by ID
// @Description  Retrieves an invoice by its unique ID.
// @Tags         invoices
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Invoice ID"
// @Success      200 {object} dtos.GetInvoiceDTO "Invoice details"
// @Failure      400 {object} models.ErrorResponse "Invalid invoice ID"
// @Failure      404 {object} models.ErrorResponse "Invoice not found"
// @Failure      403 {object} models.ErrorResponse "Access denied"
// @Security     ApiKeyAuth
// @Router       /invoices/{id} [get]
func (ic *InvoiceController) GetInvoiceByID(c *gin.Context) {
	idParam := c.Param("id")
	if ic.Log.RegisterLog(c, "Attempting to retrieve invoice with ID: "+idParam) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
		return
	}

	permissionId := config.PERMISSION_GET_INVOICE_BY_ID
	if !ic.Auth.CheckPermission(c, permissionId) {
		_ = ic.Log.RegisterLog(c, "Access denied for GetInvoiceByID")
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		_ = ic.Log.RegisterLog(c, "Invalid invoice ID: "+idParam)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid invoice ID"})
		return
	}

	invoice, err := ic.Service.GetInvoiceByID(strconv.Itoa(id))
	if err != nil {
		_ = ic.Log.RegisterLog(c, "Invoice not found with ID: "+idParam)
		c.JSON(http.StatusNotFound, gin.H{"error": "Invoice not found"})
		return
	}

	invoiceDTO := dtos.GetInvoiceDTO{
		ID:             invoice.ID,
		EnterpriseData: invoice.EnterpriseData,
		DateTime:       invoice.DateTime,
		CustomerID:     invoice.CustomerID,
		Subtotal:       invoice.Subtotal,
		Total:          invoice.Total,
		Items:          extractInvoiceBillingItems(invoice.Items),
		Discounts:      extractDiscountIds(invoice.Discounts),
		Taxes:          extractTaxIds(invoice.Taxes),
	}

	_ = ic.Log.RegisterLog(c, "Successfully retrieved invoice with ID: "+idParam)
	c.JSON(http.StatusOK, invoiceDTO)
}

// SearchInvoiceByID godoc
// @Summary      Search invoices by ID
// @Description  Search for invoices using a query parameter for the invoice ID.
// @Tags         invoices
// @Accept       json
// @Produce      json
// @Param        id   query     string  true  "Invoice ID Query"
// @Success      200 {array} dtos.GetInvoiceDTO "List of invoices found"
// @Failure      400 {object} models.ErrorResponse "Query parameter is required"
// @Failure      404 {object} models.ErrorResponse "No invoices found"
// @Failure      403 {object} models.ErrorResponse "Access denied"
// @Security     ApiKeyAuth
// @Router       /invoices/searchById [get]
func (ic *InvoiceController) SearchInvoiceByID(c *gin.Context) {
	query := c.Query("id")

	if ic.Log.RegisterLog(c, "Attempting to search invoice(s) by ID query: "+query) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
		return
	}

	permissionId := config.PERMISSION_SEARCH_INVOICE_BY_ID
	if !ic.Auth.CheckPermission(c, permissionId) {
		_ = ic.Log.RegisterLog(c, "Access denied for SearchInvoiceByID")
		return
	}

	if query == "" {
		_ = ic.Log.RegisterLog(c, "Missing query parameter for SearchInvoiceByID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter is required"})
		return
	}

	invoices, err := ic.Service.SearchInvoiceByID(query)
	if err != nil {
		_ = ic.Log.RegisterLog(c, "Error searching invoices by ID query "+query+": "+err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error searching invoices"})
		return
	}

	var invoiceDTOs []dtos.GetInvoiceDTO
	for _, invoice := range invoices {
		invoiceDTOs = append(invoiceDTOs, dtos.GetInvoiceDTO{
			ID:             invoice.ID,
			EnterpriseData: invoice.EnterpriseData,
			DateTime:       invoice.DateTime,
			CustomerID:     invoice.CustomerID,
			Subtotal:       invoice.Subtotal,
			Total:          invoice.Total,
			Items:          extractInvoiceBillingItems(invoice.Items),
			Discounts:      extractDiscountIds(invoice.Discounts),
			Taxes:          extractTaxIds(invoice.Taxes),
		})
	}

	_ = ic.Log.RegisterLog(c, "Successfully retrieved "+strconv.Itoa(len(invoiceDTOs))+" invoice(s) for search ID: "+query)
	c.JSON(http.StatusOK, invoiceDTOs)
}

// SearchInvoiceByCustomerPersonalId godoc
// @Summary      Search invoices by customer personal ID
// @Description  Search for invoices using a query parameter for the customer's personal ID.
// @Tags         invoices
// @Accept       json
// @Produce      json
// @Param        personal_id   query     string  true  "Customer Personal ID Query"
// @Success      200 {array} dtos.GetInvoiceDTO "List of invoices found"
// @Failure      400 {object} models.ErrorResponse "Query parameter 'personal_id' is required"
// @Failure      404 {object} models.ErrorResponse "No invoices found"
// @Failure      403 {object} models.ErrorResponse "Access denied"
// @Security     ApiKeyAuth
// @Router       /invoices/searchByPersonalId [get]
func (ic *InvoiceController) SearchInvoiceByCustomerPersonalId(c *gin.Context) {
	query := c.Query("personal_id")

	if ic.Log.RegisterLog(c, "Attempting to search invoice(s) by customer personal ID: "+query) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
		return
	}

	permissionId := config.PERMISSION_SEARCH_INVOICE_BY_CUSTOMER_PERSONAL_ID
	if !ic.Auth.CheckPermission(c, permissionId) {
		_ = ic.Log.RegisterLog(c, "Access denied for SearchInvoiceByCustomerPersonalId")
		return
	}

	if query == "" {
		_ = ic.Log.RegisterLog(c, "Missing query parameter 'personal_id' for SearchInvoiceByCustomerPersonalId")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'personal_id' is required"})
		return
	}

	invoices, err := ic.Service.SearchInvoiceByCustomerPersonalId(query)
	if err != nil {
		_ = ic.Log.RegisterLog(c, "Error searching invoices by customer personal ID "+query+": "+err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error searching invoices by customer personal ID"})
		return
	}

	var invoiceDTOs []dtos.GetInvoiceDTO
	for _, invoice := range invoices {
		invoiceDTOs = append(invoiceDTOs, dtos.GetInvoiceDTO{
			ID:             invoice.ID,
			EnterpriseData: invoice.EnterpriseData,
			DateTime:       invoice.DateTime,
			CustomerID:     invoice.CustomerID,
			Subtotal:       invoice.Subtotal,
			Total:          invoice.Total,
			Items:          extractInvoiceBillingItems(invoice.Items),
			Discounts:      extractDiscountIds(invoice.Discounts),
			Taxes:          extractTaxIds(invoice.Taxes),
		})
	}

	_ = ic.Log.RegisterLog(c, "Successfully retrieved "+strconv.Itoa(len(invoiceDTOs))+" invoice(s) for customer personal ID: "+query)
	c.JSON(http.StatusOK, invoiceDTOs)
}

// CreateInvoice godoc
// @Summary      Create a new invoice
// @Description  Create a new invoice based on the provided JSON data. Requires appropriate permissions.
// @Tags         invoices
// @Accept       json
// @Produce      json
// @Param        invoice_body  body      dtos.CreateInvoiceDTO  true  "Invoice data"
// @Success      201 {object} dtos.GetInvoiceDTO "Created invoice"
// @Failure      400 {object} models.ErrorResponse "Invalid request data"
// @Failure      403 {object} models.ErrorResponse "Access denied"
// @Failure      500 {object} models.ErrorResponse "Error creating invoice"
// @Security     ApiKeyAuth
// @Router       /invoices [post]
func (ic *InvoiceController) CreateInvoice(c *gin.Context) {
	if ic.Log.RegisterLog(c, "Attempting to create new invoice") != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
		return
	}

	permissionId := config.PERMISSION_CREATE_INVOICE
	if !ic.Auth.CheckPermission(c, permissionId) {
		_ = ic.Log.RegisterLog(c, "Access denied for CreateInvoice")
		return
	}

	var dto dtos.CreateInvoiceDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		_ = ic.Log.RegisterLog(c, "Invalid invoice creation request data: "+err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	invoice, err := ic.Service.CreateInvoice(&dto)
	if err != nil {
		_ = ic.Log.RegisterLog(c, "Error creating invoice: "+err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	invoiceDTO := dtos.GetInvoiceDTO{
		ID:             invoice.ID,
		EnterpriseData: invoice.EnterpriseData,
		DateTime:       invoice.DateTime,
		CustomerID:     invoice.CustomerID,
		Subtotal:       invoice.Subtotal,
		Total:          invoice.Total,
		Items:          extractInvoiceBillingItems(invoice.Items),
		Discounts:      extractDiscountIds(invoice.Discounts),
		Taxes:          extractTaxIds(invoice.Taxes),
	}

	_ = ic.Log.RegisterLog(c, "Successfully created invoice with ID: "+strconv.Itoa(invoice.ID))
	c.JSON(http.StatusCreated, invoiceDTO)
}

func extractInvoiceBillingItems(items []models.InvoiceItem) []dtos.BillingItemDTO {
	var billingItems []dtos.BillingItemDTO
	for _, item := range items {
		billingItems = append(billingItems, dtos.BillingItemDTO{
			ID:    item.ItemID,
			Stock: item.Amount,
		})
	}
	return billingItems
}

func extractDiscountIds(discounts []models.DiscountType) []int {
	var ids []int
	for _, discount := range discounts {
		ids = append(ids, discount.ID)
	}
	return ids
}

func extractTaxIds(taxes []models.TaxType) []int {
	var ids []int
	for _, tax := range taxes {
		ids = append(ids, tax.ID)
	}
	return ids
}
