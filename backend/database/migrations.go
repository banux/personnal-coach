package database

import (
	"fmt"
	"log"
)

// migration represents a versioned SQL migration
type migration struct {
	version int
	sql     string
}

// migrations is the ordered list of schema migrations.
// NEVER modify an existing migration — always add a new one.
var migrations = []migration{
	{
		version: 1,
		sql: `
-- Programs table: stores each generated workout program as JSON
CREATE TABLE IF NOT EXISTS programs (
	id          TEXT    PRIMARY KEY,
	person_id   TEXT    NOT NULL DEFAULT '',
	person_name TEXT    NOT NULL DEFAULT '',
	week_number INTEGER NOT NULL DEFAULT 1,
	total_weeks INTEGER NOT NULL DEFAULT 1,
	objective   TEXT    NOT NULL DEFAULT '',
	payload     TEXT    NOT NULL,  -- full program JSON
	created_at  TEXT    NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_programs_person_id   ON programs(person_id);
CREATE INDEX IF NOT EXISTS idx_programs_person_name ON programs(person_name);
CREATE INDEX IF NOT EXISTS idx_programs_created_at  ON programs(created_at DESC);
`,
	},
	{
		version: 2,
		sql: `
-- Profiles table: named users sharing the same app password
CREATE TABLE IF NOT EXISTS profiles (
	id         TEXT PRIMARY KEY,
	name       TEXT NOT NULL UNIQUE,
	created_at TEXT NOT NULL DEFAULT (datetime('now'))
);

-- Add profile_id to programs for per-user filtering
ALTER TABLE programs ADD COLUMN profile_id TEXT NOT NULL DEFAULT '';
CREATE INDEX IF NOT EXISTS idx_programs_profile_id ON programs(profile_id);
`,
	},
}

// migrate applies all pending migrations using a simple version-tracking table.
func (db *DB) migrate() error {
	// Ensure the schema_version table exists
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_version (
			version INTEGER PRIMARY KEY,
			applied_at TEXT NOT NULL DEFAULT (datetime('now'))
		)
	`)
	if err != nil {
		return fmt.Errorf("create schema_version: %w", err)
	}

	for _, m := range migrations {
		var count int
		if err := db.QueryRow(`SELECT COUNT(*) FROM schema_version WHERE version = ?`, m.version).Scan(&count); err != nil {
			return fmt.Errorf("check migration %d: %w", m.version, err)
		}
		if count > 0 {
			continue // already applied
		}

		log.Printf("Applying migration v%d", m.version)

		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("begin tx for migration %d: %w", m.version, err)
		}

		if _, err := tx.Exec(m.sql); err != nil {
			tx.Rollback()
			return fmt.Errorf("exec migration %d: %w", m.version, err)
		}

		if _, err := tx.Exec(`INSERT INTO schema_version (version) VALUES (?)`, m.version); err != nil {
			tx.Rollback()
			return fmt.Errorf("record migration %d: %w", m.version, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit migration %d: %w", m.version, err)
		}

		log.Printf("Migration v%d applied", m.version)
	}

	return nil
}
