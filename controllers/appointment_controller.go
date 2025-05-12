package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"totesbackend/config"
	"totesbackend/controllers/utilities"
	"totesbackend/models"
	"totesbackend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AppointmentController struct {
	Service *services.AppointmentService
	Auth    *utilities.AuthorizationUtil
	Log     *utilities.LogUtil
}

func NewAppointmentController(service *services.AppointmentService, auth *utilities.AuthorizationUtil,
	log *utilities.LogUtil) *AppointmentController {
	return &AppointmentController{Service: service, Auth: auth, Log: log}
}

// GetAppointmentByID godoc
// @Summary      Get Appointment by ID
// @Description  Retrieves an appointment by its unique ID. Requires permission.
// @Tags         appointments
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Appointment ID"
// @Success      200  {object}  models.Appointment           "The appointment object"
// @Failure      400  {object}  models.ErrorResponse         "Invalid appointment ID"
// @Failure      401  {object}  models.ErrorResponse         "Unauthorized or permission denied"
// @Failure      404  {object}  models.ErrorResponse         "Appointment not found"
// @Failure      500  {object}  models.ErrorResponse         "Internal server error"
// @Security     ApiKeyAuth
// @Router       /appointments/{id} [get]
func (ac *AppointmentController) GetAppointmentByID(c *gin.Context) {
	if ac.Log.RegisterLog(c, "Attempting to get appointment by ID") != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
		return
	}

	permissionId := config.PERMISSION_GET_APPOINTMENT_BY_ID
	if !ac.Auth.CheckPermission(c, permissionId) {
		_ = ac.Log.RegisterLog(c, "Access denied for GetAppointmentByID")
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		_ = ac.Log.RegisterLog(c, "Invalid appointment ID: "+c.Param("id"))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid appointment ID"})
		return
	}

	appointment, err := ac.Service.GetAppointmentByID(id)
	if err != nil {
		_ = ac.Log.RegisterLog(c, "Appointment not found for ID: "+strconv.Itoa(id))
		c.JSON(http.StatusNotFound, gin.H{"error": "Appointment not found"})
		return
	}

	_ = ac.Log.RegisterLog(c, "Appointment retrieved successfully for ID: "+strconv.Itoa(id))
	c.JSON(http.StatusOK, appointment)
}

// GetAllAppointments godoc
// @Summary      Get all appointments
// @Description  Retrieves a list of all appointments. Requires proper permission.
// @Tags         appointments
// @Accept       json
// @Produce      json
// @Success      200  {array}   models.Appointment       "List of all appointments"
// @Failure      401  {object}  models.ErrorResponse     "Unauthorized or permission denied"
// @Failure      500  {object}  models.ErrorResponse     "Error retrieving appointments or logging"
// @Security     ApiKeyAuth
// @Router       /appointments [get]
func (ac *AppointmentController) GetAllAppointments(c *gin.Context) {
	if ac.Log.RegisterLog(c, "Attempting to retrieve all appointments") != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
		return
	}

	permissionId := config.PERMISSION_GET_ALL_APPOINTMENTS
	if !ac.Auth.CheckPermission(c, permissionId) {
		_ = ac.Log.RegisterLog(c, "Access denied for GetAllAppointments")
		return
	}

	appointments, err := ac.Service.GetAllAppointments()
	if err != nil {
		_ = ac.Log.RegisterLog(c, "Error retrieving appointments")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving appointments"})
		return
	}

	_ = ac.Log.RegisterLog(c, "All appointments retrieved successfully")
	c.JSON(http.StatusOK, appointments)
}

