package dbutil

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// OpenSqliteDB opens a SQLite database with WAL mode and busy timeout for concurrency.
func OpenSqliteDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "file:../../internal/db/bookmarker_db1.db?_journal_mode=WAL&_busy_timeout=5000&_sync=NORMAL")
	if err != nil {
		return nil, err
	}
	// Set WAL mode for better concurrency
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	return db, nil
}

// ShutdownSqliteDB cleanly closes the SQLite database connection.
func ShutdownSqliteDB(db *sql.DB) error {
	if db != nil {
		return db.Close()
	}
	return nil
}
