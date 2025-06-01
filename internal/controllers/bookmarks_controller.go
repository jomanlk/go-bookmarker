package controllers

import (
	"bookmarker/internal/repositories"
	"bookmarker/internal/services"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BookmarksController struct {
    DB *pgxpool.Pool
}

func NewBookmarksController(db *pgxpool.Pool) *BookmarksController {
    return &BookmarksController{DB: db}
}

func (bc *BookmarksController) GetBookmarks(c *gin.Context) {
    // Initialize the repository and service
    bookmarkRepo := repositories.NewBookmarkRepository(bc.DB)
    tagRepo := repositories.NewTagRepository(bc.DB)
    bookmarkService := services.NewBookmarkServiceWithTags(bookmarkRepo, tagRepo)

    // Extract pagination parameters from the request
    pageStr := c.DefaultQuery("page", "1")
    pageSizeStr := c.DefaultQuery("limit", "50")

    // Convert page and limit to integers
    page, err := strconv.Atoi(pageStr)
    if err != nil || page < 1 {
        page = 1
    }

    limit, err := strconv.Atoi(pageSizeStr)
    if err != nil || limit < 1 {
        limit = 10
    }

    // Fetch bookmarks with tags using the service layer
    bookmarks, err := bookmarkService.ListBookmarksWithTags(page, limit)
    if err != nil {
        log.Printf("Failed to list bookmarks: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list bookmarks"})
        return
    }

    // Respond with a JSON array of the bookmarks
    c.JSON(http.StatusOK, gin.H{
        "bookmarks": bookmarks,
    })
}

func (bc *BookmarksController) CreateBookmark(c *gin.Context) {
    var input struct {
        URL         string   `json:"url" binding:"required"`
        Title       string   `json:"title"`
        Description string   `json:"description"`
        Thumbnail   string   `json:"thumbnail"`
        Tags        []string `json:"tags"`
    }

    // Bind JSON input
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }

    // Initialize the repositories and service
    bookmarkRepo := repositories.NewBookmarkRepository(bc.DB)
    tagRepo := repositories.NewTagRepository(bc.DB)
    bookmarkService := services.NewBookmarkServiceWithTags(bookmarkRepo, tagRepo)
    

    // Create the bookmark with tags
    bookmark, err := bookmarkService.CreateBookmarkWithTags(input.URL, input.Title, input.Description, input.Thumbnail, input.Tags, time.Now())
    if err != nil {
        log.Printf("Failed to create bookmark: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create bookmark"})
        return
    }
    // get the tags for the bookmark
    tags, err := tagRepo.GetTagsForBookmark(int(bookmark.ID))
    if err != nil {
        log.Printf("Failed to fetch tags for bookmark: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tags for bookmark"})
        return
    }
    bookmark.Tags = tags

    // Respond with the created bookmark
    c.JSON(http.StatusOK, gin.H{"bookmark": bookmark})
}

func (bc *BookmarksController) GetBookmark(c *gin.Context) {
    // Extract the bookmark ID from the URL
    id := c.Param("id")

    // Initialize the repository and service
    bookmarkRepo := repositories.NewBookmarkRepository(bc.DB)
    tagRepo := repositories.NewTagRepository(bc.DB)
    bookmarkService := services.NewBookmarkServiceWithTags(bookmarkRepo, tagRepo)

    // Convert the bookmark ID to an integer
    bookmarkID, err := strconv.Atoi(id)
    if err != nil {
        log.Printf("Invalid bookmark ID: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bookmark ID"})
        return
    }

    // Fetch the bookmark by ID with tags
    bookmark, err := bookmarkService.GetBookmarkWithTags(bookmarkID)
    if err != nil {
        log.Printf("Failed to fetch bookmark: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bookmark"})
        return
    }

    // Respond with the bookmark wrapped in the 'bookmark' key
    c.JSON(http.StatusOK, gin.H{"bookmark": bookmark})
}

func (bc *BookmarksController) UpdateBookmark(c *gin.Context) {
    id := c.Param("id")
    bookmarkID, err := strconv.Atoi(id)
    if err != nil {
        log.Printf("Invalid bookmark ID: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bookmark ID"})
        return
    }

    var input struct {
        URL         *string   `json:"url"`
        Title       *string   `json:"title"`
        Description *string   `json:"description"`
        Thumbnail   *string   `json:"thumbnail"`
        Tags        *[]string `json:"tags"`
    }
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }

    updateFields := make(map[string]interface{})
    if input.URL != nil {
        updateFields["url"] = *input.URL
    }
    if input.Title != nil {
        updateFields["title"] = *input.Title
    }
    if input.Description != nil {
        updateFields["description"] = *input.Description
    }
    if input.Thumbnail != nil {
        updateFields["thumbnail"] = *input.Thumbnail
    }

    bookmarkRepo := repositories.NewBookmarkRepository(bc.DB)
    tagRepo := repositories.NewTagRepository(bc.DB)
    bookmarkService := services.NewBookmarkServiceWithTags(bookmarkRepo, tagRepo)

    var updatedBookmark interface{}
    if input.Tags != nil {
        // Update tags as well
        updatedBookmark, err = bookmarkService.UpdateBookmarkWithTags(bookmarkID, updateFields, *input.Tags)
    } else {
        updatedBookmark, err = bookmarkService.UpdateBookmark(bookmarkID, updateFields)
    }
    if err != nil {
        log.Printf("Failed to update bookmark: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update bookmark"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"bookmark": updatedBookmark})
}