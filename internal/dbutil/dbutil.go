package dbutil

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// OpenPostgresDB opens a PostgreSQL database using pgx and a connection string from the environment.
func OpenPostgresDB() (*pgxpool.Pool, error) {
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")
	if user == "" || pass == "" || host == "" || port == "" || dbname == "" {
		return nil, errors.New("database environment variables not set")
	}
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=require", user, pass, host, port, dbname)
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, err
	}
	return pool, nil
}

// ShutdownPostgresDB cleanly closes the PostgreSQL database connection.
func ShutdownPostgresDB(pool *pgxpool.Pool) error {
	if pool != nil {
		pool.Close()
	}
	return nil
}
