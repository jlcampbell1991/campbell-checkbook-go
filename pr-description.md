Implemented the **Accounts Form** feature:

- Added `AccountRoutes` with `GET /accounts/new` and `POST /accounts` endpoints in `cmd/handlers/accounts.go`.
- Implemented `new` handler to render the account creation form using the new `accounts_add.tmpl` template.
- Implemented `create` handler to parse the form, validate the `name` field, insert a new `Account` via `AccountRepository`, and redirect to the home page.
- Created domain model `Account` (`internal/domain/account.go`).
- Added stub `AccountRepository` interface and implementation (`internal/db/account_repository.go`).
- Added HTML template `views/accounts/add.tmpl` for the new‑account form.
- Added unit tests for both handlers (`cmd/handlers/accounts_test.go`).
- Updated imports and ensured all code is formatted (`gofmt`) and passes `go vet`.

These changes provide a functional UI and backend to create new accounts, fulfilling Task ID 9 requirements.