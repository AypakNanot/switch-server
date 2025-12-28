## ADDED Requirements

### Requirement: Static File Embedding
The system SHALL embed frontend static files into the Go binary at compile time using Go's `embed` directive.

#### Scenario: Frontend files are embedded
- **GIVEN** a `web/dist/` directory containing compiled frontend assets
- **AND** a `web/static.go` file with `//go:embed dist/*` directive
- **WHEN** the project is compiled with `go build`
- **THEN** all static files from `web/dist/` are embedded into the binary executable
- **AND** the binary contains all HTML, CSS, JavaScript, fonts, and images
- **AND** files are accessed in the embed filesystem with `dist/` prefix (e.g., `dist/css/app.css`)

#### Scenario: Build excludes webdist when tag is set
- **GIVEN** the build tag `exclude_webdist` is set
- **WHEN** the project is compiled with that tag
- **THEN** the static files are NOT embedded into the binary
- **AND** the system can run without frontend assets for API-only deployments

### Requirement: Static File Routes
The system SHALL serve embedded static files through HTTP routes with appropriate content-type headers.

#### Scenario: CSS files are served
- **GIVEN** the system is running with embedded static files
- **WHEN** a client requests `/css/app.fdc3c312.css`
- **THEN** the system returns the CSS file content
- **AND** the `Content-Type` header is set to `text/css; charset=utf-8`

#### Scenario: JavaScript files are served
- **GIVEN** the system is running with embedded static files
- **WHEN** a client requests `/js/app.1f40785a.js`
- **THEN** the system returns the JavaScript file content
- **AND** the `Content-Type` header is set to `application/javascript; charset=utf-8`

#### Scenario: Font files are served
- **GIVEN** the system is running with embedded static files
- **WHEN** a client requests `/fonts/font.woff2`
- **THEN** the system returns the font file content
- **AND** the `Content-Type` header is set to `font/woff2`

#### Scenario: Image files are served
- **GIVEN** the system is running with embedded static files
- **WHEN** a client requests `/img/logo.png`
- **THEN** the system returns the image file content
- **AND** the `Content-Type` header is set to `image/png`

#### Scenario: Favicon is served
- **GIVEN** the system is running with embedded static files
- **WHEN** a client requests `/favicon.ico`
- **THEN** the system returns the favicon file content
- **AND** the `Content-Type` header is set to `image/x-icon`

### Requirement: Root Route Serves SPA
The system SHALL serve the `index.html` file when the root path `/` is requested.

#### Scenario: Root path returns HTML
- **GIVEN** the system is running with embedded static files
- **WHEN** a client requests `GET /`
- **THEN** the system returns the `index.html` content
- **AND** the `Content-Type` header is set to `text/html; charset=utf-8`

### Requirement: SPA Fallback Routing
The system SHALL return `index.html` for all non-API routes to support client-side routing.

#### Scenario: Unknown routes return index.html
- **GIVEN** the system is running with embedded static files
- **WHEN** a client requests any path that is NOT `/api/*`, `/swagger/*`, `/static/*`, `/info`, or a static asset
- **THEN** the system returns the `index.html` content
- **AND** the Vue.js router can handle the client-side routing

#### Scenario: API routes are not affected
- **GIVEN** the system is running with embedded static files
- **WHEN** a client requests `/api/v1/login`
- **THEN** the request is passed to the API handler
- **AND** the SPA fallback does NOT interfere

#### Scenario: Swagger routes are not affected
- **GIVEN** the system is running with embedded static files
- **WHEN** a client requests `/swagger/admin/index.html`
- **THEN** the request is passed to the Swagger handler
- **AND** the SPA fallback does NOT interfere

### Requirement: Build Artifact Exclusion
The system SHALL ignore frontend build artifacts in version control while preserving the embed source file.

#### Scenario: web/dist is ignored by git
- **GIVEN** the `.gitignore` file contains `web/dist/`
- **WHEN** a developer runs `git status`
- **THEN** files in `web/dist/` are NOT shown as untracked
- **AND** the embed source file `web/dist/static.go` is preserved (with exception pattern if needed)

#### Scenario: Frontend build process
- **GIVEN** a frontend project is configured to build to `web/dist/`
- **WHEN** the build process completes
- **THEN** the generated files are NOT tracked by git
- **AND** only the `static.go` embed source is committed

### Requirement: Content-Type Detection
The system SHALL automatically detect and set correct MIME types based on file extensions.

#### Scenario: HTML files
- **GIVEN** a file with `.html` extension
- **WHEN** served via static file handler
- **THEN** `Content-Type: text/html; charset=utf-8` is set

#### Scenario: CSS files
- **GIVEN** a file with `.css` extension
- **WHEN** served via static file handler
- **THEN** `Content-Type: text/css; charset=utf-8` is set

#### Scenario: JavaScript files
- **GIVEN** a file with `.js` extension
- **WHEN** served via static file handler
- **THEN** `Content-Type: application/javascript; charset=utf-8` is set

#### Scenario: JSON files
- **GIVEN** a file with `.json` extension
- **WHEN** served via static file handler
- **THEN** `Content-Type: application/json; charset=utf-8` is set

#### Scenario: PNG images
- **GIVEN** a file with `.png` extension
- **WHEN** served via static file handler
- **THEN** `Content-Type: image/png` is set

#### Scenario: JPEG images
- **GIVEN** a file with `.jpg` or `.jpeg` extension
- **WHEN** served via static file handler
- **THEN** `Content-Type: image/jpeg` is set

#### Scenario: SVG files
- **GIVEN** a file with `.svg` extension
- **WHEN** served via static file handler
- **THEN** `Content-Type: image/svg+xml` is set

#### Scenario: WOFF2 fonts
- **GIVEN** a file with `.woff2` extension
- **WHEN** served via static file handler
- **THEN** `Content-Type: font/woff2` is set

### Requirement: Single Binary Deployment
The system SHALL be deployable as a single executable file without requiring separate static file hosting.

#### Scenario: Deployment without external files
- **GIVEN** a compiled `go-admin.exe` binary with embedded static files
- **WHEN** deployed to a fresh server
- **THEN** the application serves the complete frontend
- **AND** no additional files or directories are required
- **AND** no static file hosting configuration is needed

#### Scenario: Cross-platform compatibility
- **GIVEN** the same source code compiled on different platforms
- **WHEN** the binary is executed on Windows, Linux, or macOS
- **THEN** the embedded static files work identically on all platforms
- **AND** no platform-specific file system issues occur
