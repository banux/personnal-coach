# Agent Build Instructions

## Project Overview
Personal Coach AI - Generates personalized workout programs using Claude AI.
- **Backend**: Go (Gin framework) on port 8080
- **Frontend**: VueJS + Tailwind CSS (Vite) on port 5173
- **MCP Server**: Integrated in Go binary (pass `mcp` arg)

## Environment Setup

### Critical Go Environment
The Go installation uses a non-standard path. Always use these env vars:
```bash
export GOROOT=/home/banux/go
export GOPATH=/home/banux/go
export GOPROXY=direct
export GONOSUMDB="*"
```

### Required Environment Variables
```bash
export ANTHROPIC_API_KEY="your-api-key-here"
```

## Quick Start

```bash
# 1. Copy and configure env
cp .env.example .env
# Edit .env: set ANTHROPIC_API_KEY and APP_PASSWORD

# 2. Full build (frontend + backend with embedded dist)
./scripts/build.sh

# 3. Run
ANTHROPIC_API_KEY=your-key APP_PASSWORD=yourpassword ./personal-coach
# Open http://localhost:8080
```

## Docker

```bash
# Build and run with docker-compose
cp .env.example .env  # Set ANTHROPIC_API_KEY
docker-compose up -d

# Rebuild after code changes
docker-compose up -d --build

# Run MCP server (for Claude Desktop integration)
docker-compose --profile mcp up personal-coach-mcp
```

## Running the Backend

```bash
cd /home/banux/personal-coach/backend

# IMPORTANT: Copy frontend dist before running
cp -r ../frontend/dist ./dist

# Development
GOROOT=/home/banux/go GOPATH=/home/banux/go ANTHROPIC_API_KEY=your-key go run main.go

# Build binary
GOROOT=/home/banux/go GOPATH=/home/banux/go go build -o ../personal-coach .

# Run MCP server mode
ANTHROPIC_API_KEY=your-key ./personal-coach mcp
```

## Running the Frontend (dev only)

The frontend is normally served by the Go server (embedded in binary).
For development with hot-reload:

```bash
cd /home/banux/personal-coach/frontend

# Development (hot reload) - connects to Go backend on :8080
npm run dev

# Production build + copy to backend
npm run build && cp -r dist ../backend/dist
```

## Build Commands

```bash
# Build backend
cd backend && GOROOT=/home/banux/go GOPATH=/home/banux/go go build ./...

# Build frontend
cd frontend && npm run build

# Build all
cd backend && GOROOT=/home/banux/go GOPATH=/home/banux/go go build ./... && cd ../frontend && npm run build
```

## Running Tests

```bash
# Backend tests (when written)
cd backend && GOROOT=/home/banux/go GOPATH=/home/banux/go go test ./...

# Frontend tests (when written)
cd frontend && npm test
```

## Project Structure

```
personal-coach/
├── backend/
│   ├── main.go              # Entry point (HTTP or MCP server)
│   ├── go.mod
│   ├── handlers/
│   │   └── program.go       # REST API handlers
│   ├── models/
│   │   └── models.go        # Data structures
│   ├── services/
│   │   ├── claude.go        # Claude AI integration
│   │   ├── pdf.go           # PDF generation
│   │   └── timer.go         # Timer sequence builder
│   └── mcp/
│       └── server.go        # MCP server (stdio transport)
└── frontend/
    ├── src/
    │   ├── App.vue           # Root component (nav + router)
    │   ├── main.js           # App bootstrap (Pinia + Router)
    │   ├── style.css         # Tailwind + custom components
    │   ├── stores/
    │   │   └── program.js    # Pinia store (API calls)
    │   ├── views/
    │   │   ├── HomeView.vue  # Form to create program
    │   │   └── ProgramView.vue # Display program + timer
    │   └── components/
    │       └── TimerModal.vue # Workout timer modal
    ├── package.json
    └── vite.config.js
```

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| POST | /api/programs/generate | Generate a new program via Claude |
| GET | /api/programs | List all programs |
| GET | /api/programs/:id | Get a specific program |
| GET | /api/programs/:id/pdf | Download program as PDF |
| GET | /api/programs/:id/timer/:day | Get timer for a day |
| GET | /health | Health check |

## MCP Tools

- `generate_workout_program` - Generate personalized program
- `get_workout_timer` - Build timer sequence from a program

## Key Learnings
- Go GOROOT and GOPATH must both be set to `/home/banux/go`
- Use GOPROXY=direct and GONOSUMDB="*" for dependency installation
- Frontend API URL configurable via VITE_API_URL env var
- gofpdf produces PDFs with structured tables (exercise per row)
- Timer uses AudioContext for beep sounds in browser
- MCP server uses JSON-RPC 2.0 over stdio

## Feature Development Quality Standards

**CRITICAL**: All new features MUST meet the following mandatory requirements before being considered complete.

### Testing Requirements

- **Minimum Coverage**: 85% code coverage ratio required for all new code
- **Test Pass Rate**: 100% - all tests must pass, no exceptions
- **Coverage Validation**: Run coverage reports before marking features complete

### Git Workflow Requirements

Before moving to the next feature, ALL changes must be committed with clear messages:
```bash
git add .
git commit -m "feat(module): descriptive message following conventional commits"
```

### Feature Completion Checklist

Before marking ANY feature as complete, verify:

- [ ] All tests pass
- [ ] Code coverage meets 85% minimum threshold
- [ ] Code formatted according to project standards
- [ ] All changes committed with conventional commit messages
- [ ] .ralph/fix_plan.md task marked as complete
- [ ] Implementation documentation updated
- [ ] .ralph/AGENT.md updated (if new patterns introduced)
