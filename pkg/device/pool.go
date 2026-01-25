package device

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// Connection represents a device connection
type Connection struct {
	ID        string
	Adapter   ProtocolAdapter
	CreatedAt time.Time
	LastUsed  time.Time
	InUse     int32 // atomic
}

// CommandTask represents a command execution task
type CommandTask struct {
	Commands []string
	Timeout  time.Duration
	ResultCh chan *CommandResult
	UserID   string
}

// ConnectionPool manages device connections using semaphore pattern
type ConnectionPool struct {
	config    *DeviceConfig
	adapter   ProtocolAdapter
	semaphore chan struct{} // Semaphore for connection limiting
	queue     chan *CommandTask
	connections map[string]*Connection
	mu        sync.RWMutex
	running   int32
	wg        sync.WaitGroup
}

// NewConnectionPool creates a new connection pool
func NewConnectionPool(config *DeviceConfig) (*ConnectionPool, error) {
	if config == nil {
		return nil, NewInvalidConfigError("config is nil")
	}

	// Create protocol adapter based on type
	var adapter ProtocolAdapter
	switch ProtocolType(config.Connection.Protocol) {
	case ProtocolSSH:
		adapter = NewSSHAdapterFunc()
	case ProtocolTelnet:
		adapter = NewTelnetAdapterFunc()
	default:
		return nil, NewInvalidConfigError(fmt.Sprintf("unsupported protocol: %s", config.Connection.Protocol))
	}

	pool := &ConnectionPool{
		config:      config,
		adapter:     adapter,
		semaphore:   make(chan struct{}, config.Pool.MaxConnections),
		queue:       make(chan *CommandTask, config.Pool.MaxQueueSize),
		connections: make(map[string]*Connection),
	}

	// Initialize semaphore with available slots
	for i := 0; i < config.Pool.MinConnections; i++ {
		pool.semaphore <- struct{}{}
	}

	return pool, nil
}

// Start starts the connection pool workers
func (p *ConnectionPool) Start(ctx context.Context) error {
	if !atomic.CompareAndSwapInt32(&p.running, 0, 1) {
		return nil // Already running
	}

	// Start worker goroutines
	for i := 0; i < p.config.Pool.MaxConnections; i++ {
		p.wg.Add(1)
		go p.worker(ctx, i)
	}

	// Establish initial connections
	for i := 0; i < p.config.Pool.MinConnections; i++ {
		if _, err := p.getConnection(ctx); err != nil {
			// Log error but continue
			fmt.Printf("Failed to establish initial connection: %v\n", err)
		}
	}

	return nil
}

// Stop stops the connection pool
func (p *ConnectionPool) Stop() error {
	if !atomic.CompareAndSwapInt32(&p.running, 1, 0) {
		return nil // Already stopped
	}

	// Close queue
	close(p.queue)

	// Wait for workers to finish
	p.wg.Wait()

	// Close all connections
	p.mu.Lock()
	for _, conn := range p.connections {
		_ = conn.Adapter.Disconnect(context.Background())
	}
	p.connections = make(map[string]*Connection)
	p.mu.Unlock()

	return nil
}

// Execute submits a command for execution
func (p *ConnectionPool) Execute(ctx context.Context, commands []string, timeout time.Duration) ([]*CommandResult, error) {
	if !p.IsRunning() {
		return nil, NewConnectionClosed()
	}

	// Create result channel
	resultCh := make(chan *CommandResult, len(commands))

	// Create task
	task := &CommandTask{
		Commands: commands,
		Timeout:  timeout,
		ResultCh: resultCh,
	}

	// Try to submit to queue
	select {
	case p.queue <- task:
		// Task submitted successfully
	case <-time.After(time.Duration(p.config.Pool.QueueTimeout) * time.Second):
		return nil, NewQueueTimeoutError()
	}

	// Collect results
	results := make([]*CommandResult, 0, len(commands))
	for i := 0; i < len(commands); i++ {
		select {
		case result := <-resultCh:
			results = append(results, result)
		case <-time.After(timeout + time.Duration(p.config.Pool.CommandTimeout)*time.Second):
			return results, NewCommandTimeoutError()
		}
	}

	return results, nil
}

