package handlers

import (
	"net/http"
	"time"

	"fitness-market/internal/database"
	"fitness-market/internal/models"

	"github.com/gin-gonic/gin"
)

type CreateEntryRequest struct {
	ExerciseID uint    `json:"exercise_id" binding:"required"`
	Weight     float64 `json:"weight" binding:"required,gte=0"`
	Reps       int     `json:"reps" binding:"required,gte=1"`
	Sets       int     `json:"sets" binding:"required,gte=1"`
	Notes      string  `json:"notes"`
}

// CreateEntry creates a new workout entry for the authenticated user
func CreateEntry(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req CreateEntryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	// Verify the exercise exists and belongs to the user
	var exercise models.Exercise
	if err := database.DB.Where("id = ? AND user_id = ?", req.ExerciseID, userID).First(&exercise).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Exercise not found"})
		return
	}

	entry := models.WorkoutEntry{
		UserID:     userID.(uint),
		ExerciseID: req.ExerciseID,
		Date:       time.Now(),
		Sets:       req.Sets,
		Reps:       req.Reps,
		Weight:     req.Weight,
		Notes:      req.Notes,
	}

	if err := database.DB.Create(&entry).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create workout entry"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Workout entry created successfully",
		"entry":   entry,
	})
}
