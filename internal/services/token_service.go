package services

import (
	"bookmarker/internal/repositories"
	"time"
)

type TokenService struct {
	TokenRepo *repositories.TokenRepository
}

func NewTokenService(tokenRepo *repositories.TokenRepository) *TokenService {
	return &TokenService{TokenRepo: tokenRepo}
}

func (s *TokenService) CreateToken(userID int, token string, expiresAt time.Time) error {
	return s.TokenRepo.CreateToken(userID, token, expiresAt)
}

func (s *TokenService) FindByToken(token string) (int, error) {
	return s.TokenRepo.FindByToken(token)
}

func (s *TokenService) DeleteToken(token string) error {
	return s.TokenRepo.DeleteToken(token)
}
