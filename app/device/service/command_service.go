package service

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"

	"opt-switch/app/device/service/dto"
	"opt-switch/pkg/device"

	"github.com/go-admin-team/go-admin-core/sdk/service"
)

// CommandService handles command execution business logic
type CommandService struct {
	service.Service
}

// NewCommandService creates a new command service
func NewCommandService() *CommandService {
	return &CommandService{}
}

// ExecuteCommand executes a single command
func (s *CommandService) ExecuteCommand(c *gin.Context, req *dto.CommandExecuteReq) (*dto.CommandExecuteResp, error) {
	// Get timeout from request or config
	timeout := time.Duration(device.GetConfig().Pool.CommandTimeout) * time.Second
	if req.Timeout > 0 {
		timeout = time.Duration(req.Timeout) * time.Second
	}

	// Extract user info for logging
	userID, username, clientIP := s.extractUserInfo(c)

	// Execute command
	ctx := context.Background()
	results, err := device.GetPool().Execute(ctx, []string{req.Command}, timeout)
	if err != nil {
		s.Log.Errorf("Failed to execute command: %v", err)
		return nil, err
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no result returned")
	}

	result := results[0]

	// Log execution asynchronously
	go func() {
		if device.GetLogger() != nil {
			_ = device.GetLogger().LogFromResult(result, userID, username, clientIP)
		}
	}()

	// Map to response
	resp := &dto.CommandExecuteResp{
		Command:  result.Command,
		Output:   result.Output,
		Success:  result.Success,
		Duration: result.Duration,
		Error:    result.Error,
	}

	return resp, nil
}

// ExecuteBatch executes multiple commands
func (s *CommandService) ExecuteBatch(c *gin.Context, req *dto.BatchCommandReq) (*dto.BatchCommandResp, error) {
	// Get timeout from request or config
	timeout := time.Duration(device.GetConfig().Pool.CommandTimeout) * time.Second
	if req.Timeout > 0 {
		timeout = time.Duration(req.Timeout) * time.Second
	}

	// Extract user info for logging
	userID, username, clientIP := s.extractUserInfo(c)

	// Execute commands
	ctx := context.Background()
	results, err := device.GetPool().Execute(ctx, req.Commands, timeout)
	if err != nil {
		s.Log.Errorf("Failed to execute batch commands: %v", err)
		return nil, err
	}

	// Log each execution asynchronously
	go func() {
		if device.GetLogger() != nil {
			for _, result := range results {
				_ = device.GetLogger().LogFromResult(result, userID, username, clientIP)
			}
		}
	}()

	// Map results
	respResults := make([]dto.CommandExecuteResp, len(results))
	successCount := 0
	failedCount := 0

	for i, result := range results {
		respResults[i] = dto.CommandExecuteResp{
			Command:  result.Command,
			Output:   result.Output,
			Success:  result.Success,
			Duration: result.Duration,
			Error:    result.Error,
		}
		if result.Success {
			successCount++
		} else {
			failedCount++
		}
	}

	return &dto.BatchCommandResp{
		Results: respResults,
		Total:   len(respResults),
		Success: successCount,
		Failed:  failedCount,
	}, nil
}

// GetHistory retrieves command execution history
func (s *CommandService) GetHistory(req *dto.CommandHistoryReq) (*dto.CommandHistoryResp, error) {
	logger := device.GetLogger()
	if logger == nil {
		return &dto.CommandHistoryResp{
			History: []dto.CommandHistoryItem{},
			Total:   0,
			Limit:   req.Limit,
			Offset:  req.Offset,
		}, nil
	}

	logs, err := logger.GetHistory(req.Limit, req.Offset)
	if err != nil {
		s.Log.Errorf("Failed to get command history: %v", err)
		return nil, err
	}

	// Map logs to response
	history := make([]dto.CommandHistoryItem, len(logs))
	for i, log := range logs {
		history[i] = dto.CommandHistoryItem{
			Timestamp: log.Timestamp,
			UserID:    log.UserID,
			Username:  log.Username,
			Command:   log.Command,
			Success:   log.Success,
			Duration:  log.Duration,
			ClientIP:  log.ClientIP,
		}
	}

	return &dto.CommandHistoryResp{
		History: history,
		Total:   len(history),
		Limit:   req.Limit,
		Offset:  req.Offset,
	}, nil
}

// GetStatus returns the device connection status
func (s *CommandService) GetStatus() *dto.DeviceStatusResp {
	pool := device.GetPool()
	if pool == nil {
		return &dto.DeviceStatusResp{
			Connected: false,
		}
	}

	status := pool.GetStatus()

	connected := false
	if running, ok := status["running"].(bool); ok && running {
		connected = true
	}

	totalConns := 0
	if v, ok := status["total_connections"].(int); ok {
		totalConns = v
	}

	activeConns := 0
	if v, ok := status["active_connections"].(int); ok {
		activeConns = v
	}

	queueSize := 0
	if v, ok := status["queue_size"].(int); ok {
		queueSize = v
	}

	maxConns := 0
	if v, ok := status["max_connections"].(int); ok {
		maxConns = v
	}

	maxQueue := 0
	if v, ok := status["max_queue_size"].(int); ok {
		maxQueue = v
	}

	return &dto.DeviceStatusResp{
		Connected:         connected,
		TotalConnections:  totalConns,
		ActiveConnections: activeConns,
		QueueSize:         queueSize,
		MaxConnections:    maxConns,
		MaxQueueSize:      maxQueue,
	}
}

// extractUserInfo extracts user information from gin context
func (s *CommandService) extractUserInfo(c *gin.Context) (userID, username, clientIP string) {
	// Get user ID from context (set by auth middleware)
	if userIDVal, exists := c.Get("user_id"); exists {
		if id, ok := userIDVal.(string); ok {
			userID = id
		}
	}

	// Get username from context
	if usernameVal, exists := c.Get("username"); exists {
		if name, ok := usernameVal.(string); ok {
			username = name
		}
	}

	// Get client IP
	clientIP = c.ClientIP()

	return
}

// MapError maps device errors to response messages
func (s *CommandService) MapError(err error) (int, string) {
	if deviceErr, ok := err.(*device.DeviceError); ok {
		switch deviceErr.Code {
		case device.ErrConnectionFailed, device.ErrAuthFailed, device.ErrConnectionClosed:
			return 503, "Device connection failed: " + err.Error()
		case device.ErrQueueFull, device.ErrQueueTimeout:
			return 429, "Service busy, please try again later"
		case device.ErrCommandTimeout:
			return 504, "Command execution timeout"
		case device.ErrCommandFailed:
			return 500, "Command execution failed"
		case device.ErrInvalidConfig, device.ErrDeviceNotConfigured:
			return 500, "Device configuration error"
		}
	}

	return 500, "Internal server error"
}
