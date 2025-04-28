package main

import (
	"bookmarker/internal/controllers"
	"bookmarker/internal/dbutil"
	"bookmarker/internal/repositories"
	"bookmarker/internal/server"
	"bookmarker/internal/services"
	"log"
	"net/http"
	"os"

	"bookmarker/internal/middleware"

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

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	tokenRepo := repositories.NewTokenRepository(db)
	refreshTokenRepo := repositories.NewRefreshTokenRepository(db)

	// Initialize services
	userService := services.NewUserService(userRepo)
	tokenService := services.NewTokenService(tokenRepo)
	refreshTokenService := services.NewRefreshTokenService(refreshTokenRepo)
	authService := services.NewAuthService(userService, tokenService, refreshTokenService)

	// Initialize controllers
	bookmarksController := controllers.NewBookmarksController(db)
	searchController := controllers.NewSearchController(db)
	tagsController := controllers.NewTagsController(db)
	userController := controllers.NewUserController(authService)

	// Define routes
	// Public routes
	r.POST("/login", userController.Login)
	r.POST("/refresh", userController.Refresh)

	// Protected routes
	r.Use(middleware.AuthMiddleware(authService))
	r.GET("/bookmarks", bookmarksController.GetBookmarks)
	r.POST("/bookmarks", bookmarksController.CreateBookmark)
	r.GET("/bookmarks/:id", bookmarksController.GetBookmark)
	r.PATCH("/bookmarks/:id", bookmarksController.UpdateBookmark)
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
	if len(os.Args) > 3 && os.Args[1] == "create-user" {
		username := os.Args[2]
		password := os.Args[3]
		createUserCommand(username, password)
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
	log.Fatalf("No command provided. Use 'start-server', 'import-pinboard <filename>', or 'create-user <username> <password>'")
}

