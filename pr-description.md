# Pay Day Deposits Form Implementation

Implemented the UI and backend for creating Pay Day Deposits.

## Changes
- **Embedded templates**: Updated `main.go` embed directive and parsing to include `views/pay_day_deposits/*.tmpl`.
- **Routes**: Added mount for `/pay-day-deposits` in `main.go` and created `PayDayDepositRoutes` with `GET /new` and `POST /` endpoints.
- **Handler**: New `cmd/handlers/pay_day_deposits.go` handling form rendering and deposit creation using stub repositories.
- **Repository**: Utilized existing `PayDayRepository` and `PayDayDepositRepository` stubs.
- **Templates**: Added `form.tmpl`, `add.tmpl`, and `edit.tmpl` under `views/pay_day_deposits` for the deposit form.
- **Tests**: Added unit tests `pay_day_deposits_test.go` covering form rendering and submission with redirect.

All new code follows project conventions, passes `go fmt`, `go vet`, and existing unit tests.
