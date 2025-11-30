package handlers

import (
	"net/http"
	"time"

	"fitness-market/internal/database"
	"fitness-market/internal/models"
	"fitness-market/internal/services"

	"github.com/gin-gonic/gin"
)

type CreateEntryRequest struct {
	ExerciseID uint    `json:"exercise_id" binding:"required"`
	Weight     float64 `json:"weight" binding:"required"`
	Reps       int     `json:"reps" binding:"required"`
	Sets       int     `json:"sets" binding:"required"`
	Notes      string  `json:"notes"`
	Date       string  `json:"date"`
}

type EntryResponse struct {
	ID         uint      `json:"id"`
	UserID     uint      `json:"user_id"`
	ExerciseID uint      `json:"exercise_id"`
	Weight     float64   `json:"weight"`
	Reps       int       `json:"reps"`
	Sets       int       `json:"sets"`
	Notes      string    `json:"notes"`
	Date       time.Time `json:"date"`
	Score      float64   `json:"score"`
	IsPR       bool      `json:"is_pr"`
	CreatedAt  time.Time `json:"created_at"`
}

// CreateEntry handles POST /api/v1/entries
func CreateEntry(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req CreateEntryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate exercise exists
	var exercise models.Exercise
	if err := database.DB.First(&exercise, req.ExerciseID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Exercise not found"})
		return
	}

	// Parse date or use current time
	entryDate := time.Now()
	if req.Date != "" {
		parsedDate, err := time.Parse("2006-01-02", req.Date)
		if err == nil {
			entryDate = parsedDate
		}
	}

	uid := userID.(uint)

	// Calculate score using scoring engine
	score := services.CalculateEntryScore(req.Weight, req.Reps, req.Sets, uid)

	// Detect if this is a PR
	isPR := services.DetectPR(uid, req.ExerciseID, req.Weight, req.Reps)

	// Create the workout entry
	entry := models.WorkoutEntry{
		UserID:     uid,
		ExerciseID: req.ExerciseID,
		Date:       entryDate,
		Sets:       req.Sets,
		Reps:       req.Reps,
		Weight:     req.Weight,
		Notes:      req.Notes,
		Score:      score,
		IsPR:       isPR,
	}

	if err := database.DB.Create(&entry).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create entry"})
		return
	}

	// Update PR record if this is a new PR
	if isPR {
		services.UpdatePRIfNeeded(uid, req.ExerciseID, req.Weight, req.Reps, entryDate)
	}

	response := EntryResponse{
		ID:         entry.ID,
		UserID:     entry.UserID,
		ExerciseID: entry.ExerciseID,
		Weight:     entry.Weight,
		Reps:       entry.Reps,
		Sets:       entry.Sets,
		Notes:      entry.Notes,
		Date:       entry.Date,
		Score:      entry.Score,
		IsPR:       entry.IsPR,
		CreatedAt:  entry.CreatedAt,
	}

	c.JSON(http.StatusCreated, response)
}
