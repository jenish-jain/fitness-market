package main

import (
	"fitness-market/internal/auth"
	"fitness-market/internal/database"
	"fitness-market/internal/middleware"
	"fitness-market/internal/models"
	"log"
	"net/http"
)

func main() {
	log.Println("Starting fitness-market server...")
	// Server implementation will go here
	http.ListenAndServe(":8080", nil)
}