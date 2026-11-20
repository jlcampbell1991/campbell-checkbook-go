Implemented the Transactions page:

- Added route registration for `/transactions` in `main.go`.
- Created domain model with full transaction fields.
- Added stub transaction repository (`internal/db/transaction_repo.go`).
- Implemented HTTP handler `TransactionsIndex` with pagination, template rendering, and stub data retrieval.
- Provided `SetTemplate` to share parsed templates with the handler.
- Added transaction templates (`views/transactions/layout.tmpl` and `views/transactions/index.tmpl`).
- Updated embed directive and template parsing to include new templates.
- Fixed utility function `IsWithdrawal` and cleaned unused imports.
- Ran go vet, go test, and gofmt to ensure code quality.