// worker processes commands from the queue
func (p *ConnectionPool) worker(ctx context.Context, workerID int) {
	defer p.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case task, ok := <-p.queue:
			if !ok {
				// Queue closed
				return
			}
			if task == nil {
				continue
			}

			// Acquire semaphore (wait for available connection slot)
			p.semaphore <- struct{}{}
			func() {
				defer func() { <-p.semaphore }() // Release semaphore

				// Execute commands
				for _, cmd := range task.Commands {
					result, err := p.executeCommand(ctx, cmd, task.Timeout)
					if err != nil {
						result = &CommandResult{
							Command:   cmd,
							Error:     err.Error(),
							Success:   false,
							Timestamp: time.Now().Unix(),
						}
					}
					task.ResultCh <- result
				}
			}()
		}
	}
}

// executeCommand executes a single command
func (p *ConnectionPool) executeCommand(ctx context.Context, cmd string, timeout time.Duration) (*CommandResult, error) {
	// Get or create connection
	conn, err := p.getConnection(ctx)
	if err != nil {
		return nil, err
	}

	// Execute command
	execCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	result, err := conn.Adapter.ExecuteCommand(execCtx, cmd)
	if err != nil {
		// Mark connection as potentially stale
		atomic.StoreInt32(&conn.InUse, 0)
		return nil, err
	}

	// Update last used time
	conn.LastUsed = time.Now()
	atomic.StoreInt32(&conn.InUse, 0)

	return result, nil
}

// getConnection gets or creates a connection
func (p *ConnectionPool) getConnection(ctx context.Context) (*Connection, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Try to find an available connection
	for _, conn := range p.connections {
		if atomic.CompareAndSwapInt32(&conn.InUse, 0, 1) {
			// Verify connection is still alive
			if conn.Adapter.IsConnected() {
				return conn, nil
			}
			// Connection is dead, recreate it
			p.recreateConnection(ctx, conn)
			return conn, nil
		}
	}

	// No available connection, create a new one if under limit
	if len(p.connections) < p.config.Pool.MaxConnections {
		conn := &Connection{
			ID:        fmt.Sprintf("conn-%d", time.Now().UnixNano()),
			Adapter:   p.adapter,
			CreatedAt: time.Now(),
			LastUsed:  time.Now(),
		}

		// Establish connection
		connCtx, cancel := context.WithTimeout(ctx, time.Duration(p.config.Connection.Timeout)*time.Second)
		defer cancel()

		if err := conn.Adapter.Connect(connCtx, &p.config.Connection); err != nil {
			return nil, NewConnectionError(err)
		}

		atomic.StoreInt32(&conn.InUse, 1)
		p.connections[conn.ID] = conn
		return conn, nil
	}

	// Wait for an available connection
	return nil, NewQueueFullError()
}

// recreateConnection recreates a stale connection
func (p *ConnectionPool) recreateConnection(ctx context.Context, conn *Connection) {
	conn.Adapter.Disconnect(context.Background())

	connCtx, cancel := context.WithTimeout(ctx, time.Duration(p.config.Connection.Timeout)*time.Second)
	defer cancel()

	if err := conn.Adapter.Connect(connCtx, &p.config.Connection); err != nil {
		fmt.Printf("Failed to recreate connection %s: %v\n", conn.ID, err)
	}
}

// IsRunning returns whether the pool is running
func (p *ConnectionPool) IsRunning() bool {
	return atomic.LoadInt32(&p.running) == 1
}

// GetStatus returns the pool status
func (p *ConnectionPool) GetStatus() map[string]interface{} {
	p.mu.RLock()
	defer p.mu.RUnlock()

	activeConnections := 0
	for _, conn := range p.connections {
		if atomic.LoadInt32(&conn.InUse) == 1 {
			activeConnections++
		}
	}

	return map[string]interface{}{
		"running":            p.IsRunning(),
		"total_connections":  len(p.connections),
		"active_connections": activeConnections,
		"queue_size":         len(p.queue),
		"max_connections":    p.config.Pool.MaxConnections,
		"max_queue_size":     p.config.Pool.MaxQueueSize,
	}
}

// ReloadConfig reloads the configuration
func (p *ConnectionPool) ReloadConfig(newConfig *DeviceConfig) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Update config
	p.config = newConfig

	// Recreate adapter if protocol changed
	// (For simplicity, we just update the config, actual reconnection happens on next use)

	return nil
}
