Implemented full CRUD UI and backend for transactions.

- Added TransactionRoutes with GET `/`, GET `/new`, POST `/` handlers.
- Created `transactionService` struct to encapsulate repository and templates.
- Implemented index, new, and create handlers (form rendering, validation, repository interaction, redirects).
- Extended `TransactionRepository` interface with Create, Show, Update, Delete methods and provided stub implementations in `transaction_repo.go`.
- Updated `main.go` to mount the new transaction router instead of the legacy handler.
- Added HTML templates: `add.tmpl`, `edit.tmpl`, `csv_form.tmpl` for transaction forms.
- Added unit tests for the new handlers and stub repository.
- Ensured `.gitignore` already excludes `task-description.md`.
