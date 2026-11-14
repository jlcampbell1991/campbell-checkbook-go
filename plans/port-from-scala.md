# Agent Context — Campbell Checkbook Go

## Project Overview

Port `campbell-checkbook-v2` (Scala) to `campbell-checkbook` (Go), preserving the functional behavior of the original application.

### Repository Roots

- **Scala reference:** `/Users/jamble-campbell/Development/scala/campbell-checkbook-v2`
- **Go target:** `/Users/jamble-campbell/Development/go/campbell-checkbook`

The Scala repository is the authoritative reference for existing business logic, API behavior, data models, and UI functionality. The Go repository is the implementation target.

## 1. Objectives and Scope

Produce a Go implementation that mirrors the functional behavior of the Scala Checkbook service, including:

- Backend HTTP API.
- Domain entities and business rules.
- PostgreSQL persistence.
- Server-rendered HTML templates.
- Minimal UI necessary to reproduce existing functionality.
- Automated unit and integration tests.

### In Scope

- Accounts.
- Transactions.
- Transfers.
- Pay-days.
- Line-items.
- Reconciliation.
- HTTP routes corresponding to Scala `Routes.scala` and individual route files.
- PostgreSQL persistence.
- Server-side HTML rendering with Go `html/template`.
- Static CSS and JavaScript assets.
- Unit and integration tests.

### Out of Scope

- Experimental Scala LLM integration in `src/main/scala/jamble/campbell/checkbook/llm`.
- UI redesign beyond reproducing existing functionality.
- Containerization; deployment uses a systemd service.

## 2. Architecture and Engineering Conventions

Follow the established conventions of `campbell-tube-go` wherever appropriate.

### Technology Stack

- Go.
- chi HTTP router.
- PostgreSQL.
- pgx.
- golang-migrate.
- godotenv.
- Go `html/template`.
- Standard Go testing framework.

### Configuration

- Load environment variables from `.env` using `github.com/joho/godotenv`.
- Use environment variables for database configuration and application settings.
- Never hardcode credentials.
- Follow existing PostgreSQL DSN conventions.

### HTTP Routing

- Use chi.
- Preserve the Scala application's HTTP endpoint paths and behavior.
- Keep handlers organized by domain.
- Use shared route utilities where appropriate.

### Persistence

- Use PostgreSQL.
- Follow the database conventions established in `campbell-tube-go`.
- Implement migrations using golang-migrate.
- Maintain compatibility with the existing Scala database schema and behavior unless a task explicitly authorizes a change.

### HTML and Static Assets

- Use Go `html/template` for server-side rendering.
- Organize templates under `views/`.
- Serve static assets from `static/css/` and `static/javascript/`.
- Preserve existing UI structure and behavior.
- Avoid unnecessary redesigns.

### Deployment

- Build a standalone Go binary.
- Run the binary through systemd.
- Follow the deployment conventions documented in `campbell-tube-go/README.md`.

## 3. Reference Repository

Use `/Users/jamble-campbell/Development/go/campbell-tube-go` as the reference implementation for established Go conventions.

Relevant reference files:

- `cmd/main.go` — application bootstrap, chi router, environment loading, server startup.
- `go.mod` — dependencies and module conventions.
- `internal/db/migrate.go` — migration runner.
- `internal/db/buildpostgresurl.go` — PostgreSQL DSN construction.
- `utils/shufflearr.go` — utility package conventions.
- `static/javascript/common.js` — shared JavaScript.
- `static/css/common.css` — shared CSS.
- `README.md` — systemd deployment example.
- `Makefile` — build, test, and run targets.

Use these files as examples of coding style, folder structure, tooling, and deployment practices.

Do not copy reference code blindly. Adapt it to the target application's requirements and architecture.

## 4. Scala-to-Go Migration Mapping

Use the following mapping as the initial architectural guide. Verify actual Scala implementations before implementing their Go equivalents.

