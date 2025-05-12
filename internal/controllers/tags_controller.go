package controllers

import (
	"bookmarker/internal/repositories"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TagsController struct {
	DB *pgxpool.Pool
}

func NewTagsController(db *pgxpool.Pool) *TagsController {
	return &TagsController{DB: db}
}

// ListAllTags handles GET /tags and returns all tags in the database (no pagination)
func (tc *TagsController) ListAllTags(c *gin.Context) {
	tagRepo := repositories.NewTagRepository(tc.DB)
	tags, err := tagRepo.ListAllTags()
	if err != nil {
		log.Printf("Failed to list tags: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list tags"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tags": tags})
}

// ListTags handles GET /tags and returns paginated tags
func (tc *TagsController) ListTags(c *gin.Context) {
	tagRepo := repositories.NewTagRepository(tc.DB)
	// Get pagination params
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "50")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 50
	}
	tags, err := tagRepo.ListTags(page, limit)
	if err != nil {
		log.Printf("Failed to list tags: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list tags"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tags": tags})
}
