# Project Context

## Purpose
**go-admin** is a backend separation authority management system based on Gin (Go web framework) with support for Vue/Element UI and Arco Design frontends. It provides a complete RBAC (Role-Based Access Control) system with multi-organization support, data permissions, and code generation capabilities.

Key goals:
- Simplify admin system initialization with minimal configuration
- Provide comprehensive user/role/menu/permission management
- Support multiple database backends (MySQL, PostgreSQL, SQLite, SQL Server)
- Offer code generation tools to accelerate development
- Implement flexible data scope permissions

## Tech Stack
- **Language**: Go 1.24
- **Web Framework**: Gin (github.com/gin-gonic/gin)
- **ORM**: GORM (gorm.io/gorm)
- **Authorization**: Casbin (github.com/casbin/casbin/v2) for RBAC
- **CLI Framework**: Cobra (github.com/spf13/cobra)
- **Authentication**: JWT (golang.org/x/crypto, golang-jwt/jwt)
- **API Documentation**: Swagger (github.com/swaggo/gin-swagger)
- **Database Drivers**: MySQL, PostgreSQL, SQLite, SQL Server
- **Logging**: Uber Zap (go.uber.org/zap)
- **Configuration**: Viper (github.com/spf13/viper)
- **Monitoring**: Prometheus (github.com/prometheus/client_golang)
- **Rate Limiting**: Alibaba Sentinel (github.com/alibaba/sentinel-golang)
- **Scheduled Jobs**: Cron (github.com/robfig/cron/v3)
- **File Storage**: Aliyun OSS, Huawei OBS, Qiniu Kodo support
- **Build Tools**: Docker, Make

## Project Conventions

### Code Style
- **Naming Conventions**:
  - Files: `snake_case.go` (e.g., `sys_user.go`, `sys_role.go`)
  - Types: `PascalCase` for exported, `camelCase` for unexported
  - Functions/Methods: `PascalCase` for exported, `camelCase` for unexported
  - Variables: `camelCase`
  - Constants: `PascalCase` or `UPPER_SNAKE_CASE`
  - Database tables: `snake_case` with `sys_` prefix for system tables (e.g., `sys_user`, `sys_role`)
  - JSON fields: `camelCase` (e.g., `userId`, `userName`)

- **Formatting**: Standard Go formatting (`gofmt`)
  - Use tabs for indentation
  - Maximum line length: no strict limit, but prefer readability
  - Exported functions must have documentation comments
  - API endpoints use Swagger annotations

- **Package Structure**:
  - Each domain module (admin, jobs, other) has: `apis/`, `models/`, `service/`, `router/`, `service/dto/`
  - Common utilities in `common/` directory
  - Configuration in `config/` directory

- **API Design**:
  - RESTful conventions
  - Routes follow pattern: `/api/v1/{resource}` (e.g., `/api/v1/sys-user`)
  - HTTP methods: GET (query), POST (create), PUT (update), DELETE (delete)
  - Response format: `{code, data, msg}` structure
  - Use DTO (Data Transfer Objects) in service layer for request/response

### Architecture Patterns

**Three-Layer Architecture**:
```
┌─────────────────┐
│   APIs Layer    │  ← Controllers (app/*/apis/)
├─────────────────┤
│  Service Layer  │  ← Business Logic (app/*/service/)
├─────────────────┤
│   Model Layer   │  ← Data Access (app/*/models/)
└─────────────────┘
```

**Key Patterns**:
1. **Active Record Pattern**: Models implement `Generate()` and `GetId()` from `models.ActiveRecord` interface
2. **DTO Pattern**: Request/response data structures in `service/dto/` packages
3. **Middleware Chain**: Authentication → Authorization → Logging → Request Handling
4. **Context-based Request Flow**: Use `gin.Context` to pass request-scoped data
5. **Dependency Injection**: Services and dependencies injected via `MakeService()`, `MakeOrm()`, `MakeContext()`

**Directory Structure**:
- `app/admin/` - Core admin system (users, roles, menus, depts, etc.)
- `app/jobs/` - Scheduled job management
- `app/other/` - Auxiliary features (file storage, monitoring, code generation)
- `cmd/` - CLI commands (server, migrate, config)
- `common/` - Shared utilities, middleware, models
- `config/` - Configuration files (settings.yml)
- `docs/` - Swagger documentation
- `static/` - Static assets
- `template/` - Code generation templates

