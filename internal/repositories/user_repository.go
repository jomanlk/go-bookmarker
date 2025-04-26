package repositories

import (
	"bookmarker/internal/models"
	"database/sql"
	"errors"
	"time"
)

type UserRepository interface {
	CreateUser(username, passwordHash string) (models.User, error)
	GetUserByUsername(username string) (models.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(username, passwordHash string) (models.User, error) {
	createdAt := time.Now().Unix()
	result, err := r.db.Exec(`INSERT INTO users (username, password_hash, created_at, updated_at) VALUES (?, ?, ?, ?)`, username, passwordHash, createdAt, createdAt)
	if err != nil {
		return models.User{}, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return models.User{}, err
	}
	return models.User{
		ID:           id,
		Username:     username,
		PasswordHash: passwordHash,
		CreatedAt:    createdAt,
		UpdatedAt:    createdAt,
	}, nil
}

func (r *userRepository) GetUserByUsername(username string) (models.User, error) {
	row := r.db.QueryRow(`SELECT id, username, password_hash, created_at, updated_at FROM users WHERE username = ?`, username)
	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return user, errors.New("user not found")
	}
	return user, err
}
