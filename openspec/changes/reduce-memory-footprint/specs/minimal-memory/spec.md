# Spec: Minimal Memory Capability

## ADDED Requirements

### Requirement: Minimal Memory Configuration

**Given** the system is deployed on resource-constrained devices (128MB-256MB RAM)
**When** the operator selects `settings.minimal.yml` configuration
**THEN** the system SHALL:
- Limit idle memory usage to < 50MB
- Reduce database connection pool to 2 max / 1 idle
- Reduce memory queue pool to 5
- Disable frontend static files (API-only mode)
- Disable non-essential middleware (Sentinel, RequestID, Metrics)

#### Scenario: Operator selects minimal config on 128MB device

**GIVEN** a device with 128MB RAM
**AND** the operator starts opt-switch with `-c config/settings.minimal.yml`
**WHEN** the system starts
**THEN** memory usage (RSS) SHOULD be < 50MB at idle
**AND** core API endpoints SHOULD be functional
**AND** Web UI SHOULD NOT be available (404 on /)

#### Scenario: Verify memory target achievement

**GIVEN** opt-switch is running with minimal config
**WHEN** memory is measured after 60 seconds of idle time
**THEN** RSS (Resident Set Size) SHOULD be < 50MB
**AND** VSZ (Virtual Memory Size) SHOULD be < 100MB

### Requirement: Runtime Memory Tuning

**Given** the system is configured with minimal memory settings
**When** the application starts
**THEN** the system SHALL:
- Set GOMAXPROCS to 1 (configurable)
- Set GOGC to 200 (reduce GC frequency)
- Set memory limit to 60MB (Go 1.19+)
- Limit maximum threads to 100

#### Scenario: Runtime tuning applied at startup

**GIVEN** minimal configuration is selected
**WHEN** main() function calls initRuntime()
**THEN** GOMAXPROCS SHOULD be set to 1
**AND** GOGC SHOULD be set to 200
**AND** memory limit SHOULD be set to 60MB (if Go >= 1.19)

#### Scenario: Override runtime settings via config

**GIVEN** operator wants different runtime settings
**WHEN** `runtime.gomaxprocs` is set in config
**THEN** the specified value SHOULD be used instead of default

### Requirement: Conditional Feature Loading

**Given** minimal memory configuration is selected
**When** the system initializes
**THEN** the system SHALL:
- Disable frontend static file serving
- Disable Sentinel middleware
- Disable RequestID middleware
- Disable Metrics middleware
- Keep essential middleware (Recovery, Logger, Auth)

#### Scenario: Frontend disabled in minimal mode

**GIVEN** settings.minimal.yml has `enableFrontend: false`
**WHEN** user accesses http://switch:8000/
**THEN** SHOULD receive 404 or API-only response
**AND** CSS/JS endpoints SHOULD return 404

#### Scenario: API endpoints remain functional

**GIVEN** minimal mode is enabled
**WHEN** client calls `/api/v1/captcha`
**THEN** SHOULD receive valid response with status 200
**AND** `/api/v1/login` SHOULD work normally

### Requirement: Memory Measurement Tools

**Given** the operator needs to verify memory usage
**When** the operator runs measurement scripts
**THEN** the system SHALL provide:
- `scripts/measure-memory.sh` - Monitor process memory over time
- `scripts/check-memory.sh` - Verify available memory before start
- `scripts/memory-test.sh` - Automated memory testing

#### Scenario: Check available memory before deployment

**GIVEN** operator runs `./scripts/check-memory.sh`
**WHEN** available memory is >= 128MB
**THEN** script SHOULD exit with code 0 and print "OK"
**WHEN** available memory is < 128MB
**THEN** script SHOULD exit with code 1 and print "ERROR"

#### Scenario: Measure memory during operation

**GIVEN** opt-switch is running with PID 1234
**WHEN** operator runs `./scripts/measure-memory.sh 1234`
**THEN** script SHOULD print CSV output: Time,RSS(MB),VSZ(MB)
**AND** update every 5 seconds
**AND** continue until process exits or script is interrupted

### Requirement: Configuration Documentation

**Given** multiple configuration options exist
**When** the operator reviews documentation
**THEN** the system SHALL provide clear guidance on:
- When to use settings.yml (full features, high memory)
- When to use settings.switch.yml (balanced, normal deployment)
- When to use settings.minimal.yml (minimal memory, low concurrency)

#### Scenario: Configuration selection guide

**GIVEN** operator has device with 256MB RAM
**AND** expects 1-3 users
**WHEN** reading deployment guide
**THEN** SHOULD recommend settings.minimal.yml
**AND** document expected memory usage (~45MB)

#### Scenario: Configuration comparison table

**GIVEN** operator compares configurations
**WHEN** viewing README or deployment guide
**THEN** SHOULD see comparison table with:
  - Memory requirements
  - Concurrent user capacity
  - Enabled features
  - Use case recommendations

## MODIFIED Requirements

### Requirement: Deployment Configuration

**Previously**: deployment spec defined single optimized configuration for switch devices

**Now**: deployment spec SHALL support multiple configuration levels:
- **settings.yml** - Full features, ~80-100MB memory
- **settings.switch.yml** - Balanced deployment, ~60-80MB memory
- **settings.minimal.yml** - Minimal memory, ~35-50MB memory

#### Scenario: Operator chooses appropriate config for hardware

**GIVEN** operator has device with specific memory constraints
**WHEN** selecting configuration file
**THEN** SHOULD choose based on:
  - **< 256MB RAM**: use settings.minimal.yml
  - **256MB - 512MB**: use settings.switch.yml
  - **> 512MB**: use settings.yml

## Cross-References

- **Deployment spec**: Defines configuration levels and selection criteria
- **Performance spec**: May need updates for minimal mode performance characteristics
- **API spec**: API functionality unchanged in minimal mode (except frontend endpoints)
