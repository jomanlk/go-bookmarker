package controllers

import (
	"bookmarker/internal/repositories"
	"bookmarker/internal/services"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SearchController struct {
	DB *pgxpool.Pool
}

func NewSearchController(db *pgxpool.Pool) *SearchController {
	return &SearchController{DB: db}
}

func (sc *SearchController) SearchBookmarks(c *gin.Context) {
	bookmarkRepo := repositories.NewBookmarkRepository(sc.DB)
	tagRepo := repositories.NewTagRepository(sc.DB)
	bookmarkService := services.NewBookmarkServiceWithTags(bookmarkRepo, tagRepo)

	searchQuery := c.DefaultQuery("q", "")
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("limit", "50")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(pageSizeStr)
	if err != nil || limit < 1 {
		limit = 10
	}
	if searchQuery == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing search query parameter 'q'"})
		return
	}

	bookmarks, err := bookmarkService.SearchBookmarks(searchQuery, page, limit)
	if err != nil {
		log.Printf("Failed to search bookmarks: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search bookmarks"})
		return
	}

	for i := range bookmarks {
		tags, err := tagRepo.GetTagsForBookmark(int(bookmarks[i].ID))
		if err == nil {
			bookmarks[i].Tags = tags
		}
	}

	c.JSON(http.StatusOK, gin.H{"bookmarks": bookmarks})
}

// GetBookmarksByTag fetches bookmarks by tag name with pagination
func (sc *SearchController) GetBookmarksByTag(c *gin.Context) {
	tagName := c.Query("tag")
	if tagName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing tag parameter"})
		return
	}

	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "50")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}
	tagRepo := repositories.NewTagRepository(sc.DB)
	bookmarkRepo := repositories.NewBookmarkRepository(sc.DB)

	// Find tag by name
	tag, err := tagRepo.GetTagByName(tagName)
	if err != nil  {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tag not found"})
		return
	}

	offset := (page - 1) * limit
	bookmarks, err := bookmarkRepo.ListBookmarksByTag(int(tag.ID), offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bookmarks for tag"})
		return
	}
	// Attach tags to each bookmark
	for i := range bookmarks {
		tags, err := tagRepo.GetTagsForBookmark(int(bookmarks[i].ID))
		if err == nil {
			bookmarks[i].Tags = tags
		}
	}
	c.JSON(http.StatusOK, gin.H{"bookmarks": bookmarks})
}
