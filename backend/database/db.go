// Package database handles SQLite persistence for the personal coach application.
// It uses modernc.org/sqlite (pure Go, CGO-free) so it works in Alpine containers.
package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

// DB wraps the sql.DB connection
type DB struct {
	*sql.DB
}

// Open initialises the SQLite database at the given path (or :memory: for tests).
// It runs all pending migrations before returning.
func Open(dataDir string) (*DB, error) {
	var dsn string
	if dataDir == ":memory:" {
		dsn = ":memory:"
	} else {
		if err := os.MkdirAll(dataDir, 0755); err != nil {
			return nil, fmt.Errorf("create data dir: %w", err)
		}
		dsn = filepath.Join(dataDir, "coach.db")
	}

	sqlDB, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}

	// SQLite should only have one writer at a time
	sqlDB.SetMaxOpenConns(1)

	db := &DB{sqlDB}
	if err := db.migrate(); err != nil {
		sqlDB.Close()
		return nil, fmt.Errorf("migrate: %w", err)
	}

	log.Printf("Database opened: %s", dsn)
	return db, nil
}
