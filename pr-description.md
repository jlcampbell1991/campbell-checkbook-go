Implemented navigation menu and application layout enhancements:

- Added `.env.example` with default `PORT=8080`.
- Copied Scala CSS stylesheet to `static/stylesheets/application.css`.
- Updated `main.go` to load environment variables via `internal/config`, read `PORT` (default 8080), and log the listening address.
- Imported `os` and config package, invoked `config.Load()`.
- Fixed menu template to use double‑quoted `id="menu"` for exact HTML parity with Scala version.
- Adjusted server start logic to use the configurable address.
- Ensured static assets are served correctly.

These changes provide configurable server port, correct navigation menu markup, and proper styling matching the original Scala implementation.