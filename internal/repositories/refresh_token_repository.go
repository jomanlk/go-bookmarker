package repositories

import (
	"bookmarker/internal/models"
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RefreshTokenRepository struct {
	DB *pgxpool.Pool
}

func NewRefreshTokenRepository(db *pgxpool.Pool) *RefreshTokenRepository {
	return &RefreshTokenRepository{DB: db}
}

func (r *RefreshTokenRepository) CreateRefreshToken(userID int, refreshToken string, expiresAt time.Time) error {
	createdAt := time.Now().UTC()
	_, err := r.DB.Exec(context.Background(),
		"INSERT INTO refresh_tokens (user_id, refresh_token, expires_at, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)",
		userID, refreshToken, expiresAt, createdAt, createdAt,
	)
	return err
}

func (r *RefreshTokenRepository) FindByToken(refreshToken string) (*models.RefreshToken, error) {
	var t models.RefreshToken
	err := r.DB.QueryRow(context.Background(),
		"SELECT id, user_id, refresh_token, expires_at, created_at, updated_at FROM refresh_tokens WHERE refresh_token = $1",
		refreshToken,
	).Scan(&t.ID, &t.UserID, &t.RefreshToken, &t.ExpiresAt, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *RefreshTokenRepository) DeleteRefreshToken(refreshToken string) error {
	_, err := r.DB.Exec(context.Background(), "DELETE FROM refresh_tokens WHERE refresh_token = $1", refreshToken)
	return err
}
