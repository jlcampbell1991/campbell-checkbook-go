# Plan – Port `campbell-checkbook-v2` (Scala) → `campbell-checkbook` (Go)

**Repository roots**
- Scala source (reference): `/Users/jamble-campbell/Development/scala/campbell-checkbook-v2`
- Go target (this repo): `/Users/jamble-campbell/Development/go/campbell-checkbook`

## 1. Goal & Scope
| Item | Description |
|------|-------------|
| **Goal** | Produce a Go implementation that mirrors the functional behaviour of the Scala Checkbook service (backend API, data model, business rules, and the minimal UI needed for the existing front‑end). |
| **In‑scope** | • All domain entities (accounts, transactions, transfers, pay‑days, line‑items, reconciliation).<br>• HTTP API routes (same endpoints as Scala `*Routes.scala`).<br>• Persistence layer (PostgreSQL via the `internal/db` package used in *campbell‑tube‑go*).<br>• Server‑side rendered HTML templates (using Go `html/template`).<br>• Unit‑ and integration‑tests for each migrated component. |
| **Out‑of‑scope** | • Experimental LLM integration present in Scala (`src/main/scala/jamble/campbell/checkbook/llm`).<br>• UI redesign beyond reproducing existing pages.<br>• Containerisation – we will run the binary as a systemd service. |

## 2. Assumptions (validated by you)
1. **Router** – we will use **chi** (already used in `campbell-tube-go`).
2. **Configuration** – environment variables loaded from a `.env` file (using `github.com/joho/godotenv`).
3. **HTML templates** – same server‑side rendering approach as `campbell-tube-go` (`views/*.go` with `html/template`).
4. **Deployment** – binary run via a systemd unit (see `campbell-tube-go/README.md` for example). 

## 3. High‑level Mapping (Scala → Go)

**Reference files in `campbell-tube-go`** (used for conventions, router setup, DB utilities, static file serving, and systemd deployment):
- `cmd/main.go` – chi router initialization, `.env` loading, server start.
- `go.mod` – module name, required dependencies (chi, godotenv, pgx, golang-migrate, etc.).
- `internal/db/migrate.go` – migration runner used as a model for our own DB migration command.
- `internal/db/buildpostgresurl.go` – helper to construct the PostgreSQL DSN from env vars.
- `utils/shufflearr.go` – example of a utility package with gofmt‑compliant code.
- `static/javascript/common.js` & `static/css/common.css` – static asset directories and chi file‑server configuration.
- `README.md` – example systemd service definition (see lines 1‑15).
- `Makefile` – build, test, and run targets we will replicate/adapt.

These files illustrate the coding style, folder layout, and tooling (GitHub Actions workflow, linting, testing) that the ported project should follow.

**Mapping from Scala files to target Go files**
| Scala file (path) | Target Go file (suggested path) | Mapping notes |
|-------------------|-------------------------------|--------------|
| `src/main/scala/jamble/campbell/checkbook/accounts/models.scala` | `internal/domain/account.go` | Struct definitions, JSON tags. |
| `src/main/scala/jamble/campbell/checkbook/accounts/AccountRepository.scala` | `internal/db/account_repo.go` | CRUD using `database/sql` / `pgx`. |
| `src/main/scala/jamble/campbell/checkbook/accounts/AccountsRoutes.scala` | `cmd/handlers/accounts.go` | chi router, same endpoint paths (`/accounts`). |
| `src/main/scala/jamble/campbell/checkbook/transactions/Transaction.scala` | `internal/domain/transaction.go` | Structs, validation. |
| `src/main/scala/jamble/campbell/checkbook/transactions/TransactionRepository.scala` | `internal/db/transaction_repo.go` | DB layer. |
| `src/main/scala/jamble/campbell/checkbook/transactions/TransactionRoutes.scala` | `cmd/handlers/transactions.go` | API routes. |
| `src/main/scala/jamble/campbell/checkbook/transfers/Transfer.scala` | `internal/domain/transfer.go` | Domain model. |
| `src/main/scala/jamble/campbell/checkbook/transfers/TransferRoutes.scala` | `cmd/handlers/transfers.go` | API routes. |
| `src/main/scala/jamble/campbell/checkbook/pay_days/PayDay.scala` | `internal/domain/payday.go` | Domain model. |
| `src/main/scala/jamble/campbell/checkbook/pay_days/PayDayRoutes.scala` | `cmd/handlers/payday.go` | API routes. |
| `src/main/scala/jamble/campbell/checkbook/line_items/LineItem.scala` | `internal/domain/lineitem.go` | Model. |
| `src/main/scala/jamble/campbell/checkbook/line_items/LineItemsRoutes.scala` | `cmd/handlers/lineitems.go` | API routes. |
| `src/main/scala/jamble/campbell/checkbook/reconcile/JobProcessor.scala` (business logic) | `internal/service/reconcile.go` | Core reconciliation service. |
| `src/main/scala/jamble/campbell/checkbook/reconcile/ReconcileRoutes.scala` | `cmd/handlers/reconcile.go` | API routes. |
| `src/main/scala/jamble/campbell/checkbook/utils/RoutesUtils.scala` | `utils/routes.go` | Helper functions used by routes. |
| `src/main/scala/jamble/campbell/checkbook/Main.scala` | `cmd/main.go` (already exists, will be expanded) | Server bootstrap, chi router, env config. |
| `src/main/resources/application.conf` | `.env` (key‑value) | Translate config keys to env vars (e.g., `POSTGRES_URL`). |
| `src/main/resources/assets/*.css` & `*.js` | `static/css/` & `static/javascript/` (already present) | Serve as static files via chi's file server. |
| `src/main/scala/jamble/campbell/checkbook/views/*.scala` (template logic) | `views/*.go` (html/template) | Convert Scala template fragments to Go templates, preserving structure. |


## 5. Verification Strategy
- **Automated tests** for each milestone (unit + integration) must pass on the CI runner.
- **Static analysis**: `go vet`, `golangci-lint` (if present), and `gofmt -s -w` after every change.
- **Manual validation** for UI milestones – the user will be asked to open `http://localhost:8080` and confirm rendering.
- **CI gate**: any lint or test failure blocks merging.

## 6. Risks & Mitigations
| Risk | Impact | Mitigation |
|------|--------|------------|
| Schema drift between Scala migrations and Go migrations | Runtime DB errors | Run migration verification (M02) against a fresh DB; diff tables.
| JSON field naming differences | Client breakage | Use explicit `json:"field_name"` tags; add contract tests comparing sample Scala responses.
| Missing third‑party Scala libs (e.g., LLM) | Feature gaps | Defer those features to a later optional milestone; document as out‑of‑scope.
| Milestone creep (trying to do too much) | Hard to verify, longer feedback | Enforce “one milestone per session”; split oversized work.
