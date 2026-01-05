# Local Development Setup Guide

## Quick Start: In-Memory Mode (No Dependencies)

This is the simplest way to run LearnForge locally - no Redis, no PostgreSQL needed!

```bash
# 1. Set your AI API key
export AI_API_KEY=your-api-key-here

# 2. Run the application
make run

# OR directly:
ENV=development go run cmd/api/main.go
```

The app will:
- ✅ Use in-memory storage (no database)
- ✅ Use in-memory cache (no Redis)
- ✅ Start on http://localhost:8080

## Installing Redis (Optional)

Redis is **NOT required** for in-memory mode, but if you want to test Redis functionality:

### Check if Redis is installed:
```bash
redis-cli --version
```

### Install Redis on macOS:

**Option 1: Using Homebrew (Recommended)**
```bash
# Install Homebrew if you don't have it
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# Install Redis
brew install redis

# Start Redis
brew services start redis

# Verify it's running
redis-cli ping
# Should return: PONG
```

**Option 2: Using Docker**
```bash
# Start Redis using Docker Compose (already configured)
make up

# This starts both PostgreSQL and Redis
# Redis will be available at redis://localhost:6379
```

### Using Redis with LearnForge:

Once Redis is installed, update your config or environment:

```bash
# Option 1: Set environment variable
export REDIS_URL=redis://localhost:6379/0

# Option 2: Update config/development.yml
redis_url: "redis://localhost:6379/0"

# Then run
make run
```

## Installing PostgreSQL (Optional)

PostgreSQL is **NOT required** for in-memory mode, but if you want persistent storage:

### Using Docker (Easiest):
```bash
# Start PostgreSQL and Redis
make up

# This will start:
# - PostgreSQL on port 5432
# - Redis on port 6379
# - Adminer (database admin UI) on http://localhost:8081
```

### Using Homebrew:
```bash
# Install PostgreSQL
brew install postgresql@15

# Start PostgreSQL
brew services start postgresql@15

# Create database
createdb learnforge
```

## Testing the Application

### 1. Health Check
```bash
curl http://localhost:8080/healthz
```

### 2. Process Text
```bash
curl -X POST http://localhost:8080/v1/process \
  -H "Content-Type: application/json" \
  -d '{
    "text": "The water cycle is the continuous movement of water on Earth.",
    "topic": "Science",
    "level": "beginner"
  }'
```

### 3. Get Result
```bash
# Use the ID from the previous response
curl http://localhost:8080/v1/process/{id}
```

## Troubleshooting

### Redis Connection Issues
If Redis is not available, the app will automatically fall back to in-memory cache. You'll see:
```
{"level":"warn","msg":"Failed to connect to Redis, using in-memory cache"}
```

### Port Already in Use
If port 8080 is busy:
```bash
export PORT=8081
make run
```

### Missing AI API Key
The app requires an AI API key:
```bash
export AI_API_KEY=your-key-here
```

## Summary

- **In-Memory Mode**: No dependencies needed, just set `AI_API_KEY` and run `make run`
- **With Redis**: Install via `brew install redis` or use `make up` for Docker
- **With PostgreSQL**: Use `make up` for Docker or install via Homebrew