// SearchAppointmentsByID godoc
// @Summary      Search appointments by ID
// @Description  Search appointments using a partial or complete ID string. Requires permission.
// @Tags         appointments
// @Accept       json
// @Produce      json
// @Param        id     query     string  true  "Appointment ID to search"
// @Success      200    {array}   models.Appointment
// @Failure      401    {object} models.ErrorResponse  "Unauthorized or permission denied"
// @Failure      404    {object}  models.ErrorResponse   "No appointments found"
// @Failure      500    {object}  models.ErrorResponse  "Error retrieving appointments or logging"
// @Security     ApiKeyAuth
// @Router       /appointments/searchByID [get]
func (ac *AppointmentController) SearchAppointmentsByID(c *gin.Context) {

	if ac.Log.RegisterLog(c, "Attempting to search appointments by ID") != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
		return
	}

	permissionId := config.PERMISSION_SEARCH_APPOINTMENTS_BY_ID
	if !ac.Auth.CheckPermission(c, permissionId) {
		_ = ac.Log.RegisterLog(c, "Access denied for SearchAppointmentsByID")
		return
	}

	query := c.Query("id")
	fmt.Println("Searching appointments by ID with:", query)

	appointments, err := ac.Service.SearchAppointmentsByID(query)
	if err != nil {
		_ = ac.Log.RegisterLog(c, "Error retrieving appointments")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving appointments"})
		return
	}

	if len(appointments) == 0 {
		_ = ac.Log.RegisterLog(c, "No appointments found for given ID")
		c.JSON(http.StatusNotFound, gin.H{"message": "No appointments found"})
		return
	}

	_ = ac.Log.RegisterLog(c, "Appointments found by ID successfully")
	c.JSON(http.StatusOK, appointments)
}

// SearchAppointmentsByCustomerID godoc
// @Summary      Search appointments by Customer ID
// @Description  Search appointments by the customer's ID. Requires permission to access this data.
// @Tags         appointments
// @Accept       json
// @Produce      json
// @Param        id     query     string  true  "Customer ID to search appointments"
// @Success      200    {array}   models.Appointment   "List of appointments found"
// @Failure      401    {object} models.ErrorResponse   "Unauthorized or permission denied"
// @Failure      404    {object}  models.ErrorResponse   "No appointments found for the given customer ID"
// @Failure      500    {object}  models.ErrorResponse  "Error retrieving appointments or logging"
// @Security     ApiKeyAuth
// @Router       /appointments/searchByCustomerID [get]
func (ac *AppointmentController) SearchAppointmentsByCustomerID(c *gin.Context) {

	if ac.Log.RegisterLog(c, "Attempting to search appointments by customer ID") != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
		return
	}

	permissionId := config.PERMISSION_GET_APPOINTMENT_BY_CUSTOMER_ID
	if !ac.Auth.CheckPermission(c, permissionId) {
		_ = ac.Log.RegisterLog(c, "Access denied for SearchAppointmentsByCustomerID")
		return
	}

	query := c.Query("id")
	fmt.Println("Searching appointments by Customer ID with:", query)

	appointments, err := ac.Service.SearchAppointmentsByCustomerID(query)
	if err != nil {
		_ = ac.Log.RegisterLog(c, "Error retrieving appointments by customer ID")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving appointments"})
		return
	}

	if len(appointments) == 0 {
		_ = ac.Log.RegisterLog(c, "No appointments found for given customer ID")
		c.JSON(http.StatusNotFound, gin.H{"message": "No appointments found"})
		return
	}

	_ = ac.Log.RegisterLog(c, "Appointments found by customer ID successfully")
	c.JSON(http.StatusOK, appointments)
}

// SearchAppointmentsByState godoc
// @Summary      Search appointments by state
// @Description  Search appointments based on their state (e.g., confirmed, pending). Requires permission to access this data.
// @Tags         appointments
// @Accept       json
// @Produce      json
// @Param        state   query     bool    true  "State of the appointment (true for confirmed, false for pending)"
// @Success      200     {array}   models.Appointment   "List of appointments found based on state"
// @Failure      400     {object}  models.ErrorResponse   "Invalid state value provided"
// @Failure      401     {object}  models.ErrorResponse   "Unauthorized or permission denied"
// @Failure      500     {object}  models.ErrorResponse   "Error retrieving appointments or logging"
// @Security     ApiKeyAuth
// @Router       /appointments/searchByState [get]
func (ac *AppointmentController) SearchAppointmentsByState(c *gin.Context) {

	if ac.Log.RegisterLog(c, "Attempting to search appointments by state") != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
		return
	}

	permissionId := config.PERMISSION_SEARCH_APPOINTMENT_BY_STATE
	if !ac.Auth.CheckPermission(c, permissionId) {
		_ = ac.Log.RegisterLog(c, "Access denied for SearchAppointmentsByState")
		return
	}

	state, err := strconv.ParseBool(c.Query("state"))
	if err != nil {
		_ = ac.Log.RegisterLog(c, "Invalid state value provided for appointment search")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid state value"})
		return
	}

	appointments, err := ac.Service.SearchAppointmentsByState(state)
	if err != nil {
		_ = ac.Log.RegisterLog(c, "Error retrieving appointments by state")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving appointments"})
		return
	}

	_ = ac.Log.RegisterLog(c, "Appointments retrieved successfully by state")
	c.JSON(http.StatusOK, appointments)
}

