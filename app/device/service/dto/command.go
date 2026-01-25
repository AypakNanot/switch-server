package dto

// CommandExecuteReq is the request for executing a single command
type CommandExecuteReq struct {
	Command string `json:"command" binding:"required"`
	Timeout int    `json:"timeout"` // seconds, default from config
}

// CommandExecuteResp is the response for executing a command
type CommandExecuteResp struct {
	Command  string `json:"command"`
	Output   string `json:"output,omitempty"`
	Success  bool   `json:"success"`
	Duration int64  `json:"duration_ms"`
	Error    string `json:"error,omitempty"`
}

// BatchCommandReq is the request for executing multiple commands
type BatchCommandReq struct {
	Commands []string `json:"commands" binding:"required,min=1,max=50"`
	Timeout  int      `json:"timeout"` // seconds, default from config
}

// BatchCommandResp is the response for executing multiple commands
type BatchCommandResp struct {
	Results []CommandExecuteResp `json:"results"`
	Total   int                  `json:"total"`
	Success int                  `json:"success"`
	Failed  int                  `json:"failed"`
}

// CommandHistoryReq is the request for querying command history
type CommandHistoryReq struct {
	Limit  int `form:"limit" binding:"required,min=1,max=1000"`
	Offset int `form:"offset" binding:"required,min=0"`
}

// CommandHistoryResp is the response for command history
type CommandHistoryResp struct {
	History []CommandHistoryItem `json:"history"`
	Total   int                  `json:"total"`
	Limit   int                  `json:"limit"`
	Offset  int                  `json:"offset"`
}

// CommandHistoryItem represents a single history item
type CommandHistoryItem struct {
	Timestamp int64  `json:"timestamp"`
	UserID    string `json:"user_id,omitempty"`
	Username  string `json:"username,omitempty"`
	Command   string `json:"command"`
	Success   bool   `json:"success"`
	Duration  int64  `json:"duration_ms"`
	ClientIP  string `json:"client_ip,omitempty"`
}

// DeviceStatusResp is the response for device status
type DeviceStatusResp struct {
	Connected         bool `json:"connected"`
	TotalConnections  int  `json:"total_connections"`
	ActiveConnections int  `json:"active_connections"`
	QueueSize         int  `json:"queue_size"`
	MaxConnections    int  `json:"max_connections"`
	MaxQueueSize      int  `json:"max_queue_size"`
}
