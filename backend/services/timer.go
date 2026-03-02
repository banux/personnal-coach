package services

import (
	"fmt"
	"personal-coach/models"
)

// BuildTimer generates a timer sequence from a day program
func BuildTimer(program models.Program, dayIndex int) (*models.TimerProgram, error) {
	if dayIndex < 0 || dayIndex >= len(program.Days) {
		return nil, fmt.Errorf("invalid day index: %d (program has %d days)", dayIndex, len(program.Days))
	}

	day := program.Days[dayIndex]
	var sets []models.TimerSet
	totalTime := 0

	for _, block := range day.Blocks {
		for _, exercise := range block.Exercises {
			for setNum := 1; setNum <= exercise.Sets; setNum++ {
				timerSet := models.TimerSet{
					ExerciseName: exercise.Name,
					SetNumber:    setNum,
					Reps:         exercise.Reps,
					RestSeconds:  exercise.RestSeconds,
				}

				// If exercise has duration (timed), use it; otherwise estimate from reps
				if exercise.DurationSecs > 0 {
					timerSet.WorkSeconds = exercise.DurationSecs
				} else {
					// Estimate ~3 seconds per rep for time calculation
					timerSet.WorkSeconds = estimateWorkDuration(exercise.Reps, exercise.Tempo)
				}

				sets = append(sets, timerSet)
				totalTime += timerSet.WorkSeconds

				// Add rest time (except after last set of last exercise in a block)
				if setNum < exercise.Sets || isLastExercise(block.Exercises, exercise) {
					totalTime += exercise.RestSeconds
				}
			}
		}
	}

	// Add estimated warm-up time
	if day.WarmupNotes != "" {
		totalTime += 600 // 10 min warm-up
	}
	if day.CooldownNotes != "" {
		totalTime += 300 // 5 min cooldown
	}

	return &models.TimerProgram{
		ProgramID: program.ID,
		DayIndex:  dayIndex,
		DayName:   day.Name,
		Sets:      sets,
		TotalTime: totalTime,
	}, nil
}

// estimateWorkDuration estimates work duration from reps string and tempo
func estimateWorkDuration(reps string, tempo string) int {
	// Parse reps - handle ranges like "8-12"
	maxReps := 10 // default
	fmt.Sscanf(reps, "%d", &maxReps)

	// If tempo is specified, calculate from it
	if tempo != "" {
		var eccentric, pauseBottom, concentric, pauseTop int
		n, _ := fmt.Sscanf(tempo, "%d-%d-%d-%d", &eccentric, &pauseBottom, &concentric, &pauseTop)
		if n == 4 {
			secPerRep := eccentric + pauseBottom + concentric + pauseTop
			return secPerRep * maxReps
		}
	}

	// Default: ~3 seconds per rep
	return maxReps * 3
}

// isLastExercise checks if the exercise is the last in the slice
func isLastExercise(exercises []models.Exercise, ex models.Exercise) bool {
	if len(exercises) == 0 {
		return true
	}
	last := exercises[len(exercises)-1]
	return last.Name == ex.Name
}
