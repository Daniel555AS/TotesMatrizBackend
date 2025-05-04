package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ControllerHealthCheck godoc
// @Summary      Health Check
// @Description  Returns a 200 status if the server is running correctly.
// @Tags         Health
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /health [get]
func ControllerHealthCheck(c *gin.Context) {
	// IndentedJSON serializes the given struct as pretty JSON (indented + endlines) into the response body.
	// Includes de status onf http and the struct
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Successful Health Check."})
}
