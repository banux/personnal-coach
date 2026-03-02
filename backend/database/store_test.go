package database_test

import (
	"testing"
	"time"

	"personal-coach/database"
	"personal-coach/models"
)

func openTestDB(t *testing.T) *database.DB {
	t.Helper()
	db, err := database.Open(":memory:")
	if err != nil {
		t.Fatalf("open test db: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

func sampleProgram(id string) models.Program {
	return models.Program{
		ID:          id,
		PersonID:    "person-1",
		PersonName:  "Alice",
		WeekNumber:  1,
		TotalWeeks:  4,
		Objective:   "Prise de masse",
		GeneratedAt: time.Now().UTC().Truncate(time.Second),
		Days: []models.DayProgram{
			{
				Day:    1,
				Name:   "Lundi",
				Focus:  "Pectoraux",
				Duration: 60,
				Blocks: []models.WorkoutBlock{
					{
						Name: "Bloc principal",
						Exercises: []models.Exercise{
							{
								Name:        "Développé couché",
								Sets:        4,
								Reps:        "8-10",
								Intensity:   "75%RM",
								RestSeconds: 120,
								Tempo:       "3-1-2-0",
								Notes:       "Garder les omoplates rétractées",
								MuscleGroups: []string{"Pectoraux", "Triceps"},
							},
						},
					},
				},
			},
		},
	}
}

func TestSaveAndGetProgram(t *testing.T) {
	db := openTestDB(t)
	p := sampleProgram("prog-1")

	if err := db.SaveProgram(p); err != nil {
		t.Fatalf("save: %v", err)
	}

	got, err := db.GetProgram("prog-1")
	if err != nil {
		t.Fatalf("get: %v", err)
	}

	if got.ID != p.ID {
		t.Errorf("id: got %q, want %q", got.ID, p.ID)
	}
	if got.PersonName != p.PersonName {
		t.Errorf("person_name: got %q, want %q", got.PersonName, p.PersonName)
	}
	if len(got.Days) != 1 {
		t.Errorf("days: got %d, want 1", len(got.Days))
	}
	if len(got.Days[0].Blocks[0].Exercises) != 1 {
		t.Errorf("exercises: got %d, want 1", len(got.Days[0].Blocks[0].Exercises))
	}
}

func TestGetProgramNotFound(t *testing.T) {
	db := openTestDB(t)

	_, err := db.GetProgram("nonexistent")
	if err != database.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestListPrograms(t *testing.T) {
	db := openTestDB(t)

	// Empty list
	programs, err := db.ListPrograms()
	if err != nil {
		t.Fatalf("list empty: %v", err)
	}
	if len(programs) != 0 {
		t.Errorf("expected 0 programs, got %d", len(programs))
	}

	// Save two programs
	if err := db.SaveProgram(sampleProgram("p1")); err != nil {
		t.Fatal(err)
	}
	if err := db.SaveProgram(sampleProgram("p2")); err != nil {
		t.Fatal(err)
	}

	programs, err = db.ListPrograms()
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(programs) != 2 {
		t.Errorf("expected 2 programs, got %d", len(programs))
	}
}

func TestDeleteProgram(t *testing.T) {
	db := openTestDB(t)

	if err := db.SaveProgram(sampleProgram("del-1")); err != nil {
		t.Fatal(err)
	}

	if err := db.DeleteProgram("del-1"); err != nil {
		t.Fatalf("delete: %v", err)
	}

	_, err := db.GetProgram("del-1")
	if err != database.ErrNotFound {
		t.Errorf("after delete: expected ErrNotFound, got %v", err)
	}
}

func TestDeleteProgramNotFound(t *testing.T) {
	db := openTestDB(t)
	err := db.DeleteProgram("ghost")
	if err != database.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestMigrationsIdempotent(t *testing.T) {
	// Opening the same :memory: db twice would be different connections,
	// so test idempotency by reopening a file-based db in a temp dir.
	t.TempDir() // just verify Open doesn't fail if run twice on the same dir
	tmp := t.TempDir()
	db1, err := database.Open(tmp)
	if err != nil {
		t.Fatal(err)
	}
	db1.Close()

	db2, err := database.Open(tmp)
	if err != nil {
		t.Fatalf("second open: %v", err)
	}
	db2.Close()
}
