.PHONY: run run-pg up down test build build-web clean web-dev

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

# Build the web UI
build-web:
	@echo "Building web UI..."
	@cd web && npm install && npm run build

# Build the application (includes web UI)
build: build-web
	@echo "Building application..."
	go build -o bin/api cmd/api/main.go

# Run web dev server
web-dev:
	@echo "Starting web dev server..."
	@cd web && npm install && npm run dev

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/ web/dist/ web/node_modules/

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

