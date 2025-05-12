package controllers

import (
	"errors"
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

type EmployeeController struct {
	Service *services.EmployeeService
	Auth    *utilities.AuthorizationUtil
	Log     *utilities.LogUtil
}

func NewEmployeeController(service *services.EmployeeService, auth *utilities.AuthorizationUtil, log *utilities.LogUtil) *EmployeeController {
	return &EmployeeController{Service: service, Auth: auth, Log: log}
}

// GetEmployeeByID godoc
// @Summary      Get employee by ID
// @Description  Retrieves a specific employee by their unique ID from the system. Requires the appropriate permissions.
// @Tags         employees
// @Accept       json
// @Produce      json
// @Param        id path string true "Employee ID"
// @Success      200 {object} dtos.GetEmployeeDTO "Successfully retrieved employee details"
// @Failure      400 {object} models.ErrorResponse "Invalid employee ID"
// @Failure      403 {object} models.ErrorResponse "Permission denied"
// @Failure      404 {object} models.ErrorResponse "Employee not found"
// @Security     ApiKeyAuth
// @Router       /employees/{id} [get]
func (ec *EmployeeController) GetEmployeeByID(c *gin.Context) {
	permissionId := config.PERMISSION_GET_EMPLOYEE_BY_ID

	if ec.Log.RegisterLog(c, "Attempting to get employee by ID") != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
		return
	}

	if !ec.Auth.CheckPermission(c, permissionId) {
		_ = ec.Log.RegisterLog(c, "Access denied for GetEmployeeByID")
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	id := c.Param("id")

	employee, err := ec.Service.GetEmployeeByID(id)
	if err != nil {
		_ = ec.Log.RegisterLog(c, "Employee not found with ID: "+id)
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}

	employeeDTO := dtos.GetEmployeeDTO{
		ID:               employee.ID,
		Names:            employee.Names,
		LastNames:        employee.LastNames,
		PersonalID:       employee.PersonalID,
		Address:          employee.Address,
		PhoneNumbers:     employee.PhoneNumbers,
		UserID:           employee.UserID,
		IdentifierTypeID: employee.IdentifierTypeID,
	}

	_ = ec.Log.RegisterLog(c, "Successfully retrieved employee with ID: "+id)
	c.JSON(http.StatusOK, employeeDTO)
}

// GetAllEmployees godoc
// @Summary      Get all employees
// @Description  Retrieves a list of all employees in the system. Requires the appropriate permissions.
// @Tags         employees
// @Accept       json
// @Produce      json
// @Success      200 {array} dtos.GetEmployeeDTO "Successfully retrieved list of employees"
// @Failure      403 {object} models.ErrorResponse "Permission denied"
// @Failure      500 {object} models.ErrorResponse "Error retrieving employees"
// @Security     ApiKeyAuth
// @Router       /employees [get]
func (ec *EmployeeController) GetAllEmployees(c *gin.Context) {
	permissionId := config.PERMISSION_GET_ALL_EMPLOYEES

	if ec.Log.RegisterLog(c, "Attempting to get all employees") != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
		return
	}

	if !ec.Auth.CheckPermission(c, permissionId) {
		_ = ec.Log.RegisterLog(c, "Access denied for GetAllEmployees")
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	employees, err := ec.Service.GetAllEmployees()
	if err != nil {
		_ = ec.Log.RegisterLog(c, "Error retrieving employees: "+err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving employees"})
		return
	}

	var employeesDTO []dtos.GetEmployeeDTO
	for _, employee := range employees {
		employeeDTO := dtos.GetEmployeeDTO{
			ID:               employee.ID,
			Names:            employee.Names,
			LastNames:        employee.LastNames,
			PersonalID:       employee.PersonalID,
			Address:          employee.Address,
			PhoneNumbers:     employee.PhoneNumbers,
			UserID:           employee.UserID,
			IdentifierTypeID: employee.IdentifierTypeID,
		}
		employeesDTO = append(employeesDTO, employeeDTO)
	}

	_ = ec.Log.RegisterLog(c, "Successfully retrieved all employees")
	c.JSON(http.StatusOK, employeesDTO)
}

// SearchEmployeesByID godoc
// @Summary      Search employees by ID
// @Description  Searches for employees by a given ID. Returns a list of employees that match the ID.
// @Tags         employees
// @Accept       json
// @Produce      json
// @Param        id query string true "Employee ID to search for"
// @Success      200 {array} dtos.GetEmployeeDTO "Successfully found employees matching ID"
// @Failure      403 {object} models.ErrorResponse "Permission denied"
// @Failure      404 {object} models.ErrorResponse "No employees found"
// @Failure      500 {object} models.ErrorResponse "Error retrieving employees"
// @Security     ApiKeyAuth
// @Router       /employees/searchByID [get]
func (ec *EmployeeController) SearchEmployeesByID(c *gin.Context) {
	query := c.Query("id")
	permissionId := config.PERMISSION_SEARCH_EMPLOYEES_BY_ID

	if ec.Log.RegisterLog(c, "Attempting to search employees by ID: "+query) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
		return
	}

	if !ec.Auth.CheckPermission(c, permissionId) {
		_ = ec.Log.RegisterLog(c, "Access denied for SearchEmployeesByID")
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	employees, err := ec.Service.SearchEmployeesByID(query)
	if err != nil {
		_ = ec.Log.RegisterLog(c, "Error retrieving employees by ID: "+query+" - "+err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving employees"})
		return
	}

	if len(employees) == 0 {
		_ = ec.Log.RegisterLog(c, "No employees found with ID: "+query)
		c.JSON(http.StatusNotFound, gin.H{"message": "No employees found"})
		return
	}

	var employeesDTO []dtos.GetEmployeeDTO
	for _, employee := range employees {
		employeesDTO = append(employeesDTO, dtos.GetEmployeeDTO{
			ID:               employee.ID,
			Names:            employee.Names,
			LastNames:        employee.LastNames,
			PhoneNumbers:     employee.PhoneNumbers,
			UserID:           employee.UserID,
			IdentifierTypeID: employee.IdentifierTypeID,
		})
	}

	_ = ec.Log.RegisterLog(c, "Successfully found employees matching ID: "+query)
	c.JSON(http.StatusOK, employeesDTO)
}

// SearchEmployeesByName godoc
// @Summary      Search employees by name
// @Description  Searches for employees by their name. Returns a list of employees that match the name.
// @Tags         employees
// @Accept       json
// @Produce      json
// @Param        names query string true "Employee name to search for"
// @Success      200 {array} dtos.GetEmployeeDTO "Successfully found employees matching name"
// @Failure      400 {object} models.ErrorResponse "Search query is required"
// @Failure      403 {object} models.ErrorResponse "Permission denied"
// @Failure      404 {object} models.ErrorResponse "No employees found"
// @Failure      500 {object} models.ErrorResponse "Error retrieving employees"
// @Security     ApiKeyAuth
// @Router       /employees/searchByName [get]
func (ec *EmployeeController) SearchEmployeesByName(c *gin.Context) {
	permissionId := config.PERMISSION_SEARCH_EMPLOYEES_BY_NAME

	query := c.Query("names")

	if ec.Log.RegisterLog(c, "Attempting to search employees by name: "+query) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
		return
	}

	if !ec.Auth.CheckPermission(c, permissionId) {
		_ = ec.Log.RegisterLog(c, "Access denied for SearchEmployeesByName")
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	if query == "" {
		_ = ec.Log.RegisterLog(c, "Empty name query provided in SearchEmployeesByName")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
		return
	}

	employees, err := ec.Service.SearchEmployeesByName(query)
	if err != nil {
		_ = ec.Log.RegisterLog(c, "Error retrieving employees by name: "+query+" - "+err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving employees"})
		return
	}

	if len(employees) == 0 {
		_ = ec.Log.RegisterLog(c, "No employees found with name: "+query)
		c.JSON(http.StatusNotFound, gin.H{"message": "No employees found"})
		return
	}

	var employeesDTO []dtos.GetEmployeeDTO
	for _, employee := range employees {
		employeesDTO = append(employeesDTO, dtos.GetEmployeeDTO{
			ID:               employee.ID,
			Names:            employee.Names,
			LastNames:        employee.LastNames,
			PersonalID:       employee.PersonalID,
			Address:          employee.Address,
			PhoneNumbers:     employee.PhoneNumbers,
			UserID:           employee.UserID,
			IdentifierTypeID: employee.IdentifierTypeID,
		})
	}

	_ = ec.Log.RegisterLog(c, "Successfully found employees with name: "+query)
	c.JSON(http.StatusOK, employeesDTO)
}

// CreateEmployee godoc
// @Summary      Create a new employee
// @Description  Creates a new employee. The request body must contain the employee details.
// @Tags         employees
// @Accept       json
// @Produce      json
// @Param        employee body dtos.CreateEmployeeDTO true "Employee information"
// @Success      201 {object} dtos.GetEmployeeDTO "Successfully created employee"
// @Failure      400 {object} models.ErrorResponse "Invalid JSON format, or missing fields"
// @Failure      403 {object} models.ErrorResponse "Permission denied"
// @Failure      409 {object} models.ErrorResponse "Employee with this Personal ID already exists"
// @Failure      500 {object} models.ErrorResponse "Error creating employee"
// @Security     ApiKeyAuth
// @Router       /employees [post]
func (ec *EmployeeController) CreateEmployee(c *gin.Context) {
	permissionId := config.PERMISSION_CREATE_EMPLOYEE

	if ec.Log.RegisterLog(c, "Attempting to create an employee") != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
		return
	}

	if !ec.Auth.CheckPermission(c, permissionId) {
		_ = ec.Log.RegisterLog(c, "Permission denied for CreateEmployee")
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	var dto dtos.CreateEmployeeDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		_ = ec.Log.RegisterLog(c, "Invalid JSON format: "+err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format", "details": err.Error()})
		return
	}

	existingEmployee, _ := ec.Service.GetEmployeeByID(dto.PersonalID)
	if existingEmployee != nil {
		_ = ec.Log.RegisterLog(c, "Attempt to create duplicate employee with PersonalID: "+dto.PersonalID)
		c.JSON(http.StatusConflict, gin.H{"error": "An employee with this Personal ID already exists"})
		return
	}

	if dto.UserID <= 0 || dto.IdentifierTypeID <= 0 {
		_ = ec.Log.RegisterLog(c, "Invalid UserID or IdentifierTypeID: UserID="+strconv.Itoa(dto.UserID)+", IdentifierTypeID="+strconv.Itoa(dto.IdentifierTypeID))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid User ID or Identifier Type ID"})
		return
	}

	employee := &models.Employee{
		Names:            dto.Names,
		LastNames:        dto.LastNames,
		PersonalID:       dto.PersonalID,
		Address:          dto.Address,
		PhoneNumbers:     dto.PhoneNumbers,
		UserID:           dto.UserID,
		IdentifierTypeID: dto.IdentifierTypeID,
	}

	createdEmployee, err := ec.Service.CreateEmployee(employee)
	if err != nil {
		_ = ec.Log.RegisterLog(c, "Error creating employee: "+err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating employee", "details": err.Error()})
		return
	}

	employeeDTO := dtos.GetEmployeeDTO{
		ID:               createdEmployee.ID,
		Names:            createdEmployee.Names,
		LastNames:        createdEmployee.LastNames,
		PersonalID:       createdEmployee.PersonalID,
		Address:          createdEmployee.Address,
		PhoneNumbers:     createdEmployee.PhoneNumbers,
		UserID:           createdEmployee.UserID,
		IdentifierTypeID: createdEmployee.IdentifierTypeID,
	}

	_ = ec.Log.RegisterLog(c, "Successfully created employee with PersonalID: "+createdEmployee.PersonalID)
	c.JSON(http.StatusCreated, employeeDTO)
}

// UpdateEmployee godoc
// @Summary      Update an existing employee
// @Description  Updates an employee's details by ID. The request body must contain the updated employee information.
// @Tags         employees
// @Accept       json
// @Produce      json
// @Param        id path string true "Employee ID"
// @Param        employee body dtos.UpdateEmployeeDTO true "Updated employee information"
// @Success      200 {object} dtos.GetEmployeeDTO "Successfully updated employee"
// @Failure      400 {object} models.ErrorResponse "Invalid JSON format"
// @Failure      403 {object} models.ErrorResponse "Permission denied"
// @Failure      404 {object} models.ErrorResponse "Employee not found"
// @Failure      500 {object} models.ErrorResponse "Error updating employee"
// @Security     ApiKeyAuth
// @Router       /employees/{id} [put]
func (ec *EmployeeController) UpdateEmployee(c *gin.Context) {
	permissionId := config.PERMISSION_UPDATE_EMPLOYEE

	if err := ec.Log.RegisterLog(c, "Attempting to update an employee"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
		return
	}

	if !ec.Auth.CheckPermission(c, permissionId) {
		_ = ec.Log.RegisterLog(c, "Permission denied for UpdateEmployee")
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	id := c.Param("id")

	var dto dtos.UpdateEmployeeDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		_ = ec.Log.RegisterLog(c, "Invalid JSON in UpdateEmployee: "+err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	employee, err := ec.Service.GetEmployeeByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			_ = ec.Log.RegisterLog(c, "Employee not found in UpdateEmployee: ID = "+id)
			c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
			return
		}
		_ = ec.Log.RegisterLog(c, "Error retrieving employee in UpdateEmployee: "+err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	employee.Names = dto.Names
	employee.LastNames = dto.LastNames
	employee.PersonalID = dto.PersonalID
	employee.Address = dto.Address
	employee.PhoneNumbers = dto.PhoneNumbers
	employee.UserID = dto.UserID
	employee.IdentifierTypeID = dto.IdentifierTypeID

	err = ec.Service.UpdateEmployee(employee)
	if err != nil {
		_ = ec.Log.RegisterLog(c, "Error updating employee: "+err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	employeeDTO := dtos.GetEmployeeDTO{
		ID:               employee.ID,
		Names:            employee.Names,
		LastNames:        employee.LastNames,
		PersonalID:       employee.PersonalID,
		Address:          employee.Address,
		PhoneNumbers:     employee.PhoneNumbers,
		UserID:           employee.UserID,
		IdentifierTypeID: employee.IdentifierTypeID,
	}

	_ = ec.Log.RegisterLog(c, "Successfully updated Employee with ID: "+id)

	c.JSON(http.StatusOK, employeeDTO)
}
