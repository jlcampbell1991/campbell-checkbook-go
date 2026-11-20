Implemented line‑item creation UI and backend.

- Added template parsing for `views/line_items/*.tmpl`.
- Created `line_items_form.tmpl` providing a form for new line items.
- Introduced GET `/accounts/{id}/line-items/new` and POST `/accounts/{id}/line-items` routes in `main.go` with handlers that render the form, validate input, insert via repository, and redirect to the account page.
- Extended `LineItemRepository` with an `Insert` method and implemented it in the stub repository, returning a placeholder ID.
- Updated imports to include `fmt`, `strconv`, and the domain package.
- Added unit test for the stub `Insert` implementation.

All code passes `go vet` and `go test ./...`.
