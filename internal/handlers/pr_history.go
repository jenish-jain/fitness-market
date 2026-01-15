package handlers

import (
	"net/http"
	"strconv"

	"fitness-market/internal/services"

	"github.com/gin-gonic/gin"
)

// GetPRHistoryByExercise returns PR history for a specific exercise
func GetPRHistoryByExercise(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	exerciseIDStr := c.Param("exercise_id")
	exerciseID, err := strconv.ParseUint(exerciseIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid exercise ID"})
		return
	}

	history, err := services.GetPRHistory(userID.(uint), uint(exerciseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch PR history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"exercise_id": exerciseID,
		"pr_history":  history,
		"count":       len(history),
	})
}

// GetAllUserPRs returns the current best PR for each exercise
func GetAllUserPRs(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	prs, err := services.GetAllPRs(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch PRs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"personal_records": prs,
		"count":            len(prs),
	})
}
