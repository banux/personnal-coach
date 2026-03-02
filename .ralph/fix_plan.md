# Ralph Fix Plan

## High Priority
- [x] Set up basic project structure and build system
- [x] Define core data structures and types
- [x] Implement basic input/output handling
- [ ] Create test framework and initial tests

## Medium Priority
- [ ] Add error handling and validation
- [x] Implement core business logic (Claude SDK integration)
- [ ] Add configuration management (env vars / config file)
- [ ] Create user documentation

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

## Next Steps
1. Add Go tests for services (claude, timer, pdf)
2. Add env configuration (ANTHROPIC_API_KEY, PORT validation)
3. Add persistence (SQLite or file-based storage)
4. Deploy configuration (docker-compose or scripts)

## Notes
- Go env: GOROOT=/home/banux/go, GOPATH=/home/banux/go (same dir, expected warning)
- Build: cd backend && GOROOT=/home/banux/go GOPATH=/home/banux/go go build ./...
- Frontend: cd frontend && npm run build
- Backend runs on port 8080, frontend dev on 5173
