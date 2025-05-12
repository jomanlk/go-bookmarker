package repositories

import (
	"bookmarker/internal/models"
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// TokenRepository handles token-related DB operations
type TokenRepository struct {
	DB *pgxpool.Pool
}

func NewTokenRepository(db *pgxpool.Pool) *TokenRepository {
	return &TokenRepository{DB: db}
}

func (r *TokenRepository) CreateToken(userID int, token string, expiresAt time.Time) error {
	createdAt := time.Now().UTC()
	_, err := r.DB.Exec(context.Background(),
		"INSERT INTO tokens (user_id, token, expires_at, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)",
		userID, token, expiresAt, createdAt, createdAt,
	)
	return err
}

func (r *TokenRepository) FindByToken(token string) (*models.Token, error) {
	var t models.Token
	err := r.DB.QueryRow(context.Background(),
		"SELECT id, user_id, token, expires_at, created_at, updated_at FROM tokens WHERE token = $1",
		token,
	).Scan(&t.ID, &t.UserID, &t.Token, &t.ExpiresAt, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TokenRepository) DeleteToken(token string) error {
	_, err := r.DB.Exec(context.Background(), "DELETE FROM tokens WHERE token = $1", token)
	return err
}
