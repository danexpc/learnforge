.PHONY: run run-pg up down test build clean

# Run in in-memory mode
run:
	@echo "Running in in-memory mode..."
	@STORAGE=inmem go run cmd/api/main.go

# Run with PostgreSQL (requires DATABASE_URL)
run-pg:
	@echo "Running with PostgreSQL..."
	@STORAGE=postgres DATABASE_URL=postgres://learnforge:learnforge@localhost:5432/learnforge?sslmode=disable go run cmd/api/main.go

# Start Docker Compose services
up:
	@echo "Starting Docker Compose services..."
	docker-compose up -d

# Stop Docker Compose services
down:
	@echo "Stopping Docker Compose services..."
	docker-compose down

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Build the application
build:
	@echo "Building application..."
	go build -o bin/api cmd/api/main.go

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

