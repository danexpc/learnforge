# LearnForge Web UI

Modern Preact-based web interface for LearnForge.

## Development

```bash
# Install dependencies
npm install

# Start dev server (runs on http://localhost:5173)
npm run dev
```

The dev server proxies API requests to `http://localhost:8080`.

## Building

```bash
# Build for production
npm run build
```

Outputs optimized files to `dist/` directory, which are then embedded into the Go binary.

## Project Structure

```
web/
├── src/
│   ├── components/     # Reusable UI components
│   ├── pages/         # Page components
│   ├── utils/         # Utilities (API client, etc.)
│   ├── App.jsx        # Main app component
│   └── main.jsx       # Entry point
├── index.html         # HTML template
├── package.json       # Dependencies
├── vite.config.js     # Vite configuration
└── tailwind.config.js # Tailwind CSS configuration
```

## Tech Stack

- **Preact** - Lightweight React alternative
- **Vite** - Fast build tool and dev server
- **Tailwind CSS** - Utility-first CSS framework

