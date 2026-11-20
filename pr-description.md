## PR Description

Implemented database support package for **campbell-checkbook-go**:

- Added `.env` loader (`internal/config/config.go`).
- Created DB configuration struct and loader (`internal/db/config.go`).
- Implemented PostgreSQL DSN builder (`internal/db/url.go`).
- Added migration runner stub (`internal/db/migrate.go`).
- Defined generic `DB` interface and concrete wrapper (`internal/db/db.go`).
- Provided unit test for DSN construction (`internal/db/db_test.go`).
- Created migrations directory placeholder.
- Updated module dependencies to include `godotenv`, `golang-migrate`, and `pgx`.
- Ensured `.gitignore` already contains `task-description.md`.

These changes lay the groundwork for database interactions and future migrations.
