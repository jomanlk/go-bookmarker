package controllers

import (
	"bookmarker/internal/clients"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UrlController struct{}

func NewUrlController() *UrlController {
	return &UrlController{}
}

// UrlPreviewHandler handles GET /url/preview?url=...
func (uc *UrlController) UrlPreviewHandler(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing url parameter"})
		return
	}
	client := clients.NewURLPreviewApiClient()
	preview, err := client.Fetch(url)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"url": preview.URL,
		"title": preview.Title,
		"description": preview.Description,
		"image": preview.Image,
	})
}