| Scala Reference | Go Target |
| --- | --- |
| `accounts/models.scala` | `internal/domain/account.go` |
| `accounts/AccountRepository.scala` | `internal/db/account_repo.go` |
| `accounts/AccountsRoutes.scala` | `cmd/handlers/accounts.go` |
| `transactions/Transaction.scala` | `internal/domain/transaction.go` |
| `transactions/TransactionRepository.scala` | `internal/db/transaction_repo.go` |
| `transactions/TransactionRoutes.scala` | `cmd/handlers/transactions.go` |
| `transfers/Transfer.scala` | `internal/domain/transfer.go` |
| `transfers/TransferRoutes.scala` | `cmd/handlers/transfers.go` |
| `pay_days/PayDay.scala` | `internal/domain/payday.go` |
| `pay_days/PayDayRoutes.scala` | `cmd/handlers/payday.go` |
| `line_items/LineItem.scala` | `internal/domain/lineitem.go` |
| `line_items/LineItemsRoutes.scala` | `cmd/handlers/lineitems.go` |
| `reconcile/JobProcessor.scala` | `internal/service/reconcile.go` |
| `reconcile/ReconcileRoutes.scala` | `cmd/handlers/reconcile.go` |
| `utils/RoutesUtils.scala` | `utils/routes.go` |
| `Main.scala` | `cmd/main.go` |
| `application.conf` | `.env` configuration |
| `resources/assets/` | `static/css/` and `static/javascript/` |
| `views/*.scala` | `views/*.go` |

## 5. Behavioral Compatibility

The Scala implementation is the source of truth for migrated functionality.

For every assigned migration task:

1. Inspect the relevant Scala implementation before implementing the Go equivalent.
2. Identify existing behavior, validation, error handling, and data contracts.
3. Preserve existing API paths, HTTP methods, status codes, JSON field names, and response structures.
4. Preserve business rules and persistence semantics.
5. Use explicit JSON struct tags.
6. Add tests covering the behavior being migrated.
7. Do not introduce unrelated changes or redesign existing behavior.

If the Scala behavior is ambiguous, investigate related code and tests before making assumptions.

Document unavoidable differences.

## 6. Testing and Verification

Every implementation task must include appropriate verification.

### Automated Tests

- Write unit tests for domain logic.
- Write integration tests for persistence and HTTP behavior where appropriate.
- Add API contract tests when migrating endpoints.
- Ensure tests pass before declaring a task complete.

### Static Analysis and Formatting

Run applicable checks:

- `gofmt -s -w` on modified Go files.
- `go vet ./...`.
- `golangci-lint run` if configured.

Do not leave lint or test failures unresolved.

### CI

- CI failures block merging.
- Do not claim verification succeeded unless the relevant checks actually ran and passed.
- Report any checks that could not be executed and why.

### UI Validation

For UI-related tasks:

- Verify templates compile and render.
- Verify relevant routes and static assets.
- Ask the user to manually validate rendering at `http://localhost:8080` when visual confirmation is necessary.

## 7. Risks and Mitigations

### Schema Drift

Risk: Scala and Go database schemas diverge.

Mitigation:

- Verify migrations against a fresh PostgreSQL database.
- Compare tables, columns, constraints, and indexes.
- Preserve compatibility with existing data.

### JSON Contract Differences

Risk: Go responses differ from Scala responses.

Mitigation:

- Use explicit JSON tags.
- Add contract tests against representative Scala responses.
- Preserve field names, nullability, and response structure.

### Missing Scala Dependencies

Risk: Experimental Scala features have no Go equivalent.

Mitigation:

- Keep LLM integration out of scope.
- Document deferred features.

### Milestone Creep

Risk: Workers implement unrelated features or exceed the assigned scope.

Mitigation:

- Implement only the assigned Vikunja task.
- Respect task acceptance criteria.
- Split oversized tasks into smaller milestones.
- Avoid unrelated refactoring.
- Do not implement future milestones unless explicitly assigned.

## 8. Worker Execution Requirements

Every worker operating in this project must:

1. Follow the assigned Vikunja task as the implementation scope.
2. Use this project description as persistent engineering context.
3. Inspect the relevant Scala reference implementation.
4. Follow established Go conventions and architecture.
5. Respect the supplied Git branch, parent branch, and parent commit instructions.
6. Implement only the assigned task.
7. Add or update tests.
8. Run appropriate formatting, linting, and tests.
9. Report implementation details and verification results.
10. Identify blockers, assumptions, or unresolved issues explicitly.

The task description and its latest comments provide the authoritative task-specific requirements. They supplement, but do not replace, this project context.

**Priority:** Preserve functional compatibility with Scala, maintain consistency with established Go conventions, and keep implementation changes narrowly scoped to the assigned task.