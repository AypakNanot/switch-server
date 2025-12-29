# Implementation Tasks

## 1. Core Module Changes
- [ ] 1.1 Update `go.mod` module name from `go-admin` to `opt-switch`
- [ ] 1.2 Run `go mod tidy` to update dependencies
- [ ] 1.3 Update all import paths in Go source files
- [ ] 1.4 Test compilation after module rename

## 2. Configuration Files
- [ ] 2.1 Update `config/settings.yml` - system name
- [ ] 2.2 Update `config/settings.switch.yml` - switch config name
- [ ] 2.3 Update `config/db.sql` - sys_config values (sys_app_name, sys_app_logo)
- [ ] 2.4 Update database filename references

## 3. Binary Output
- [ ] 3.1 Update Makefile BINARY_NAME to `opt-switch`
- [ ] 3.2 Update Dockerfile binary names
- [ ] 3.3 Update deployment scripts

## 4. Documentation
- [ ] 4.1 Update README.md project name
- [ ] 4.2 Update README.Zh-cn.md project name
- [ ] 4.3 Update code comments containing "go-admin"
- [ ] 4.4 Update OpenSpec documentation references

## 5. Build Verification
- [ ] 5.1 Clean build: `go clean && go build`
- [ ] 5.2 Verify binary name is `opt-switch`
- [ ] 5.3 Test application startup
- [ ] 5.4 Verify API endpoints work
- [ ] 5.5 Test login page displays correctly

## 6. Docker Verification
- [ ] 6.1 Rebuild Docker image
- [ ] 6.2 Test container startup
- [ ] 6.3 Verify web interface accessible

## 7. Code Search & Replace
- [ ] 7.1 Search for remaining "go-admin" references in code
- [ ] 7.2 Update any remaining hardcoded strings
- [ ] 7.3 Verify no broken references

## 8. Testing
- [ ] 8.1 Test basic CRUD operations
- [ ] 8.2 Test authentication flow
- [ ] 8.3 Test database operations
- [ ] 8.4 Verify static file serving
