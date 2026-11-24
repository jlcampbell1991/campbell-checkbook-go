Implemented account creation UI and backend.

### Changes
- Added `internal/domain/account.go` defining the `Account` struct.
- Added `internal/db/account_repository.go` with `AccountRepository` interface and a stub implementation.
- Created `cmd/handlers/accounts.go` with routes for rendering the new‑account form (`GET /accounts/new`) and handling form submissions (`POST /accounts`).
- Added HTML template `views/accounts/add.tmpl` for the account creation form.
- Updated `main.go` to initialize the stub account repository and mount the new account routes.
- Added unit tests `cmd/handlers/accounts_test.go` covering rendering, successful creation, and validation errors.

These updates provide a functional UI for adding accounts and integrate it into the existing server architecture.
