## ADDED Requirements

### Requirement: Cross-Platform Build Support
The system SHALL support cross-compilation to multiple CPU architectures commonly found in network switches and embedded devices.

#### Scenario: Build for ARMv7
- **GIVEN** the project uses pure Go dependencies (no CGO)
- **WHEN** developer runs `make build-armv7`
- **THEN** a binary for ARMv7 architecture is created
- **AND** the binary is statically linked (no external glibc dependencies)
- **AND** the binary can run on ARMv7 devices (e.g., Raspberry Pi 2/3, many switches)

#### Scenario: Build for MIPS Little-Endian
- **GIVEN** many network switches use MIPS architecture
- **WHEN** developer runs `make build-mipsle`
- **THEN** a binary for MIPS little-endian architecture is created
- **AND** the binary is statically linked
- **AND** the binary can run on common MIPS switches/routers

#### Scenario: Build for ARM64
- **GIVEN** newer switches use ARM64 (aarch64) architecture
- **WHEN** developer runs `make build-arm64`
- **THEN** a binary for ARM64 architecture is created
- **AND** the binary is statically linked
- **AND** the binary can run on ARM64 devices (e.g., Raspberry Pi 4/5)

#### Scenario: Build All Architectures
- **GIVEN** developer wants to build binaries for all supported architectures
- **WHEN** developer runs `make build-switch`
- **THEN** binaries are created for common switch architectures (armv7, arm64, mipsle, mips)
- **AND** each binary is named with architecture suffix (e.g., `go-admin-armv7`)

### Requirement: Low-Memory Configuration
The system SHALL provide a configuration file optimized for low-memory environments (256MB-512MB RAM).

#### Scenario: Configuration uses minimal resources
- **GIVEN** the switch has limited memory (256MB-512MB)
- **WHEN** the system starts with `config/settings.switch.yml`
- **THEN** database connection pool is limited (maxOpenConns: 5, maxIdleConns: 2)
- **AND** memory queue pool size is reduced (poolSize: 20)
- **AND** debug logging is disabled (mode: prod, level: warn)
- **AND** database query logging is disabled (enableddb: false)

#### Scenario: System runs within memory constraints
- **GIVEN** a switch with 256MB RAM
- **WHEN** the system is started with switch configuration
- **THEN** the application runs without out-of-memory errors
- **AND** basic API operations work correctly
- **AND** static file serving works correctly

### Requirement: Automated Deployment Script
The system SHALL provide a script for automated deployment to switch devices via SSH.

#### Scenario: Deploy binary to switch
- **GIVEN** a compiled binary for the target architecture
- **AND** SSH access to the switch device
- **WHEN** developer runs `./scripts/deploy-to-switch.sh --host=192.168.1.1 --user=root --arch=armv7`
- **THEN** the script connects to the switch via SSH
- **AND** backs up any existing binary
- **AND** uploads the new binary to `/usr/bin/go-admin`
- **AND** sets executable permissions (chmod +x)
- **AND** starts the service
- **AND** verifies the service is running (health check)

#### Scenario: Deployment fails with rollback
- **GIVEN** a deployment is in progress
- **WHEN** the health check fails after starting the new binary
- **THEN** the script stops the new binary
- **AND** restores the previous binary from backup
- **AND** restarts the old binary
- **AND** reports the failure to the user

#### Scenario: Restart existing deployment
- **GIVEN** the application is already deployed on a switch
- **WHEN** developer runs `./scripts/deploy-to-switch.sh --host=192.168.1.1 --user=root --action=restart`
- **THEN** the script connects to the switch
- **AND** restarts the application (stop then start)
- **AND** verifies the service is running

### Requirement: Static Binary Verification
The system SHALL provide tools to verify that compiled binaries are statically linked and target the correct architecture.

#### Scenario: Verify ARM binary is static
- **GIVEN** a compiled binary for ARM architecture
- **WHEN** developer runs `file go-admin-armv7`
- **THEN** output shows "ARM" as the architecture
- **AND** output shows "statically linked" (no dynamically linked dependencies)

#### Scenario: Verify no external dependencies
- **GIVEN** a compiled binary
- **WHEN** developer runs `ldd go-admin-armv7` on Linux or `otool -L` on macOS
- **THEN** command reports "not a dynamic executable" or similar
- **AND** the binary has no external shared library dependencies

### Requirement: Architecture Documentation
The system SHALL document supported architectures and corresponding device types.

#### Scenario: Architecture reference guide
- **GIVEN** developer needs to know which architecture to compile for
- **WHEN** developer reads `docs/switch-deployment.md`
- **THEN** the document lists all supported architectures (armv5, armv6, armv7, arm64, mips, mipsle, mips64, mips64le, ppc64)
- **AND** for each architecture, typical device examples are provided
- **AND** compiler flags (GOARM version) are explained

#### Scenario: Device-specific deployment notes
- **GIVEN** developer is deploying to a specific switch model
- **WHEN** developer reads `docs/switch-deployment.md`
- **THEN** the document provides model-specific notes if available
- **AND** common issues and solutions are documented
- **AND** performance expectations are listed

### Requirement: Binary Size Optimization
The system SHALL optimize binary size for constrained storage environments.

#### Scenario: Strip debugging symbols
- **GIVEN** the build process creates a binary
- **WHEN** any build target is executed
- **THEN** the build uses `-ldflags="-w -s"` to remove debug info
- **AND** binary size is reduced by ~30-50% compared to unstripped binary

#### Scenario: Report binary sizes
- **GIVEN** developer runs `make build-switch`
- **WHEN** build completes
- **THEN** the output displays the size of each compiled binary
- **AND** developer can verify the binary fits in available storage

### Requirement: Service Management
The system SHALL provide service management scripts for common init systems.

#### Scenario: Systemd service file
- **GIVEN** a switch running Linux with systemd
- **WHEN** developer copies `scripts/go-admin.service` to `/etc/systemd/system/`
- **THEN** the service can be managed with `systemctl start/stop/restart go-admin`
- **AND** the service starts automatically on boot if enabled

#### Scenario: Init.d script
- **GIVEN** a switch running Linux with traditional init.d
- **WHEN** developer copies `scripts/go-admin.init` to `/etc/init.d/go-admin`
- **THEN** the service can be managed with `/etc/init.d/go-admin start/stop/restart`
- **AND** the service starts automatically on boot if enabled

### Requirement: Configuration File Location
The system SHALL support configurable configuration file paths for different deployment scenarios.

#### Scenario: Custom config file location
- **GIVEN** the switch has limited writable storage
- **WHEN** starting the application with `-c /tmp/go-admin-config.yml`
- **THEN** the application loads configuration from the specified path
- **AND** the application runs normally with the custom configuration

#### Scenario: Default config fallback
- **GIVEN** no config file is specified
- **WHEN** the application starts
- **THEN** the application tries to load `config/settings.yml` in the current directory
- **AND** if not found, tries to load `/etc/go-admin/settings.yml`
