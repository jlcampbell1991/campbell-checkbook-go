Implemented Pay Day Deposits UI and backend.

- Added embed of new view templates and parsing in `main.go`.
- Created `cmd/handlers/pay_day_deposits.go` providing routes for creating a deposit (`GET /new` and `POST /`).
- Implemented `payDayDepositService` with dependencies on `PayDayRepository`, `PayDayDepositRepository`, and `LineItemRepository`.
- Added HTML templates (`form.tmpl`, `add.tmpl`, `edit.tmpl`) under `views/pay_day_deposits/` for rendering the deposit form.
- Updated routing in `main.go` to mount the new deposit routes.
- Added unit tests `pay_day_deposits_test.go` covering form rendering and deposit creation handling.
- Ensured all new code follows project conventions and compiles with existing build pipeline.
