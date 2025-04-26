package main

import (
	"bookmarker/internal/controllers"
	"bookmarker/internal/dbutil"
	"bookmarker/internal/server"
	"log"
	"net/http"
	"os"

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
	db, err := dbutil.OpenSqliteDB("../../internal/db/bookmarker_db1.db")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize controllers
	bookmarksController := controllers.NewBookmarksController(db)
	searchController := controllers.NewSearchController(db)
	tagsController := controllers.NewTagsController(db)

	// Define routes
	r.GET("/bookmarks", bookmarksController.GetBookmarks)
	r.POST("/bookmarks", bookmarksController.CreateBookmark)
	r.GET("/bookmarks/:id", bookmarksController.GetBookmark)
	r.PATCH("/bookmarks/:id", bookmarksController.UpdateBookmark)
	// Use SearchController for /search
	r.GET("/search", searchController.SearchBookmarks)
	r.GET("/bookmarks/tag", searchController.GetBookmarksByTag)
	r.GET("/tags", tagsController.ListTags)

	return r
}

func main() {
	if len(os.Args) > 2 && os.Args[1] == "import-pinboard" {
		filename := os.Args[2]
		importPinboard(filename)
		return
	}
	if len(os.Args) > 1 && os.Args[1] == "start-server" {
		r := setupRouter()
		r.Run(":8080")
		return
	}
	if len(os.Args) > 1 {
		log.Fatalf("Unrecognized command: %s", os.Args[1])
	}
	log.Fatalf("No command provided. Use 'start-server' or 'import-pinboard <filename>'")
}
