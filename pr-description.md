Implemented line‑item creation form and backend.

- Added template parsing for `views/line_items/*.tmpl`.
- Created new GET `/accounts/{id}/line-items/new` and POST `/accounts/{id}/line-items` routes in `main.go` with handlers rendering a form and inserting a line item via repository.
- Extended `LineItemRepository` with `Insert` method and provided stub implementation returning a placeholder ID.
- Added `views/line_items/form.tmpl` for the new form.
- Updated imports to include required packages (`fmt`, `strconv`, `domain`).
- Added unit test for stub `Insert` implementation.
- Updated handler tests with mock repository to satisfy new interface.
- Ensured all code passes `go test ./...` and follows formatting conventions.
