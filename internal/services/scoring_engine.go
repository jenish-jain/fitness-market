package services

import (
	"fitness-market/internal/database"
	"fitness-market/internal/models"
)

// CalculateEntryScore calculates the score for a workout entry
// Score formula: weight * reps * sets * bodyweight_factor
func CalculateEntryScore(weight float64, reps int, sets int, userID uint) float64 {
	if weight <= 0 || reps <= 0 || sets <= 0 {
		return 0
	}

	rawScore := weight * float64(reps) * float64(sets)
	normalizedScore := NormalizeScore(rawScore, userID)

	return normalizedScore
}

// DetectPR checks if the current entry is a personal record for the exercise
func DetectPR(userID uint, exerciseID uint, weight float64, reps int) bool {
	if weight <= 0 || reps <= 0 {
		return false
	}

	var existingPR models.ExercisePR
	var exercise models.Exercise

	// Get exercise name
	if err := database.DB.First(&exercise, exerciseID).Error; err != nil {
		return false
	}

	// Check if there's an existing PR for this exercise
	err := database.DB.Where("user_id = ? AND exercise_name = ?", userID, exercise.Name).First(&existingPR).Error

	if err != nil {
		// No existing PR, this is a new PR
		return true
	}

	// Compare: PR if weight is higher at same or more reps, or same weight with more reps
	if weight > existingPR.Weight {
		return true
	}
	if weight == existingPR.Weight && reps > existingPR.Reps {
		return true
	}

	return false
}

// UpdatePRIfNeeded updates the PR record if the entry is a new PR
func UpdatePRIfNeeded(userID uint, exerciseID uint, weight float64, reps int, date interface{}) error {
	var exercise models.Exercise
	if err := database.DB.First(&exercise, exerciseID).Error; err != nil {
		return err
	}

	var existingPR models.ExercisePR
	err := database.DB.Where("user_id = ? AND exercise_name = ?", userID, exercise.Name).First(&existingPR).Error

	if err != nil {
		// Create new PR
		newPR := models.ExercisePR{
			UserID:       userID,
			ExerciseName: exercise.Name,
			Weight:       weight,
			Reps:         reps,
			Unit:         "kg",
		}
		if t, ok := date.(interface{ Unix() int64 }); ok {
			_ = t
		}
		return database.DB.Create(&newPR).Error
	}

	// Update existing PR
	existingPR.Weight = weight
	existingPR.Reps = reps
	return database.DB.Save(&existingPR).Error
}
