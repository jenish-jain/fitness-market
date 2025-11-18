package auth

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/supabase-community/supabase-go"
)

type SupabaseClient struct {
	client    *supabase.Client
	jwtSecret string
}

type SupabaseUser struct {
	ID    string `json:"sub"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

func NewSupabaseClient() (*SupabaseClient, error) {
	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_ANON_KEY")
	jwtSecret := os.Getenv("SUPABASE_JWT_SECRET")

	if supabaseURL == "" || supabaseKey == "" || jwtSecret == "" {
		return nil, errors.New("missing required Supabase environment variables")
	}

	client, err := supabase.NewClient(supabaseURL, supabaseKey, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Supabase client: %w", err)
	}

	return &SupabaseClient{
		client:    client,
		jwtSecret: jwtSecret,
	}, nil
}

func (s *SupabaseClient) ValidateToken(ctx context.Context, tokenString string) (*jwt.Token, error) {
	// Parse the token without verifying the signature first to extract header info
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Return the JWT secret for HMAC validation
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if !token.Valid {
		return nil, errors.New("token is not valid")
	}

	return token, nil
}

func (s *SupabaseClient) GetUserFromToken(token *jwt.Token) (*SupabaseUser, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	user := &SupabaseUser{}

	if sub, ok := claims["sub"].(string); ok {
		user.ID = sub
	} else {
		return nil, errors.New("missing user ID in token")
	}

	if email, ok := claims["email"].(string); ok {
		user.Email = email
	}

	if role, ok := claims["role"].(string); ok {
		user.Role = role
	} else {
		user.Role = "authenticated" // default role
	}

	return user, nil
}