Implemented the accounts creation UI and backend:

- Added routes `GET /accounts/new` and `POST /accounts` in `cmd/handlers/accounts.go`.
- Implemented `new` handler to render the account creation form.
- Implemented `create` handler to process form submissions, validate the `name` field, insert a new account via `AccountRepository`, and redirect to home.
- Updated route registration order to include the new endpoints.
- Added `views/accounts/add.tmpl` template for the new account form.
- Added unit tests `TestAccountNewHandler` and `TestAccountCreateHandler` covering rendering and submission handling.
- Ensured code passes `go test`, `go vet`, and is formatted with `gofmt`.

These changes fulfill the task requirement to provide a UI and backend for adding new accounts.