package services_test

import (
	"testing"

	"personal-coach/models"
	"personal-coach/services"
)

// helper: builds a minimal program with one day
func programWith(day models.DayProgram) models.Program {
	return models.Program{
		ID:   "prog-1",
		Days: []models.DayProgram{day},
	}
}

// helper: builds one exercise
func exercise(name string, sets int, reps string, rest int, tempo string, durationSecs int) models.Exercise {
	return models.Exercise{
		Name:         name,
		Sets:         sets,
		Reps:         reps,
		RestSeconds:  rest,
		Tempo:        tempo,
		DurationSecs: durationSecs,
	}
}

// TestBuildTimer_InvalidDay verifies out-of-range day returns an error
func TestBuildTimer_InvalidDay(t *testing.T) {
	prog := programWith(models.DayProgram{Day: 1, Name: "Lundi"})

	if _, err := services.BuildTimer(prog, -1); err == nil {
		t.Error("expected error for day -1")
	}
	if _, err := services.BuildTimer(prog, 1); err == nil {
		t.Error("expected error for day 1 (only day 0 exists)")
	}
}

// TestBuildTimer_SetCount verifies the correct number of timer sets is produced
func TestBuildTimer_SetCount(t *testing.T) {
	day := models.DayProgram{
		Day:  1,
		Name: "Lundi",
		Blocks: []models.WorkoutBlock{
			{
				Name: "Bloc principal",
				Exercises: []models.Exercise{
					exercise("Squat", 4, "8", 120, "", 0),
					exercise("Presse", 3, "10", 90, "", 0),
				},
			},
		},
	}
	prog := programWith(day)

	timer, err := services.BuildTimer(prog, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 4 sets for Squat + 3 sets for Presse = 7
	if len(timer.Sets) != 7 {
		t.Errorf("expected 7 sets, got %d", len(timer.Sets))
	}
}

// TestBuildTimer_SetNames verifies exercise names are correctly propagated
func TestBuildTimer_SetNames(t *testing.T) {
	day := models.DayProgram{
		Day:  1,
		Name: "Test",
		Blocks: []models.WorkoutBlock{
			{
				Name: "Main",
				Exercises: []models.Exercise{
					exercise("Développé couché", 2, "10", 90, "", 0),
				},
			},
		},
	}

	timer, err := services.BuildTimer(programWith(day), 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for i, s := range timer.Sets {
		if s.ExerciseName != "Développé couché" {
			t.Errorf("set %d: wrong exercise name %q", i, s.ExerciseName)
		}
		if s.SetNumber != i+1 {
			t.Errorf("set %d: wrong set number %d", i, s.SetNumber)
		}
	}
}

// TestBuildTimer_DurationSecs verifies timed exercises use DurationSecs
func TestBuildTimer_DurationSecs(t *testing.T) {
	day := models.DayProgram{
		Day:  1,
		Name: "Cardio",
		Blocks: []models.WorkoutBlock{
			{
				Name: "Bloc",
				Exercises: []models.Exercise{
					exercise("Planche", 3, "1", 60, "", 45),
				},
			},
		},
	}

	timer, err := services.BuildTimer(programWith(day), 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, s := range timer.Sets {
		if s.WorkSeconds != 45 {
			t.Errorf("expected WorkSeconds=45, got %d", s.WorkSeconds)
		}
	}
}

// TestBuildTimer_TempoCalculation verifies tempo-based work duration
func TestBuildTimer_TempoCalculation(t *testing.T) {
	// tempo "3-1-2-0" = 6 sec/rep; 10 reps = 60 sec
	day := models.DayProgram{
		Day:  1,
		Name: "Force",
		Blocks: []models.WorkoutBlock{
			{
				Name: "Main",
				Exercises: []models.Exercise{
					exercise("Romanian DL", 1, "10", 120, "3-1-2-0", 0),
				},
			},
		},
	}

	timer, err := services.BuildTimer(programWith(day), 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(timer.Sets) != 1 {
		t.Fatalf("expected 1 set, got %d", len(timer.Sets))
	}
	// 3+1+2+0 = 6 sec/rep × 10 reps = 60
	if timer.Sets[0].WorkSeconds != 60 {
		t.Errorf("expected WorkSeconds=60 for tempo 3-1-2-0 × 10, got %d", timer.Sets[0].WorkSeconds)
	}
}

// TestBuildTimer_DefaultWorkDuration verifies 3s/rep default when no tempo
func TestBuildTimer_DefaultWorkDuration(t *testing.T) {
	day := models.DayProgram{
		Day:  1,
		Name: "Day",
		Blocks: []models.WorkoutBlock{
			{
				Name: "B",
				Exercises: []models.Exercise{
					exercise("Curl", 1, "12", 60, "", 0),
				},
			},
		},
	}

	timer, err := services.BuildTimer(programWith(day), 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// default: 12 reps × 3s = 36
	if timer.Sets[0].WorkSeconds != 36 {
		t.Errorf("expected WorkSeconds=36, got %d", timer.Sets[0].WorkSeconds)
	}
}

// TestBuildTimer_WarmupCooldownAdded verifies warmup/cooldown adds to totalTime
func TestBuildTimer_WarmupCooldownAdded(t *testing.T) {
	day := models.DayProgram{
		Day:           1,
		Name:          "Day",
		WarmupNotes:   "5 min vélo",
		CooldownNotes: "Étirements",
		Blocks: []models.WorkoutBlock{
			{
				Name: "B",
				Exercises: []models.Exercise{
					exercise("Squat", 1, "5", 60, "", 0),
				},
			},
		},
	}

	// Without warmup/cooldown
	dayNoNotes := day
	dayNoNotes.WarmupNotes = ""
	dayNoNotes.CooldownNotes = ""
	timerNoNotes, _ := services.BuildTimer(programWith(dayNoNotes), 0)

	// With warmup + cooldown
	timerWithNotes, _ := services.BuildTimer(programWith(day), 0)

	diff := timerWithNotes.TotalTime - timerNoNotes.TotalTime
	if diff != 900 { // 600 + 300
		t.Errorf("expected 900s extra for warmup+cooldown, got %d", diff)
	}
}

// TestBuildTimer_DayName verifies DayName is set correctly
func TestBuildTimer_DayName(t *testing.T) {
	day := models.DayProgram{Day: 1, Name: "Mercredi", Blocks: []models.WorkoutBlock{}}
	timer, err := services.BuildTimer(programWith(day), 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if timer.DayName != "Mercredi" {
		t.Errorf("expected DayName=Mercredi, got %q", timer.DayName)
	}
	if timer.ProgramID != "prog-1" {
		t.Errorf("expected ProgramID=prog-1, got %q", timer.ProgramID)
	}
}
