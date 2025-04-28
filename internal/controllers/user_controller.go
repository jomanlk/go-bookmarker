package controllers

import (
	"bookmarker/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserLoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserController struct {
	AuthService *services.AuthService
}

func NewUserController(authService *services.AuthService) *UserController {
	return &UserController{AuthService: authService}
}

func (uc *UserController) Login(c *gin.Context) {
	var req UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	result, err := uc.AuthService.Authenticate(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}
	// Set access and refresh tokens as HTTP-only cookies
	cookieMaxAge := 60 * 30 // 30 minutes for access token
	refreshCookieMaxAge := 60 * 60 * 24 * 30 // 30 days for refresh token
	domain := "" // set to your domain if needed
	secure := false // set to true if using HTTPS
	c.SetCookie("access_token", result.AccessToken, cookieMaxAge, "/", domain, secure, true)
	c.SetCookie("refresh_token", result.RefreshToken, refreshCookieMaxAge, "/", domain, secure, true)
	c.JSON(http.StatusOK, UserLoginResponse{AccessToken: result.AccessToken, RefreshToken: result.RefreshToken})
}

// RefreshRequest and Refresh handler

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (uc *UserController) Refresh(c *gin.Context) {
	var req RefreshRequest

	// Try to get refresh token from cookie first
	refreshToken, errCookie := c.Cookie("refresh_token")
	if errCookie == nil && refreshToken != "" {
		req.RefreshToken = refreshToken
	} else {
		// Fallback: get from request body
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}
	}

	result, err := uc.AuthService.RefreshTokens(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		return
	}
	// Set new tokens as HTTP-only cookies
	cookieMaxAge := 60 * 30 // 30 minutes for access token
	refreshCookieMaxAge := 60 * 60 * 24 * 30 // 30 days for refresh token
	domain := "" // set to your domain if needed
	secure := false // set to true if using HTTPS
	c.SetCookie("access_token", result.AccessToken, cookieMaxAge, "/", domain, secure, true)
	c.SetCookie("refresh_token", result.RefreshToken, refreshCookieMaxAge, "/", domain, secure, true)
	c.JSON(http.StatusOK, UserLoginResponse{AccessToken: result.AccessToken, RefreshToken: result.RefreshToken})
}

// LogoutRequest for logout endpoint
// Accepts both access and refresh tokens
type LogoutRequest struct {
	AccessToken  string `json:"access_token" binding:"required"`
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// Logout handler
func (uc *UserController) Logout(c *gin.Context) {
	var req LogoutRequest

	// Try to get tokens from cookies first
	accessToken, errAccess := c.Cookie("access_token")
	refreshToken, errRefresh := c.Cookie("refresh_token")

	if errAccess == nil && errRefresh == nil && accessToken != "" && refreshToken != "" {
		req.AccessToken = accessToken
		req.RefreshToken = refreshToken
	} else {
		// Fallback: get from request body
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}
	}

	err := uc.AuthService.Logout(req.AccessToken, req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Remove the cookies by setting MaxAge to -1
	domain := "" // set to your domain if needed
	secure := false // set to true if using HTTPS
	c.SetCookie("access_token", "", -1, "/", domain, secure, true)
	c.SetCookie("refresh_token", "", -1, "/", domain, secure, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// Me handler to return the currently logged-in user's info
func (uc *UserController) Me(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userID, ok := userIDVal.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}
	user, err := uc.AuthService.UserService.GetUserByID(int64(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
        "user": user,
    })
}
