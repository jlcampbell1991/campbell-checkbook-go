# Plan: Port **campbell‑checkbook‑v2** (Scala) → **campbell‑checkbook‑go** (Go)

**Slug:** `port-scala-to-go`
**Root directory (future repo):** `/Users/jamble-campbell/Development/go/campbell-checkbook`

## Goal
Create a fully‑functional Go web application that mirrors the existing Scala service, using:

* **Router:** `github.com/go-chi/chi/v5` (as in the reference project `campbell‑tube‑go`)
* **Database:** plain `database/sql` with the PostgreSQL driver (`github.com/jackc/pgx/v5/stdlib`) – no schema changes, same tables as the Scala version.
* **LLM client:** identical Hermes HTTP API contract.
* **Testing:** Go `testing` + `net/http/httptest` for each milestone.
* **Conventions:** Global Go conventions (`go‑conventions.md`) – run `gofmt`, `go vet`, `go test` after every change.

## Overall Milestone Strategy
Milestones are sized to be comfortably implementable within a single 90 k‑token session. Each milestone adds a self‑contained domain (API routes, data models, repository, and tests) and can be compiled & tested independently.

| ID | Title | Objective | Prereqs | Approx. token budget | Scala source files | Go reference files (campbell‑tube‑go) |
|----|-------|-----------|----------|----------------------|--------------------|-----------------------------------|
| **M01** | **Project scaffolding & core server** | Set up Go module, main entry point, chi router, shared DB init, config loading. | None | ≤ 15 k | `src/main/scala/jamble/campbell/checkbook/Main.scala`, `src/main/scala/jamble/campbell/checkbook/CheckbookServer.scala` | `cmd/main.go`, `internal/db/buildpostgresurl.go`, `internal/db/migrate.go` |
| **M02** | **Accounts domain** | Define `Account` model, repository (CRUD), routes (`/accounts`), request/response structs, and unit tests. | M01 | ≤ 25 k | `src/main/scala/jamble/campbell/checkbook/accounts/AccountRepository.scala`, `src/main/scala/jamble/campbell/checkbook/accounts/AccountsRoutes.scala`, `src/main/scala/jamble/campbell/checkbook/accounts/models.scala` | `internal/db/buildpostgresurl.go` |
| **M03** | **Line‑items domain** | Model `LineItem`, repository, routes (`/line‑items`), validation, tests. | M02 | ≤ 25 k | `src/main/scala/jamble/campbell/checkbook/line_items/LineItem.scala`, `src/main/scala/jamble/campbell/checkbook/line_items/LineItemRepository.scala`, `src/main/scala/jamble/campbell/checkbook/line_items/LineItemsRoutes.scala` | `utils/shufflearr.go` |
| **M04** | **Transactions domain** | Model `Transaction`, CSV import endpoint, repository, routes, tests. | M03 | ≤ 30 k | `src/main/scala/jamble/campbell/checkbook/transactions/Transaction.scala`, `src/main/scala/jamble/campbell/checkbook/transactions/TransactionRepository.scala`, `src/main/scala/jamble/campbell/checkbook/transactions/TransactionRoutes.scala` | `internal/db/migrate.go` |
| **M05** | **Pay‑days domain** | Models `PayDay` & `PayDayDeposit`, repositories, routes, tests. | M04 | ≤ 25 k | `src/main/scala/jamble/campbell/checkbook/pay_days/PayDay.scala`, `src/main/scala/jamble/campbell/checkbook/pay_days/PayDayDepositRepository.scala`, `src/main/scala/jamble/campbell/checkbook/pay_days/PayDayRepository.scala`, `src/main/scala/jamble/campbell/checkbook/pay_days/PayDayRoutes.scala` | `utils/shufflearr.go` |
| **M06** | **Reconcile import‑job domain** | Models `ImportJob`, repository, background worker stub, routes, tests. | M05 | ≤ 30 k | `src/main/scala/jamble/campbell/checkbook/reconcile/ImportJobRepository.scala`, `src/main/scala/jamble/campbell/checkbook/reconcile/ImportJobRoutes.scala`, `src/main/scala/jamble/campbell/checkbook/reconcile/JobProcessor.scala`, `src/main/scala/jamble/campbell/checkbook/reconcile/JobWorker.scala`, `src/main/scala/jamble/campbell/checkbook/reconcile/MatchingService.scala`, `src/main/scala/jamble/campbell/checkbook/reconcile/Models.scala`, `src/main/scala/jamble/campbell/checkbook/reconcile/NormalisationService.scala`, `src/main/scala/jamble/campbell/checkbook/reconcile/ReconcileRoutes.scala` | `internal/db/migrate.go` |
| **M07** | **LLM (Hermes) client** | Go client mirroring `HermesClient.scala`, request/response structs, integration tests (mock server). | M06 | ≤ 25 k | `src/main/scala/jamble/campbell/checkbook/llm/HermesClient.scala`, `src/main/scala/jamble/campbell/checkbook/llm/Prompts.scala` | `internal/db/buildpostgresurl.go` |
| **M08** | **Utility package** | `Currency` handling, DB helper (`DBUtils`), logger, error types, tests. | M01 | ≤ 15 k | `src/main/scala/jamble/campbell/checkbook/utils/Currency.scala`, `src/main/scala/jamble/campbell/checkbook/utils/DBUtils.scala`, `src/main/scala/jamble/campbell/checkbook/utils/Debugger.scala`, `src/main/scala/jamble/campbell/checkbook/utils/RoutesUtils.scala` | `utils/shufflearr.go` |
| **M09** | **Session management & auth** | Session model, JWT middleware (if present in Scala), routes, tests. | M02 | ≤ 20 k | `src/main/scala/jamble/campbell/checkbook/sessions/Session.scala`, `src/main/scala/jamble/campbell/checkbook/sessions/SessionRoutes.scala` | `cmd/main.go` |
| **M10** | **Full‑stack sanity check** | Wire all domains together, run end‑to‑end smoke test (create account → add line‑item → transaction → reconcile). | M01‑M09 | ≤ 30 k | *(All of the above)* | `cmd/main.go`, `internal/db/buildpostgresurl.go`, `internal/db/migrate.go` |

