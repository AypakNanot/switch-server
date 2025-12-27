## ADDED Requirements

### Requirement: SQLite3 as Default Database
The system SHALL use SQLite3 as the default database driver for new installations and development environments.

#### Scenario: Fresh project start
- **WHEN** a developer clones and runs the project without modifying configuration
- **THEN** the system SHALL start successfully using SQLite3
- **AND** data SHALL be stored in `./go-admin-db.db` file by default

#### Scenario: First-time database initialization
- **WHEN** running `./go-admin migrate -c config/settings.yml` with default config
- **THEN** the system SHALL create all required tables in SQLite database
- **AND** initialization SHALL complete without requiring external database server

### Requirement: Alternative Database Configuration
The system SHALL support easy configuration switching to MySQL, PostgreSQL, or SQL Server.

#### Scenario: Switch to MySQL
- **WHEN** user modifies `settings.yml` to use MySQL driver and connection string
- **THEN** the system SHALL connect to MySQL database on next restart
- **AND** all features SHALL work identically to SQLite3

#### Scenario: Switch to PostgreSQL
- **WHEN** user modifies `settings.yml` to use PostgreSQL driver and connection string
- **THEN** the system SHALL connect to PostgreSQL database on next restart
- **AND** all features SHALL work identically to SQLite3

## MODIFIED Requirements

### Requirement: Build Configuration
The system build process SHALL include SQLite3 support by default for all platforms.

#### Scenario: Standard build
- **WHEN** running `make build` or `go build`
- **THEN** the resulting binary SHALL support SQLite3 database
- **AND** CGO SHALL be enabled for SQLite support

#### Scenario: Docker build
- **WHEN** running `docker build -t go-admin .`
- **THEN** the Docker image SHALL include SQLite3 library support
- **AND** the containerized application SHALL use SQLite3 successfully

### Requirement: Quick Start Documentation
The README documentation SHALL provide SQLite3-based quick start instructions as the primary method.

#### Scenario: New developer onboarding
- **WHEN** a new developer reads the README
- **THEN** they SHALL see SQLite3 as the default database option
- **AND** they SHALL be able to start the application without installing MySQL
- **AND** alternative database options SHALL be clearly documented

### Requirement: Database Configuration Examples
The system SHALL provide example configuration files for all supported database types.

#### Scenario: Explore configuration options
- **WHEN** user explores the `config/` directory
- **THEN** they SHALL find `settings.yml` (SQLite3 default)
- **AND** they SHALL find `settings.mysql.yml` example
- **AND** they SHALL find `settings.postgres.yml` example
- **AND** they SHALL find `settings.sqlite.yml` reference