// GetAppointmentsByCustomerID godoc
// @Summary      Get appointments by customer ID
// @Description  Retrieves a list of appointments associated with a specific customer. Requires permission to view appointments by customer ID.
// @Tags         appointments
// @Accept       json
// @Produce      json
// @Param        customerID  path      int                          true  "ID of the customer"
// @Success      200         {array}   models.Appointment           "List of appointments"
// @Failure      400         {object}  models.ErrorResponse       "Invalid customer ID"
// @Failure      401         {object} models.ErrorResponse        "Unauthorized or permission denied"
// @Failure      500         {object}  models.ErrorResponse       "Error retrieving appointments"
// @Router       /appointments/customer/{customerID} [get]
func (ac *AppointmentController) GetAppointmentsByCustomerID(c *gin.Context) {

	if ac.Log.RegisterLog(c, "Attempting to retrieve appointments by customer ID") != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
		return
	}

	permissionId := config.PERMISSION_GET_APPOINTMENT_BY_CUSTOMER_ID
	if !ac.Auth.CheckPermission(c, permissionId) {
		_ = ac.Log.RegisterLog(c, "Access denied for GetAppointmentsByCustomerID")
		return
	}

	customerID, err := strconv.Atoi(c.Param("customerID"))
	if err != nil {
		_ = ac.Log.RegisterLog(c, "Invalid customer ID provided")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID"})
		return
	}

	appointments, err := ac.Service.GetAppointmentsByCustomerID(customerID)
	if err != nil {
		_ = ac.Log.RegisterLog(c, "Error retrieving appointments by customer ID")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving appointments"})
		return
	}

	_ = ac.Log.RegisterLog(c, "Appointments retrieved successfully by customer ID")
	c.JSON(http.StatusOK, appointments)
}

// CreateAppointment godoc
// @Summary      Create a new appointment
// @Description  Create a new appointment in the system. Requires permission to create appointments.
// @Tags         appointments
// @Accept       json
// @Produce      json
// @Param        appointment  body      models.Appointment  true  "Appointment data to create"
// @Success      201          {object}  models.Appointment  "Appointment successfully created"
// @Failure      400          {object}  models.ErrorResponse   "Invalid JSON format or appointment limit reached"
// @Failure      403          {object}  models.ErrorResponse   "Forbidden, no permission to create appointments"
// @Failure      500          {object}  models.ErrorResponse   "Error creating appointment or logging"
// @Security     ApiKeyAuth
// @Router       /appointments [post]
func (ac *AppointmentController) CreateAppointment(c *gin.Context) {
	if ac.Log.RegisterLog(c, "Attempting to create appointment") != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
		return
	}

	permissionId := config.PERMISSION_CREATE_APPOINTMENT
	if !ac.Auth.CheckPermission(c, permissionId) {
		_ = ac.Log.RegisterLog(c, "Access denied for CreateAppointment")
		c.JSON(http.StatusForbidden, gin.H{"error": "No tienes permisos para crear citas"})
		return
	}

	var appointment models.Appointment
	if err := c.ShouldBindJSON(&appointment); err != nil {
		_ = ac.Log.RegisterLog(c, "Invalid JSON format when creating appointment")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Formato JSON inválido"})
		return
	}

	createdAppointment, err := ac.Service.CreateAppointment(appointment)
	if err != nil {
		if err.Error() == "ya existen 3 citas agendadas para esta fecha y hora" {
			_ = ac.Log.RegisterLog(c, "limite de citas alcanzado :v")
			c.JSON(http.StatusBadRequest, gin.H{"error": "no se puede crear la cita. Ya hay 3 citas agendadas para esta fecha y hora."})
		} else {
			_ = ac.Log.RegisterLog(c, "Error creando cita")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear la cita"})
		}
		return
	}

	_ = ac.Log.RegisterLog(c, "Cita creada exitosamente")
	c.JSON(http.StatusCreated, createdAppointment)
}

