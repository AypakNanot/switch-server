# Spec: Project Rename

## Scope
此规格说明项目重命名的相关要求。

## MODIFIED Requirements

### Requirement: Project Name Consistency
The system SHALL use consistent naming across all components.

#### Scenario: Module name matches project name
- **GIVEN** the project is renamed to opt-switch
- **WHEN** the go.mod file is examined
- **THEN** the module name SHALL be `opt-switch`
- **AND** all import paths SHALL use `opt-switch` as the base

#### Scenario: Binary name matches project name
- **GIVEN** the project is compiled
- **WHEN** the build process completes
- **THEN** the output binary SHALL be named `opt-switch`
- **AND** the binary name SHALL be consistent across platforms (Windows: opt-switch.exe, Linux: opt-switch)

#### Scenario: System name in configuration
- **GIVEN** the application starts
- **WHEN** the system configuration is loaded
- **THEN** the system name SHALL be "opt-switch管理系统"
- **AND** this name SHALL be displayed in the UI

### Requirement: Backward Compatibility
The system SHALL maintain compatibility with existing deployments.

#### Scenario: Database compatibility
- **GIVEN** an existing database from go-admin
- **WHEN** the system is upgraded to opt-switch
- **THEN** the database SHALL continue to work without modifications
- **AND** existing data SHALL be preserved

#### Scenario: API compatibility
- **GIVEN** existing clients using the API
- **WHEN** the project is renamed
- **THEN** all API endpoints SHALL remain unchanged
- **AND** existing API clients SHALL continue to work
