Implemented the Transactions page:

- Added full `Transaction` domain model with comprehensive fields and `time` import.
- Created stub transaction repository (`internal/db/transaction_repo.go`) exposing `Index` and `GetAll` functions.
- Developed HTTP handler `TransactionsIndex` to render paginated transaction listings, parse `page` query param, and build a next‑page URL.
- Embedded transaction templates via `//go:embed` and added `views/transactions/layout.tmpl` and `views/transactions/index.tmpl` reflecting the Scala UI.
- Registered the `/transactions` route in `main.go`.
- Updated imports and ensured all new code complies with Go formatting and vet checks.
