# Task API

Cloud-native REST API for task management, built with Go and PostgreSQL, following the [12-Factor App](https://12factor.net/) methodology.

---

## Stack

- **Go 1.22** — runtime
- **chi** — HTTP router
- **lib/pq** — PostgreSQL driver
- **swaggo/swag** — OpenAPI/Swagger docs
- **Docker** — containerization

---

## Local Development

### Prerequisites

- Docker & Docker Compose
- Go 1.22+
- [`swag`](https://github.com/swaggo/swag) CLI (for doc generation)

### Quick start with Docker Compose

```bash
cp .env.example .env
docker compose up --build
```

API available at `http://localhost:8080`
Swagger UI at `http://localhost:8080/swagger/index.html`

### Running without Docker

```bash
# 1. Start only postgres
docker compose up postgres -d

# 2. Copy and source env file
cp .env.example .env
export $(grep -v '^#' .env | xargs)

# 3. Generate swagger docs
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g cmd/api/main.go -o docs

# 4. Run
go run ./cmd/api
```

---

## Environment Variables

All configuration is done via environment variables (12-Factor III).

| Variable              | Required | Default       | Description                           |
|-----------------------|----------|---------------|---------------------------------------|
| `APP_ENV`             | No       | `development` | Environment name (dev/staging/prod)   |
| `APP_NAME`            | No       | `task-api`    | App name used in logs                 |
| `SERVER_PORT`         | No       | `8080`        | HTTP port                             |
| `SERVER_READ_TIMEOUT` | No       | `15s`         | Max time to read a request            |
| `SERVER_WRITE_TIMEOUT`| No       | `15s`         | Max time to write a response          |
| `DB_HOST`             | No       | `localhost`   | PostgreSQL hostname                   |
| `DB_PORT`             | No       | `5432`        | PostgreSQL port                       |
| `DB_NAME`             | **Yes**  | —             | Database name                         |
| `DB_USER`             | **Yes**  | —             | Database user                         |
| `DB_PASSWORD`         | **Yes**  | —             | Database password                     |
| `DB_SSLMODE`          | No       | `disable`     | SSL mode (disable/require/verify-full)|
| `DB_MAX_OPEN_CONNS`   | No       | `25`          | Max open DB connections               |
| `DB_MAX_IDLE_CONNS`   | No       | `25`          | Max idle DB connections               |
| `DB_CONN_MAX_LIFETIME`| No       | `5m`          | Connection max lifetime               |

---

## API Reference

Base path: `/api/v1`

| Method | Endpoint                              | Description               |
|--------|---------------------------------------|---------------------------|
| POST   | `/tasks`                              | Create a new task         |
| GET    | `/tasks`                              | List all tasks            |
| GET    | `/tasks/status?completed=true\|false` | Filter tasks by status    |
| GET    | `/tasks/{id}`                         | Get a task by ID          |
| PUT    | `/tasks/{id}`                         | Update a task             |

### Task schema

```json
{
  "id": "uuid",
  "description": "string",
  "isCompleted": false,
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

### Create task — `POST /api/v1/tasks`

```json
{ "description": "Buy groceries" }
```

### Update task — `PUT /api/v1/tasks/{id}`

All fields are optional:

```json
{ "description": "Buy groceries and cook", "isCompleted": true }
```

### Health probes

| Endpoint      | Description      |
|---------------|------------------|
| `GET /healthz` | Liveness check  |
| `GET /readyz`  | Readiness check |

Full interactive docs at `/swagger/index.html`.

---

## CI/CD (GitHub Actions)

Three jobs run on every push:

1. **lint** — runs `golangci-lint`
2. **build** — compiles the binary (lint must pass)
3. **docker** — builds and pushes `juanozorio/task-api` to Docker Hub (main branch only)

### Required GitHub Secrets

| Secret               | Value                                   |
|----------------------|-----------------------------------------|
| `DOCKERHUB_USERNAME` | Your Docker Hub username (`juanozorio`) |
| `DOCKERHUB_TOKEN`    | Docker Hub access token (not password)  |

Create the token at: https://hub.docker.com/settings/security

---

## Project Structure

```
.
├── cmd/api/            # Application entrypoint
├── internal/
│   ├── config/         # 12-factor env-based config
│   ├── database/       # DB connection & migrations
│   ├── domain/         # Models + repository interfaces
│   ├── handler/        # HTTP handlers & middleware
│   ├── repository/     # PostgreSQL implementations
│   └── service/        # Business logic
├── docs/               # Generated swagger files (git-ignored)
├── .github/workflows/  # CI/CD pipelines
├── docker-compose.yml  # Local development stack
├── Dockerfile          # Multi-stage production image
├── .env.example        # Environment variable reference
└── .golangci.yml       # Linter configuration
```

---

## 12-Factor Compliance

| Factor | Implementation |
|--------|---------------|
| I. Codebase | Single repo, multiple deploys via image tags |
| II. Dependencies | Explicit in `go.mod` / `go.sum` |
| III. Config | All config via environment variables |
| IV. Backing services | PostgreSQL treated as attached resource via env |
| V. Build/release/run | Separated via Dockerfile stages |
| VI. Processes | Stateless, share-nothing HTTP server |
| VII. Port binding | Exports service via `SERVER_PORT` |
| VIII. Concurrency | Horizontal scaling via container orchestration |
| IX. Disposability | Graceful shutdown with 30s drain |
| X. Dev/prod parity | Same Docker image across environments |
| XI. Logs | Structured JSON to stdout via `slog` |
| XII. Admin processes | Migrations run at startup automatically |
