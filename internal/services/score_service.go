package services

import (
	"fitness-market/internal/database"
	"fitness-market/internal/models"
)

const (
	DefaultBodyweight = 70.0
	ReferenceBodyweight = 75.0
)

func GetUserCurrentBodyweight(userID uint) float64 {
	var entry models.BodyweightEntry
	if err := database.DB.Where("user_id = ?", userID).Order("recorded_at desc").First(&entry).Error; err != nil {
		return DefaultBodyweight
	}

	weight := entry.Weight
	if entry.Unit == "lb" || entry.Unit == "lbs" {
		weight = weight * 0.453592
	}

	return weight
}

func CalculateBodyweightFactor(userID uint) float64 {
	bodyweight := GetUserCurrentBodyweight(userID)
	if bodyweight <= 0 {
		return 1.0
	}
	return ReferenceBodyweight / bodyweight
}

func NormalizeScore(rawScore float64, userID uint) float64 {
	factor := CalculateBodyweightFactor(userID)
	return rawScore * factor
}
