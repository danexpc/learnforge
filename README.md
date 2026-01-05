# LearnForge

LearnForge is a small AI-powered service that turns raw text into structured learning content. You send it any text, and it generates a clear summary, key points, flashcards, and quiz questions tailored to a specific topic and difficulty level.

The service can use a topic you provide (like engineering, biology, or business), or automatically detect the topic on its own. It's designed to support personalized learning experiences and show how AI can make educational content more engaging and easier to understand.

## Overview

A production-style Go microservice that transforms raw text into structured learning content using LLMs (OpenAI-compatible APIs). LearnForge solves the problem of converting unstructured educational text into structured learning materials (summaries, key points, flashcards, and quizzes) automatically. It's designed as a clean, maintainable microservice with proper separation of concerns, observability, and production-ready features.

## Features

- **Text Processing**: Converts raw text into structured learning content (summary, key points, flashcards, quizzes)
- **Topic Detection**: Auto-detects topics when not provided, with confidence scores
- **Flexible Storage**: Supports in-memory (default) or PostgreSQL storage
- **Observability**: Structured logging, Prometheus metrics, request tracing
- **Production Ready**: Timeouts, retries, graceful shutdown, error handling
- **Clean Architecture**: Dependency inversion, interface-based design
- **Slack Integration**: Error logging and daily summaries to separate Slack channels
- **Daily Summaries**: Automated daily summaries with Redis caching (in-memory fallback)
- **Manual Summary Trigger**: Secure API endpoint to generate and send summaries on-demand

## Architecture

```
┌─────────────────┐
│   HTTP Client   │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  HTTP Transport │  (handlers, middleware)
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│    Service      │  (business logic)
└────────┬────────┘
         │
    ┌────┴────┐
    ▼         ▼
┌────────┐ ┌────────┐
│   AI   │ │ Store  │  (interfaces)
└────────┘ └────────┘
    │         │
    ▼         ▼
┌────────┐ ┌────────┐
│OpenAI  │ │InMem/  │
│Client  │ │Postgres│
└────────┘ └────────┘
```

### Directory Structure

```
learnforge/
├── cmd/
│   └── api/              # Application entry point
├── config/               # YAML configuration files
│   ├── development.yml
│   ├── local.example.yml
│   └── production.yml
├── internal/
│   ├── domain/           # Domain models, errors, and validation
│   ├── service/          # Business logic
│   ├── transport/
│   │   └── http/         # HTTP handlers and middleware
│   ├── store/            # Storage interface, implementations, and migrations
│   ├── ai/               # AI client interface and implementation
│   └── config/           # Configuration management
├── go.mod
├── go.sum
├── Makefile
├── docker-compose.yml
└── README.md
```

## Getting Started

### Prerequisites

- Go 1.22+
- Docker and Docker Compose (for PostgreSQL mode)
- OpenAI-compatible API key

### Configuration

The service supports YAML configuration files with environment variable overrides. Configuration files are located in the `config/` directory:

- `development.yml` - Default development settings (in-memory storage)
- `local.example.yml` - Example for local PostgreSQL setup
- `production.yml` - Production settings template

The service loads configuration based on the `ENV` environment variable (defaults to `development`). You can also specify a custom path via `CONFIG_PATH`.

Environment variables always override YAML values, allowing you to override sensitive values like API keys.

### Running in In-Memory Mode (Zero Dependencies)

1. Copy the example config (optional):
   ```bash
   cp config/local.example.yml config/local.yml
   # Edit config/local.yml and set your AI_API_KEY
   ```

2. Set your AI API key (if not in config file):
   ```bash
   export AI_API_KEY=your-api-key-here
   ```

3. Run the service:
   ```bash
   make run
   # or
   ENV=development go run cmd/api/main.go
   ```

The service will start on `http://localhost:8080` (configurable via `PORT` env var or config file).

### Running with PostgreSQL and Redis

1. Start services:
   ```bash
   make up
   ```

   This starts:
   - PostgreSQL on port 5432
   - Redis on port 6379
   - Adminer (database admin) on port 8081

2. Create a local config file:
   ```bash
   cp config/local.example.yml config/local.yml
   # Edit config/local.yml with your settings
   ```

3. Or set environment variables:
   ```bash
   export ENV=local
   export STORAGE=postgres
   export DATABASE_URL=postgres://learnforge:learnforge@localhost:5432/learnforge?sslmode=disable
   export AI_API_KEY=your-api-key-here
   ```

4. Run the service:
   ```bash
   make run-pg
   # or
   ENV=local go run cmd/api/main.go
   ```

**Note**: Migrations run automatically on startup when using PostgreSQL. The in-memory store doesn't require migrations.

### Database Access

If using Docker Compose, you can access the database via Adminer at `http://localhost:8081`:
- System: PostgreSQL
- Server: postgres
- Username: learnforge
- Password: learnforge
- Database: learnforge

## API Usage

### Process Text

```bash
curl -X POST http://localhost:8080/v1/process \
  -H "Content-Type: application/json" \
  -d '{
    "text": "The water cycle is the continuous movement of water on, above, and below the Earth's surface. Water evaporates from oceans, lakes, and rivers, forms clouds, and falls back as precipitation.",
    "mode": "lesson",
    "topic": "Water Cycle",
    "level": "beginner",
    "language": "en"
  }'
```

