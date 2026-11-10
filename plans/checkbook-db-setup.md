# Checkbook DB Setup Plan

**Repository:** `github.com/jlcampbell1991/campbell-checkbook-go`

## Summary

Add a robust PostgreSQL database layer to the Go version of the Checkbook application, mirroring the architecture used in `campbell-tube-go`. This includes:

- A migration runner using `golang-migrate`.
- A helper to build a PostgreSQL connection URL.
- Database schema migrations for the `checkbook_accounts` table.
- An `Account` model and repository interface with basic CRUD operations.
- Unit tests for the repository using `go-sqlmock`.
- Integration of migration execution into the existing `InitDB` flow.

All work will be broken into small, independently testable milestones.

## Goals

- Provide a clean, reusable DB layer that other packages can depend on.
- Ensure the schema is created/updated automatically on startup.
- Maintain high test coverage for DB interactions.
- Follow the project’s Go conventions (gofmt, go vet, go test, lint).

## Non‑Goals

- Expose a full REST/GraphQL API (out of scope for this plan).
- Implement UI or CLI commands beyond what already exists.
- Add production‑grade migrations beyond the initial table creation.

## Requirements & Acceptance Criteria

| Requirement | Acceptance Criteria |
|-------------|----------------------|
| Migration runner | `RunMigrations(dsn string)` exists, logs the source URL, and applies all migrations without error (or reports `no change`). |
| PostgreSQL URL builder | `BuildPostgresURL(cfg PostgresConfig) string` correctly URL‑escapes user/password and supports both TCP host and Unix socket formats. |
| Schema migrations | `000001__create-accounts-table.up.sql` creates `checkbook_accounts` with columns matching the `Account` struct. |
| Account model | `type Account struct { ID uuid.UUID; Name string; CreatedAt time.Time; UpdatedAt time.Time }` with appropriate `db` tags. |
| Repository interface | `AccountRepository` defines `Create`, `GetByID`, `List`, `Update`, `Delete` methods. |
| Repository implementation | `accountRepository` implements the interface using `database/sql` and parameterised queries. |
| Unit tests | All repository methods have tests using `go-sqlmock`; tests compile and pass (`go test ./...`). |
| Integration with InitDB | `InitDB` calls `RunMigrations` after establishing the connection. |
| CI verification | `go test ./...`, `go vet ./...`, and `gofmt -l` report no issues. |

## Technical Approach & Design Decisions

- Reuse the migration runner from `campbell-tube-go` (golang‑migrate with file source).
- Place migrations under `internal/db/migrations/`.
- Keep the URL builder in `internal/db/buildpostgresurl.go`.
- Define the repository in a new file `internal/db/account.go`.
- Add tests under `internal/db/account_test.go`.
- Use the existing `go-sqlmock` dependency already present in the repo.
- Follow the global Go conventions defined in `~/.config/opencode/go-conventions.md`.

## Milestones

| ID | Title | Objective | Prerequisites | Acceptance Criteria | Status |
|----|-------|-----------|---------------|---------------------|--------|
| **M01** | Add migration runner | Implement `RunMigrations(dsn string)` and a basic migration file for the accounts table. | None | Function exists, logs source URL, runs migrations without error; migration file present. | USER_VALIDATION_PENDING |
| **M02** | Add PostgreSQL URL builder | Implement `BuildPostgresURL` and unit test. | M01 | Helper correctly builds URL for TCP and socket hosts; test passes. | NOT_STARTED |
| **M03** | Create `checkbook_accounts` migration | Add `000001__create-accounts-table.up.sql` with appropriate columns. | M01 | Migration creates the table matching the `Account` struct; `RunMigrations` applies it. | NOT_STARTED |
| **M04** | Define `Account` model | Add `Account` struct with `db` tags. | M03 | Struct compiled; fields match migration columns. | NOT_STARTED |
| **M05** | Define repository interface | Add `AccountRepository` with CRUD signatures. | M04 | Interface compiles; no implementation yet. | NOT_STARTED |
| **M06** | Implement repository | Add `accountRepository` with SQL queries for all CRUD methods. | M05 | Methods execute queries using `*sql.DB`; compile without errors. | NOT_STARTED |
| **M07** | Add repository unit tests | Write tests for each CRUD method using `go-sqlmock`. | M06 | Tests cover all methods; `go test ./...` passes. | NOT_STARTED |
| **M08** | Integrate migrations into `InitDB` | Update `internal/db/db.go` to call `RunMigrations` after connection. | M01, M03 | Application starts, logs migration execution, and DB schema is ensured. | NOT_STARTED |
| **M09** | CI verification & cleanup | Run `go fmt`, `go vet`, `golangci-lint` (if configured); ensure no failures. | All previous milestones | All commands pass; repository is ready for production use. | NOT_STARTED |

### Modified Files
- internal/db/migrate.go
- internal/db/migrations/000001__init.up.sql

**Automated verification:** `go fmt ./...`, `go vet ./...`, and `go test ./...` all passed.

## Manual Validation Steps

1. Run the application (`go run ./cmd/main.go`) with a local PostgreSQL instance and verify that the `checkbook_accounts` table is created.
2. Use a simple script or `psql` to insert, retrieve, update, and delete an account via the repository methods.
3. Confirm that the unit tests exercise each path and pass.

## Open Questions & Risks

- The project currently has no CI configuration; adding `golangci-lint` may require additional setup.
- No existing `Account` model; introducing it may affect other packages – we will keep the scope limited to this repository.
- Migration files must be placed relative to the repository root; ensure the path used by `RunMigrations` matches the actual location.

---

*This plan is saved in `plans/checkbook-db-setup.md`. The next step is to implement Milestone M01.*