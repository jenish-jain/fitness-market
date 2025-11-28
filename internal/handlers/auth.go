package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strings"
	"time"

	"fitness-market/internal/auth"
	"fitness-market/internal/database"
	"fitness-market/internal/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type ResetPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordConfirmRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

type AuthResponse struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}

func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	var existingUser models.User
	if err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
		return
	}

	user := models.User{
		Email:    req.Email,
		Password: string(hashedPassword),
		Name:     req.Name,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	token, err := auth.GenerateToken(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, AuthResponse{
		Token: token,
		User:  &user,
	})
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	token, err := auth.GenerateToken(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		Token: token,
		User:  &user,
	})
}

func Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

func RequestPasswordReset(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "If an account with that email exists, a password reset link has been sent"})
		return
	}

	resetToken := generateSecureToken(32)
	expiryTime := time.Now().Add(1 * time.Hour)

	user.ResetToken = resetToken
	user.ResetTokenExpiry = &expiryTime

	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "If an account with that email exists, a password reset link has been sent",
		"reset_token": resetToken,
	})
}

func ResetPassword(c *gin.Context) {
	var req ResetPasswordConfirmRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	var user models.User
	if err := database.DB.Where("reset_token = ?", req.Token).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired reset token"})
		return
	}

	if user.ResetTokenExpiry == nil || time.Now().After(*user.ResetTokenExpiry) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired reset token"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
		return
	}

	user.Password = string(hashedPassword)
	user.ResetToken = ""
	user.ResetTokenExpiry = nil

	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password successfully reset"})
}

func generateSecureToken(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)
}
