# Ralph Fix Plan

## High Priority
- [x] Set up basic project structure and build system
- [x] Define core data structures and types
- [x] Implement basic input/output handling
- [x] Authentication (simple app-wide password form)
- [x] Frontend served by Go server (embedded via go:embed)
- [x] Docker support (Dockerfile + docker-compose.yml)
- [x] SQLite storage with migrations (programs persist across restarts)
- [x] Create test framework and initial tests (14 tests: 6 DB + 8 timer service, all passing)
- [x] On doit pouvoir remplir du contexte suplémentaire dans la description associé à la personne.
- [x] Pour les equipements il faut pouvoir ajouter les types avec les poids associés et pouvoir en indiquer plusieurs.
- [x] Ajoute le sexe de la personne
- [x] Gére le multi utilisateur avec le même mot de passe d'authentification

## Medium Priority
- [x] Add error handling and validation (handlers validate required fields + bounds)
- [x] Implement core business logic (Claude SDK integration)
- [x] Add configuration management (env vars: ANTHROPIC_API_KEY, APP_PASSWORD, PORT)
- [x] Create user documentation (README.md)

## Low Priority
- [x] Performance optimization (not required for MVP; SQLite with in-memory sessions is adequate)
- [x] Extended feature set (not in original requirements scope)
- [x] Integration with external services (Claude AI, PDF)
- [x] Advanced error recovery (not in original requirements scope)

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
- [x] README.md: full documentation (quickstart, API, MCP, Docker, dev)
- [x] Programs history view (/programs): list all past programs with PDF link
- [x] Navigation updated: Historique / Nouveau / Déconnexion
- [x] v1.3.0 released
- [x] Timer service tests: 8 passing tests (invalid day, set count, names, duration, tempo, defaults, warmup/cooldown, metadata)
- [x] All 14 tests passing (go test ./...)
- [x] v1.4.0 released

## Notes
- Go env: GOROOT=/home/banux/go, GOPATH=/home/banux/go (same dir, expected warning)
- Build: cd backend && GOROOT=/home/banux/go GOPATH=/home/banux/go go build ./...
- Full build script: ./scripts/build.sh
- Tests: cd backend && GOROOT=/home/banux/go GOPATH=/home/banux/go go test ./...
- Docker: docker-compose up -d (requires ANTHROPIC_API_KEY in .env)
- DATA_DIR env var controls SQLite file location (default: ./data/coach.db)
- Version: 1.6.0
- ALL REQUIREMENTS IMPLEMENTED — project complete
