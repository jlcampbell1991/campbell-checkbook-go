Implemented Pay Day CRUD functionality:

- Added domain model `PayDay` and request struct.
- Created repository interfaces and stub implementations for PayDay and PayDayDeposit.
- Implemented `PayDayRoutes` with handlers for list, create, show, edit, update, deactivate/reactivate.
- Added HTML templates for Pay Day pages (form, index, show, inactive index).
- Integrated Pay Day routes into the main server setup.
- Added unit tests for the Pay Day form handler.
- Updated imports and template embedding.

All new code follows project conventions, passes `go test ./...`, and `go vet` without issues.