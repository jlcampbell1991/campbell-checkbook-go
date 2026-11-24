Implemented the Line‑Item "show" feature:

- Added domain models for LineItem, PayDayDeposit, and Transaction.
- Defined `LineItemRepository` interface with a stub implementation.
- Created HTTP handler for `GET /line-items/{id}` including pagination support and template rendering.
- Registered template helper functions (`IsWithdrawal`, `CurrencyFromDouble`).
- Integrated the new routes into the main server and instantiated the stub repository.
- Added Go HTML templates for the line‑item show page and supporting pay‑day deposit and transaction layouts.
- Provided utility functions and registered them with the template engine.
- Added unit tests for the handler and utility functions.
- Updated imports and build configuration accordingly.
