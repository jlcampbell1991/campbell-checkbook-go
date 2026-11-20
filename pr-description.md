Implemented navigation menu and layout enhancements:

- Added `.env.example` with `PORT=8080` for configurable server port.
- Copied Scala stylesheet to `static/stylesheets/application.css` and ensured static serving.
- Updated `views/layouts/menu.tmpl` to use double‑quoted `id="menu"` matching Scala template.
- Modified `main.go`:
  * Imported `os` and the internal config package.
  * Loaded environment variables via `config.Load()`.
  * Determined listening port from `PORT` env var (default 8080).
  * Logged the actual address using `log.Printf`.
- Adjusted imports and formatting per Go conventions.

These changes bring the Go UI into visual parity with the Scala implementation and make the server port configurable via environment variables.
