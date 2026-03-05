package database

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"personal-coach/models"
)

// ErrNotFound is returned when a requested record does not exist.
var ErrNotFound = errors.New("not found")

// SaveProgram inserts or replaces a program in the database.
func (db *DB) SaveProgram(p models.Program) error {
	payload, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("marshal program: %w", err)
	}

	_, err = db.Exec(`
		INSERT OR REPLACE INTO programs
			(id, profile_id, person_id, person_name, week_number, total_weeks, objective, payload, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		p.ID,
		p.ProfileID,
		p.PersonID,
		p.PersonName,
		p.WeekNumber,
		p.TotalWeeks,
		p.Objective,
		string(payload),
		p.GeneratedAt.UTC().Format(time.RFC3339),
	)
	return err
}

// GetProgram fetches a program by ID. Returns ErrNotFound if absent.
func (db *DB) GetProgram(id string) (*models.Program, error) {
	var payload string
	err := db.QueryRow(`SELECT payload FROM programs WHERE id = ?`, id).Scan(&payload)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("query program: %w", err)
	}

	var p models.Program
	if err := json.Unmarshal([]byte(payload), &p); err != nil {
		return nil, fmt.Errorf("unmarshal program: %w", err)
	}
	return &p, nil
}

// ListPrograms returns programs ordered by creation date (newest first).
// If profileID is non-empty, only programs for that profile are returned.
func (db *DB) ListPrograms(profileID string) ([]models.Program, error) {
	var rows *sql.Rows
	var err error
	if profileID != "" {
		rows, err = db.Query(`
			SELECT payload FROM programs
			WHERE profile_id = ?
			ORDER BY created_at DESC
		`, profileID)
	} else {
		rows, err = db.Query(`
			SELECT payload FROM programs
			ORDER BY created_at DESC
		`)
	}
	if err != nil {
		return nil, fmt.Errorf("query programs: %w", err)
	}
	defer rows.Close()

	var programs []models.Program
	for rows.Next() {
		var payload string
		if err := rows.Scan(&payload); err != nil {
			return nil, fmt.Errorf("scan program: %w", err)
		}
		var p models.Program
		if err := json.Unmarshal([]byte(payload), &p); err != nil {
			return nil, fmt.Errorf("unmarshal program: %w", err)
		}
		programs = append(programs, p)
	}
	return programs, rows.Err()
}

// DeleteProgram removes a program by ID.
func (db *DB) DeleteProgram(id string) error {
	res, err := db.Exec(`DELETE FROM programs WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete program: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

// SaveProfile inserts a new profile. Returns ErrNotFound if a profile with the same name exists.
func (db *DB) SaveProfile(p *models.Profile) error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	if p.CreatedAt.IsZero() {
		p.CreatedAt = time.Now()
	}
	_, err := db.Exec(
		`INSERT INTO profiles (id, name, created_at) VALUES (?, ?, ?)`,
		p.ID, p.Name, p.CreatedAt.UTC().Format(time.RFC3339),
	)
	return err
}

// GetProfile fetches a profile by ID.
func (db *DB) GetProfile(id string) (*models.Profile, error) {
	var p models.Profile
	var createdAt string
	err := db.QueryRow(`SELECT id, name, created_at FROM profiles WHERE id = ?`, id).
		Scan(&p.ID, &p.Name, &createdAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("query profile: %w", err)
	}
	p.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	return &p, nil
}

// ListProfiles returns all profiles ordered by creation date.
func (db *DB) ListProfiles() ([]models.Profile, error) {
	rows, err := db.Query(`SELECT id, name, created_at FROM profiles ORDER BY created_at ASC`)
	if err != nil {
		return nil, fmt.Errorf("query profiles: %w", err)
	}
	defer rows.Close()

	var profiles []models.Profile
	for rows.Next() {
		var p models.Profile
		var createdAt string
		if err := rows.Scan(&p.ID, &p.Name, &createdAt); err != nil {
			return nil, fmt.Errorf("scan profile: %w", err)
		}
		p.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		profiles = append(profiles, p)
	}
	return profiles, rows.Err()
}
