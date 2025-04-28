package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// CorsMiddleware is a Gin middleware that sets CORS headers and handles preflight requests.
func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Use APP_URL as the allowed origin if ALLOWED_CORS_ORIGINS is not set
		allowedOrigin := os.Getenv("ALLOWED_CORS_ORIGINS")
		if allowedOrigin == "" {
			allowedOrigin = "*"
		}
		c.Writer.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
