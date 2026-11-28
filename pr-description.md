Implemented database support package for campbell-checkbook-go.

- Added configuration loader (`internal/config/config.go`) using godotenv.
- Created DB configuration struct and loader (`internal/db/config.go`).
- Implemented PostgreSQL DSN builder (`internal/db/url.go`).
- Added migration runner (`internal/db/migrate.go`).
- Defined placeholder DB interface and constructor (`internal/db/db.go`).
- Provided unit test for DSN construction (`internal/db/db_test.go`).
- Created empty migrations directory.
- Updated .gitignore to ignore task-description.md.
- Added required dependencies to go.mod.
