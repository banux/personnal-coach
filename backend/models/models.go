package models

import "time"

// EquipmentItem represents a piece of equipment with optional available weights
type EquipmentItem struct {
	Type    string    `json:"type"`
	Weights []float64 `json:"weights,omitempty"` // available weights in kg
}

// Person represents a user of the personal coach app
type Person struct {
	ID             string          `json:"id"`
	Name           string          `json:"name"`
	Age            int             `json:"age"`
	Sex            string          `json:"sex,omitempty"`     // homme, femme, autre
	Weight         float64         `json:"weight"`            // kg
	Height         float64         `json:"height"`            // cm
	Level          string          `json:"level"`             // beginner, intermediate, advanced
	Goals          []string        `json:"goals"`             // weight_loss, muscle_gain, endurance, etc.
	Equipment      []string        `json:"equipment,omitempty"` // backward compat: type names only
	EquipmentItems []EquipmentItem `json:"equipment_items,omitempty"` // equipment with weight details
	Description    string          `json:"description,omitempty"` // additional context (injuries, preferences, etc.)
	CreatedAt      time.Time       `json:"created_at"`
}

// Profile represents a named user in the multi-user system (all share the same app password)
type Profile struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// Exercise represents a single exercise in a workout
type Exercise struct {
	Name           string  `json:"name"`
	Sets           int     `json:"sets"`
	Reps           string  `json:"reps"`           // e.g., "8-12" or "10"
	Intensity      string  `json:"intensity"`      // e.g., "70%RM" or "RPE 7"
	RestSeconds    int     `json:"rest_seconds"`   // rest between sets in seconds
	Tempo          string  `json:"tempo"`          // e.g., "3-1-2-0" (eccentric-pause-concentric-pause)
	Notes          string  `json:"notes"`          // technical instructions
	DurationSecs   int     `json:"duration_secs"`  // for timed exercises (seconds)
	MuscleGroups   []string `json:"muscle_groups"` // targeted muscles
}

// WorkoutBlock groups exercises (e.g., warm-up, main block, accessory)
type WorkoutBlock struct {
	Name      string     `json:"name"`
	Exercises []Exercise `json:"exercises"`
}

// DayProgram represents a single day's workout
type DayProgram struct {
	Day          int            `json:"day"`
	Name         string         `json:"name"`
	Focus        string         `json:"focus"`     // e.g., "Upper body", "Legs", "Full body"
	Duration     int            `json:"duration"`  // estimated minutes
	Blocks       []WorkoutBlock `json:"blocks"`
	WarmupNotes  string         `json:"warmup_notes"`
	CooldownNotes string        `json:"cooldown_notes"`
}

// WeeklyFeedback captures how the user felt the previous week
type WeeklyFeedback struct {
	EnergyLevel    int    `json:"energy_level"`    // 1-10
	SorenessLevel  int    `json:"soreness_level"`  // 1-10
	MotivationLevel int   `json:"motivation_level"` // 1-10
	Notes          string `json:"notes"`
	CompletedDays  int    `json:"completed_days"`
}

// Program represents a complete workout program
type Program struct {
	ID          string         `json:"id"`
	ProfileID   string         `json:"profile_id,omitempty"`
	PersonID    string         `json:"person_id"`
	PersonName  string         `json:"person_name"`
	WeekNumber  int            `json:"week_number"`
	TotalWeeks  int            `json:"total_weeks"`
	Objective   string         `json:"objective"`
	Days        []DayProgram   `json:"days"`
	Feedback    *WeeklyFeedback `json:"feedback,omitempty"`
	GeneratedAt time.Time      `json:"generated_at"`
	Notes       string         `json:"notes"`
}

// TimerSet represents one set in the timer
type TimerSet struct {
	ExerciseName string `json:"exercise_name"`
	SetNumber    int    `json:"set_number"`
	WorkSeconds  int    `json:"work_seconds"`   // 0 for reps-based
	Reps         string `json:"reps"`
	RestSeconds  int    `json:"rest_seconds"`
}

// TimerProgram represents the timer sequence for a day's workout
type TimerProgram struct {
	ProgramID string      `json:"program_id"`
	DayIndex  int         `json:"day_index"`
	DayName   string      `json:"day_name"`
	Sets      []TimerSet  `json:"sets"`
	TotalTime int         `json:"total_time"` // estimated total seconds
}

// GenerateRequest is the API request for generating a program
type GenerateRequest struct {
	Person   Person          `json:"person"`
	Weeks    int             `json:"weeks"`
	DaysPerWeek int          `json:"days_per_week"`
	Feedback *WeeklyFeedback `json:"feedback,omitempty"`
}

// GenerateResponse is the API response after generating a program
type GenerateResponse struct {
	Program Program `json:"program"`
	Message string  `json:"message"`
}
