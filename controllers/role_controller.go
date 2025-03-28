package controllers

import (
	"fmt"
	"net/http"
	"totesbackend/config"
	"totesbackend/controllers/utilities"
	"totesbackend/dtos"
	"totesbackend/services"

	"github.com/gin-gonic/gin"
)

type RoleController struct {
	Service *services.RoleService
	Auth    *utilities.AuthorizationUtil
}

func NewRoleController(service *services.RoleService, auth *utilities.AuthorizationUtil) *RoleController {
	return &RoleController{Service: service, Auth: auth}
}

func (rc *RoleController) GetRoleByID(c *gin.Context) {

	permissionId := config.PERMISSION_GET_ROLE_BY_ID

	if !rc.Auth.CheckPermission(c, permissionId) {
		return
	}

	idParam := c.Param("id")
	var id uint
	if _, err := fmt.Sscanf(idParam, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	role, err := rc.Service.GetRoleByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}

	permissionIDs, err := rc.Service.GetRolePermissions(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving role permissions"})
		return
	}

	roleDTO := dtos.RoleDTO{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		Permissions: make([]string, len(permissionIDs)),
	}

	for i, permissionID := range permissionIDs {
		roleDTO.Permissions[i] = fmt.Sprintf("%d", permissionID)
	}

	c.JSON(http.StatusOK, roleDTO)
}

func (rc *RoleController) GetAllRoles(c *gin.Context) {

	permissionId := config.PERMISSION_GET_ALL_ROLES

	if !rc.Auth.CheckPermission(c, permissionId) {
		return
	}

	roles, err := rc.Service.GetAllRoles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving roles"})
		return
	}

	var rolesDTO []dtos.RoleDTO
	for _, role := range roles {
		permissionIDs, err := rc.Service.GetRolePermissions(role.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving role permissions"})
			return
		}

		roleDTO := dtos.RoleDTO{
			ID:          role.ID,
			Name:        role.Name,
			Description: role.Description,
			Permissions: make([]string, len(permissionIDs)),
		}

		for i, permissionID := range permissionIDs {
			roleDTO.Permissions[i] = fmt.Sprintf("%d", permissionID)
		}

		rolesDTO = append(rolesDTO, roleDTO)
	}

	c.JSON(http.StatusOK, rolesDTO)
}

func (rc *RoleController) GetAllPermissionsOfRole(c *gin.Context) {

	permissionId := config.PERMISSION_GET_ALL_PERMISSIONS_OF_ROLE

	if !rc.Auth.CheckPermission(c, permissionId) {
		return
	}
	roleIDParam := c.Param("id")
	var roleID uint
	if _, err := fmt.Sscanf(roleIDParam, "%d", &roleID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	permissions, err := rc.Service.GetAllPermissionsOfRole(roleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving permissions for role"})
		return
	}

	c.JSON(http.StatusOK, permissions)
}

func (rc *RoleController) ExistRole(c *gin.Context) {

	permissionId := config.PERMISSION_EXIST_ROLE

	if !rc.Auth.CheckPermission(c, permissionId) {
		return
	}

	idParam := c.Param("id")
	var id uint
	if _, err := fmt.Sscanf(idParam, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	exists, err := rc.Service.ExistRole(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking role existence"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"exists": exists})
}

func (rc *RoleController) SearchRolesByID(c *gin.Context) {

	permissionId := config.PERMISSION_SEARCH_ROLE_BY_ID

	if !rc.Auth.CheckPermission(c, permissionId) {
		return
	}

	username := c.GetHeader("Username")
	fmt.Println("Request made by user:", username)

	query := c.Query("id")
	roles, err := rc.Service.SearchRolesByID(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error searching roles by ID"})
		return
	}

	c.JSON(http.StatusOK, roles)
}

func (rc *RoleController) SearchRolesByName(c *gin.Context) {

	permissionId := config.PERMISSION_SEARCH_ROLE_BY_NAME

	if !rc.Auth.CheckPermission(c, permissionId) {
		return
	}

	username := c.GetHeader("Username")
	fmt.Println("Request made by user:", username)

	query := c.Query("name")
	roles, err := rc.Service.SearchRolesByName(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error searching roles by name"})
		return
	}

	c.JSON(http.StatusOK, roles)
}
