package dbutil

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// OpenSqliteDB opens a SQLite database with WAL mode and busy timeout for concurrency.
func OpenSqliteDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	// Set WAL mode for better concurrency
	_, err = db.Exec("PRAGMA journal_mode=WAL; PRAGMA synchronous=1; PRAGMA mmap_size=134217728; PRAGMA journal_size_limit=67108864; PRAGMA cache_size=2000;")
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to set WAL mode: %w", err)
	}
	// Set busy timeout to 5 seconds
	_, err = db.Exec("PRAGMA busy_timeout=5000;")
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to set busy timeout: %w", err)
	}
	// Optional: tune connection pool
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	return db, nil
}
