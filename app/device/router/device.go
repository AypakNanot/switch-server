package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"opt-switch/app/device/apis"
	"opt-switch/pkg/device"
)

var (
	logger *zap.Logger
)

// InitDeviceService initializes the device service and API
func InitDeviceService(log *zap.Logger) error {
	logger = log

	// Initialize device layer
	if err := device.Initialize(logger); err != nil {
		logger.Warn("Failed to initialize device layer", zap.Error(err))
		// Don't fail if device layer doesn't initialize
		// (e.g., device not configured)
		return nil
	}

	logger.Info("Device service initialized")
	return nil
}

// InitDeviceRouter initializes device routes
func InitDeviceRouter(router *gin.RouterGroup) {
	if !device.IsInitialized() {
		// Device layer not initialized, skip routes
		return
	}

	commandAPI := &apis.CommandAPI{}

	deviceRouter := router.Group("/device")
	{
		// Device info endpoint (no auth required for health check)
		deviceRouter.GET("", commandAPI.GetDeviceInfo)

		// Command execution routes (require authentication)
		commandGroup := deviceRouter.Group("/command")
		{
			commandGroup.POST("/execute", commandAPI.ExecuteCommand)
			commandGroup.POST("/batch", commandAPI.ExecuteBatch)
			commandGroup.GET("/history", commandAPI.GetHistory)
		}

		// Device status route (require authentication)
		deviceRouter.GET("/status", commandAPI.GetStatus)
	}
}

// ShutdownDeviceService shuts down the device service
func ShutdownDeviceService() error {
	if logger != nil {
		return device.Shutdown(logger)
	}
	return nil
}
