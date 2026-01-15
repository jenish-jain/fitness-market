package handlers

import (
	"errors"
	"net/http"
	"regexp"
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

// validateTicker validates the ticker symbol format
func validateTicker(ticker string) error {
	if len(ticker) < 2 || len(ticker) > 10 {
		return errors.New("ticker symbol must be between 2 and 10 characters")
	}
	// Ticker must be uppercase alphanumeric only
	matched, _ := regexp.MatchString("^[A-Z0-9]+$", ticker)
	if !matched {
		return errors.New("ticker symbol must contain only uppercase letters and numbers")
	}
	return nil
}

// CreateExercise creates a new exercise for the authenticated user
func CreateExercise(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req CreateExerciseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Normalize ticker to uppercase
	req.Ticker = strings.ToUpper(strings.TrimSpace(req.Ticker))
	req.Name = strings.TrimSpace(req.Name)

	// Validate ticker format
	if err := validateTicker(req.Ticker); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check ticker uniqueness for this user
	var existingExercise models.Exercise
	err := database.DB.Where("user_id = ? AND ticker = ?", userID, req.Ticker).First(&existingExercise).Error
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Ticker symbol already exists for this user"})
		return
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	exercise := models.Exercise{
		UserID:      userID.(uint),
		Ticker:      req.Ticker,
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
		StockPrice:  req.StockPrice,
	}

	if err := database.DB.Create(&exercise).Error; err != nil {
		if strings.Contains(err.Error(), "ticker symbol already exists") {
			c.JSON(http.StatusConflict, gin.H{"error": "Ticker symbol already exists for this user"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create exercise"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Exercise created successfully",
		"exercise": exercise,
	})
}

// GetExercises returns all exercises for the authenticated user
func GetExercises(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var exercises []models.Exercise
	if err := database.DB.Where("user_id = ?", userID).Find(&exercises).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch exercises"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"exercises": exercises})
}

// UpdateExercise updates an exercise for the authenticated user
func UpdateExercise(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	var req UpdateExerciseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update fields if provided
	if req.Ticker != "" {
		req.Ticker = strings.ToUpper(strings.TrimSpace(req.Ticker))
		if err := validateTicker(req.Ticker); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// Check ticker uniqueness for this user (excluding current exercise)
		var existingExercise models.Exercise
		err := database.DB.Where("user_id = ? AND ticker = ? AND id != ?", userID, req.Ticker, exerciseID).First(&existingExercise).Error
		if err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Ticker symbol already exists for this user"})
			return
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		exercise.Ticker = req.Ticker
	}
	if req.Name != "" {
		exercise.Name = strings.TrimSpace(req.Name)
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
		if strings.Contains(err.Error(), "ticker symbol already exists") {
			c.JSON(http.StatusConflict, gin.H{"error": "Ticker symbol already exists for this user"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update exercise"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Exercise updated successfully",
		"exercise": exercise,
	})
}

// DeleteExercise deletes an exercise for the authenticated user
func DeleteExercise(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if err := database.DB.Delete(&exercise).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete exercise"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Exercise deleted successfully"})
}
