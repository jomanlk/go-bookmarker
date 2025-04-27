package repositories

import (
	"database/sql"
	"time"
)

type RefreshTokenRepository struct {
	DB *sql.DB
}

func NewRefreshTokenRepository(db *sql.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{DB: db}
}

func (r *RefreshTokenRepository) CreateRefreshToken(userID int, refreshToken string, expiresAt time.Time) error {
	_, err := r.DB.Exec("INSERT INTO refresh_tokens (user_id, refresh_token, expires_at, created_at, updated_at) VALUES (?, ?, ?, ?, ?)", userID, refreshToken, expiresAt.Unix(), time.Now().Unix(), time.Now().Unix())
	return err
}

func (r *RefreshTokenRepository) FindByToken(refreshToken string) (int, error) {
	var userID int
	err := r.DB.QueryRow("SELECT user_id FROM refresh_tokens WHERE refresh_token = ? AND expires_at > ?", refreshToken, time.Now().Unix()).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func (r *RefreshTokenRepository) DeleteRefreshToken(refreshToken string) error {
	_, err := r.DB.Exec("DELETE FROM refresh_tokens WHERE refresh_token = ?", refreshToken)
	return err
}
