package device

import "fmt"

// ErrorCode represents the device error code
type ErrorCode int

const (
	// Connection errors 1000-1099
	ErrConnectionFailed ErrorCode = 1001
	ErrAuthFailed       ErrorCode = 1002
	ErrConnectionClosed ErrorCode = 1003

	// Queue errors 1100-1199
	ErrQueueFull    ErrorCode = 1101
	ErrQueueTimeout ErrorCode = 1102

	// Execution errors 1200-1299
	ErrCommandFailed  ErrorCode = 1201
	ErrCommandTimeout ErrorCode = 1202
	ErrOutputTooLarge ErrorCode = 1203

	// Config errors 1300-1399
	ErrInvalidConfig        ErrorCode = 1301
	ErrDeviceNotConfigured  ErrorCode = 1302
)

// Error messages mapping
var errorMessages = map[ErrorCode]string{
	ErrConnectionFailed:   "Failed to connect to device",
	ErrAuthFailed:         "Authentication failed",
	ErrConnectionClosed:   "Connection closed",
	ErrQueueFull:          "Command queue is full, please try again later",
	ErrQueueTimeout:       "Queue wait timeout",
	ErrCommandFailed:      "Command execution failed",
	ErrCommandTimeout:     "Command execution timeout",
	ErrOutputTooLarge:     "Command output too large, truncated",
	ErrInvalidConfig:      "Invalid device configuration",
	ErrDeviceNotConfigured: "Device not configured",
}

// DeviceError represents a device operation error
type DeviceError struct {
	Code    ErrorCode
	Message string
	Cause   error
}

// Error implements the error interface
func (e *DeviceError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("[%d] %s: %s", e.Code, errorMessages[e.Code], e.Message)
	}
	if e.Cause != nil {
		return fmt.Sprintf("[%d] %s: %v", e.Code, errorMessages[e.Code], e.Cause)
	}
	return fmt.Sprintf("[%d] %s", e.Code, errorMessages[e.Code])
}

// Unwrap returns the underlying cause
func (e *DeviceError) Unwrap() error {
	return e.Cause
}

// NewConnectionError creates a new connection error
func NewConnectionError(cause error) *DeviceError {
	return &DeviceError{
		Code:  ErrConnectionFailed,
		Cause: cause,
	}
}

// NewAuthError creates a new authentication error
func NewAuthError(cause error) *DeviceError {
	return &DeviceError{
		Code:  ErrAuthFailed,
		Cause: cause,
	}
}

// NewQueueTimeoutError creates a new queue timeout error
func NewQueueTimeoutError() *DeviceError {
	return &DeviceError{
		Code: ErrQueueTimeout,
	}
}

// NewQueueFullError creates a new queue full error
func NewQueueFullError() *DeviceError {
	return &DeviceError{
		Code: ErrQueueFull,
	}
}

// NewCommandTimeoutError creates a new command timeout error
func NewCommandTimeoutError() *DeviceError {
	return &DeviceError{
		Code: ErrCommandTimeout,
	}
}

// NewCommandFailedError creates a new command failed error
func NewCommandFailedError(cause error) *DeviceError {
	return &DeviceError{
		Code:  ErrCommandFailed,
		Cause: cause,
	}
}

// NewInvalidConfigError creates a new invalid config error
func NewInvalidConfigError(message string) *DeviceError {
	return &DeviceError{
		Code:    ErrInvalidConfig,
		Message: message,
	}
}

// NewConnectionClosed creates a new connection closed error
func NewConnectionClosed() *DeviceError {
	return &DeviceError{
		Code: ErrConnectionClosed,
	}
}
