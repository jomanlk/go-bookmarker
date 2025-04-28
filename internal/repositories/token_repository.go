package repositories

import (
	"bookmarker/internal/models"
	"database/sql"
	"time"
)

// TokenRepository handles token-related DB operations
type TokenRepository struct {
	DB *sql.DB
}

func NewTokenRepository(db *sql.DB) *TokenRepository {
	return &TokenRepository{DB: db}
}

func (r *TokenRepository) CreateToken(userID int, token string, expiresAt time.Time) error {
	_, err := r.DB.Exec("INSERT INTO tokens (user_id, token, expires_at, created_at, updated_at) VALUES (?, ?, ?, ?, ?)", userID, token, expiresAt.Unix(), time.Now().Unix(), time.Now().Unix())
	return err
}

func (r *TokenRepository) FindByToken(token string) (*models.Token, error) {
	var t models.Token
	err := r.DB.QueryRow("SELECT id, user_id, token, expires_at, created_at, updated_at FROM tokens WHERE token = ?", token).Scan(
		&t.ID, &t.UserID, &t.Token, &t.ExpiresAt, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TokenRepository) DeleteToken(token string) error {
	_, err := r.DB.Exec("DELETE FROM tokens WHERE token = ?", token)
	return err
}
