package main

import (
	"bookmarker/internal/controllers"
	"bookmarker/internal/dbutil"
	"bookmarker/internal/repositories"
	"bookmarker/internal/services"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"bookmarker/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)


func main() {
	_ = godotenv.Load("../../.env") // Loads .env file if present

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
		db, err := dbutil.OpenPostgresDB()
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
		r := setupRouter(db)

		// Graceful shutdown setup
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		server := &http.Server{
			Addr:    ":8080",
			Handler: r,
		}
		go func() {
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("listen: %s\n", err)
			}
		}()
		<-quit
		log.Println("Shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Server forced to shutdown: %v", err)
		}
		if err := dbutil.ShutdownPostgresDB(db); err != nil {
			log.Printf("Error closing database: %v", err)
		}
		log.Println("Server exiting")
		return
	}
	if len(os.Args) > 1 {
		log.Fatalf("Unrecognized command: %s", os.Args[1])
	}
	log.Fatalf("No command provided. Use 'start-server', 'import-pinboard <filename>', or 'create-user <username> <password>'")
}


func setupRouter(db *pgxpool.Pool) *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	// Apply the Gin CORS middleware
	r.Use(middleware.CorsMiddleware())

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

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
	telegramController := controllers.NewTelegramController(db)
	urlController := controllers.NewUrlController()

	// Define routes
	// Public routes
	r.POST("/login", userController.Login)
	r.POST("/refresh", userController.Refresh)
	r.POST("/logout", userController.Logout)
	// Telegram webhook route
	r.POST("/telegram/listen", telegramController.TelegramWebhookHandler)

	// Protected routes
	r.Use(middleware.AuthMiddleware(authService))
	r.GET("/bookmarks", bookmarksController.GetBookmarks)
	r.POST("/bookmarks", bookmarksController.CreateBookmark)
	r.GET("/bookmarks/:id", bookmarksController.GetBookmark)
	r.PATCH("/bookmarks/:id", bookmarksController.UpdateBookmark)
	r.DELETE("/bookmarks/:id", bookmarksController.DeleteBookmark)
	r.GET("/search", searchController.SearchBookmarks)
	r.GET("/bookmarks/tag", searchController.GetBookmarksByTag)
	r.GET("/tags", tagsController.ListTags)
	r.GET("/me", userController.Me)
	r.GET("/url/preview", urlController.UrlPreviewHandler)

	return r
}