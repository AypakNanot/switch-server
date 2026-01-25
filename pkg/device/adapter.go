package device

import "context"

// ProtocolType represents the protocol type
type ProtocolType string

const (
	ProtocolSSH    ProtocolType = "ssh"
	ProtocolTelnet ProtocolType = "telnet"
	ProtocolNETCONF ProtocolType = "netconf"
)

// ConnectionConfig holds the connection configuration
type ConnectionConfig struct {
	Protocol string
	Host     string
	Port     int
	Username string
	Password string
	Timeout  int // seconds
}

// DeviceConfig holds the device configuration from config file
type DeviceConfig struct {
	Connection ConnectionConfig `yaml:"connection" mapstructure:"connection"`
	Pool       PoolConfig       `yaml:"pool" mapstructure:"pool"`
	Log        LogConfig        `yaml:"log" mapstructure:"log"`
}

// PoolConfig holds the connection pool configuration
type PoolConfig struct {
	MaxConnections int  `yaml:"max_connections" mapstructure:"max_connections"`
	MinConnections int  `yaml:"min_connections" mapstructure:"min_connections"`
	IdleTimeout    int  `yaml:"idle_timeout" mapstructure:"idle_timeout"`    // seconds
	CommandTimeout int  `yaml:"command_timeout" mapstructure:"command_timeout"` // seconds
	QueueTimeout   int  `yaml:"queue_timeout" mapstructure:"queue_timeout"`   // seconds
	MaxQueueSize   int  `yaml:"max_queue_size" mapstructure:"max_queue_size"`
}

// LogConfig holds the logging configuration
type LogConfig struct {
	Enabled       bool   `yaml:"enabled" mapstructure:"enabled"`
	File          string `yaml:"file" mapstructure:"file"`
	MaxSize       int    `yaml:"max_size" mapstructure:"max_size"`       // MB
	MaxBackups    int    `yaml:"max_backups" mapstructure:"max_backups"`
	MaxAge        int    `yaml:"max_age" mapstructure:"max_age"`         // days
	Compress      bool   `yaml:"compress" mapstructure:"compress"`
	IncludeOutput bool   `yaml:"include_output" mapstructure:"include_output"`
	MaxOutputSize int    `yaml:"max_output_size" mapstructure:"max_output_size"` // bytes
}

// CommandResult represents the result of a command execution
type CommandResult struct {
	Command   string        `json:"command"`
	Output    string        `json:"output,omitempty"`
	Error     string        `json:"error,omitempty"`
	Duration  int64         `json:"duration_ms"`
	Success   bool          `json:"success"`
	Timestamp int64         `json:"timestamp"`
}

// ProtocolAdapter defines the interface for protocol adapters
type ProtocolAdapter interface {
	// Connect establishes a connection to the device
	Connect(ctx context.Context, config *ConnectionConfig) error

	// Disconnect closes the connection
	Disconnect(ctx context.Context) error

	// ExecuteCommand executes a single command
	ExecuteCommand(ctx context.Context, cmd string) (*CommandResult, error)

	// IsConnected returns whether the connection is active
	IsConnected() bool

	// ProtocolType returns the protocol type
	ProtocolType() ProtocolType
}
