# Implementation Tasks

## 1. Embed Package Setup
- [x] 1.1 Create `web/static.go` with `//go:embed dist/*` directive
- [x] 1.2 Add build tag to exclude when needed (`//go:build !exclude_webdist`)
- [x] 1.3 Export `WebFS embed.FS` variable for use by server
- [x] 1.4 Move from `web/dist/` to `web/` to prevent deletion during frontend rebuilds

## 2. Static File Serving Routes
- [x] 2.1 Import `go-admin/web` package in `cmd/api/server.go`
- [x] 2.2 Add `strings` package import for path handling
- [x] 2.3 Create `serveStaticFile` helper function with Content-Type detection
- [x] 2.4 Add `/css/*filepath` route (reads from `dist/css/*`)
- [x] 2.5 Add `/js/*filepath` route (reads from `dist/js/*`)
- [x] 2.6 Add `/fonts/*filepath` route (reads from `dist/fonts/*`)
- [x] 2.7 Add `/img/*filepath` route (reads from `dist/img/*`)
- [x] 2.8 Add `/favicon.ico` route (reads from `dist/favicon.ico`)

## 3. Root and SPA Routes
- [x] 3.1 Add `/` route to serve `dist/index.html`
- [x] 3.2 Add `NoRoute` handler for SPA fallback (serves `dist/index.html`)
- [x] 3.3 Exclude `/api/*` paths from SPA fallback
- [x] 3.4 Exclude `/swagger/*` paths from SPA fallback
- [x] 3.5 Exclude `/info` path from SPA fallback

## 4. Welcome Page Removal
- [x] 4.1 Comment out default welcome page route in `app/admin/router/sys_router.go`
- [x] 4.2 Verify no conflicts with new static routes

## 5. Content-Type Handling
- [x] 5.1 Add HTML content type (`.html`)
- [x] 5.2 Add CSS content type (`.css`)
- [x] 5.3 Add JavaScript content type (`.js`)
- [x] 5.4 Add JSON content type (`.json`)
- [x] 5.5 Add image content types (`.png`, `.jpg`, `.jpeg`, `.svg`)
- [x] 5.6 Add favicon content type (`.ico`)
- [x] 5.7 Add font content types (`.woff`, `.woff2`)

## 6. Git Configuration
- [x] 6.1 Add `web/dist/` to `.gitignore`
- [x] 6.2 Add exception for `!web/dist/embed/*.go` (not needed, using different approach)
- [x] 6.3 Finalize `.gitignore` with `web/dist/` pattern

## 7. Build Verification
- [x] 7.1 Run `go build` successfully with no errors
- [x] 7.2 Verify binary includes embedded files
- [x] 7.3 Check binary size reflects embedded assets

## 8. Runtime Testing
- [x] 8.1 Start server and verify homepage returns Vue SPA HTML
- [x] 8.2 Test `/favicon.ico` returns 3628 bytes
- [x] 8.3 Test CSS file serves with correct content
- [x] 8.4 Test SPA fallback returns `index.html` for unknown routes
- [x] 8.5 Verify API routes still work (`/api/v1/*`)
- [x] 8.6 Verify Swagger still works (`/swagger/admin/*`)

## ✅ Implementation Complete

### Test Results

| Test Case | Expected | Result |
|-----------|----------|--------|
| `GET /` | Vue SPA HTML | ✅ Pass |
| `GET /favicon.ico` | 3628 bytes, image/x-icon | ✅ Pass |
| `GET /css/app.*.css` | CSS content | ✅ Pass |
| `GET /js/app.*.js` | JS content | ✅ Pass |
| `GET /unknown` | Falls back to index.html | ✅ Pass |
| `GET /api/v1/login` | API endpoint works | ✅ Pass |
| `GET /swagger/admin/*` | Swagger works | ✅ Pass |

### Files Modified

**Created**:
- `web/static.go` - Embed package (placed in parent directory to prevent deletion)

**Modified**:
- `cmd/api/server.go` - Static file routes with `dist/` prefix paths
- `app/admin/router/sys_router.go` - Removed welcome page
- `.gitignore` - Ignore `web/dist/` but preserve `web/static.go`

### Deployed Binary Size
- **Before**: ~48MB (without frontend)
- **After**: ~48MB (frontend embedded)
- Note: Frontend build files are relatively small (~500KB compressed)

### Usage
```bash
# 1. Build frontend (do this in your frontend project)
cd web && npm run build

# 2. Build Go binary (this embeds web/dist/*)
go build -o go-admin.exe main.go

# 3. Run - single executable, no other files needed!
./go-admin.exe server
```

Access at: `http://localhost:8000/`
