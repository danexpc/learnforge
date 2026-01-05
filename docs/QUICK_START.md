# Quick Start Guide - Running LearnForge Locally

## Prerequisites

✅ You have an OpenAI API key  
✅ Go 1.22+ installed  
✅ That's it! No Redis, no PostgreSQL needed for basic testing.

## Step-by-Step

### 1. Set Your API Key

```bash
export AI_API_KEY=sk-proj-your-key-here
```

**Note**: Replace `sk-proj-your-key-here` with your actual OpenAI API key.

### 2. Run the Application

```bash
make run
```

Or directly:
```bash
go run cmd/api/main.go
```

### 3. Verify It's Running

You should see:
```
{"level":"info","msg":"Using in-memory storage"}
{"level":"info","msg":"Using in-memory cache"}
{"level":"info","msg":"Using OpenAI AI provider"}
{"level":"info","msg":"Starting server","port":"8080"}
```

### 4. Test the Application

**Option A: Use the Web UI (Recommended)**
1. Open `http://localhost:8080` in your browser
2. Enter text, select mode and level
3. Click "Generate Learning Content"
4. View results with shareable URL

**Option B: Test the API**

**Health Check:**
```bash
curl http://localhost:8080/healthz
```

**Process Text:**
```bash
curl -X POST http://localhost:8080/v1/process \
  -H "Content-Type: application/json" \
  -d '{
    "text": "Photosynthesis is how plants make food using sunlight, water, and carbon dioxide.",
    "topic": "Biology",
    "level": "beginner",
    "generate_meme": false
  }'
```

**View API Documentation:**
- Open `http://localhost:8080/docs` in your browser for interactive API docs

## What You Get

- ✅ **In-memory storage** - No database needed
- ✅ **In-memory cache** - No Redis needed  
- ✅ **Full API functionality** - All endpoints work
- ✅ **Web UI** - Modern interface at `http://localhost:8080`
- ✅ **OpenAPI docs** - Interactive docs at `/docs`
- ✅ **OpenAI integration** - Uses your API key
- ✅ **Meme generation** - Optional beta feature
- ❌ **No Slack notifications** (optional)
- ❌ **No persistent storage** (data lost on restart)

## Alternative: Using Config File

Instead of environment variables, you can edit `config/development.yml`:

```yaml
ai_api_key: "sk-proj-your-key-here"
```

Then run:
```bash
ENV=development go run cmd/api/main.go
```

## Troubleshooting

### "AI_API_KEY is required" error
- Make sure you set the environment variable: `export AI_API_KEY=...`
- Or add it to `config/development.yml`

### Port 8080 already in use
```bash
export PORT=8081
make run
```

### API key not working
- Verify your key starts with `sk-` or `sk-proj-`
- Check if you have credits/quota in OpenAI dashboard
- Make sure there are no extra spaces in the key

## Next Steps

Once it's running:
1. **Try the Web UI** - Open `http://localhost:8080` in your browser
2. **Test the API** - Use `/v1/process` endpoint or the web UI
3. **Try different modes** - `lesson`, `flashcards`, `quiz`
4. **Test topic inference** - Don't provide `topic` field and let AI detect it
5. **Try meme generation** - Enable the beta meme feature in the web UI
6. **Check API docs** - Visit `/docs` for interactive OpenAPI documentation
7. **View metrics** - Check `/metrics` for Prometheus metrics
8. **Share content** - Copy the URL with `?id=...` to share generated content

## Using Gemini Instead

If you want to use Gemini:
```bash
export AI_PROVIDER=gemini
export AI_API_KEY=AIzaSy-your-gemini-key
export AI_MODEL=gemini-2.0-flash-exp
make run
```

## Stopping the App

Press `Ctrl+C` in the terminal to stop the server gracefully.

