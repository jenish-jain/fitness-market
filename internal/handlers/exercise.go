package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"fitness-market/internal/database"
	"fitness-market/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreateExerciseRequest struct {
	Ticker      string  `json:"ticker" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Category    string  `json:"category" binding:"required"`
	StockPrice  float64 `json:"stock_price"`
}

type UpdateExerciseRequest struct {
	Ticker      string  `json:"ticker"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	StockPrice  float64 `json:"stock_price"`
}

// CreateExercise creates a new exercise for the authenticated user
func CreateExercise(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req CreateExerciseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Normalize ticker to uppercase
	ticker := strings.ToUpper(strings.TrimSpace(req.Ticker))

	exercise := models.Exercise{
		UserID:      userID.(uint),
		Ticker:      ticker,
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
		StockPrice:  req.StockPrice,
	}

	if err := database.DB.Create(&exercise).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			c.JSON(http.StatusConflict, gin.H{"error": "Ticker already exists for this user"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create exercise"})
		return
	}

	c.JSON(http.StatusCreated, exercise)
}

// GetExercises returns all exercises for the authenticated user
func GetExercises(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var exercises []models.Exercise
	if err := database.DB.Where("user_id = ?", userID).Find(&exercises).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch exercises"})
		return
	}

	c.JSON(http.StatusOK, exercises)
}

// UpdateExercise updates an exercise for the authenticated user
func UpdateExercise(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	exerciseID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid exercise ID"})
		return
	}

	var exercise models.Exercise
	if err := database.DB.Where("id = ? AND user_id = ?", exerciseID, userID).First(&exercise).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Exercise not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch exercise"})
		return
	}

	var req UpdateExerciseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update fields if provided
	if req.Ticker != "" {
		exercise.Ticker = strings.ToUpper(strings.TrimSpace(req.Ticker))
	}
	if req.Name != "" {
		exercise.Name = req.Name
	}
	if req.Description != "" {
		exercise.Description = req.Description
	}
	if req.Category != "" {
		exercise.Category = req.Category
	}
	if req.StockPrice != 0 {
		exercise.StockPrice = req.StockPrice
	}

	if err := database.DB.Save(&exercise).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			c.JSON(http.StatusConflict, gin.H{"error": "Ticker already exists for this user"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update exercise"})
		return
	}

	c.JSON(http.StatusOK, exercise)
}

// DeleteExercise deletes an exercise for the authenticated user
func DeleteExercise(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	exerciseID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid exercise ID"})
		return
	}

	result := database.DB.Where("id = ? AND user_id = ?", exerciseID, userID).Delete(&models.Exercise{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete exercise"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Exercise not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Exercise deleted successfully"})
}
