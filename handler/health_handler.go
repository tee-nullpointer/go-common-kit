package common_handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HealthCheck(c *gin.Context) {
	healthData := gin.H{
		"status": "healthy",
	}
	c.JSON(http.StatusOK, healthData)
}
