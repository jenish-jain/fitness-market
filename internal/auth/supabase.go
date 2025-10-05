package auth

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/supabase-community/supabase-go"
)

type SupabaseClient struct {
	client *supabase.Client
}

func NewSupabaseClient() (*SupabaseClient, error) {
	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_ANON_KEY")

	if supabaseURL == "" || supabaseKey == "" {
		return nil, errors.New("supabase configuration not found")
	}

	client, err := supabase.NewClient(supabaseURL, supabaseKey, nil)
	if err != nil {
		return nil, err
	}

	return &SupabaseClient{client: client}, nil
}

func (s *SupabaseClient) ValidateToken(ctx context.Context, token string) (*jwt.Token, error) {
	// Remove Bearer prefix if present
	token = strings.TrimPrefix(token, "Bearer ")

	// Parse JWT token
	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	if !jwtToken.Valid {
		return nil, errors.New("invalid token")
	}

	return jwtToken, nil
}

func (s *SupabaseClient) GetUserFromToken(token *jwt.Token) (map[string]interface{}, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}