**Response:**
```json
{
  "id": "abc123...",
  "topic": "Water Cycle",
  "topic_source": "user",
  "topic_confidence": 1.0,
  "summary": "...",
  "key_points": ["...", "..."],
  "flashcards": [
    {"q": "What is the water cycle?", "a": "..."}
  ],
  "quiz": [
    {
      "q": "Where does water evaporate from?",
      "choices": ["Oceans", "Mountains", "Deserts"],
      "answer": "Oceans"
    }
  ],
  "meta": {
    "model": "gpt-3.5-turbo",
    "provider": "openai-compatible",
    "processing_ms": 1234
  },
  "created_at": "2024-01-01T12:00:00Z"
}
```

### Get Result by ID

```bash
curl http://localhost:8080/v1/process/{id}
```

### Health Endpoints

```bash
# Health check
curl http://localhost:8080/healthz

# Readiness check
curl http://localhost:8080/readyz

# Prometheus metrics
curl http://localhost:8080/metrics
```

## Configuration

### YAML Configuration Files

Configuration files support the following fields:

```yaml
port: "8080"
storage: "inmem"  # or "postgres"
database_url: ""
ai_base_url: "https://api.openai.com"
ai_api_key: ""
ai_model: "gpt-3.5-turbo"
log_level: "info"
```

### Environment Variables

Environment variables override YAML values:

| Variable | Default | Description |
|----------|---------|-------------|
| `ENV` | `development` | Environment name (determines which YAML file to load) |
| `CONFIG_PATH` | - | Custom path to config file (overrides ENV-based path) |
| `PORT` | `8080` | HTTP server port |
| `STORAGE` | `inmem` | Storage backend: `inmem` or `postgres` |
| `DATABASE_URL` | - | PostgreSQL connection string (required for postgres mode) |
| `AI_BASE_URL` | `https://api.openai.com` | OpenAI-compatible API base URL |
| `AI_API_KEY` | - | API key for AI service (required) |
| `AI_MODEL` | `gpt-3.5-turbo` | Model to use |
| `LOG_LEVEL` | `info` | Logging level |
| `SLACK_WEBHOOK_URL` | - | Slack webhook URL for daily summaries |
| `SLACK_ERROR_WEBHOOK_URL` | - | Slack webhook URL for error notifications |
| `SUMMARY_API_KEY` | - | API key for manual summary generation endpoint |
| `REDIS_URL` | - | Redis connection URL (optional, falls back to in-memory cache) |

## Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage
```

## Development

### Building

```bash
make build
```

Binary will be in `bin/api`.

### Project Structure Principles

- **Domain Layer**: Core business entities and errors (no dependencies)
- **Service Layer**: Business logic orchestration
- **Transport Layer**: HTTP-specific concerns (handlers, middleware)
- **Infrastructure**: External dependencies (AI, database) behind interfaces

### Key Design Decisions

1. **Interface-based Design**: Store and AI clients are interfaces, enabling easy testing and swapping implementations
2. **Dependency Inversion**: High-level modules don't depend on low-level modules
3. **Error Handling**: Consistent error format with proper HTTP status code mapping
4. **Observability**: Structured logging with request IDs, Prometheus metrics
5. **Graceful Shutdown**: Proper cleanup on SIGINT/SIGTERM
6. **Validation**: Centralized validation for request fields (mode, level) with clear error messages
7. **Migrations**: Automatic database migrations for PostgreSQL (skipped for in-memory mode)
8. **Configuration**: YAML-based config with environment variable overrides for flexibility

## Trade-offs and Future Improvements

### Current Trade-offs

1. **In-Memory Storage**: Fast but not persistent (data lost on restart)
2. **Simple Retry Logic**: Only one retry for AI requests (could be more sophisticated)
3. **Basic Metrics**: Core metrics only (could add more detailed business metrics)
4. **No Rate Limiting**: Currently no rate limiting (mentioned as nice-to-have)

### Slack Integration

#### Setting Up Slack Webhooks

1. Create a Slack app at https://api.slack.com/apps
2. Enable "Incoming Webhooks"
3. Create webhooks for:
   - **Summary Channel**: For daily summaries
   - **Error Channel**: For error notifications
4. Add webhook URLs to your config file or environment variables

#### Daily Summaries

Daily summaries are automatically generated at midnight UTC and sent to the summary Slack channel. The summary includes:
- Total requests processed
- Top topics
- Error count (if any)

Summaries are cached in Redis (or in-memory) for 7 days to avoid regeneration.

#### Error Logging

Errors are automatically sent to the error Slack channel with context information including:
- Request ID
- Endpoint
- Error message

### Future Improvements

- [ ] Rate limiting per IP (token bucket)
- [ ] OpenAPI/Swagger documentation
- [ ] Batch processing support
- [ ] Webhook notifications for async processing
- [ ] Multi-language support improvements
- [ ] Custom prompt templates
- [ ] Analytics dashboard
- [ ] Summary email notifications

## License

MIT

