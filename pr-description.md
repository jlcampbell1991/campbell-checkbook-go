Implemented full CRUD UI and backend for transactions:
- Added `transactionService` with repository and template handling.
- Created `TransactionRoutes` mounting CRUD endpoints (`/`, `/new`, POST `/`).
- Implemented handlers for listing, rendering add form, and creating transactions with validation and redirects.
- Extended `TransactionRepository` interface with Create, Show, Update, Delete methods and provided stub implementations.
- Added `NewStubTransactionRepo` factory.
- Updated `main.go` to mount the new transaction router using the stub repo.
- Added HTML templates: `add.tmpl`, `edit.tmpl`, `csv_form.tmpl`.
- Added unit tests for transaction routes and stub repository.
- Updated `.gitignore` (already contained `task-description.md`).