**Data Permission Model**:
- Tree-structured departments (SysDept)
- Users belong to departments and posts
- Roles assigned to users with menu permissions
- Data scope permissions: All, Custom, Department, Department+Children, Self Only

### Testing Strategy
- **Current State**: Limited test coverage (marked as TODO in README)
- **Existing Tests**: File storage operations (OSS, OBS, Kodo) in `common/file_store/`
- **Test Framework**: Standard Go `testing` package
- **Coverage Goal**: Unit tests planned but not yet implemented
- **Integration Tests**: Manual testing via Swagger UI and front-end applications

### Git Workflow

**Branching Strategy**:
- **master**: Main production branch, protected
- **feature/***: New feature branches
- **bugfix/***: Bug fix branches
- PRs must be based on `master` branch

**Commit Conventions**:
- PR template requires categorizing changes (feature, bugfix, refactor, etc.)
- Multi-language changelog (English + Chinese)
- Self-checklist before merge (docs, demos, TypeScript definitions)

**CI/CD Pipeline** (`.github/workflows/build.yml`):
1. Triggers: Push/PR to `master` branch
2. Steps:
   - Checkout code
   - Set up Go 1.24
   - Run `go mod tidy`
   - Build with CGO_ENABLED=1, SQLite tags
   - Build Docker image
   - Push to Aliyun Container Registry
   - Deploy to production server via SSH
3. Additional workflows: CodeQL analysis, issue automation

**Makefile Targets**:
- `make build` - Standard Go build (CGO_ENABLED=0)
- `make build-linux` - Docker build
- `make build-sqlite` - Build with SQLite support
- `make run` - Docker Compose deployment
- `make stop` - Stop containers

## Domain Context

**Core Entities**:
- **SysUser**: System users with authentication, department/role assignments
- **SysRole**: Roles with menu permissions and data scope settings
- **SysMenu**: Menu tree with API and button permissions
- **SysDept**: Hierarchical department structure for data permissions
- **SysPost**: Job positions within the organization
- **SysApi**: API endpoints registered for permission control
- **SysDictType/SysDictData**: Dictionary management for system configurations
- **SysConfig**: Runtime configuration parameters
- **SysLoginLog/SysOperaLog**: Audit logs for login and operations

**Permission Model**:
- Casbin-based RBAC with data scope control
- Permissions can be at menu, API, or button level
- Data scopes: All Data, Custom Department, Department Only, Department+Children, Self Only
- Supports multiple roles per user with combined permissions

**Code Generation**:
- Reads database table structure
- Generates: Models, APIs, Services, DTOs, Routers
- Templates in `template/` directory
- Supports visual configuration and customization

## Important Constraints

**Technical Constraints**:
- Go version requirement: 1.24+
- Database support: MySQL 5.7+, PostgreSQL 9.6+, SQLite 3.x, SQL Server 2012+
- Windows builds require CGO for SQLite (GCC compiler needed)
- HTTP server timeout: Read/write timeout configured (recently added)
- Rate limiting via Sentinel for API protection

**Business Constraints**:
- Admin system default credentials: admin / 123456 (must be changed in production)
- Multi-language support: Chinese and English
- Default data port: 8000
- Swagger docs available at `/swagger/index.html`

**Security Considerations**:
- Passwords hashed with bcrypt
- JWT token authentication
- CORS middleware configurable
- Secure headers middleware enabled
- Casbin enforces API-level authorization
- Data scope permissions prevent unauthorized data access

## External Dependencies

**Storage Services**:
- Aliyun OSS (Object Storage Service)
- Huawei OBS (Object Storage Service)
- Qiniu Kodo (Cloud Storage)

**Infrastructure**:
- Docker & Docker Compose for containerization
- Aliyun Container Registry (registry.ap-northeast-1.aliyuncs.com)
- Prometheus metrics endpoint (configurable)

**Development Tools**:
- GoLand IDE (JetBrains open source license)
- Swagger UI for API testing
- Code generation tools (built-in)

**Documentation**:
- Official docs: https://www.go-admin.dev
- Frontend repo: https://github.com/go-admin-team/go-admin-ui
- Video tutorials: Bilibili channel