> **Milestone size note:** The token estimates include source code, comments, and tests. The largest milestones (M04, M06, M10) stay well under the 90 k limit.

## Common Implementation Pattern (taken from `campbell‑tube‑go`)

* `cmd/` – entry‑point (`main.go`) starts HTTP server.
* `internal/` – contains sub‑packages per domain (`accounts`, `lineitems`, `transactions`, …).
* `internal/db/` – DB connection pool, `InitDB` helper.
* `internal/config/` – environment‑based configuration (port, DB URL).
* `internal/middleware/` – request logging, panic recovery, optional auth.
* `internal/llm/` – Hermes client implementation.

Each domain package will expose:

```go
type Service struct {
    repo Repository
}
func NewService(repo Repository) *Service { … }

func (s *Service) RegisterRoutes(r chi.Router) { … }
```

## Verification Commands per Milestone

| Milestone | Commands to run after implementation |
|-----------|---------------------------------------|
| M01 | `go test ./...` (should pass – only config & router init) |
| M02 | `go test ./internal/accounts -run TestCreateAccount` |
| M03 | `go test ./internal/lineitems -run TestCreateLineItem` |
| M04 | `go test ./internal/transactions -run TestImportCSV` |
| M05 | `go test ./internal/paydays -run TestCreatePayDay` |
| M06 | `go test ./internal/reconcile -run TestCreateImportJob` |
| M07 | `go test ./internal/llm -run TestHermesClient` |
| M08 | `go test ./internal/utils` |
| M09 | `go test ./internal/sessions -run TestLogin` |
| M10 | `go test ./...` (full suite) + manual curl smoke test (documented in plan) |

## Acceptance Criteria per Milestone (example for M02)

* All Go files compile (`go build ./...`).
* `go vet` and `gofmt -l` report no issues.
* Unit tests cover ≥ 80 % of the service’s public methods and all pass.
* HTTP routes behave as per Scala counterpart (status codes, JSON schema).
* No database schema changes – queries succeed against the existing PostgreSQL DB.

(Analogous criteria apply to each subsequent milestone.)

## Status Overview (initial)

| Milestone | Status |
|-----------|--------|
| M01 | **COMPLETE** |
| M02 | **COMPLETE** |
| M03 | **COMPLETE** |
| M04 | **COMPLETE** |
| M05 | **COMPLETE** |
| M06 | **COMPLETE** |
| M07 | **COMPLETE** |
| M08 | **COMPLETE** |
| M09 | **COMPLETE** |
| M10 | **COMPLETE** |

---

**Next step:** Once you give the go‑ahead, we will start **Milestone M01** (project scaffolding & core server).