// UpdateAppointment godoc
// @Summary      Update an existing appointment
// @Description  Update the details of an existing appointment. Requires permission to update appointments.
// @Tags         appointments
// @Accept       json
// @Produce      json
// @Param        id          path      int                 true  "Appointment ID to update"
// @Param        appointment body      models.Appointment   true  "Appointment data to update"
// @Success      200         {object}  models.Appointment   "Appointment successfully updated"
// @Failure      400         {object}  models.ErrorResponse   "Invalid appointment ID or JSON format"
// @Failure      403         {object} models.ErrorResponse   "Forbidden, no permission to update appointments"
// @Failure      404         {object} models.ErrorResponse   "Appointment not found for update"
// @Failure      500         {object}  models.ErrorResponse   "Error updating appointment or logging"
// @Security     ApiKeyAuth
// @Router       /appointments/{id} [put]
func (ac *AppointmentController) UpdateAppointment(c *gin.Context) {

	if ac.Log.RegisterLog(c, "Attempting to update appointment") != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
		return
	}

	permissionId := config.PERMISSION_UPDATE_APPOINTMENT

	if !ac.Auth.CheckPermission(c, permissionId) {
		_ = ac.Log.RegisterLog(c, "Access denied for UpdateAppointment")
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		_ = ac.Log.RegisterLog(c, "Invalid appointment ID format")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid appointment ID"})
		return
	}

	var appointment models.Appointment
	if err := c.ShouldBindJSON(&appointment); err != nil {
		_ = ac.Log.RegisterLog(c, "Invalid JSON format on update appointment")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	appointment.ID = id

	err = ac.Service.UpdateAppointment(&appointment)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			_ = ac.Log.RegisterLog(c, "Appointment not found for update")
			c.JSON(http.StatusNotFound, gin.H{"error": "Appointment not found"})
			return
		}
		_ = ac.Log.RegisterLog(c, "Error updating appointment")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating appointment"})
		return
	}

	_ = ac.Log.RegisterLog(c, "Appointment updated successfully")
	c.JSON(http.StatusOK, appointment)
}

// GetAppointmentByCustomerIDAndDate godoc
// @Summary      Get appointment by customer ID and date
// @Description  Retrieve an appointment based on the customer ID and appointment date. Requires permission to get appointments.
// @Tags         appointments
// @Accept       json
// @Produce      json
// @Param        customerId  query     int    true  "Customer ID"
// @Param        dateTime    query     string true  "Appointment date and time (format: YYYY-MM-DD HH:MM:SS)"
// @Success      200         {object}  models.Appointment  "Appointment successfully retrieved"
// @Failure      400         {object}  models.ErrorResponse   "Invalid customer ID or date format"
// @Failure      401         {object}  models.ErrorResponse   "Unauthorized or permission denied"
// @Failure      404         {object} models.ErrorResponse   "Appointment not found for the given customer ID and date"
// @Failure      500         {object}  models.ErrorResponse  "Error retrieving appointment or logging"
// @Security     ApiKeyAuth
// @Router       /appointments/byCustomerIdAndDate [get]
func (ac *AppointmentController) GetAppointmentByCustomerIDAndDate(c *gin.Context) {

	if ac.Log.RegisterLog(c, "Attempting to get appointment by customer ID and date") != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
		return
	}

	permissionId := config.PERMISSION_GET_APPOINTMENTS_BY_CUSTOMERID_AND_DATE

	if !ac.Auth.CheckPermission(c, permissionId) {
		_ = ac.Log.RegisterLog(c, "Access denied for GetAppointmentByCustomerIDAndDate")
		return
	}

	customerID, err := strconv.Atoi(c.Query("customerId"))
	if err != nil {
		_ = ac.Log.RegisterLog(c, "Invalid customer ID format")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID"})
		return
	}

	dateTime, err := time.Parse("2006-01-02 15:04:05", c.Query("dateTime"))
	if err != nil {
		_ = ac.Log.RegisterLog(c, "Invalid date format for GetAppointmentByCustomerIDAndDate")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format, use 'YYYY-MM-DD HH:MM:SS'"})
		return
	}

	appointment, err := ac.Service.GetAppointmentByCustomerIDAndDate(customerID, dateTime)
	if err != nil {
		_ = ac.Log.RegisterLog(c, "Appointment not found for given customer ID and date")
		c.JSON(http.StatusNotFound, gin.H{"error": "Appointment not found"})
		return
	}

	_ = ac.Log.RegisterLog(c, "Appointment successfully retrieved by customer ID and date")

	response := gin.H{
		"id":           appointment.ID,
		"dateTime":     appointment.DateTime,
		"customerName": appointment.CustomerName,
		"lastName":     appointment.LastName,
		"email":        appointment.Email,
		"customerID":   appointment.CustomerID,
	}

	c.JSON(http.StatusOK, response)
}

