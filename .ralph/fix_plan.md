# Ralph Fix Plan

## High Priority
- [x] Set up basic project structure and build system
- [x] Define core data structures and types
- [x] Implement basic input/output handling
- [x] Authentication (simple app-wide password form)
- [x] Frontend served by Go server (embedded via go:embed)
- [x] Docker support (Dockerfile + docker-compose.yml)
- [ ] Create test framework and initial tests

## Medium Priority
- [ ] Add error handling and validation (input sanitization)
- [x] Implement core business logic (Claude SDK integration)
- [x] Add configuration management (env vars: ANTHROPIC_API_KEY, APP_PASSWORD, PORT)
- [ ] Create user documentation (README.md)

## Low Priority
- [ ] Performance optimization
- [ ] Extended feature set (multi-person programs)
- [x] Integration with external services (Claude AI, PDF)
- [ ] Advanced error recovery

## Completed
- [x] Project initialization
- [x] Go backend: models, services, handlers, MCP server
- [x] VueJS frontend with Tailwind CSS
- [x] Claude SDK integration for program generation
- [x] PDF generation with gofpdf
- [x] Exercise timer service
- [x] MCP server with generate_workout_program and get_workout_timer tools
- [x] REST API: POST /api/programs/generate, GET /api/programs/:id, GET /api/programs/:id/pdf, GET /api/programs/:id/timer/:day
- [x] Authentication: session cookie (APP_PASSWORD env var, default: coach2024)
- [x] Frontend embedded in Go binary via go:embed (single binary deployment)
- [x] Docker: multi-stage Dockerfile + docker-compose.yml + .env.example
- [x] Build script: scripts/build.sh

## Next Steps
1. Write README.md with setup and usage instructions
2. Add Go unit tests for timer service
3. Add persistence (SQLite file-based) for programs across restarts

## Notes
- Go env: GOROOT=/home/banux/go, GOPATH=/home/banux/go (same dir, expected warning)
- Build: cd backend && GOROOT=/home/banux/go GOPATH=/home/banux/go go build ./...
- Frontend build: cd frontend && npm run build && cp -r dist ../backend/dist
- Full build script: ./scripts/build.sh
- Docker: docker-compose up -d (requires ANTHROPIC_API_KEY in .env)
- Version: 1.1.0
