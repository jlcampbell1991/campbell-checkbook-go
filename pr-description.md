# Pull Request Description

Implemented the **Account "show"** page and the associated **delete** endpoint for the `campbell-checkbook-go` project.

### Key Additions
- Added route registration for accounts in `main.go` (mounting the new accounts router).
- Created `cmd/handlers/accounts.go` with:
  - `show` handler that loads an account, its paginated transactions, dislocated transactions, computes the dislocated total, builds a next‑page URL, and renders the `accounts_show.tmpl` template.
  - `delete` handler that destroys an account and redirects to the home page.
- Defined domain models:
  - `internal/domain/account.go` – `Account` struct with JSON tags (`id`, `name`, `balance`, `budget`).
  - `internal/domain/transaction.go` – `Transaction` struct with fields for ID, account ID, amount, category, date, etc.
- Added repository interface in `internal/db/account_repository.go` outlining CRUD operations and transaction queries required by the handlers.
- Created view templates:
  - `views/accounts/show.tmpl` – Go `html/template` version of the Scala show view, rendering account details, line items, transactions, dislocated transactions, and navigation links.
  - Helper templates for line items and transactions (`account_items.tmpl` and `transactions_layout.tmpl`).
- Added utility function `isWithdrawal` to determine CSS classes based on amount sign and registered it with the template set.
- Updated `.gitignore` to exclude the generated `task-description.md`.

### Testing
- Added unit tests for the new handlers (`cmd/handlers/accounts_test.go`) using a mock repository to verify rendering, pagination, and delete behavior.
- Ensured all new Go files pass `go vet` and are formatted with `gofmt`.

These changes bring the Go implementation in line with the reference Scala behavior for displaying an account and deleting it.
