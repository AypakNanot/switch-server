package device

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// ExecutionLog represents a command execution log entry
type ExecutionLog struct {
	Timestamp   int64  `json:"timestamp"`
	UserID      string `json:"user_id,omitempty"`
	Username    string `json:"username,omitempty"`
	Command     string `json:"command"`
	Output      string `json:"output,omitempty"`
	OutputSize  int    `json:"output_size"`
	Success     bool   `json:"success"`
	Error       string `json:"error,omitempty"`
	Duration    int64  `json:"duration_ms"`
	ClientIP    string `json:"client_ip,omitempty"`
}

// ExecutionLogger handles command execution logging
type ExecutionLogger struct {
	logger *zap.Logger
	config *LogConfig
	mu     sync.RWMutex
}

// NewExecutionLogger creates a new execution logger
func NewExecutionLogger(config *LogConfig) (*ExecutionLogger, error) {
	if !config.Enabled {
		return &ExecutionLogger{
			config: config,
		}, nil
	}

	// Ensure log directory exists
	logDir := filepath.Dir(config.File)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	// Set up lumberjack for log rotation
	lumberjackLogger := &lumberjack.Logger{
		Filename:   config.File,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   config.Compress,
	}

	// Create encoder config
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Create core
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(lumberjackLogger),
		zapcore.InfoLevel,
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return &ExecutionLogger{
		logger: logger,
		config: config,
	}, nil
}

// Log logs a command execution
func (l *ExecutionLogger) Log(log *ExecutionLog) error {
	if !l.config.Enabled {
		return nil
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	// Truncate output if needed
	output := log.Output
	if l.config.MaxOutputSize > 0 && len(output) > l.config.MaxOutputSize {
		output = output[:l.config.MaxOutputSize] + "... (truncated)"
	}

	// Create log entry
	logEntry := map[string]interface{}{
		"timestamp":   time.Unix(log.Timestamp, 0).Format(time.RFC3339),
		"user_id":     log.UserID,
		"username":    log.Username,
		"command":     log.Command,
		"output_size": log.OutputSize,
		"success":     log.Success,
		"duration_ms": log.Duration,
		"client_ip":   log.ClientIP,
	}

	if l.config.IncludeOutput && len(output) > 0 {
		logEntry["output"] = output
	}

	if log.Error != "" {
		logEntry["error"] = log.Error
	}

	// Log as JSON
	if l.logger != nil {
		l.logger.Info("command_execution", zap.Any("data", logEntry))
	}

	return nil
}

// LogFromResult logs a command execution from CommandResult
func (l *ExecutionLogger) LogFromResult(result *CommandResult, userID, username, clientIP string) error {
	return l.Log(&ExecutionLog{
		Timestamp:  result.Timestamp,
		UserID:     userID,
		Username:   username,
		Command:    result.Command,
		Output:     result.Output,
		OutputSize: len(result.Output),
		Success:    result.Success,
		Error:      result.Error,
		Duration:   result.Duration,
		ClientIP:   clientIP,
	})
}

// GetHistory retrieves execution history from log file
func (l *ExecutionLogger) GetHistory(limit, offset int) ([]ExecutionLog, error) {
	if !l.config.Enabled {
		return []ExecutionLog{}, nil
	}

	l.mu.RLock()
	defer l.mu.RUnlock()

	// Open log file
	file, err := os.Open(l.config.File)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}
	defer file.Close()

	// Read file line by line from the end
	var logs []ExecutionLog
	scanner := newReverseScanner(file)

	count := 0
	skipped := 0

	for scanner.Scan() {
		line := scanner.Text()

		// Skip empty lines
		if line == "" {
			continue
		}

		// Parse JSON log entry
		var logEntry struct {
			Data ExecutionLog `json:"data"`
		}

		if err := json.Unmarshal([]byte(line), &logEntry); err != nil {
			continue // Skip invalid JSON
		}

		// Skip entries before offset
		if skipped < offset {
			skipped++
			continue
		}

		logs = append(logs, logEntry.Data)
		count++

		if count >= limit {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read log file: %w", err)
	}

	return logs, nil
}

// Close closes the logger
func (l *ExecutionLogger) Close() error {
	if l.logger != nil {
		return l.logger.Sync()
	}
	return nil
}

// reverseScanner reads a file line by line from the end
type reverseScanner struct {
	file   *os.File
	pos    int64
	buffer []byte
	line   string
	err    error
}

func newReverseScanner(file *os.File) *reverseScanner {
	// Get file size
	stat, _ := file.Stat()
	size := stat.Size()

	return &reverseScanner{
		file:   file,
		pos:    size,
		buffer: make([]byte, 4096),
	}
}

func (s *reverseScanner) Scan() bool {
	if s.pos <= 0 {
		return false
	}

	var lines []string
	currentLine := ""

	// Read file backwards in chunks
	for s.pos > 0 {
		// Determine chunk size
		chunkSize := int64(len(s.buffer))
		if s.pos < chunkSize {
			chunkSize = s.pos
		}

		// Seek to chunk position
		s.pos -= chunkSize
		if _, err := s.file.Seek(s.pos, 0); err != nil {
			s.err = err
			return false
		}

		// Read chunk
		n, err := s.file.Read(s.buffer[:chunkSize])
		if err != nil {
			s.err = err
			return false
		}

		// Process chunk backwards
		for i := n - 1; i >= 0; i-- {
			if s.buffer[i] == '\n' {
				lines = append(lines, currentLine)
				currentLine = ""
			} else {
				currentLine = string(s.buffer[i]) + currentLine
			}
		}

		// If we have lines, return the first one
		if len(lines) > 0 {
			s.line = lines[len(lines)-1]
			return true
		}
	}

	// Return the last line if exists
	if currentLine != "" {
		s.line = currentLine
		return true
	}

	return false
}

func (s *reverseScanner) Text() string {
	return s.line
}

func (s *reverseScanner) Err() error {
	return s.err
}
