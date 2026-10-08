# Sezzle Calculator

[![CI](https://github.com/diego-654/sezzle-calculator/actions/workflows/ci.yml/badge.svg)](https://github.com/diego-654/sezzle-calculator/actions/workflows/ci.yml)

A full-stack calculator: a React + TypeScript frontend that calls a Go REST microservice.
It supports add, subtract, multiply, divide, power, square root and percent.

## Project structure

```
.
├── backend/                     # Go microservice (standard library only)
│   ├── cmd/server/              # entry point: config, HTTP server, graceful shutdown
│   └── internal/
│       ├── calculator/          # pure domain logic, returns (result, error)
│       └── transport/http/      # handlers, JSON, error mapping, middleware
├── frontend/                    # Vite + React + TypeScript, served by nginx in Docker
├── docker-compose.yml           # runs both services together
├── Makefile                     # shortcuts for tests, coverage and Docker
└── .github/workflows/ci.yml     # lint, tests and build on every push and PR
```

## Requirements

- Go 1.27+
- Node.js 22+
- Docker (optional, to run the full stack with one command)
- make (optional; every target is a plain command you can also run by hand)

## Running the app

### Option 1: Docker Compose (recommended)

```bash
docker compose up --build
```

Open http://localhost:3000. nginx serves the frontend and forwards `/api` to the backend,
which is not exposed to the host.

### Option 2: Run each part locally

Backend (listens on port 8080, change it with the `PORT` environment variable):

```bash
cd backend
go run ./cmd/server
```

Frontend, in a second terminal:

```bash
cd frontend
npm install
npm run dev
```

Open http://localhost:5173. The Vite dev server proxies `/api` to `http://localhost:8080`.

## Tests and coverage

```bash
make test        # backend + frontend unit tests
make coverage    # same, with coverage reports
```

Without make:

```bash
cd backend && go test -cover ./...
cd frontend && npm run coverage
```

Current coverage: 94% for the HTTP layer and 65% for the calculator package (backend),
and about 96% of statements in the frontend. The frontend HTML report is written to
`frontend/coverage/`. CI runs the same checks plus `gofmt`, `go vet`, ESLint and a
production build.

## API

All operations use `POST /api/v1/{operation}` with a JSON body.

| Operation  | Body             | Result         |
|------------|------------------|----------------|
| `add`      | `{"a": 2, "b": 3}` | `a + b`      |
| `subtract` | `{"a": 2, "b": 3}` | `a - b`      |
| `multiply` | `{"a": 2, "b": 3}` | `a * b`      |
| `divide`   | `{"a": 6, "b": 3}` | `a / b`      |
| `power`    | `{"a": 2, "b": 3}` | `a ^ b`      |
| `percent`  | `{"a": 200, "b": 15}` | `b% of a` |
| `sqrt`     | `{"a": 16}`       | `√a`          |

`GET /health` returns `{"status":"ok"}` for health checks.

### Examples

```bash
curl -X POST http://localhost:8080/api/v1/add \
  -H "Content-Type: application/json" -d '{"a": 2, "b": 3}'
# 200 {"result":5}

curl -X POST http://localhost:8080/api/v1/sqrt \
  -H "Content-Type: application/json" -d '{"a": 16}'
# 200 {"result":4}

curl -X POST http://localhost:8080/api/v1/divide \
  -H "Content-Type: application/json" -d '{"a": 1, "b": 0}'
# 422 {"error":{"code":"DIVISION_BY_ZERO","message":"division by zero"}}
```

(With Docker Compose, use `http://localhost:3000/api/v1/...` instead.)

### Errors

Every error has the same shape: `{"error": {"code": "...", "message": "..."}}`.

| Status | When | Codes |
|--------|------|-------|
| 400 | Malformed request: invalid JSON, missing operand, unknown field, NaN/Infinity | `INVALID_INPUT`, `INVALID_OPERAND` |
| 404 | Unknown operation | `INVALID_OPERATION` |
| 405 | Wrong HTTP method | (plain response from Go's router) |
| 422 | Valid request the math cannot answer | `DIVISION_BY_ZERO`, `NEGATIVE_SQUARE_ROOT`, `INVALID_EXPONENT`, `OVERFLOW`, `UNDEFINED_RESULT` |
| 500 | Unexpected error (details are logged, not returned) | `INTERNAL_ERROR` |

## Design decisions

- **One microservice, not one per operation.** The backend is a single stateless service with
  its own `go.mod` and Dockerfile, a health check, configuration through environment variables,
  structured logs (`log/slog`) and graceful shutdown on SIGTERM. Splitting every operation into
  its own service would add network calls and deployment work with no benefit.
- **Standard library only.** Go 1.22+ routing (`POST /api/v1/{op}`) covers everything this API
  needs, so there is no framework to learn or update.
- **Domain separated from HTTP.** `internal/calculator` knows nothing about HTTP and returns
  errors instead of panicking. The HTTP layer maps each domain error to a status code in one
  table, which keeps handlers small and makes both layers easy to test.
- **Validation in both layers.** The frontend validates for a good user experience; the backend
  validates again because it is the source of truth and can be called directly.
- **Same origin, no CORS.** nginx (in Docker) and the Vite proxy (in development) send `/api` to
  the backend, so the browser only talks to one origin.
- **Edge cases.** NaN and ±Infinity inputs are rejected, overflowing results return `OVERFLOW`,
  `0` raised to a negative power is a division by zero, a negative base with a fractional
  exponent is rejected, and `-0` is normalized to `0`.

### Numeric precision

The service uses `float64`, so results follow IEEE 754: `0.1 + 0.2` returns
`0.30000000000000004`. This is a known trade-off kept on purpose for a simple calculator. For
money, I would switch to a decimal type (for example `shopspring/decimal`) or integer cents,
and define a rounding rule.

## What I would do with more time

- Decimal arithmetic and an explicit rounding policy.
- More unit tests for the calculator package edge cases.
- An OpenAPI spec for the API contract.
- Rate limiting, Prometheus metrics and tracing.
- Deployment to Kubernetes using the existing health check.

## AI usage

See [AI_USAGE.md](AI_USAGE.md) for how I used AI tools in this project.
