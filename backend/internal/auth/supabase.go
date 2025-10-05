package auth

import (
	"context"
	"errors"
)

// Supabase auth implementation
type SupabaseAuth struct {
	// Configuration fields
}

func NewSupabaseAuth() *SupabaseAuth {
	return &SupabaseAuth{}
}

func (s *SupabaseAuth) Authenticate(ctx context.Context, token string) error {
	return errors.New("not implemented")
}