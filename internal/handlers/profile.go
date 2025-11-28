package handlers

import (
	"fitness-market/internal/database"
	"fitness-market/internal/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type UpdateBodyweightRequest struct {
	Weight     float64   `json:"weight" binding:"required,gt=0"`
	Unit       string    `json:"unit"`
	RecordedAt time.Time `json:"recorded_at"`
}

type AddExercisePRRequest struct {
	ExerciseName string    `json:"exercise_name" binding:"required"`
	Weight       float64   `json:"weight" binding:"required,gt=0"`
	Unit         string    `json:"unit"`
	Reps         int       `json:"reps"`
	RecordedAt   time.Time `json:"recorded_at"`
}

func GetUserProfile(c *gin.Context) {
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}
	user := userInterface.(models.User)

	var profile models.UserProfile
	if err := database.DB.Where("user_id = ?", user.ID).First(&profile).Error; err != nil {
		profile = models.UserProfile{UserID: user.ID}
		database.DB.Create(&profile)
	}

	var latestBodyweight models.BodyweightEntry
	database.DB.Where("user_id = ?", user.ID).Order("recorded_at desc").First(&latestBodyweight)

	var exercisePRs []models.ExercisePR
	database.DB.Where("user_id = ?", user.ID).Find(&exercisePRs)

	c.JSON(http.StatusOK, gin.H{
		"profile":          profile,
		"current_bodyweight": latestBodyweight,
		"exercise_prs":     exercisePRs,
	})
}

func AddBodyweight(c *gin.Context) {
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}
	user := userInterface.(models.User)

	var req UpdateBodyweightRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	unit := req.Unit
	if unit == "" {
		unit = "kg"
	}

	recordedAt := req.RecordedAt
	if recordedAt.IsZero() {
		recordedAt = time.Now()
	}

	entry := models.BodyweightEntry{
		UserID:     user.ID,
		Weight:     req.Weight,
		Unit:       unit,
		RecordedAt: recordedAt,
	}

	if err := database.DB.Create(&entry).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save bodyweight"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"bodyweight": entry})
}

func GetBodyweightHistory(c *gin.Context) {
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}
	user := userInterface.(models.User)

	var entries []models.BodyweightEntry
	database.DB.Where("user_id = ?", user.ID).Order("recorded_at desc").Find(&entries)

	c.JSON(http.StatusOK, gin.H{"bodyweight_history": entries})
}

func AddExercisePR(c *gin.Context) {
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}
	user := userInterface.(models.User)

	var req AddExercisePRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	unit := req.Unit
	if unit == "" {
		unit = "kg"
	}

	reps := req.Reps
	if reps <= 0 {
		reps = 1
	}

	recordedAt := req.RecordedAt
	if recordedAt.IsZero() {
		recordedAt = time.Now()
	}

	pr := models.ExercisePR{
		UserID:       user.ID,
		ExerciseName: req.ExerciseName,
		Weight:       req.Weight,
		Unit:         unit,
		Reps:         reps,
		RecordedAt:   recordedAt,
	}

	if err := database.DB.Create(&pr).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save exercise PR"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"exercise_pr": pr})
}

func GetExercisePRs(c *gin.Context) {
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}
	user := userInterface.(models.User)

	var prs []models.ExercisePR
	database.DB.Where("user_id = ?", user.ID).Order("exercise_name asc, recorded_at desc").Find(&prs)

	c.JSON(http.StatusOK, gin.H{"exercise_prs": prs})
}

func UpdateExercisePR(c *gin.Context) {
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}
	user := userInterface.(models.User)

	prID := c.Param("id")

	var pr models.ExercisePR
	if err := database.DB.Where("id = ? AND user_id = ?", prID, user.ID).First(&pr).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Exercise PR not found"})
		return
	}

	var req AddExercisePRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pr.ExerciseName = req.ExerciseName
	pr.Weight = req.Weight
	if req.Unit != "" {
		pr.Unit = req.Unit
	}
	if req.Reps > 0 {
		pr.Reps = req.Reps
	}
	if !req.RecordedAt.IsZero() {
		pr.RecordedAt = req.RecordedAt
	}

	if err := database.DB.Save(&pr).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update exercise PR"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"exercise_pr": pr})
}

func DeleteExercisePR(c *gin.Context) {
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}
	user := userInterface.(models.User)

	prID := c.Param("id")

	var pr models.ExercisePR
	if err := database.DB.Where("id = ? AND user_id = ?", prID, user.ID).First(&pr).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Exercise PR not found"})
		return
	}

	if err := database.DB.Delete(&pr).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete exercise PR"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Exercise PR deleted"})
}
