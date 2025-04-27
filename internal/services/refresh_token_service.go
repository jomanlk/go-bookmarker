package services

import (
	"bookmarker/internal/repositories"
	"time"
)

type RefreshTokenService struct {
	RefreshTokenRepo *repositories.RefreshTokenRepository
}

func NewRefreshTokenService(refreshTokenRepo *repositories.RefreshTokenRepository) *RefreshTokenService {
	return &RefreshTokenService{RefreshTokenRepo: refreshTokenRepo}
}

func (s *RefreshTokenService) CreateRefreshToken(userID int, refreshToken string, expiresAt time.Time) error {
	return s.RefreshTokenRepo.CreateRefreshToken(userID, refreshToken, expiresAt)
}

func (s *RefreshTokenService) FindByToken(refreshToken string) (int, error) {
	return s.RefreshTokenRepo.FindByToken(refreshToken)
}

func (s *RefreshTokenService) DeleteRefreshToken(refreshToken string) error {
	return s.RefreshTokenRepo.DeleteRefreshToken(refreshToken)
}
