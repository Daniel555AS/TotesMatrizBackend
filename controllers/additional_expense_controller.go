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
	"gorm.io/gorm"
)

type AdditionalExpenseController struct {
	Service *services.AdditionalExpenseService
	Auth    *utilities.AuthorizationUtil
	Log     *utilities.LogUtil
}

func NewAdditionalExpenseController(service *services.AdditionalExpenseService,
	auth *utilities.AuthorizationUtil, log *utilities.LogUtil) *AdditionalExpenseController {
	return &AdditionalExpenseController{Service: service, Auth: auth, Log: log}
}

// GetAdditionalExpenseByID godoc
// @Summary      Get additional expense by ID
// @Description  Retrieves an additional expense record by its ID. Requires permission to view additional expenses by ID.
// @Tags         additional-expenses
// @Accept       json
// @Produce      json
// @Param        id   path      string                   true  "ID of the additional expense"
// @Success      200  {object}  models.AdditionalExpense "The additional expense record"
// @Failure      401  {object}  models.ErrorResponse      "Unauthorized or permission denied"
// @Failure      404  {object}  models.ErrorResponse      "Additional expense not found"
// @Failure      500  {object}  models.ErrorResponse      "Internal server error (log registration or DB error)"
// @Router       /additional-expenses/{id} [get]
func (aec *AdditionalExpenseController) GetAdditionalExpenseByID(c *gin.Context) {
	idParam := c.Param("id")

	if aec.Log.RegisterLog(c, "Attempting to retrieve AdditionalExpense with ID: "+idParam) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
		return
	}

	permissionId := config.PERMISSION_GET_ADDITIONAL_EXPENSE_BY_ID
	if !aec.Auth.CheckPermission(c, permissionId) {
		_ = aec.Log.RegisterLog(c, "Access denied for GetAdditionalExpenseByID")
		return
	}

	additionalExpense, err := aec.Service.GetAdditionalExpenseByID(idParam)
	if err != nil {
		_ = aec.Log.RegisterLog(c, "Error retrieving AdditionalExpense with ID "+idParam+": "+err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving Additional Expense"})
		return
	}

	if additionalExpense == nil {
		_ = aec.Log.RegisterLog(c, "AdditionalExpense with ID "+idParam+" not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "Additional Expense not found"})
		return
	}

	_ = aec.Log.RegisterLog(c, "Successfully retrieved AdditionalExpense with ID: "+idParam)

	c.JSON(http.StatusOK, additionalExpense)
}

// GetAllAdditionalExpenses godoc
// @Summary      Get all additional expenses
// @Description  Retrieves all additional expense records. Requires permission to view all additional expenses.
// @Tags         additional-expenses
// @Accept       json
// @Produce      json
// @Success      200  {array}   models.AdditionalExpense   "A list of all additional expenses"
// @Failure      401  {object}  models.ErrorResponse       "Unauthorized or permission denied"
// @Failure      500  {object}  models.ErrorResponse       "Error retrieving additional expenses"
// @Security     ApiKeyAuth
// @Router       /additional-expenses [get]
func (aec *AdditionalExpenseController) GetAllAdditionalExpenses(c *gin.Context) {
	if aec.Log.RegisterLog(c, "Attempting to retrieve all AdditionalExpenses") != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
		return
	}

	permissionId := config.PERMISSION_GET_ALL_ADDITIONAL_EXPENSE
	if !aec.Auth.CheckPermission(c, permissionId) {
		_ = aec.Log.RegisterLog(c, "Access denied for GetAllAdditionalExpenses")
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	additionalExpenses, err := aec.Service.GetAllAdditionalExpenses()
	if err != nil {
		_ = aec.Log.RegisterLog(c, "Error retrieving all AdditionalExpenses: "+err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving additional expenses"})
		return
	}

	_ = aec.Log.RegisterLog(c, "Successfully retrieved all AdditionalExpenses")

	c.JSON(http.StatusOK, additionalExpenses)
}

// CreateAdditionalExpense godoc
// @Summary      Create a new additional expense
// @Description  Creates a new additional expense record. Requires permission to create an additional expense.
// @Tags         additional-expenses
// @Accept       json
// @Produce      json
// @Param        expense  body      dtos.UpdateAdditionalExpenseDTO  true  "Additional Expense DTO"
// @Success      201      {object}  models.AdditionalExpense         "The created additional expense"
// @Failure      400      {object}  models.ErrorResponse             "Invalid JSON format"
// @Failure      401      {object}  models.ErrorResponse             "Unauthorized or permission denied"
// @Failure      500      {object}  models.ErrorResponse             "Error creating additional expense"
// @Security     ApiKeyAuth
// @Router       /additional-expenses [post]
func (aec *AdditionalExpenseController) CreateAdditionalExpense(c *gin.Context) {
	if aec.Log.RegisterLog(c, "Attempting to create a new AdditionalExpense") != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
		return
	}

	permissionId := config.PERMISSION_CREATE_ADDITIONAL_EXPENSE
	if !aec.Auth.CheckPermission(c, permissionId) {
		_ = aec.Log.RegisterLog(c, "Access denied for CreateAdditionalExpense")
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	var dto dtos.UpdateAdditionalExpenseDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		_ = aec.Log.RegisterLog(c, "Invalid JSON format for CreateAdditionalExpense: "+err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	newExpense := &models.AdditionalExpense{
		Name:        dto.Name,
		ItemID:      dto.ItemID,
		Expense:     dto.Expense,
		Description: dto.Description,
	}

	createdExpense, err := aec.Service.CreateAdditionalExpense(newExpense)
	if err != nil {
		_ = aec.Log.RegisterLog(c, "Error creating AdditionalExpense: "+err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating additional expense"})
		return
	}

	_ = aec.Log.RegisterLog(c, "Successfully created AdditionalExpense with ID: "+strconv.Itoa(createdExpense.ID))

	c.JSON(http.StatusCreated, createdExpense)
}

// DeleteAdditionalExpense godoc
// @Summary      Delete an additional expense by ID
// @Description  Deletes an additional expense record by its ID. Requires permission to delete an additional expense.
// @Tags         additional-expenses
// @Accept       json
// @Produce      json
// @Param        id    path      string                  true  "Additional Expense ID"
// @Success      200   {object}  models.MessageResponse       "Message indicating successful deletion"
// @Failure      400   {object}  models.ErrorResponse    "Invalid ID format or request"
// @Failure      401   {object}  models.ErrorResponse    "Unauthorized or permission denied"
// @Failure      404   {object}  models.ErrorResponse    "Additional expense not found"
// @Failure      500   {object}  models.ErrorResponse    "Error deleting additional expense"
// @Security     ApiKeyAuth
// @Router       /additional-expenses/{id} [delete]
func (aec *AdditionalExpenseController) DeleteAdditionalExpense(c *gin.Context) {
	id := c.Param("id")

	if aec.Log.RegisterLog(c, "Attempting to delete AdditionalExpense with ID: "+id) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
		return
	}

	permissionId := config.PERMISSION_DELETE_ADDITIONAL_EXPENSE
	if !aec.Auth.CheckPermission(c, permissionId) {
		_ = aec.Log.RegisterLog(c, "Access denied for DeleteAdditionalExpense with ID: "+id)
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	err := aec.Service.DeleteAdditionalExpense(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			_ = aec.Log.RegisterLog(c, "AdditionalExpense with ID "+id+" not found")
			c.JSON(http.StatusNotFound, gin.H{"error": "Additional Expense not found"})
			return
		}
		_ = aec.Log.RegisterLog(c, "Error deleting AdditionalExpense with ID "+id+": "+err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting Additional Expense"})
		return
	}

	_ = aec.Log.RegisterLog(c, "Successfully deleted AdditionalExpense with ID: "+id)

	c.JSON(http.StatusOK, gin.H{"message": "Additional Expense deleted successfully"})
}

// UpdateAdditionalExpense godoc
// @Summary      Update an additional expense by ID
// @Description  Updates the details of an existing additional expense identified by its ID. Requires permission to update an additional expense.
// @Tags         additional-expenses
// @Accept       json
// @Produce      json
// @Param        id    path      string                            true  "Additional Expense ID"
// @Param        body  body      dtos.UpdateAdditionalExpenseDTO   true  "Updated Additional Expense details"
// @Success      200   {object}  models.AdditionalExpense           "The updated additional expense"
// @Failure      400   {object}  models.ErrorResponse               "Invalid request or JSON format"
// @Failure      401   {object}  models.ErrorResponse               "Unauthorized or permission denied"
// @Failure      404   {object}  models.ErrorResponse               "Additional expense not found"
// @Failure      500   {object}  models.ErrorResponse               "Internal server error"
// @Security     ApiKeyAuth
// @Router       /additional-expenses/{id} [put]
func (aec *AdditionalExpenseController) UpdateAdditionalExpense(c *gin.Context) {
	id := c.Param("id")

	if aec.Log.RegisterLog(c, "Attempting to update AdditionalExpense with ID: "+id) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
		return
	}

	permissionId := config.PERMISSION_UPDATE_ADDITIONAL_EXPENSE
	if !aec.Auth.CheckPermission(c, permissionId) {
		_ = aec.Log.RegisterLog(c, "Access denied for UpdateAdditionalExpense with ID: "+id)
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	var dto dtos.UpdateAdditionalExpenseDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		_ = aec.Log.RegisterLog(c, "Invalid JSON format for UpdateAdditionalExpense with ID: "+id)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	expense, err := aec.Service.GetAdditionalExpenseByID(id)
	if err != nil {
		_ = aec.Log.RegisterLog(c, "AdditionalExpense with ID "+id+" not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "AdditionalExpense not found"})
		return
	}

	expense.Name = dto.Name
	expense.ItemID = dto.ItemID
	expense.Expense = dto.Expense
	expense.Description = dto.Description

	updatedExpense, err := aec.Service.UpdateAdditionalExpense(expense)
	if err != nil {
		_ = aec.Log.RegisterLog(c, "Error updating AdditionalExpense with ID "+id+": "+err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating AdditionalExpense"})
		return
	}

	_ = aec.Log.RegisterLog(c, "Successfully updated AdditionalExpense with ID: "+id)

	c.JSON(http.StatusOK, updatedExpense)
}
