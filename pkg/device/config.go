package device

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

// ConfigManager manages device configuration
type ConfigManager struct {
	config      *DeviceConfig
	encryptionKey []byte
}

// NewConfigManager creates a new config manager
func NewConfigManager(config *DeviceConfig) *ConfigManager {
	// Load encryption key from environment
	key := []byte(os.Getenv("DEVICE_ENCRYPTION_KEY"))

	return &ConfigManager{
		config:      config,
		encryptionKey: key,
	}
}

// LoadConfig loads and validates the device configuration
func (m *ConfigManager) LoadConfig() (*DeviceConfig, error) {
	if m.config == nil {
		return nil, NewInvalidConfigError("device config is nil")
	}

	// Validate required fields
	if err := m.validateConfig(m.config); err != nil {
		return nil, err
	}

	// Decrypt password if needed
	if strings.HasPrefix(m.config.Connection.Password, "encrypted:") {
		decrypted, err := m.decryptPassword(m.config.Connection.Password)
		if err != nil {
			return nil, NewInvalidConfigError(fmt.Sprintf("failed to decrypt password: %v", err))
		}
		m.config.Connection.Password = decrypted
	}

	return m.config, nil
}

// validateConfig validates the device configuration
func (m *ConfigManager) validateConfig(config *DeviceConfig) error {
	// Validate connection config
	if config.Connection.Host == "" {
		return NewInvalidConfigError("connection.host is required")
	}
	if config.Connection.Port <= 0 {
		return NewInvalidConfigError("connection.port must be positive")
	}
	if config.Connection.Protocol == "" {
		return NewInvalidConfigError("connection.protocol is required")
	}
	if config.Connection.Username == "" {
		return NewInvalidConfigError("connection.username is required")
	}
	if config.Connection.Password == "" {
		return NewInvalidConfigError("connection.password is required")
	}
	if config.Connection.Timeout <= 0 {
		config.Connection.Timeout = 30 // Default 30 seconds
	}

	// Validate pool config
	if config.Pool.MaxConnections <= 0 {
		config.Pool.MaxConnections = 3 // Default
	}
	if config.Pool.MinConnections < 0 {
		config.Pool.MinConnections = 1 // Default
	}
	if config.Pool.MinConnections > config.Pool.MaxConnections {
		config.Pool.MinConnections = config.Pool.MaxConnections
	}
	if config.Pool.IdleTimeout <= 0 {
		config.Pool.IdleTimeout = 300 // Default 5 minutes
	}
	if config.Pool.CommandTimeout <= 0 {
		config.Pool.CommandTimeout = 30 // Default 30 seconds
	}
	if config.Pool.QueueTimeout <= 0 {
		config.Pool.QueueTimeout = 60 // Default 60 seconds
	}
	if config.Pool.MaxQueueSize <= 0 {
		config.Pool.MaxQueueSize = 100 // Default
	}

	// Validate log config
	if config.Log.File == "" {
		config.Log.File = "logs/command.log" // Default
	}
	if config.Log.MaxSize <= 0 {
		config.Log.MaxSize = 100 // Default 100MB
	}
	if config.Log.MaxBackups < 0 {
		config.Log.MaxBackups = 3 // Default
	}
	if config.Log.MaxAge <= 0 {
		config.Log.MaxAge = 7 // Default 7 days
	}
	if config.Log.MaxOutputSize < 0 {
		config.Log.MaxOutputSize = 10240 // Default 10KB
	}

	return nil
}

// decryptPassword decrypts an encrypted password
func (m *ConfigManager) decryptPassword(encrypted string) (string, error) {
	if len(m.encryptionKey) == 0 {
		return "", fmt.Errorf("encryption key not set (DEVICE_ENCRYPTION_KEY environment variable)")
	}

	// Remove "encrypted:" prefix
	encrypted = strings.TrimPrefix(encrypted, "encrypted:")

	// Decode base64
	ciphertext, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	// Create cipher
	block, err := aes.NewCipher(m.encryptionKey)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	// Check ciphertext length
	if len(ciphertext) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	// Extract IV and actual ciphertext
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	// Decrypt
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}

// EncryptPassword encrypts a password (utility function for generating encrypted passwords)
func EncryptPassword(password string, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create IV
	iv := make([]byte, aes.BlockSize)
	// In production, use crypto/rand for IV generation
	// For simplicity here, we use a fixed IV (NOT SECURE for production)

	// Encrypt
	ciphertext := make([]byte, len(password))
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext, []byte(password))

	// Combine IV and ciphertext
	result := append(iv, ciphertext...)

	// Encode to base64 with prefix
	return "encrypted:" + base64.StdEncoding.EncodeToString(result), nil
}

// GetConfig returns the current configuration
func (m *ConfigManager) GetConfig() *DeviceConfig {
	return m.config
}

// UpdateConfig updates the configuration
func (m *ConfigManager) UpdateConfig(newConfig *DeviceConfig) error {
	if err := m.validateConfig(newConfig); err != nil {
		return err
	}

	m.config = newConfig
	return nil
}

// ValidateSSHConnection validates the SSH connection configuration
func (m *ConfigManager) ValidateSSHConnection() error {
	config := m.config.Connection

	sshConfig := &ssh.ClientConfig{
		User: config.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(config.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:        5 * time.Second,
	}

	address := fmt.Sprintf("%s:%d", config.Host, config.Port)
	client, err := ssh.Dial("tcp", address, sshConfig)
	if err != nil {
		return NewConnectionError(err)
	}
	defer client.Close()

	return nil
}
