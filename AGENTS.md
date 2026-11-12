# campbell-checkbook (Go)

## Overview
A Go implementation of the Checkbook service. Provides HTTP API for accounts, transactions, line items, pay days, transfers, and server‑side HTML rendering using `html/template`. Uses Chi router and in‑memory repositories for demo.

## Repository layout
- `main.go` – entry point, sets up Chi router and registers handlers.
- `cmd/handlers/*.go` – HTTP route registrations and handler functions.
- `internal/domain/*.go` – domain structs (Account, Transaction, Transfer, PayDay, LineItem, etc.).
- `internal/db/*.go` – simple in‑memory repository implementations.
- `views/*.html` – HTML templates rendered by handlers.
- `static/` – (optional) static assets.
- `Makefile` – basic `run` target.
- `README.md` – systemd service setup instructions.
- `go.mod` – module definition and dependencies.

## Development commands
```sh
# Run the server (default port 8080)
go run main.go

# Or use the Makefile shortcut
make run

# Run tests (none currently, but `go test ./...` will discover any)
go test ./...

# Format code
gofmt -s -w .

# Lint / vet
go vet ./...
```

## Build & deployment
- Build binary: `go build -o campbell-checkbook`.
- Systemd service defined in `README.md`; adjust `ExecStart` path accordingly.

## Coding conventions
- Follow standard Go formatting (`gofmt`).
- Errors are returned as `error` values; HTTP handlers translate to appropriate status codes.
- Templates are parsed with `template.ParseFiles` using relative paths from project root.
- Use `chi` middleware: Logger, Recoverer, URLFormat.

## Testing strategy
- Unit tests should be placed alongside packages (`*_test.go`) and use table‑driven style.
- Integration tests can start the server with `httptest.NewServer` and exercise routes.

## Known limitations
- Currently uses in‑memory repositories; replace with PostgreSQL implementation when needed.
- No authentication or authorization layers.
- No external static asset handling beyond the placeholder comment.
