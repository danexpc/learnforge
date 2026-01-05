package http

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
)

// RegisterWebRoutes serves the web UI from filesystem or placeholder
func (h *Handler) RegisterWebRoutes(r chi.Router) {
	// Check if web/dist exists
	webDistPath := "web/dist"
	if _, err := os.Stat(webDistPath); os.IsNotExist(err) {
		// If web/dist doesn't exist, serve a placeholder
		r.Get("/", h.serveWebPlaceholder)
		return
	}

	// Serve static files from filesystem
	workDir, _ := os.Getwd()
	webDistFullPath := filepath.Join(workDir, webDistPath)
	fileServer := http.FileServer(http.Dir(webDistFullPath))

	// Serve static files with SPA fallback
	// This catch-all route should be registered last to not interfere with API routes
	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// Skip API routes - these should already be handled by other handlers
		// This is a safety check in case the route order changes
		if strings.HasPrefix(path, "/v1") ||
			strings.HasPrefix(path, "/docs") ||
			strings.HasPrefix(path, "/openapi") ||
			path == "/metrics" ||
			path == "/healthz" ||
			path == "/readyz" {
			// Let chi router continue to next handler (shouldn't happen if routes are registered correctly)
			return
		}

		// Try to serve the requested file
		filePath := strings.TrimPrefix(path, "/")
		if filePath == "" {
			filePath = "index.html"
		}

		// Check if file exists
		fullPath := filepath.Join(webDistFullPath, filePath)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			// If file doesn't exist, serve index.html for SPA routing
			r.URL.Path = "/index.html"
		}

		// Serve the file
		fileServer.ServeHTTP(w, r)
	})
}

// serveWebPlaceholder serves a placeholder page when web/dist is not built
func (h *Handler) serveWebPlaceholder(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html>
<head>
    <title>LearnForge - Web UI</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            display: flex;
            align-items: center;
            justify-content: center;
            min-height: 100vh;
            margin: 0;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
        }
        .container {
            text-align: center;
            padding: 2rem;
        }
        h1 { font-size: 2.5rem; margin-bottom: 1rem; }
        p { font-size: 1.2rem; opacity: 0.9; }
        code {
            background: rgba(0,0,0,0.2);
            padding: 0.2rem 0.5rem;
            border-radius: 4px;
            font-family: monospace;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>🚀 LearnForge</h1>
        <p>Web UI is not built yet.</p>
        <p>Run <code>cd web && npm install && npm run build</code> to build the UI.</p>
        <p style="margin-top: 2rem;">
            <a href="/docs" style="color: white; text-decoration: underline;">View API Documentation →</a>
        </p>
    </div>
</body>
</html>`
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

