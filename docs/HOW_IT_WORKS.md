# How LearnForge Works Without External Dependencies

## 1. How In-Memory Mode Works (No Redis Needed)

When you run in in-memory mode without Redis:

### Storage (Data Persistence)
- **In-Memory Store**: Uses a Go `map[string]*StoredResult` with `sync.RWMutex` for thread safety
- Data is stored in RAM only - lost when the app restarts
- Perfect for development and testing
- No database connection needed

### Caching (Summary Caching)
- **In-Memory Cache**: Uses a Go `map[string]cacheItem` with automatic TTL (Time-To-Live)
- Stores daily summaries in memory for 7 days
- Has a background goroutine that cleans up expired entries every minute
- Works exactly like Redis but in-process memory

**Code Location**: `internal/cache/inmem.go`

```go
// In-memory cache structure
type InMemCache struct {
    mu    sync.RWMutex
    items map[string]cacheItem  // Simple Go map
}

// Automatic cleanup every minute
func (c *InMemCache) cleanup() {
    ticker := time.NewTicker(1 * time.Minute)
    // Removes expired entries
}
```

**Result**: The app works perfectly without Redis - it just uses Go's built-in data structures in memory.

---

## 2. What Happens Without Slack?

The app is **fully functional** without Slack - it just won't send notifications.

### Error Logging
- If `SLACK_ERROR_WEBHOOK_URL` is not set, `slackError` is `nil`
- Error logging checks: `if s.slackErr != nil` before sending
- **Result**: Errors are still logged to console, just not sent to Slack

**Code**: `internal/summary/service.go:112-120`
```go
func (s *Service) LogError(ctx context.Context, err error, errorContext map[string]string) {
    if s.slackErr != nil {  // Only sends if Slack is configured
        // Send to Slack
    }
    // Otherwise, errors are just logged to console
}
```

### Daily Summaries
- If `SLACK_WEBHOOK_URL` is not set, `slackSummary` is `nil`
- The scheduler only starts if Slack is configured: `if slackSummary != nil`
- **Result**: Daily summaries won't run automatically, but you can still:
  - Generate summaries manually via API endpoint
  - View summaries in the API response
  - They just won't be sent to Slack

**Code**: `cmd/api/main.go:81-84`
```go
if slackSummary != nil {
    summaryScheduler.Start()  // Only starts if Slack configured
    defer summaryScheduler.Stop()
}
```

**Summary**: 
- ✅ App works perfectly without Slack
- ✅ Errors still logged to console
- ✅ All API endpoints work
- ❌ No Slack notifications
- ❌ No automatic daily summaries (but manual generation works)

---

## 3. AI API Key - Required, But Free Options Available

### Current Status
**AI API Key IS REQUIRED** - the app will exit if not provided:
```go
if cfg.AIApiKey == "" {
    log.Fatal("AI_API_KEY is required")
}
```

### Free/Open Source AI Options

#### Option 1: Local AI Models (Free, No API Key)
You can run AI models locally using:

**Ollama** (Recommended - Easy Setup)
```bash
# Install Ollama
brew install ollama

# Start Ollama
ollama serve

# Pull a model (in another terminal)
ollama pull llama2
ollama pull mistral

# Update config to use local Ollama
export AI_BASE_URL=http://localhost:11434
export AI_API_KEY=ollama  # Ollama doesn't require real auth, but we need a value
export AI_MODEL=llama2
```

**Note**: You'll need to modify the OpenAI client to work with Ollama's API format, or use Ollama's OpenAI-compatible endpoint.

#### Option 2: Free Tier APIs

**OpenAI** (Limited Free Tier)
- Sign up at https://platform.openai.com
- Get $5 free credit (usually enough for testing)
- Set: `AI_BASE_URL=https://api.openai.com`
- Set: `AI_API_KEY=sk-...` (your key)

**Anthropic Claude** (Free Trial)
- Sign up at https://console.anthropic.com
- Free tier available
- Would need API adapter (currently supports OpenAI format)

**Hugging Face Inference API** (Free Tier)
- Sign up at https://huggingface.co
- Free inference API available
- Would need API adapter

#### Option 3: Mock AI Client (For Testing)

You could create a mock AI client that returns sample data for testing:

```go
// internal/ai/mock.go
type MockAIClient struct{}

func (m *MockAIClient) ProcessText(ctx context.Context, req *ProcessRequest) (*domain.ProcessResponse, error) {
    return &domain.ProcessResponse{
        Topic: "Mock Topic",
        Summary: "This is a mock response for testing",
        KeyPoints: []string{"Point 1", "Point 2"},
        // ... sample data
    }, nil
}
```

---

## Minimum Requirements to Run

### Absolute Minimum (In-Memory Mode)
```bash
# Only need:
export AI_API_KEY=your-key-here  # REQUIRED
make run
```

### What You Get:
- ✅ Full API functionality
- ✅ Text processing
- ✅ In-memory storage
- ✅ In-memory caching
- ✅ All endpoints work
- ❌ No Slack notifications
- ❌ No persistent storage (data lost on restart)
- ❌ No Redis caching (but in-memory cache works)

### Recommended for Development
```bash
# With Docker (PostgreSQL + Redis)
make up
export STORAGE=postgres
export DATABASE_URL=postgres://learnforge:learnforge@localhost:5432/learnforge?sslmode=disable
export REDIS_URL=redis://localhost:6379/0
export AI_API_KEY=your-key-here
make run-pg
```

---

## Summary Table

| Component | Required? | What Happens Without It? |
|-----------|-----------|--------------------------|
| **AI API Key** | ✅ **YES** | App won't start (`log.Fatal`) |
| **Redis** | ❌ No | Uses in-memory cache (works fine) |
| **PostgreSQL** | ❌ No | Uses in-memory store (data lost on restart) |
| **Slack Webhooks** | ❌ No | No notifications, but app works |
| **Summary API Key** | ❌ No | Manual summary endpoint disabled |

---

## Quick Test Without AI Key (Development Only)

If you want to test the app structure without an AI key, you could temporarily comment out the check:

```go
// Temporarily disable for testing
// if cfg.AIApiKey == "" {
//     log.Fatal("AI_API_KEY is required")
// }
```

But you'll get errors when trying to process text since the AI client needs the key.

