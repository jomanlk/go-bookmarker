package controllers

import (
	"bookmarker/internal/services"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type UtilityController struct{}

func NewUtilityController() *UtilityController {
	return &UtilityController{}
}

// BackupDBHandler triggers the backup process
func (uc *UtilityController) BackupDBHandler(c *gin.Context) {
	token := c.Query("token")
	webhookSecret := os.Getenv("WEBHOOK_SECRET")
	if token == "" || token != webhookSecret {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	err := services.BackupPostgresDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Backup completed and uploaded successfully"})
}
