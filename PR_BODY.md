### Summary
Implemented navigation menu and application layout templates matching the Scala reference. Added server setup to render these templates.

### Requirements addressed
- Added `views/layouts/menu.html` and `application.html` with exact style (tables, not lists).
- Updated `main.go` to parse templates and serve a home page using the new layout.
- Embedded template files via `embed.FS`.
- Added necessary imports and static asset handling.

### Implementation details
- Templates defined with `{{ define "menu" }}` and `{{ define "application.tmpl" }}`.
- Server uses chi router, loads env, runs DB migrations (existing).
- Home handler renders `application.tmpl` with title "Home" and simple body.

### Validation
- `go test ./...` passes.
- `go vet ./...` passes.
- `gofmt -l .` shows no formatting issues.

### Related work
- Vikunja Task ID: 7 - add nav and menu
