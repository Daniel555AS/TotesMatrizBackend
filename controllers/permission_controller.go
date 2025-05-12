package controllers

import (
	"fmt"
	"net/http"
	"totesbackend/config"
	"totesbackend/controllers/utilities"
	"totesbackend/services"

	"github.com/gin-gonic/gin"
)

type PermissionController struct {
	Service *services.PermissionService
	Auth    *utilities.AuthorizationUtil
	Log     *utilities.LogUtil
}

func NewPermissionController(service *services.PermissionService, auth *utilities.AuthorizationUtil, log *utilities.LogUtil) *PermissionController {
	return &PermissionController{Service: service, Auth: auth, Log: log}
}

// GetPermissionByID godoc
// @Summary      Get permission by ID
// @Description  Retrieves a specific permission by its unique ID.
// @Tags         permissions
// @Produce      json
// @Param        id   path      int                           true  "Permission ID"
// @Success      200  {object}  models.Permission             "Permission data"
// @Failure      400  {object}  models.ErrorResponse          "Invalid permission ID"
// @Failure      403  {object}  models.ErrorResponse          "Access denied"
// @Failure      404  {object}  models.ErrorResponse          "Permission not found"
// @Failure      500  {object}  models.ErrorResponse          "Internal server error"
// @Security     ApiKeyAuth
// @Router       /permissions/{id} [get]
func (pc *PermissionController) GetPermissionByID(c *gin.Context) {
	permissionId := config.PERMISSION_GET_PERMISSION_BY_ID

	if !pc.Auth.CheckPermission(c, permissionId) {
		if pc.Log.RegisterLog(c, "Access denied for GetPermissionByID") != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
			return
		}
		return
	}

	idParam := c.Param("id")
	var id uint
	if _, err := fmt.Sscanf(idParam, "%d", &id); err != nil {
		if pc.Log.RegisterLog(c, "Invalid permission ID: "+idParam) != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid permission ID"})
		return
	}

	if pc.Log.RegisterLog(c, "Attempting to retrieve Permission with ID: "+idParam) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
		return
	}

	permission, err := pc.Service.GetPermissionByID(id)
	if err != nil {
		if pc.Log.RegisterLog(c, "Permission with ID "+idParam+" not found") != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "Permission not found"})
		return
	}

	if pc.Log.RegisterLog(c, "Successfully retrieved Permission with ID: "+idParam) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
		return
	}

	c.JSON(http.StatusOK, permission)
}

// GetAllPermissions godoc
// @Summary      Get all permissions
// @Description  Retrieves a list of all permissions available in the system.
// @Tags         permissions
// @Produce      json
// @Success      200  {array}   models.Permission             "List of permissions"
// @Failure      403  {object}  models.ErrorResponse          "Access denied"
// @Failure      500  {object}  models.ErrorResponse          "Internal server error"
// @Security     ApiKeyAuth
// @Router       /permissions [get]
func (pc *PermissionController) GetAllPermissions(c *gin.Context) {
	permissionId := config.PERMISSION_GET_ALL_PERMISSIONS

	if !pc.Auth.CheckPermission(c, permissionId) {
		if pc.Log.RegisterLog(c, "Access denied for GetAllPermissions") != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
			return
		}
		return
	}

	if pc.Log.RegisterLog(c, "Attempting to retrieve all permissions") != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
		return
	}

	permissions, err := pc.Service.GetAllPermissions()
	if err != nil {
		if pc.Log.RegisterLog(c, "Error retrieving all permissions: "+err.Error()) != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving permissions"})
		return
	}

	if pc.Log.RegisterLog(c, "Successfully retrieved all permissions") != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
		return
	}

	c.JSON(http.StatusOK, permissions)
}

// SearchPermissionsByID godoc
// @Summary      Search permissions by ID
// @Description  Retrieves a list of permissions that match the given ID pattern.
// @Tags         permissions
// @Produce      json
// @Param        id   query     string  true  "ID to search for (partial or full match)"
// @Success      200  {array}   models.Permission             "List of matching permissions"
// @Failure      400  {object}  models.ErrorResponse          "Missing or invalid query parameter"
// @Failure      403  {object}  models.ErrorResponse          "Access denied"
// @Failure      500  {object}  models.ErrorResponse          "Internal server error"
// @Security     ApiKeyAuth
// @Router       /permissions/searchByID [get]
func (pc *PermissionController) SearchPermissionsByID(c *gin.Context) {
	permissionId := config.PERMISSION_SEARCH_PERMISSION_BY_ID

	if !pc.Auth.CheckPermission(c, permissionId) {
		if pc.Log.RegisterLog(c, "Access denied for SearchPermissionsByID") != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
			return
		}
		return
	}

	query := c.Query("id")
	if query == "" {
		if pc.Log.RegisterLog(c, "SearchPermissionsByID: missing 'id' query parameter") != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
		return
	}

	if pc.Log.RegisterLog(c, "Attempting to search permissions by ID: "+query) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
		return
	}

	permissions, err := pc.Service.SearchPermissionsByID(query)
	if err != nil {
		if pc.Log.RegisterLog(c, "Error retrieving permissions by ID: "+err.Error()) != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving permissions"})
		return
	}

	if pc.Log.RegisterLog(c, "Successfully retrieved permissions by ID: "+query) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
		return
	}

	c.JSON(http.StatusOK, permissions)
}

// SearchPermissionsByName godoc
// @Summary      Search permissions by name
// @Description  Retrieves a list of permissions that match the given name pattern.
// @Tags         permissions
// @Produce      json
// @Param        name   query     string  true  "Name to search for (partial or full match)"
// @Success      200    {array}   models.Permission             "List of matching permissions"
// @Failure      400    {object}  models.ErrorResponse          "Missing or invalid query parameter"
// @Failure      403    {object}  models.ErrorResponse          "Access denied"
// @Failure      500    {object}  models.ErrorResponse          "Internal server error"
// @Security     ApiKeyAuth
// @Router       /permissions/searchByName [get]
func (pc *PermissionController) SearchPermissionsByName(c *gin.Context) {
	permissionId := config.PERMISSION_SEARCH_PERMISSION_BY_NAME

	if !pc.Auth.CheckPermission(c, permissionId) {
		if pc.Log.RegisterLog(c, "Access denied for SearchPermissionsByName") != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
			return
		}
		return
	}

	query := c.Query("name")
	if query == "" {
		if pc.Log.RegisterLog(c, "SearchPermissionsByName: missing 'name' query parameter") != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
		return
	}

	if pc.Log.RegisterLog(c, "Attempting to search permissions by name: "+query) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
		return
	}

	permissions, err := pc.Service.SearchPermissionsByName(query)
	if err != nil {
		if pc.Log.RegisterLog(c, "Error retrieving permissions by name: "+err.Error()) != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving permissions"})
		return
	}

	if pc.Log.RegisterLog(c, "Successfully retrieved permissions by name: "+query) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
		return
	}

	c.JSON(http.StatusOK, permissions)
}
