package repositories

import (
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

func (r *TokenRepository) FindByToken(token string) (int, error) {
	var userID int
	err := r.DB.QueryRow("SELECT user_id FROM tokens WHERE token = ? AND expires_at > ?", token, time.Now().Unix()).Scan(&userID)
	if (err != nil) {
		return 0, err
	}
	return userID, nil
}

func (r *TokenRepository) DeleteToken(token string) error {
	_, err := r.DB.Exec("DELETE FROM tokens WHERE token = ?", token)
	return err
}
