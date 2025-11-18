package middleware

import (
	"fitness-market/internal/auth"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Extract token from Bearer format
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString := tokenParts[1]

		// Initialize Supabase client
		supabaseClient, err := auth.NewSupabaseClient()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Authentication service unavailable"})
			c.Abort()
			return
		}

		// Validate token
		token, err := supabaseClient.ValidateToken(c.Request.Context(), tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Get user from token
		user, err := supabaseClient.GetUserFromToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user token"})
			c.Abort()
			return
		}

		// Set user in context
		c.Set("user", user)
		c.Next()
	}
}