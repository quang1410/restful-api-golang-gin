# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

All commands must be run from this directory (`lesson03-route-group/`).

```bash
# Run the server (port 8081)
go run main.go

# Manage dependencies
go mod tidy

# Test rate limiting (burst = 15, then 429)
ab -n 50 -c 20 -H "X-API-Key:ab2a7a8a-d601-4bf7-b0e2-dd00e5459392" http://localhost:8081/api/v1/users
```

## Architecture

This module demonstrates Gin middleware and dual API versioning (v1/v2) with no database — all handlers return mock data.

**Request pipeline (global, applied to all routes):**
```
LoggerMiddleware → ApiKeyMiddleware → RateLimitingMiddleware → Handler
```
`SimpleMiddleware` is route-level only (applied to `GET /api/v1/news/:slug`).

**Middleware:**
- `ApiKeyMiddleware` — reads `API_KEY` from `.env` (default `"secret-key"`); requires `X-API-Key` header; sets `"username"` on the context
- `RateLimitingMiddleware` — per-IP token bucket via `golang.org/x/time/rate`: 5 req/s, burst 15; background goroutine (`CleanupClients`) evicts IPs idle > 3 minutes
- `LoggerMiddleware` — structured JSON logs to `logs/http.log` via zerolog + lumberjack (1 MB, 5 backups, 5-day retention); captures request body across all content types (JSON, form, multipart) and response body
- `SimpleMiddleware` — demo middleware that writes directly to the response before and after `ctx.Next()`

**Handlers (`internal/api/`):**
- `v1/handler/` — user, product, category, news; mock responses only
- `v2/handler/` — user only (different response shape to contrast versioning)
- Input validation uses `ShouldBind` / `ShouldBindUri` with struct tags; errors are forwarded to `utils.HandleValidationErrors`

**Utils:**
- `utils.HandleValidationErrors` — maps `validator.ValidationErrors` to a `gin.H{"error": map[field]message}` with Vietnamese messages; field paths are converted from CamelCase to snake_case
- `utils.RegisterValidators` — registers custom validators: `slug`, `search`, `min_int`, `max_int`, `file_ext`; must be called at startup before Gin handles any request
- `utils.ValidateAndSaveFile` — validates and saves a single `*multipart.FileHeader` to a given directory

**File uploads:**
- `POST /api/v1/news` — single file inline (5 MB max, saved to `./uploads/`)
- `POST /api/v1/news/upload-file` — single file via `utils.ValidateAndSaveFile`
- `POST /api/v1/news/upload-multiple-file` — multiple files (`images[]`), partial-success response
- Uploaded files served statically at `/images/` → `./uploads/`

**Config:** `API_KEY` is loaded from `.env` via `godotenv`. The `.env` file must exist in this directory when running the server.
