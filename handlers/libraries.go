package handlers

import (
	"github.com/OkanUysal/go-logger"
	"github.com/OkanUysal/go-starter-api/types"
	"github.com/gin-gonic/gin"
)

// GetLibraries returns all available libraries
// @Summary      Get all available libraries
// @Description  Returns a list of all 10 production Go libraries with their metadata
// @Tags         Libraries
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}  "success, data (array of Library), count"
// @Router       /libraries [get]
func GetLibraries(c *gin.Context) {
	logger.Debug("Fetching available libraries")

	if Metrics != nil {
		Metrics.IncrementCounter("libraries_requested_total", nil)
	}

	libraries := types.GetAvailableLibraries()

	logger.Info("Libraries retrieved", logger.Int("count", len(libraries)))
	c.JSON(200, gin.H{
		"success": true,
		"data":    libraries,
		"count":   len(libraries),
	})
}
