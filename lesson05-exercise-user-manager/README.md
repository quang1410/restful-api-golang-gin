# User Manager — Layered Architecture

A Go REST API project demonstrating clean layered architecture with Gin framework.

---

## Project Structure

```
lesson05-exercise-user-manager/
├── cmd/
│   └── api/
│       └── main.go               # Entry point: init config, app, server
├── internal/
│   ├── app/
│   │   ├── app.go                # Wire all modules, start server
│   │   └── user_module.go        # Wire user-specific dependencies
│   ├── config/
│   │   └── config.go             # App config (port, env variables)
│   ├── logs/                     # Log files (access.log, error.log)
│   ├── models/
│   │   └── user.go               # User struct / data definition
│   ├── repository/
│   │   ├── interfaces.go         # Repository interface contracts
│   │   └── user_repository.go    # User data access implementation
│   ├── service/
│   │   ├── interfaces.go         # Service interface contracts
│   │   └── user_service.go       # User business logic
│   ├── handler/
│   │   └── user_handler.go       # HTTP request/response handling
│   ├── routes/
│   │   ├── routes.go             # Root router setup
│   │   └── user_routes.go        # User route definitions
│   ├── middleware/
│   │   └── auth.go               # Auth middleware
│   ├── utils/
│   │   ├── errors.go             # Error helpers
│   │   └── string.go             # String helpers
│   └── validation/
│       ├── custom_validation.go  # Custom validation rules
│       └── validator.go          # Validator setup
```

---

## Architecture

### Request Flow

```
HTTP Request
     │
     ▼
  routes        ← Map URL + HTTP method to handler
     │
     ▼
  handler       ← Parse request, call service, write response
     │
     ▼
  service       ← Business logic (validate, process, orchestrate)
     │
     ▼
 repository     ← Read / write data
     │
     ▼
   model        ← Data structure definition
```

### Layer Responsibilities

| Layer | File | Responsibility |
|---|---|---|
| **model** | `models/user.go` | Defines `User` struct and fields |
| **repository** | `repository/user_repository.go` | CRUD operations (create, read, update, delete) |
| **service** | `service/user_service.go` | Business rules, validation, orchestration |
| **handler** | `handler/user_handler.go` | Bind request body, call service, return JSON |
| **routes** | `routes/user_routes.go` | Register URL paths and HTTP methods |

---

## Dependency Injection

Dependencies are initialized **bottom-up** — each layer receives the layer below it via constructor injection:

```
model   → defines the data shape

repo    = user_repository(model)   ← depends on model
service = user_service(repo)       ← depends on repository
handler = user_handler(service)    ← depends on service
routes  = user_routes(handler)     ← depends on handler
```

This wiring is done in `internal/app/user_module.go`.

---

## Key Principles

- **Each layer only knows the layer directly below it** — handler never calls repository directly.
- **Depend on interfaces, not concrete types** — `service` and `repository` are defined as interfaces, making them easy to test and swap.
- **Data flows down, results flow up** — request travels from routes → model; response travels back from model → routes.

---

## Running

```bash
cd lesson05-exercise-user-manager
go run cmd/api/main.go
```

## API Endpoints

| Method | Path | Description |
|---|---|---|
| GET | `/users` | List all users |
| GET | `/users/:id` | Get user by ID |
| POST | `/users` | Create new user |
| PUT | `/users/:id` | Update user |
| DELETE | `/users/:id` | Delete user |
