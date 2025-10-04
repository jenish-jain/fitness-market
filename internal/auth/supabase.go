package auth

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/supabase-community/supabase-go"
)

var SupabaseClient *supabase.Client

func InitSupabase() {
	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_ANON_KEY")

	if supabaseURL == "" || supabaseKey == "" {
		panic("SUPABASE_URL and SUPABASE_ANON_KEY must be set")
	}

	var err error
	SupabaseClient, err = supabase.NewClient(supabaseURL, supabaseKey, nil)
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize Supabase client: %v", err))
	}
}

func ValidateToken(token string) (*supabase.User, error) {
	user, err := SupabaseClient.Auth.GetUser(context.Background(), token)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func ExtractTokenFromHeader(authHeader string) string {
	if authHeader == "" {
		return ""
	}
	
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}
	
	return parts[1]
}