Implemented full Pay Day CRUD UI and backend:
- Added domain model `PayDay` and request struct.
- Created repository interfaces and stub implementations for PayDays and PayDayDeposits.
- Developed `PayDayRoutes` with handlers for list, inactive list, show, new, create, edit, update, deactivate, reactivate, and get‑paid.
- Added HTML templates for Pay Day pages (index, inactive, form, show).
- Updated main application to parse new templates and mount the Pay Day router.
- Provided unit tests for Pay Day new‑form handler and adjusted existing tests to work with new templates.
All code passes `go test ./...`, `go vet`, and `gofmt` checks.
