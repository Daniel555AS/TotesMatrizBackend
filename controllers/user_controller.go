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

type UserController struct {
	Service *services.UserService
	Auth    *utilities.AuthorizationUtil
	Log     *utilities.LogUtil
}

func NewUserController(service *services.UserService, auth *utilities.AuthorizationUtil, log *utilities.LogUtil) *UserController {
	return &UserController{Service: service, Auth: auth, Log: log}
}

// GetUserByID godoc
// @Summary      Get a user by ID
// @Description  Retrieves a user by their ID. Returns user details if found, otherwise returns an error.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id  path     string  true  "User ID"
// @Success      200  {object}  dtos.GetUserDTO  "User details"
// @Failure      400  {object}  models.ErrorResponse  "Invalid ID format"
// @Failure      403  {object}  models.ErrorResponse  "Permission denied"
// @Failure      404  {object}  models.ErrorResponse  "User not found"
// @Failure      500  {object}  models.ErrorResponse  "Error retrieving user"
// @Security     ApiKeyAuth
// @Router       /users/{id} [get]
func (uc *UserController) GetUserByID(c *gin.Context) {
	id := c.Param("id")

	// Log de intento
	if uc.Log.RegisterLog(c, "Attempting to retrieve user with ID: "+id) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
		return
	}

	permissionId := config.PERMISSION_GET_USER_BY_ID
	if !uc.Auth.CheckPermission(c, permissionId) {
		_ = uc.Log.RegisterLog(c, "Access denied for GetUserByID")
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	user, err := uc.Service.GetUserByID(id)
	if err != nil {
		_ = uc.Log.RegisterLog(c, "Error retrieving user with ID "+id+": "+err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving user"})
		return
	}

	if user == nil {
		_ = uc.Log.RegisterLog(c, "User with ID "+id+" not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	userDTO := dtos.GetUserDTO{
		ID:          user.ID,
		Email:       user.Email,
		Password:    user.Password,
		UserTypeID:  user.UserTypeID,
		UserStateID: user.UserStateTypeID,
	}

	_ = uc.Log.RegisterLog(c, "Successfully retrieved user with ID: "+id)
	c.JSON(http.StatusOK, userDTO)
}

// GetAllUsers godoc
// @Summary      Get all users
// @Description  Retrieves a list of all users in the system.
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      200  {array}  dtos.GetUserDTO  "List of users"
// @Failure      403  {object}  models.ErrorResponse  "Permission denied"
// @Failure      404  {object}  models.ErrorResponse  "Users not found"
// @Failure      500  {object}  models.ErrorResponse  "Error retrieving users"
// @Security     ApiKeyAuth
// @Router       /users [get]
func (uc *UserController) GetAllUsers(c *gin.Context) {
	permissionId := config.PERMISSION_GET_ALL_USERS

	// Intento de obtener todos los usuarios
	if uc.Log.RegisterLog(c, "Attempting to retrieve all users") != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
		return
	}

	if !uc.Auth.CheckPermission(c, permissionId) {
		_ = uc.Log.RegisterLog(c, "Access denied for GetAllUsers")
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	users, err := uc.Service.GetAllUsers()
	if err != nil {
		_ = uc.Log.RegisterLog(c, "Error retrieving all users: "+err.Error())
		c.JSON(http.StatusNotFound, gin.H{"error": "Users not found"})
		return
	}

	var usersDTO []dtos.GetUserDTO
	for _, user := range users {
		userDTO := dtos.GetUserDTO{
			ID:          user.ID,
			Email:       user.Email,
			Password:    user.Password,
			UserTypeID:  user.UserTypeID,
			UserStateID: user.UserStateTypeID,
		}
		usersDTO = append(usersDTO, userDTO)
	}

	_ = uc.Log.RegisterLog(c, "Successfully retrieved all users")
	c.JSON(http.StatusOK, usersDTO)
}

// SearchUsersByID godoc
// @Summary      Search users by ID
// @Description  Searches for users by ID based on the provided query parameter.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   query     string  true  "User ID to search"
// @Success      200  {array}  dtos.GetUserDTO  "List of users matching the search criteria"
// @Failure      400  {object}  models.ErrorResponse  "Query parameter is required"
// @Failure      403  {object}  models.ErrorResponse  "Permission denied"
// @Failure      404  {object}  models.ErrorResponse  "Users not found"
// @Failure      500  {object}  models.ErrorResponse  "Error searching users"
// @Security     ApiKeyAuth
// @Router       /users/searchByID [get]
func (uc *UserController) SearchUsersByID(c *gin.Context) {
	permissionId := config.PERMISSION_SEARCH_USER_BY_ID
	query := c.Query("id")

	// Intento de búsqueda
	if uc.Log.RegisterLog(c, "Attempting to search users by ID with query: "+query) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
		return
	}

	if !uc.Auth.CheckPermission(c, permissionId) {
		_ = uc.Log.RegisterLog(c, "Access denied for SearchUsersByID")
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	if query == "" {
		_ = uc.Log.RegisterLog(c, "Query parameter 'id' is missing for SearchUsersByID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter is required"})
		return
	}

	users, err := uc.Service.SearchUsersByID(query)
	if err != nil {
		_ = uc.Log.RegisterLog(c, "Error searching users by ID "+query+": "+err.Error())
		c.JSON(http.StatusNotFound, gin.H{"error": "Users not found"})
		return
	}

	if len(users) == 0 {
		_ = uc.Log.RegisterLog(c, "No users found with ID containing: "+query)
		c.JSON(http.StatusNotFound, gin.H{"message": "No users found"})
		return
	}

	var usersDTO []dtos.GetUserDTO
	for _, user := range users {
		userDTO := dtos.GetUserDTO{
			ID:          user.ID,
			Email:       user.Email,
			Password:    user.Password,
			UserTypeID:  user.UserTypeID,
			UserStateID: user.UserStateTypeID,
		}
		usersDTO = append(usersDTO, userDTO)
	}

	_ = uc.Log.RegisterLog(c, "Successfully retrieved users with query: "+query)
	c.JSON(http.StatusOK, usersDTO)
}

// SearchUsersByEmail godoc
// @Summary      Search users by email
// @Description  Searches for users by email based on the provided query parameter.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        email   query     string  true  "User email to search"
// @Success      200     {array}  dtos.GetUserDTO  "List of users matching the search criteria"
// @Failure      400     {object}  models.ErrorResponse  "Query parameter is required"
// @Failure      403     {object}  models.ErrorResponse  "Permission denied"
// @Failure      404     {object}  models.ErrorResponse  "Users not found"
// @Failure      500     {object}  models.ErrorResponse  "Error searching users"
// @Security     ApiKeyAuth
// @Router       /users/searchByEmail [get]
func (uc *UserController) SearchUsersByEmail(c *gin.Context) {
	permissionId := config.PERMISSION_SEARCH_USERS_BY_EMAIL
	query := c.Query("email")

	// Intento de búsqueda
	if uc.Log.RegisterLog(c, "Attempting to search users by email with query: "+query) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
		return
	}

	if !uc.Auth.CheckPermission(c, permissionId) {
		_ = uc.Log.RegisterLog(c, "Access denied for SearchUsersByEmail")
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	if query == "" {
		_ = uc.Log.RegisterLog(c, "Query parameter 'email' is missing for SearchUsersByEmail")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter is required"})
		return
	}

	users, err := uc.Service.SearchUsersByEmail(query)
	if err != nil {
		_ = uc.Log.RegisterLog(c, "Error searching users by email "+query+": "+err.Error())
		c.JSON(http.StatusNotFound, gin.H{"error": "Users not found"})
		return
	}

	if len(users) == 0 {
		_ = uc.Log.RegisterLog(c, "No users found with email containing: "+query)
		c.JSON(http.StatusNotFound, gin.H{"message": "No users found"})
		return
	}

	var usersDTO []dtos.GetUserDTO
	for _, user := range users {
		userDTO := dtos.GetUserDTO{
			ID:          user.ID,
			Email:       user.Email,
			Password:    user.Password,
			UserTypeID:  user.UserTypeID,
			UserStateID: user.UserStateTypeID,
		}
		usersDTO = append(usersDTO, userDTO)
	}

	_ = uc.Log.RegisterLog(c, "Successfully retrieved users with email query: "+query)
	c.JSON(http.StatusOK, usersDTO)
}

type request struct {
	UserState int `json:"user_state"` // Correctly defines the JSON binding
}

// UpdateUserState godoc
// @Summary      Update user state
// @Description  Updates the state of a user based on the provided user ID and state value.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id      path     string  true  "User ID"
// @Param        body    body     request  true  "User state to update" // Updated 'object' to 'body' and set 'true'
// @Success      200     {object}  dtos.GetUserDTO  "Updated user information"
// @Failure      400     {object}  models.ErrorResponse  "Invalid request body"
// @Failure      403     {object}  models.ErrorResponse  "Permission denied"
// @Failure      404     {object}  models.ErrorResponse  "User not found"
// @Failure      500     {object}  models.ErrorResponse  "Error updating user state"
// @Security     ApiKeyAuth
// @Router       /users/{id}/state [patch]
func (uc *UserController) UpdateUserState(c *gin.Context) {
	var request request
	permissionId := config.PERMISSION_UPDATE_USER_STATE
	id := c.Param("id")

	// Log de intento
	if uc.Log.RegisterLog(c, "Attempting to update user state for ID: "+id) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
		return
	}

	// Check permission
	if !uc.Auth.CheckPermission(c, permissionId) {
		_ = uc.Log.RegisterLog(c, "Access denied for UpdateUserState")
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Bind JSON request body to request struct
	if err := c.ShouldBindJSON(&request); err != nil {
		_ = uc.Log.RegisterLog(c, "Invalid request body for UpdateUserState: "+err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update user state
	user, err := uc.Service.UpdateUserState(id, request.UserState)
	if err != nil {
		_ = uc.Log.RegisterLog(c, "User not found with ID "+id+" while updating state")
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Prepare userDTO to return
	userDTO := dtos.GetUserDTO{
		ID:          user.ID,
		Email:       user.Email,
		Password:    user.Password,
		UserTypeID:  user.UserTypeID,
		UserStateID: user.UserStateTypeID,
	}

	// Log success and return response
	_ = uc.Log.RegisterLog(c, "Successfully updated user state for ID: "+id)
	c.JSON(http.StatusOK, userDTO)
}

// UpdateUser godoc
// @Summary      Update user information
// @Description  Updates user details such as email, password, user type, and state.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id      path     string  true  "User ID"
// @Param        body    body     dtos.UpdateUserDTO  true  "User details to update"
// @Success      200     {object}  dtos.GetUserDTO  "Updated user information"
// @Failure      400     {object}  models.ErrorResponse  "Invalid request body"
// @Failure      403     {object}  models.ErrorResponse  "Permission denied"
// @Failure      404     {object}  models.ErrorResponse  "User not found"
// @Failure      500     {object}  models.ErrorResponse  "Error updating user information"
// @Security     ApiKeyAuth
// @Router       /users/{id} [put]
func (uc *UserController) UpdateUser(c *gin.Context) {
	permissionId := config.PERMISSION_UPDATE_USER
	id := c.Param("id")

	if uc.Log.RegisterLog(c, "Attempting to update user with ID: "+id) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
		return
	}

	if !uc.Auth.CheckPermission(c, permissionId) {
		_ = uc.Log.RegisterLog(c, "Access denied for UpdateUser")
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	var dto dtos.UpdateUserDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		_ = uc.Log.RegisterLog(c, "Invalid request body for UpdateUser: "+err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := uc.Service.GetUserByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			_ = uc.Log.RegisterLog(c, "User not found with ID: "+id)
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		_ = uc.Log.RegisterLog(c, "Error retrieving user with ID: "+id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	user.Email = dto.Email
	user.Password = dto.Password
	user.UserTypeID = dto.UserTypeID
	user.UserStateTypeID = dto.UserStateID

	err = uc.Service.UpdateUser(user)

	dtoUser := dtos.GetUserDTO{
		ID:          user.ID,
		Email:       user.Email,
		Password:    user.Password,
		UserTypeID:  user.UserTypeID,
		UserStateID: user.UserStateTypeID,
	}

	if err != nil {
		_ = uc.Log.RegisterLog(c, "Failed to update user with ID: "+id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	_ = uc.Log.RegisterLog(c, "Successfully updated user with ID: "+id)
	c.JSON(http.StatusOK, dtoUser)
}

// CreateUser godoc
// @Summary      Create a new user
// @Description  Creates a new user with the provided email, password, user type, and state.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        body    body     dtos.CreateUserDTO  true  "User details to create"
// @Success      201     {object}  dtos.GetUserDTO  "Created user information"
// @Failure      400     {object}  models.ErrorResponse  "Invalid request body"
// @Failure      403     {object}  models.ErrorResponse  "Permission denied"
// @Failure      409     {object}  models.ErrorResponse  "Email already in use"
// @Failure      500     {object}  models.ErrorResponse  "Error creating user"
// @Security     ApiKeyAuth
// @Router       /users [post]
func (uc *UserController) CreateUser(c *gin.Context) {
	permissionId := config.PERMISSION_CREATE_USER

	if uc.Log.RegisterLog(c, "Attempting to create new user") != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
		return
	}

	if !uc.Auth.CheckPermission(c, permissionId) {
		_ = uc.Log.RegisterLog(c, "Access denied for CreateUser")
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	var dto dtos.CreateUserDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		_ = uc.Log.RegisterLog(c, "Invalid request body for CreateUser: "+err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existingUser, _ := uc.Service.GetUserByEmail(dto.Email)
	if existingUser != nil {
		_ = uc.Log.RegisterLog(c, "Email already in use: "+dto.Email)
		c.JSON(http.StatusConflict, gin.H{"error": "Email already in use"})
		return
	}

	newUser := models.User{
		Email:           dto.Email,
		Password:        dto.Password,
		UserTypeID:      dto.UserTypeID,
		UserStateTypeID: dto.UserStateID,
	}

	createdUser, err := uc.Service.CreateUser(&newUser)
	if err != nil {
		_ = uc.Log.RegisterLog(c, "Failed to create user: "+err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	userDTO := dtos.GetUserDTO{
		ID:          createdUser.ID,
		Email:       createdUser.Email,
		Password:    createdUser.Password,
		UserTypeID:  createdUser.UserTypeID,
		UserStateID: createdUser.UserStateTypeID,
	}

	_ = uc.Log.RegisterLog(c, "Successfully created user with ID: "+strconv.Itoa(createdUser.ID))
	c.JSON(http.StatusCreated, userDTO)
}
