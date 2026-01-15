package services

import (
	"fitness-market/internal/database"
	"fitness-market/internal/models"
	"time"
)

// PRDetectionResult contains the result of PR detection
type PRDetectionResult struct {
	IsPR            bool    `json:"is_pr"`
	Score           float64 `json:"score"`
	PreviousBest    float64 `json:"previous_best,omitempty"`
	Improvement     float64 `json:"improvement,omitempty"`
	CelebrationText string  `json:"celebration_text,omitempty"`
}

// CalculateScore calculates a performance score based on weight, reps, and sets
// Using a standard formula: Score = Weight * Reps * Sets
func CalculateScore(weight float64, reps int, sets int) float64 {
	return weight * float64(reps) * float64(sets)
}

// DetectPR compares a new entry against historical data to detect personal records
func DetectPR(userID uint, exerciseID uint, weight float64, reps int, sets int) (*PRDetectionResult, error) {
	db := database.GetDB()
	newScore := CalculateScore(weight, reps, sets)

	result := &PRDetectionResult{
		IsPR:  false,
		Score: newScore,
	}

	// Find the current best PR for this exercise
	var currentBestPR models.PRHistory
	err := db.Where("user_id = ? AND exercise_id = ?", userID, exerciseID).
		Order("score DESC").
		First(&currentBestPR).Error

	if err != nil {
		// No previous PR exists, this is the first entry - it's a PR!
		result.IsPR = true
		result.CelebrationText = "🎉 First record set! Keep going!"
		return result, nil
	}

	// Compare against previous best
	result.PreviousBest = currentBestPR.Score

	if newScore > currentBestPR.Score {
		result.IsPR = true
		result.Improvement = newScore - currentBestPR.Score
		result.CelebrationText = "🏆 New Personal Record! You beat your previous best!"
	}

	return result, nil
}

// RecordPR saves a new PR to the PR history
func RecordPR(userID uint, exerciseID uint, workoutEntryID uint, weight float64, reps int, sets int, achievedAt time.Time) error {
	db := database.GetDB()

	prHistory := models.PRHistory{
		UserID:         userID,
		ExerciseID:     exerciseID,
		WorkoutEntryID: workoutEntryID,
		Score:          CalculateScore(weight, reps, sets),
		Weight:         weight,
		Reps:           reps,
		Sets:           sets,
		AchievedAt:     achievedAt,
	}

	return db.Create(&prHistory).Error
}

// GetPRHistory returns the PR history for a specific exercise
func GetPRHistory(userID uint, exerciseID uint) ([]models.PRHistory, error) {
	db := database.GetDB()
	var history []models.PRHistory

	err := db.Where("user_id = ? AND exercise_id = ?", userID, exerciseID).
		Order("achieved_at DESC").
		Find(&history).Error

	return history, err
}

// GetAllPRs returns the current best PR for each exercise for a user
func GetAllPRs(userID uint) ([]models.PRHistory, error) {
	db := database.GetDB()
	var prs []models.PRHistory

	// Get the best PR for each exercise using a subquery
	subquery := db.Model(&models.PRHistory{}).
		Select("exercise_id, MAX(score) as max_score").
		Where("user_id = ?", userID).
		Group("exercise_id")

	err := db.Model(&models.PRHistory{}).
		Joins("JOIN (?) as best ON pr_history.exercise_id = best.exercise_id AND pr_history.score = best.max_score", subquery).
		Where("pr_history.user_id = ?", userID).
		Preload("Exercise").
		Find(&prs).Error

	return prs, err
}
