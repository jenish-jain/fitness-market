package auth

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
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

type JWKSet struct {
	Keys []JWK `json:"keys"`
}

type JWK struct {
	Kty string `json:"kty"`
	Use string `json:"use"`
	Kid string `json:"kid"`
	N   string `json:"n"`
	E   string `json:"e"`
	Alg string `json:"alg"`
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
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			// Try RSA method for JWK validation
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			// Get JWK for RSA validation
			return s.getJWKKey(token)
		}
		// Return HMAC secret
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return token, nil
}

func (s *SupabaseClient) GetUserFromToken(token *jwt.Token) (*SupabaseUser, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	user := &SupabaseUser{}
	if id, ok := claims["sub"].(string); ok {
		user.ID = id
	}
	if email, ok := claims["email"].(string); ok {
		user.Email = email
	}
	if role, ok := claims["role"].(string); ok {
		user.Role = role
	}

	return user, nil
}

func (s *SupabaseClient) getJWKKey(token *jwt.Token) (*rsa.PublicKey, error) {
	// Get JWK URL from Supabase project
	jwkURL := fmt.Sprintf("%s/rest/v1/", strings.TrimSuffix(os.Getenv("SUPABASE_URL"), "/"))
	jwkURL = strings.Replace(jwkURL, "/rest/v1/", "/.well-known/jwks.json", 1)

	resp, err := http.Get(jwkURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JWK: %w", err)
	}
	defer resp.Body.Close()

	var jwkSet JWKSet
	if err := json.NewDecoder(resp.Body).Decode(&jwkSet); err != nil {
		return nil, fmt.Errorf("failed to decode JWK: %w", err)
	}

	// Find matching key
	kid, ok := token.Header["kid"].(string)
	if !ok {
		return nil, errors.New("missing kid in token header")
	}

	for _, key := range jwkSet.Keys {
		if key.Kid == kid {
			return s.jwkToRSAKey(key)
		}
	}

	return nil, errors.New("matching JWK not found")
}

func (s *SupabaseClient) jwkToRSAKey(jwk JWK) (*rsa.PublicKey, error) {
	n, err := base64.RawURLEncoding.DecodeString(jwk.N)
	if err != nil {
		return nil, fmt.Errorf("failed to decode N: %w", err)
	}

	e, err := base64.RawURLEncoding.DecodeString(jwk.E)
	if err != nil {
		return nil, fmt.Errorf("failed to decode E: %w", err)
	}

	pubKey := &rsa.PublicKey{
		N: new(big.Int).SetBytes(n),
		E: int(new(big.Int).SetBytes(e).Int64()),
	}

	return pubKey, nil
}