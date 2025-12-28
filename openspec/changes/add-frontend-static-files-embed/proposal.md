# Change: Embed Frontend Static Files into Go Binary

## Why
Currently, the project serves a default welcome page at the root path `/`, which is a placeholder that redirects to an external site. To provide a complete self-contained application, we need to:

1. **Eliminate external deployment complexity** - Frontend static files should be embedded directly into the Go binary
2. **Simplify deployment** - Deploy a single executable file without needing separate static file hosting
3. **Provide proper SPA routing** - Support Single Page Application routing patterns

This change enables the project to serve a pre-built Vue.js frontend application directly from the embedded filesystem, making the application truly portable and self-contained.

## What Changes
- **Create embed package** (`web/dist/static.go`) - Use Go 1.16+ `embed` directive to embed frontend static files
- **Add static file serving routes** (`cmd/api/server.go`) - Serve CSS, JS, fonts, images, and favicon from embedded filesystem
- **Remove placeholder welcome page** (`app/admin/router/sys_router.go`) - Comment out the default iframe welcome page
- **Implement SPA routing support** - All non-API routes return `index.html` for client-side routing
- **Update `.gitignore`** - Ignore `web/dist/` build artifacts while preserving the embed source file

## Impact
- **Affected specs**: frontend-static-file (new capability)
- **Affected code**:
  - `web/dist/static.go` (new) - Embed filesystem with frontend assets
  - `cmd/api/server.go` - Static file serving routes
  - `app/admin/router/sys_router.go` - Remove welcome page route
  - `.gitignore` - Ignore build artifacts

- **User-visible changes**:
  - ✅ Single executable deployment - No separate static file hosting needed
  - ✅ Proper SPA routing support - Vue.js app works correctly
  - ✅ Faster startup - No external file I/O for static assets
  - ✅ Cross-platform compatibility - Works the same on all platforms

- **Migration path**:
  - Existing deployments continue to work - API endpoints unchanged
  - Static files must be pre-built and placed in `web/dist/` before compilation
  - To update frontend: rebuild `web/dist/`, then recompile Go binary

## Implementation Status

### ✅ Completed

**Core Tasks**: 8/8 (100%)

1. **Embed Package** ✅
   - Created `web/static.go` with `//go:embed dist/*` directive
   - Files embedded at compile time from `web/dist/` directory
   - Source file placed in `web/` (not `web/dist/`) to prevent deletion during frontend rebuilds
   - Files accessed with `dist/` prefix (e.g., `dist/css/app.css`)

2. **Static File Routes** ✅
   - Added routes for `/css/*`, `/js/*`, `/fonts/*`, `/img/*`
   - Added `/favicon.ico` route
   - Root `/` route serves `index.html`
   - SPA fallback route for client-side routing
   - All file reads use `dist/` prefix to access embedded filesystem

3. **Welcome Page Removal** ✅
   - Commented out default iframe welcome page in `sys_router.go`

4. **Content-Type Handling** ✅
   - Proper MIME types for HTML, CSS, JS, images, fonts

5. **Git Configuration** ✅
   - `.gitignore` updated to exclude `web/dist/` build artifacts

6. **Build Verification** ✅
   - Compilation successful with embedded files
   - Binary size includes static assets
   - No external file dependencies

7. **Runtime Testing** ✅
   - Homepage returns Vue SPA HTML
   - Static assets (CSS, JS, favicon) served correctly
   - SPA routing works (fallback to index.html)

### Technical Details

**Project Structure**:
```
web/
├── dist/              # Frontend build output (generated)
│   ├── css/           # Stylesheets
│   ├── js/            # JavaScript bundles
│   ├── fonts/         # Font files
│   ├── img/           # Images
│   ├── index.html     # SPA entry point
│   └── favicon.ico    # Site icon
└── static.go          # Embed package source (preserved)
```

**Embed File (`web/static.go`)**:
```go
package web

import "embed"

//go:embed dist/*
var WebFS embed.FS
```

**Served Routes**:
| Route Pattern | Embed File Path | Description |
|---------------|-----------------|-------------|
| `/` | `dist/index.html` | Serves `index.html` |
| `/css/*filepath` | `dist/css/*filepath` | CSS files |
| `/js/*filepath` | `dist/js/*filepath` | JavaScript files |
| `/fonts/*filepath` | `dist/fonts/*filepath` | Font files |
| `/img/*filepath` | `dist/img/*filepath` | Image files |
| `/favicon.ico` | `dist/favicon.ico` | Favicon |
| `NoRoute` (fallback) | `dist/index.html` | Returns `index.html` for SPA routes |
| `/api/*` | N/A | API endpoints (passed through) |
| `/swagger/*` | N/A | Swagger docs (passed through) |

### Build Process

1. **Build frontend** (external step):
   ```bash
   cd web && npm run build
   ```
   This generates files in `web/dist/`

2. **Build Go binary**:
   ```bash
   go build -o go-admin.exe main.go
   ```
   This embeds `web/dist/*` into the binary

3. **Deploy**:
   ```bash
   ./go-admin.exe server
   ```
   Single executable, no additional files needed
