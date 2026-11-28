# Pull Request Description

Implemented the **Line‑Item "show"** page and its HTTP endpoint.

### Key Changes
- Added template helper functions (`IsWithdrawal`, `CurrencyFromDouble`) and registered them.
- Created domain models: `LineItem`, `PayDayDeposit`, `Transaction`.
- Defined `LineItemRepository` interface with required methods.
- Implemented handler `cmd/handlers/line_items.go` exposing `GET /line-items/{id}` with pagination support.
- Registered line‑item routes in `main.go`.
- Added Go HTML templates for the line‑item show view and supporting pay‑day‑deposit and transaction layouts.
- Updated imports and routing setup.

These changes provide the server‑side rendering of a line item, its associated pay‑day deposits, and paginated transactions, matching the behavior of the reference Scala implementation.
