package http

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) RegisterWebRoutes(r chi.Router) {
	webDistPath := "web/dist"
	if _, err := os.Stat(webDistPath); os.IsNotExist(err) {
		r.Get("/", h.serveWebPlaceholder)
		return
	}

	workDir, _ := os.Getwd()
	webDistFullPath := filepath.Join(workDir, webDistPath)
	fileServer := http.FileServer(http.Dir(webDistFullPath))

	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if strings.HasPrefix(path, "/v1") ||
			strings.HasPrefix(path, "/docs") ||
			strings.HasPrefix(path, "/openapi") ||
			path == "/metrics" ||
			path == "/healthz" ||
			path == "/readyz" {
			return
		}

		filePath := strings.TrimPrefix(path, "/")
		if filePath == "" {
			filePath = "index.html"
		}

		fullPath := filepath.Join(webDistFullPath, filePath)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			r.URL.Path = "/index.html"
		}

		fileServer.ServeHTTP(w, r)
	})
}

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

