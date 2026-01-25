package device

import (
	"context"
	"fmt"
	"sync"

	"go.uber.org/zap"

	"opt-switch/config"
)

var (
	globalPool     *ConnectionPool
	globalLogger   *ExecutionLogger
	globalConfig   *DeviceConfig
	configManager  *ConfigManager
	once           sync.Once
)

// Initialize initializes the device interaction layer
func Initialize(logger *zap.Logger) error {
	var initErr error
	once.Do(func() {
		// Load device configuration from config.ExtConfig
		extConfig := config.ExtConfig

		// Convert config.ExtConfig.Device to device.DeviceConfig
		cfg := &DeviceConfig{
			Connection: ConnectionConfig{
				Protocol: extConfig.Device.Connection.Protocol,
				Host:     extConfig.Device.Connection.Host,
				Port:     extConfig.Device.Connection.Port,
				Username: extConfig.Device.Connection.Username,
				Password: extConfig.Device.Connection.Password,
				Timeout:  extConfig.Device.Connection.Timeout,
			},
			Pool: PoolConfig{
				MaxConnections: extConfig.Device.Pool.MaxConnections,
				MinConnections: extConfig.Device.Pool.MinConnections,
				IdleTimeout:    extConfig.Device.Pool.IdleTimeout,
				CommandTimeout: extConfig.Device.Pool.CommandTimeout,
				QueueTimeout:   extConfig.Device.Pool.QueueTimeout,
				MaxQueueSize:   extConfig.Device.Pool.MaxQueueSize,
			},
			Log: LogConfig{
				Enabled:       extConfig.Device.Log.Enabled,
				File:          extConfig.Device.Log.File,
				MaxSize:       extConfig.Device.Log.MaxSize,
				MaxBackups:    extConfig.Device.Log.MaxBackups,
				MaxAge:        extConfig.Device.Log.MaxAge,
				Compress:      extConfig.Device.Log.Compress,
				IncludeOutput: extConfig.Device.Log.IncludeOutput,
				MaxOutputSize: extConfig.Device.Log.MaxOutputSize,
			},
		}

		globalConfig = cfg
		configManager = NewConfigManager(cfg)

		// Validate and load config
		if _, err := configManager.LoadConfig(); err != nil {
			initErr = fmt.Errorf("failed to validate device config: %w", err)
			return
		}

		// Create execution logger
		execLogger, err := NewExecutionLogger(&cfg.Log)
		if err != nil {
			initErr = fmt.Errorf("failed to create execution logger: %w", err)
			return
		}
		globalLogger = execLogger

		// Create connection pool
		pool, err := NewConnectionPool(cfg)
		if err != nil {
			initErr = fmt.Errorf("failed to create connection pool: %w", err)
			return
		}
		globalPool = pool

		// Start the pool
		ctx := context.Background()
		if err := pool.Start(ctx); err != nil {
			initErr = fmt.Errorf("failed to start connection pool: %w", err)
			return
		}

		logger.Info("Device interaction layer initialized",
			zap.String("protocol", cfg.Connection.Protocol),
			zap.String("host", cfg.Connection.Host),
			zap.Int("max_connections", cfg.Pool.MaxConnections),
		)
	})

	return initErr
}

// GetPool returns the global connection pool
func GetPool() *ConnectionPool {
	return globalPool
}

// GetLogger returns the global execution logger
func GetLogger() *ExecutionLogger {
	return globalLogger
}

// GetConfig returns the global device configuration
func GetConfig() *DeviceConfig {
	return globalConfig
}

// GetConfigManager returns the global config manager
func GetConfigManager() *ConfigManager {
	return configManager
}

// Shutdown shuts down the device interaction layer
func Shutdown(logger *zap.Logger) error {
	if globalPool != nil {
		if err := globalPool.Stop(); err != nil {
			logger.Error("Failed to stop connection pool", zap.Error(err))
			return err
		}
	}

	if globalLogger != nil {
		if err := globalLogger.Close(); err != nil {
			logger.Error("Failed to close execution logger", zap.Error(err))
			return err
		}
	}

	logger.Info("Device interaction layer shut down")
	return nil
}

// IsInitialized returns whether the device layer is initialized
func IsInitialized() bool {
	return globalPool != nil
}
