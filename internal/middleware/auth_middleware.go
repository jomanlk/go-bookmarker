package middleware

import (
	"bookmarker/internal/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates the access token and attaches the user ID to the context
func AuthMiddleware(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		var token string
		if header != "" && strings.HasPrefix(header, "Bearer ") {
			token = strings.TrimPrefix(header, "Bearer ")
		} else {
			// Try to get token from cookie if Authorization header is missing or invalid
			cookieToken, err := c.Cookie("access_token")
			if err == nil && cookieToken != "" {
				token = cookieToken
			}
		}
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header or access_token cookie missing or invalid"})
			return
		}
		userID, err := authService.ValidateAccessToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}
		// Attach userID to context
		c.Set("userID", userID)
		c.Next()
	}
}
