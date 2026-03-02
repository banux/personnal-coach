package database

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

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
			(id, person_id, person_name, week_number, total_weeks, objective, payload, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		p.ID,
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

// ListPrograms returns all programs ordered by creation date (newest first).
func (db *DB) ListPrograms() ([]models.Program, error) {
	rows, err := db.Query(`
		SELECT payload FROM programs
		ORDER BY created_at DESC
	`)
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
