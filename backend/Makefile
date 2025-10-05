# Database operations
.PHONY: migrate seed build run test clean

# Build the application
build:
	go build -o bin/server cmd/server/main.go

# Run the server
run:
	go run cmd/server/main.go

# Run database migrations
migrate:
	go run scripts/migrate.go

# Seed the database
seed:
	go run scripts/seed.go

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	rm -rf bin/