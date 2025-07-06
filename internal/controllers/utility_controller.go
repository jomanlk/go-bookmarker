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
	type backupRequest struct {
		Token string `json:"token"`
	}
	var req backupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	webhookSecret := os.Getenv("WEBHOOK_SECRET")
	if req.Token == "" || req.Token != webhookSecret {
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
