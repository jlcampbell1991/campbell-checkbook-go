Implemented navigation menu and layout enhancements:

- Added `.env.example` with `PORT=8080` for configurable server port.
- Updated `main.go`:
  * Imported `os` and the internal `config` package.
  * Loaded environment variables via `config.Load()`.
  * Replaced hard‑coded `:8080` with a configurable `PORT` (default 8080).
  * Adjusted logging to show the actual listening address.
- Fixed navigation menu template to use double‑quoted `id="menu"` for exact HTML parity.
- Added static CSS file `static/stylesheets/application.css` copied from the Scala reference implementation.
- Ensured templates reference the new stylesheet correctly (no code change needed).

These changes bring the Go application UI in line with the Scala version and allow the HTTP port to be configured through environment variables.
