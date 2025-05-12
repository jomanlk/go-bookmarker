package repositories

import (
	"bookmarker/internal/models"
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// UserRepository handles user-related DB operations
// Now uses pgxpool.Pool
//
type UserRepository struct {
	DB *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) CreateUser(username, passwordHash string) (models.User, error) {
	createdAt := time.Now().UTC()
	var id int64
	err := r.DB.QueryRow(context.Background(),
		`INSERT INTO users (username, password_hash, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id`,
		username, passwordHash, createdAt, createdAt,
	).Scan(&id)
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

func (r *UserRepository) GetUserByUsername(username string) (models.User, error) {
	row := r.DB.QueryRow(context.Background(),
		`SELECT id, username, password_hash, created_at, updated_at FROM users WHERE username = $1`, username)
	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return user, errors.New("user not found")
		}
	}
	return user, err
}

func (r *UserRepository) GetUserByID(id int64) (models.User, error) {
	row := r.DB.QueryRow(context.Background(),
		`SELECT id, username, password_hash, created_at, updated_at FROM users WHERE id = $1`, id)
	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return user, errors.New("user not found")
		}
	}
	return user, err
}
