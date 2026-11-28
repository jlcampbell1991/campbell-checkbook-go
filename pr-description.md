Implemented the transactions page:

- Added `Transaction` domain model with full fields.
- Created transaction repository stub (`internal/db/transaction_repo.go`).
- Added HTTP handler `TransactionsIndex` (`cmd/handlers/transactions.go`) with pagination, next link, and template rendering.
- Updated `main.go` to embed transaction templates, parse them, and register the `/transactions` route.
- Added transaction view templates (`views/transactions/layout.tmpl` and `views/transactions/index.tmpl`).
- Updated `.gitignore` (already contains `task-description.md`).
- Added `task-description.md` describing the task.
