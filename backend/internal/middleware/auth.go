package middleware

import (
	"net/http"
)

// AuthMiddleware handles authentication
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Authentication logic here
		next.ServeHTTP(w, r)
	})
}