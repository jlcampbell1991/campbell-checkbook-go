Implemented Line‑Item "show" page and endpoint.

- Added domain structs for `LineItem` and `PayDayDeposit`.
- Defined `LineItemRepository` interface with stub implementation.
- Created handler `cmd/handlers/line_items.go` providing `GET /line-items/{id}` with pagination support.
- Added Go HTML template `views/line_items/show.tmpl` mirroring the Scala view.
- Updated `main.go` to instantiate the stub repository and mount the new routes.
- Adjusted `utils.IsWithdrawal` to accept a float64 for template compatibility.
- Updated imports accordingly.

All new code compiles and the server now serves the line‑item show page.
