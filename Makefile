# Database operations
.PHONY: migrate seed db-reset

migrate:
	go run scripts/migrate.go

seed:
	go run scripts/seed.go

db-reset:
	rm -f fitness_market.db
	make migrate
	make seed

# Server operations
.PHONY: run dev

run:
	go run cmd/server/main.go

dev:
	air -c .air.toml

# Build
.PHONY: build

build:
	go build -o bin/server cmd/server/main.go

# Clean
.PHONY: clean

clean:
	rm -rf bin/
	rm -f fitness_market.db