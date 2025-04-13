package controllers

import (
	"bookmarker/internal/repositories"
	"bookmarker/internal/services"
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BookmarksController struct {
    DB *sql.DB
}

func NewBookmarksController(db *sql.DB) *BookmarksController {
    return &BookmarksController{DB: db}
}

func (bc *BookmarksController) GetBookmarks(c *gin.Context) {
    // Initialize the repository and service
    bookmarkRepo := repositories.NewBookmarkRepository(bc.DB)
    bookmarkService := services.NewBookmarkService(bookmarkRepo)

    // Pagination parameters (hardcoded for now)
    page, pageSize := 1, 10

    // Fetch bookmarks
    bookmarks, err := bookmarkService.ListBookmarks(page, pageSize)
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
        URL         string  `json:"url" binding:"required"`
        Title       string  `json:"title"`
        Description string  `json:"description"`
        Thumbnail   string  `json:"thumbnail"`
    }

    // Bind JSON input
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }

    // Initialize the repository and service
    bookmarkRepo := repositories.NewBookmarkRepository(bc.DB)
    bookmarkService := services.NewBookmarkService(bookmarkRepo)

    // Create the bookmark
    bookmark, err := bookmarkService.CreateBookmark(input.URL, input.Title, input.Description, input.Thumbnail)
    if err != nil {
        log.Printf("Failed to create bookmark: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create bookmark"})
        return
    }

    // Respond with the created bookmark
    c.JSON(http.StatusOK, gin.H{"bookmark": bookmark})
}