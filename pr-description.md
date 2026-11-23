Implemented database support for the Checkbook Go project.

- Added `.gitignore` entry for `task-description.md`.
- Created `internal/config/config.go` to load environment variables from a `.env` file.
- Implemented `internal/db/config.go` with a `Config` struct and `Load` function that reads DB connection settings.
- Added `internal/db/url.go` providing `BuildPostgresURL` to construct a PostgreSQL DSN.
- Added `internal/db/migrate.go` to run migrations using golang‑migrate.
- Added a minimal `DB` interface and concrete implementation in `internal/db/db.go`.
- Provided a unit test `TestBuildPostgresURL` in `internal/db/db_test.go`.
- Created empty `internal/db/migrations` directory for future migration files.
- Added required third‑party dependencies (`godotenv`, `golang-migrate`, `pgx`).

All code compiles, passes `go test ./...`, and conforms to formatting and vet checks.