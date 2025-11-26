package middleware

import (
	"fitness-market/internal/auth"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Check Bearer prefix
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token required"})
			c.Abort()
			return
		}

		// Extract token
		token := strings.TrimPrefix(authHeader, "Bearer ")

		// Initialize Supabase client
		supabaseClient, err := auth.NewSupabaseClient()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Authentication service unavailable"})
			c.Abort()
			return
		}

		// Validate token
		jwtToken, err := supabaseClient.ValidateToken(c.Request.Context(), token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Get user from token
		user, err := supabaseClient.GetUserFromToken(jwtToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to extract user from token"})
			c.Abort()
			return
		}

		// Set user in context
		c.Set("user", user)
		c.Next()
	}
}