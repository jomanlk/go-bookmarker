package main

import (
	"bookmarker/internal/controllers"
	"bookmarker/internal/server"
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	// Apply the CORS middleware from server.go
	r.Use(func(c *gin.Context) {
		server.CorsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c.Next()
		})).ServeHTTP(c.Writer, c.Request)
	})

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// Initialize database connection
	db, err := sql.Open("sqlite3", "../../internal/db/bookmarker_db1.db")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize controllers
	bookmarksController := controllers.NewBookmarksController(db)

	// Define routes
	r.GET("/bookmarks", bookmarksController.GetBookmarks)
	r.POST("/bookmarks", bookmarksController.CreateBookmark)
	r.GET("/bookmarks/:id", bookmarksController.GetBookmark)
	r.PATCH("/bookmarks/:id", bookmarksController.UpdateBookmark)

	return r
}

func main() {
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
