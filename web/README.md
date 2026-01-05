# LearnForge Web UI

Modern Preact-based web interface for LearnForge with visual content display, meme generation, and shareable URLs.

## Features

- **Visual Content Display**: Beautiful UI for viewing generated learning content
- **Meme Generation**: Optional beta feature to generate educational memes
- **Shareable URLs**: Content IDs in URLs for bookmarking and sharing
- **Responsive Design**: Works on desktop and mobile devices
- **Real-time Updates**: Loading states and smooth transitions

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

Outputs optimized files to `dist/` directory, which are served by the Go backend.

**Note**: The Go server automatically serves files from `web/dist/` if it exists. If not, it shows a placeholder page with build instructions.

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

