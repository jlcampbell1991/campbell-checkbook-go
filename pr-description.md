Implemented the Account "show" page and delete endpoint.

### Changes
- Added route registration in `main.go` and stub repository initialization.
- Created `cmd/handlers/accounts.go` with show and delete handlers.
- Defined domain models `Account` and `Transaction`.
- Introduced `internal/db/account_repository.go` interface and a stub implementation.
- Added utility `IsWithdrawal` and registered it for templates.
- Created `views/accounts/show.tmpl` mirroring the Scala show view.
- Updated `.gitignore` to ignore `pr-description.md` and `commit-message.txt`.

These modifications compile the project and expose the new account routes while providing placeholder logic for data access.
