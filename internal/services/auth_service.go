package services

import (
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserService          *UserService
	TokenService         *TokenService
	RefreshTokenService  *RefreshTokenService
}

func NewAuthService(userService *UserService, tokenService *TokenService, refreshTokenService *RefreshTokenService) *AuthService {
	return &AuthService{
		UserService:         userService,
		TokenService:        tokenService,
		RefreshTokenService: refreshTokenService,
	}
}

// AuthResult holds both access and refresh tokens
type AuthResult struct {
	AccessToken  string
	RefreshToken string
}

// Authenticate checks credentials, creates and stores both tokens, and returns them
func (a *AuthService) Authenticate(username, password string) (*AuthResult, error) {
	user, err := a.UserService.GetUserByUsername(username)
	if (err != nil) {
		return nil, err
	}
	
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return nil, err
	}
	accessToken := generateRandomToken()
	accessExpiresAt := time.Now().Add(30 * time.Minute) // 30 minutes expiry
	if err := a.TokenService.CreateToken(int(user.ID), accessToken, accessExpiresAt); err != nil {
		return nil, err
	}
	refreshToken := generateRandomToken()
	refreshExpiresAt := time.Now().Add(30 * 24 * time.Hour) // 30 days
	if err := a.RefreshTokenService.CreateRefreshToken(int(user.ID), refreshToken, refreshExpiresAt); err != nil {
		return nil, err
	}
	return &AuthResult{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

// RefreshTokens validates a refresh token and issues new tokens
func (a *AuthService) RefreshTokens(refreshToken string) (*AuthResult, error) {
	userID, err := a.RefreshTokenService.FindByToken(refreshToken)
	if err != nil {
		return nil, err
	}
	// Invalidate old refresh token
	err = a.RefreshTokenService.DeleteRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}
	// Issue new tokens
	accessToken := generateRandomToken()
	accessExpiresAt := time.Now().Add(30 * time.Minute) // 30 minutes expiry
	if err := a.TokenService.CreateToken(userID, accessToken, accessExpiresAt); err != nil {
		return nil, err
	}
	newRefreshToken := generateRandomToken()
	refreshExpiresAt := time.Now().Add(30 * 24 * time.Hour)
	if err := a.RefreshTokenService.CreateRefreshToken(userID, newRefreshToken, refreshExpiresAt); err != nil {
		return nil, err
	}
	return &AuthResult{AccessToken: accessToken, RefreshToken: newRefreshToken}, nil
}

// generateRandomToken generates a random string for use as a token
func generateRandomToken() string {
	return RandString(48)
}

// RandString returns a random alphanumeric string of given length
func RandString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[seededRand.Intn(len(letters))]
	}
	return string(b)
}


func dd(v interface{}) {
	fmt.Printf("%#v\n", v)
	os.Exit(1)
}

var seededRand *rand.Rand
var once sync.Once

func init() {
	once.Do(func() {
		seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
	})
}
