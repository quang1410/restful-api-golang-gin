# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Run the server (from project root)
cd cmd/api && go run main.go

# Build
go build ./...

# Run tests
go test ./...

# Run a single test
go test ./internal/... -run TestName

# Tidy dependencies
go mod tidy
```

The server port is read from the `SERVER_PORT` environment variable (`.env` file at repo root, two levels up from `cmd/api/`).

All requests require the `X-API-Key` header matching the `API_KEY` env var (falls back to `"secret-key"`).

## Architecture

The app uses a **module-based layered architecture**: each domain (e.g., User) is wired up as a self-contained `Module` in `internal/app/`, composing its own repository → service → handler → routes chain. New domains follow the same pattern: add a `*_module.go` file in `internal/app/` and register it in `app.go`.

**Request lifecycle:**
```
HTTP → LoggerMiddleware → ApiKeyMiddleware → AuthMiddleware → RateLimiterMiddleware → Handler → Service → Repository
```

**Layer responsibilities:**
- `internal/handler/v1/` — binds and validates HTTP input via `ShouldBind*`, calls service, writes response via `utils` helpers
- `internal/service/v1/` — business logic; depends on `repository` interfaces
- `internal/repository/` — data access; `interfaces.go` defines the contract, `*_repository.go` is the SQL implementation
- `internal/dto/v1/` — input structs with binding tags; DTOs have `MapInputToModel()` methods for conversion
- `internal/routes/v1/` — registers route groups under `/api/v1`

**Validation:** Custom validators are registered once at startup via `validation.InitValidator()`. Add new validators in `internal/validation/custom_validation.go` and map their error messages in `validation.HandleValidationErrors()`. Validation error messages are in Vietnamese.

**Response helpers** (`internal/utils/response.go`):
- `ResponseSuccess(ctx, status, data)` — wraps data in `{"status":"success","data":...}`
- `ResponseError(ctx, err)` — unwraps `*AppError` to get HTTP status from `ErrorCode`
- `ResponseValidator(ctx, data)` — 400 with raw `gin.H` (used after `HandleValidationErrors`)
- `ResponseStatusCode(ctx, status)` — status only (e.g., 204 No Content)

**Logging:** `LoggerMiddleware` writes structured JSON to `internal/logs/http.log` via zerolog + lumberjack (1 MB max, 5 backups, 5-day retention). It captures full request/response bodies, including multipart form data and files.

**Rate limiting:** Per-IP token bucket — 5 req/s, burst 15. Clients inactive for 3+ minutes are purged by a background goroutine (`CleanupClients()`).

**`AuthMiddleware`** is a no-op placeholder; JWT or session logic goes there.

## Adding a New Domain

1. Create `internal/repository/interfaces.go` entry or a new file for the new interface.
2. Implement the repository in `internal/repository/<domain>_repository.go`.
3. Add service interface to `internal/service/v1/interfaces.go` and implement in `user_service.go`-style file.
4. Add handler in `internal/handler/v1/<domain>_handler.go`.
5. Add routes in `internal/routes/v1/<domain>_routes.go` implementing `routes.Route`.
6. Wire everything in `internal/app/<domain>_module.go` and append to the `modules` slice in `app.go`.
