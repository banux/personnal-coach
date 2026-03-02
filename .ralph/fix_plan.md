# Ralph Fix Plan

## High Priority
- [x] Set up basic project structure and build system
- [x] Define core data structures and types
- [x] Implement basic input/output handling
- [x] Authentication (simple app-wide password form)
- [x] Frontend served by Go server (embedded via go:embed)
- [x] Docker support (Dockerfile + docker-compose.yml)
- [x] SQLite storage with migrations (programs persist across restarts)
- [ ] Create test framework and initial tests (partial: DB tests done)

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
- [x] SQLite persistence: database package with auto-migrations (modernc.org/sqlite, CGO-free)
- [x] DB tests: 6 passing tests for CRUD + migration idempotency
- [x] v1.2.0 released

## Next Steps
1. Write README.md with setup and usage instructions
2. Add person history view to frontend (list all past programs)

## Notes
- Go env: GOROOT=/home/banux/go, GOPATH=/home/banux/go (same dir, expected warning)
- Build: cd backend && GOROOT=/home/banux/go GOPATH=/home/banux/go go build ./...
- Full build script: ./scripts/build.sh
- Tests: cd backend && GOROOT=/home/banux/go GOPATH=/home/banux/go go test ./database/...
- Docker: docker-compose up -d (requires ANTHROPIC_API_KEY in .env)
- DATA_DIR env var controls SQLite file location (default: ./data/coach.db)
- Version: 1.2.0
