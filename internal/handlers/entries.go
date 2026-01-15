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
	ID              uint      `json:"id"`
	UserID          uint      `json:"user_id"`
	ExerciseID      uint      `json:"exercise_id"`
	Weight          float64   `json:"weight"`
	Reps            int       `json:"reps"`
	Sets            int       `json:"sets"`
	Notes           string    `json:"notes"`
	Date            time.Time `json:"date"`
	Score           float64   `json:"score"`
	IsPR            bool      `json:"is_pr"`
	CelebrationText string    `json:"celebration_text,omitempty"`
	PreviousBest    float64   `json:"previous_best,omitempty"`
	Improvement     float64   `json:"improvement,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
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

	db := database.GetDB()

	// Verify exercise exists and belongs to user
	var exercise models.Exercise
	if err := db.Where("id = ? AND user_id = ?", req.ExerciseID, userID).First(&exercise).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Exercise not found"})
		return
	}

	// Parse date or use current time
	entryDate := time.Now()
	if req.Date != "" {
		parsedDate, err := time.Parse("2006-01-02", req.Date)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
			return
		}
		entryDate = parsedDate
	}

	// Detect if this is a PR
	prResult, err := services.DetectPR(userID.(uint), req.ExerciseID, req.Weight, req.Reps, req.Sets)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check PR status"})
		return
	}

	// Create the workout entry
	entry := models.WorkoutEntry{
		UserID:     userID.(uint),
		ExerciseID: req.ExerciseID,
		Weight:     req.Weight,
		Reps:       req.Reps,
		Sets:       req.Sets,
		Notes:      req.Notes,
		Date:       entryDate,
		Score:      prResult.Score,
		IsPR:       prResult.IsPR,
	}

	if err := db.Create(&entry).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create entry"})
		return
	}

	// If it's a PR, record it in PR history
	if prResult.IsPR {
		if err := services.RecordPR(userID.(uint), req.ExerciseID, entry.ID, req.Weight, req.Reps, req.Sets, entryDate); err != nil {
			// Log error but don't fail the request
			// The entry was created successfully
		}
	}

	// Build response with celebration indicator
	response := EntryResponse{
		ID:              entry.ID,
		UserID:          entry.UserID,
		ExerciseID:      entry.ExerciseID,
		Weight:          entry.Weight,
		Reps:            entry.Reps,
		Sets:            entry.Sets,
		Notes:           entry.Notes,
		Date:            entry.Date,
		Score:           entry.Score,
		IsPR:            entry.IsPR,
		CelebrationText: prResult.CelebrationText,
		PreviousBest:    prResult.PreviousBest,
		Improvement:     prResult.Improvement,
		CreatedAt:       entry.CreatedAt,
	}

	c.JSON(http.StatusCreated, response)
}
