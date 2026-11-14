# AGENTS.md – Campbell Checkbook Go Service

## Project Overview
The **campbell-checkbook** repository contains a Go implementation of a personal‑finance backend and minimal UI. It serves HTML pages rendered with `html/template`, provides static assets, and (in future work) will expose a JSON HTTP API for accounts, transactions, transfers, pay‑days, line‑items, and reconciliation. The service is intended to run as a systemd unit in production.

## Architecture
- **Entry point**: `main.go` – loads environment variables (`.env` via `github.com/joho/godotenv`), configures a **chi** router, registers middleware (`Logger`, `Recoverer`), serves static files, parses templates, and starts the HTTP server.
- **Routing**: chi router (`github.com/go-chi/chi/v5`). Future route groups will be added under `cmd/handlers/`.
- **Implemented routes**:
  - Accounts show page (`GET /accounts/{id}`)
  - Account delete (`POST /accounts/{id}/delete`)
  - **Account creation form** (`GET /accounts/new`)
  - **Account creation submission** (`POST /accounts`)
  - **Line items creation form** (`GET /accounts/{id}/line-items/new`)
  - **Line items creation submission** (`POST /accounts/{id}/line-items`)
  - **Transactions list** (`GET /transactions`)
  - **Transaction creation form** (`GET /transactions/new`)
  - **Transaction creation submission** (`POST /transactions`)
   - **Transfer creation form** (`GET /transfers/new`)
   - **Transfer creation submission** (`POST /transfers`)
- **Templates**: Go `html/template` files located in `views/layouts/`. The base layout `application.html` includes the `menu.html` partial and a `{{block "content" .}}` for page‑specific content. Account creation uses `views/accounts/add.html` (template name `add`).
- **Static assets**: Served from `static/stylesheets/` and `static/javascripts/` via `http.FileServer` with URL prefixes `/stylesheets/` and `/javascripts/`.
- **Configuration**: Environment variables (`PORT`, DB credentials, etc.) loaded from a `.env` file. Values are accessed via `os.Getenv`.
- **Production deployment**: Systemd service defined in `README.md` (`campbell-checkbook-go.service`) that runs `start.sh` which invokes the compiled binary.
- **Future layers**:
  - **Domain models** under `internal/domain/` (e.g., `account.go`, `transaction.go`).
  - **Database access** under `internal/db/` using `pgx` or `database/sql`. This package contains migration helpers (`RunMigrations`) and URL builder utilities.
  - **Business services** under `internal/service/`.

## Development Commands
| Command | Description |
|---------|-------------|
| `make run` | Executes `go run main.go` – starts the server on the port defined by `PORT` (default `8080`). |
| `go build ./...` | Compiles the binary (output `campbell-checkbook-go`). |
| `go test ./...` | Runs all unit and integration tests. |
| `go vet ./...` | Static analysis for potential issues. |
| `gofmt -s -w .` | Formats the entire codebase (required before commit). |
| `golangci-lint run` | Runs lint checks if the tool is installed (optional). |
| `make db-migrate` | Runs database migrations using the `internal/db` package. Equivalent to `go run ./internal/db/migrate.go` with the `DATABASE_URL` (or individual DB env vars) supplied. |
| `make db-status` | Prints the current migration version using the same migration tooling. |
| `sudo systemctl restart campbell-checkbook-go` | Restarts the production service (see `README.md`). |

## Coding Conventions
- **Formatting**: All `.go` files must pass `gofmt -s -w`. The repository enforces this as part of the mandatory coding workflow.
- **Linting/Static analysis**: Run `go vet` and optionally `golangci-lint` on every change.
- **Error handling**: Fatal startup errors use `log.Fatalf`. Request‑level errors return via `http.Error` with appropriate status codes.
- **Imports**: Group standard library imports, third‑party imports, and internal imports separately, each block separated by a blank line.
- **Templates**: Use `{{define "name"}} … {{end}}` and `{{template "name" .}}`. Keep layout files in `views/layouts/` and page‑specific templates in `views/`.
- **Static files**: Reference CSS/JS via the `/stylesheets/` and `/javascripts/` URL prefixes.
- **Environment**: Load `.env` with `godotenv.Load()` early in `main`. Do not panic on missing optional vars; provide defaults where sensible.
- **Version control**: Never commit merge‑conflict markers (`<<<<<<<`, `>>>>>>>`). Resolve conflicts before committing.

## Testing Strategy
- **Unit tests**: Place tests in `*_test.go` alongside the code they exercise. Use table‑driven tests where appropriate.
- **HTTP tests**: Use `net/http/httptest` to spin up a test server with the chi router and verify routes, template rendering, and static file serving.
- **Integration tests**: Future DB integration tests will run against a temporary PostgreSQL instance (e.g., using Docker or a test database URL).
- **CI**: All tests must pass on the CI runner before merging. The CI pipeline runs `go test ./...`, `go vet`, and `gofmt -l`.

## Project‑Specific Constraints
- **Go conventions**: Follow the global conventions defined in `/Users/jamble-campbell/.config/opencode/go-conventions.md` (formatting, vetting, testing).
- **Systemd**: Production deployment relies on the service described in `README.md`. Ensure the binary is built and placed where `start.sh` expects it.
- **Static asset paths**: The router expects `/stylesheets/*` and `/javascripts/*`. Do not rename these directories without updating the router configuration.
- **Database constraints**:
  - Migrations must be run before the application starts in any environment that accesses the DB.
  - The `DATABASE_URL` (or the component parts `DB_USER`, `DB_PASSWORD`, `DB_HOST`, `DB_NAME`, `DB_PORT`) must be set in the `.env` file or the environment.
  - The `internal/db` package assumes PostgreSQL; other databases are not supported without code changes.
- **Future work**: The current repository is a scaffold. Adding domain models, DB repositories, and API routes should follow the mapping outlined in `plans/port-from-scala.md`.

---
*Generated by OpenCode agent.*