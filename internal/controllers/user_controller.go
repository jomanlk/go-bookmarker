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
	c.JSON(http.StatusOK, UserLoginResponse{AccessToken: result.AccessToken, RefreshToken: result.RefreshToken})
}

// RefreshRequest and Refresh handler

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (uc *UserController) Refresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	result, err := uc.AuthService.RefreshTokens(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		return
	}
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
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	err := uc.AuthService.Logout(req.AccessToken, req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