// DeleteAppointmentByID godoc
// @Summary      Delete appointment by ID
// @Description  Delete an appointment based on the appointment ID. Requires permission to delete appointments.
// @Tags         appointments
// @Accept       json
// @Produce      json
// @Param        id  path     int  true  "Appointment ID"
// @Success      200 {object} models.MessageResponse "Appointment deleted successfully"
// @Failure      400 {object} models.ErrorResponse  "Invalid appointment ID format"
// @Failure      401 {object} models.ErrorResponse  "Unauthorized or permission denied"
// @Failure      404 {object} models.ErrorResponse  "Appointment not found for the given ID"
// @Failure      500 {object} models.ErrorResponse  "Error deleting the appointment"
// @Security     ApiKeyAuth
// @Router       /appointments/deleteAppointment/{id} [delete]
func (ac *AppointmentController) DeleteAppointmentByID(c *gin.Context) {
	if ac.Log.RegisterLog(c, "Attempting to delete appointment") != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
		return
	}

	permissionId := config.PERMISSION_DELETE_APPOINTMENT
	if !ac.Auth.CheckPermission(c, permissionId) {
		_ = ac.Log.RegisterLog(c, "Access denied for DeleteAppointmentByID")
		c.JSON(http.StatusForbidden, gin.H{"error": "No tienes permisos para eliminar citas"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		_ = ac.Log.RegisterLog(c, "Invalid appointment ID: "+c.Param("id"))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid appointment ID"})
		return
	}

	err = ac.Service.DeleteAppointmentByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			_ = ac.Log.RegisterLog(c, "Appointment not found for ID: "+strconv.Itoa(id))
			c.JSON(http.StatusNotFound, gin.H{"error": "Appointment not found"})
		} else {
			_ = ac.Log.RegisterLog(c, "Error deleting appointment")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	_ = ac.Log.RegisterLog(c, "Appointment deleted successfully for ID: "+strconv.Itoa(id))
	c.JSON(http.StatusOK, gin.H{"message": "Appointment deleted successfully"})

}

// GetAppointmentsByHourRange godoc
// @Summary      Get appointment count by hourly range for a specific date
// @Description  Retrieves the number of appointments for each hour within a specified date. Requires permission to view appointments by hour.
// @Tags         appointments
// @Accept       json
// @Produce      json
// @Param        date  query     string  true  "Date in YYYY-MM-DD format"
// @Success      200  {array}  models.Appointment  "List of hourly appointment counts"
// @Failure      400  {object}  models.ErrorResponse  "Invalid date format or missing 'date' parameter"
// @Failure      401  {object}  models.ErrorResponse "Unauthorized or permission denied"
// @Failure      500  {object}  models.ErrorResponse  "Error retrieving appointment counts"
// @Security     ApiKeyAuth
// @Router       /appointments/hourly-count [get]
func (c *AppointmentController) GetAppointmentsByHourRange(ctx *gin.Context) {

	permissionId := config.PERMISSION_GET_APPOINTMENTS_BY_HOUR
	if !c.Auth.CheckPermission(ctx, permissionId) {
		_ = c.Log.RegisterLog(ctx, "Access denied for CreateAppointment")
		ctx.JSON(http.StatusForbidden, gin.H{"error": "No tienes permisos para crear citas"})
		return
	}

	dateParam := ctx.Query("date")
	if dateParam == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Se requiere el parámetro 'date' en formato YYYY-MM-DD"})
		return
	}

	date, err := time.Parse("2006-01-02", dateParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Formato de fecha inválido. Usa YYYY-MM-DD"})
		return
	}

	counts, err := c.Service.GetHourlyAppointmentCount(date)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al contar las citas: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"date": dateParam, "appointmentsPerHour": counts})
}
