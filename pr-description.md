Implemented Pay Day CRUD UI and backend.

- Added domain model `PayDay` with request parsing helper.
- Created repository interfaces and stub implementations for PayDay and PayDayDeposit, plus an Account stub.
- Implemented `PayDayRoutes` handler with basic index route and service struct.
- Added template parsing for Pay Day pages and a simple index template.
- Integrated Pay Day routes into the main server setup, including stub repository wiring.
- Provided unit test for the Pay Day index endpoint.
- Updated template parsing to include Pay Day templates.
- Adjusted imports and added necessary stub constructors.

All new code follows project conventions, passes `go test ./...`, `go vet ./...`, and is formatted with `gofmt`.