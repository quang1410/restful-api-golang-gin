# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Repository Structure

This is a Golang learning repository organized as independent modules, each in its own directory with its own `go.mod`. There is no single root module.

| Directory | Purpose |
|---|---|
| `lesson02-gin-starter/` | Basic Gin HTTP server, routes, params, query strings |
| `lesson03-route-group/` | Route groups, middleware, v1/v2 API versioning, file upload |
| `lesson05-exercise-user-manager/` | User management exercise (single file) |
| `lesson07-exercise/` – `lesson18-exercise/` | Go fundamentals: pointers, structs, interfaces, generics, arrays/slices/maps, goroutines, system monitoring |
| `Golang-basic/` | Standalone Go basic examples |
| `project-shopping-cart/` | Full layered REST API project (the main reference architecture) |

## Commands

All commands must be run from inside the relevant module directory (the one containing `go.mod`).

```bash
# Run any standalone lesson
cd lesson<N>-<name> && go run main.go

# Run/build the shopping cart project
cd project-shopping-cart
make run          # go run ./cmd/api/main.go
make build        # outputs binary to ./bin/
make test         # go test ./... -v
make test/run name=TestFoo   # run a single test by name
make tidy         # go mod tidy

# Manage dependencies in any module
go mod tidy
```

## project-shopping-cart Architecture

This is the main reference project. It uses a **module-based layered architecture** inside `internal/`:

```
HTTP → LoggerMiddleware → ApiKeyMiddleware → AuthMiddleware → RateLimiterMiddleware → Handler → Service → Repository
```

**Layers:**
- `cmd/api/main.go` — entrypoint; loads `.env`, wires the app
- `internal/app/` — one `*_module.go` per domain; each composes its own repository → service → handler → routes chain and registers itself in `app.go`
- `internal/handler/v1/` — binds/validates HTTP input, calls service, writes response via `utils` helpers
- `internal/service/v1/` — business logic; depends on repository interfaces
- `internal/repository/` — data access; `interfaces.go` defines contracts
- `internal/dto/v1/` — request structs with binding tags; each has a `MapInputToModel()` method
- `internal/routes/v1/` — registers route groups under `/api/v1`
- `internal/validation/` — custom validators registered at startup; error messages are in Vietnamese

**Adding a new domain:** create repository interface → implement repository → add service interface + implementation → add handler → add routes → wire in a new `internal/app/<domain>_module.go` and append to `modules` in `app.go`.

**Response helpers** (`internal/utils/response.go`):
- `ResponseSuccess(ctx, status, data)` — `{"status":"success","data":...}`
- `ResponseError(ctx, err)` — unwraps `*AppError` to HTTP status via `ErrorCode`
- `ResponseValidator(ctx, data)` — 400 with validation field errors
- `ResponseStatusCode(ctx, status)` — status only (e.g., 204)

**Config:** `SERVER_PORT` and `API_KEY` are read from a `.env` file. All requests require the `X-API-Key` header matching `API_KEY` (default: `"secret-key"`). Rate limiting is per-IP token bucket (5 req/s, burst 15). Logging writes structured JSON to `internal/logs/http.log` via zerolog + lumberjack.

## lesson03-route-group Architecture

Flat structure demonstrating Gin middleware and dual API versioning without a database:
- `middleware/` — API key auth, logger, rate limiter, simple auth
- `internal/api/v1/handler/` and `internal/api/v2/handler/` — handlers per resource (user, product, category, news)
- `utils/` — conversion, file handling, validation helpers

## Go Fundamentals Lessons Pattern

Lessons 07–18 are standalone CLI programs. Each follows the pattern:
- `main.go` — entry point with interactive menu or demo
- Domain packages alongside (`student/`, `library/`, `monitors/`, etc.)
- lesson18 uses goroutines + channels to concurrently collect system stats (CPU, mem, net, disk) via a `Monitor` interface pattern
