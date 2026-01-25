# Design: Extensible Device Interaction Layer (Simplified)

## Context

系统运行在交换机设备上，需要管理自身这台交换机。主要场景：
1. 配置管理：VLAN、端口、路由等配置
2. 状态监控：端口状态、流量、系统状态等
3. 多用户访问：多个管理员同时通过 Web UI 操作设备

### 核心约束
- **只管理自身设备**：无需多设备管理
- **资源受限**：50-70MB 内存
- **配置文件管理**：设备连接信息从配置文件读取
- **连接池 + 队列**：多用户操作需要排队机制
- **日志文件**：执行历史用日志文件记录，不用数据库

### Stakeholders
- 网络管理员：通过 Web UI 管理交换机
- 系统：自动化运维脚本
- 开发者：未来扩展 NETCONF 等协议

## Goals / Non-Goals

### Goals
1. **可扩展的协议抽象**：支持 SSH、Telnet，未来轻松扩展 NETCONF
2. **连接池 + 命令队列**：多用户操作排队，避免并发冲突
3. **配置文件驱动**：设备连接信息从配置文件读取，支持热加载
4. **轻量级日志**：使用日志文件记录执行历史
5. **统一的错误处理**：清晰的错误码和错误信息
6. **同步执行**：命令执行后立即返回结果

### Non-Goals (当前不实现)
1. 多设备管理（只管理自身）
2. 设备分组功能
3. 数据库存储设备信息
4. 命令模板系统（后续可扩展）
5. 设备配置备份和回滚（后续可扩展）

## Decisions

### Decision 1: 简化的三层架构

**What**: 采用简化的三层架构：
1. **API 层** (`app/device/apis/`)：HTTP 接口
2. **Service 层** (`app/device/service/`)：业务逻辑
3. **设备交互层** (`pkg/device/`)：协议抽象、连接池、命令队列

**Why**:
- 架构清晰，职责单一
- 只管理自身设备，无需复杂的设备管理层
- 符合项目现有架构模式

**Alternatives considered**:
- 两层架构（API + Device）：耦合度高
- 完整四层架构：对于单设备场景过度设计

### Decision 2: 连接池 + 命令队列设计

**What**: 实现连接池配合命令队列，支持并发控制和排队：

```go
// 连接池：管理固定数量的连接
type ConnectionPool struct {
    connections chan *Connection  // 信号量模式，限制并发数
    device      *DeviceConfig     // 设备配置（从配置文件读取）
    protocol    ProtocolAdapter   // 协议适配器
}

// 命令队列：管理待执行的命令
type CommandQueue struct {
    queue      chan *CommandTask  // 命令任务队列
    workers    int                // 工作协程数 = 连接池大小
    timeout    time.Duration      // 队列等待超时
}

// 命令任务
type CommandTask struct {
    Commands   []string           // 要执行的命令列表
    Result     chan *CommandResult// 结果通道（同步返回）
    Timeout    time.Duration      // 执行超时
    UserID     string             // 用户ID（用于日志）
}
```

**工作流程**：
```
用户请求 → API → Service → 获取连接(信号量) → 执行命令 → 释放连接 → 返回结果
                  ↓
            如果连接池满 → 排队等待 → 超时返回错误
```

**配置示例**：
```yaml
device:
  connection:
    host: "127.0.0.1"  # 本地回连
    port: 22
    protocol: "ssh"
    username: "admin"
    password: "xxx"
  pool:
    max_connections: 3        # 最大并发连接数
    idle_timeout: 300s        # 空闲超时
    command_timeout: 30s      # 命令执行超时
    queue_timeout: 60s        # 队列等待超时
```

**Why**:
- **连接池**：限制并发数，避免设备负载过高
- **信号量模式**：简单高效，Go channel 原生支持
- **命令队列**：多用户请求排队，避免并发冲突
- **可配置并发数**：根据设备性能灵活调整

**Alternatives considered**:
- 全局锁：性能差，无法并发
- 无限制并发：可能导致设备负载过高
- 异步任务队列：增加复杂度，不符合同步执行需求

### Decision 3: 协议适配器模式（保持不变）

**What**: 定义 `ProtocolAdapter` 接口，所有协议实现必须满足该接口：

```go
type ProtocolAdapter interface {
    // 连接设备
    Connect(ctx context.Context, config *ConnectionConfig) error
    // 断开连接
    Disconnect(ctx context.Context) error
    // 执行单个命令
    ExecuteCommand(ctx context.Context, cmd string) (*CommandResult, error)
    // 检查连接状态
    IsConnected() bool
    // 获取协议类型
    ProtocolType() ProtocolType
}

type ConnectionConfig struct {
    Host     string
    Port     int
    Username string
    Password string
    Timeout  time.Duration
}

type CommandResult struct {
    Command   string        // 执行的命令
    Output    string        // 命令输出
    Error     error         // 错误信息
    Duration  time.Duration // 执行耗时
    Timestamp time.Time     // 执行时间
}
```

**Why**:
- 面向接口编程，上层代码不关心具体协议
- 添加 NETCONF 只需实现接口，无需修改上层代码
- 方便单元测试和 mock

### Decision 4: 配置文件驱动

**What**: 设备连接信息从配置文件读取，支持热加载：

```yaml
# config/settings.yml
device:
  # 设备连接配置
  connection:
    host: "127.0.0.1"        # 本地回连（或实际管理IP）
    port: 22
    protocol: "ssh"          # ssh / telnet
    username: "admin"
    password: "encrypted:xxx" # 支持加密（可选）
    timeout: 30s

  # 连接池配置
  pool:
    max_connections: 3       # 最大并发连接数（可配置）
    idle_timeout: 300s       # 空闲超时
    command_timeout: 30s     # 命令执行超时
    queue_timeout: 60s       # 队列等待超时
    max_queue_size: 100      # 最大队列长度

  # 日志配置
  log:
    enabled: true
    file: "logs/command.log"       # 日志文件路径
    max_size: 100                  # 单文件最大 MB
    max_backups: 3                 # 保留历史文件数
    max_age: 7                     # 保留天数
    compress: true                 # 压缩旧文件
```

**支持热加载**：
- 通过信号（SIGHUP）或 API 触发配置重载
- 重新建立连接池（优雅关闭旧连接）

**Why**:
- 简化部署，无需数据库
- 配置修改无需重启（热加载）
- 符合 12-factor app 配置原则

**Alternatives considered**:
- 数据库存储：增加复杂度，不适合单设备场景
- 环境变量：管理不便，不支持复杂配置

### Decision 5: 日志文件记录执行历史

**What**: 使用滚动日志文件记录命令执行历史：

```go
type ExecutionLogger struct {
    logger *zap.Logger    // 使用现有 zap logger
    config *LogConfig
}

type LogConfig struct {
    Enabled     bool
    File        string
    MaxSize     int    // MB
    MaxBackups  int
    MaxAge      int    // days
    Compress    bool
}

// 日志格式（JSON，方便查询）
type ExecutionLog struct {
    Timestamp   time.Time `json:"timestamp"`
    UserID      string    `json:"user_id"`
    Username    string    `json:"username"`
    Command     string    `json:"command"`
    Output      string    `json:"output,omitempty"`       // 可选，避免日志过大
    OutputSize  int       `json:"output_size"`
    Success     bool      `json:"success"`
    Error       string    `json:"error,omitempty"`
    Duration    int64     `json:"duration_ms"`
    ClientIP    string    `json:"client_ip"`
}
```

**日志查询 API**：
```go
GET /api/v1/device/command/history?limit=100&offset=0
// 返回最近的执行历史（从日志文件尾部读取）
```

**Why**:
- 不占用数据库资源
- 日志文件天然支持历史归档
- zap lumberjack 支持滚动压缩
- JSON 格式方便后续分析

**Alternatives considered**:
- 数据库存储：增加数据库负担和复杂度
- 内存缓存：重启丢失，无法追溯历史

### Decision 6: 错误处理策略

**What**: 定义统一的错误码体系：

```go
type DeviceError struct {
    Code    ErrorCode
    Message string
    Cause   error
}

type ErrorCode int

const (
    // 连接错误 1000-1099
    ErrConnectionFailed ErrorCode = 1001
    ErrAuthFailed       ErrorCode = 1002
    ErrConnectionClosed ErrorCode = 1003

    // 队列错误 1100-1199
    ErrQueueFull        ErrorCode = 1101
    ErrQueueTimeout     ErrorCode = 1102

    // 执行错误 1200-1299
    ErrCommandFailed    ErrorCode = 1201
    ErrCommandTimeout   ErrorCode = 1202
    ErrOutputTooLarge   ErrorCode = 1203

    // 配置错误 1300-1399
    ErrInvalidConfig    ErrorCode = 1301
    ErrDeviceNotConfigured ErrorCode = 1302
)
```

**Why**:
- 统一错误处理，方便上层判断
- 保留原始错误信息，方便调试
- 支持国际化错误消息

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                         API Layer                               │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │  POST /api/v1/device/command/execute                     │  │
│  │  POST /api/v1/device/command/batch                       │  │
│  │  GET  /api/v1/device/command/history                     │  │
│  │  GET  /api/v1/device/status                              │  │
│  └──────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────────┐
│                       Service Layer                             │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │  CommandService                                           │  │
│  │  - ExecuteCommand()     获取连接 → 执行 → 释放连接       │  │
│  │  - ExecuteBatch()      循环调用 ExecuteCommand           │  │
│  │  - GetHistory()        从日志文件读取历史                │  │
│  └──────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────────┐
│                    Device Interaction Layer                     │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │ConnectionPool│  │CommandQueue  │  │ExecLogger    │          │
│  │- 信号量模式  │  │- 排队机制    │  │- 日志文件    │          │
│  └──────────────┘  └──────────────┘  └──────────────┘          │
└─────────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────────┐
│                      Protocol Adapters                          │
│  ┌──────────────────┐  ┌──────────────────┐                     │
│  │   SSH Adapter    │  │  Telnet Adapter  │                     │
│  │- golang.org/x/   │  │- telnet lib      │                     │
│  │   crypto/ssh     │  │                  │                     │
│  └──────────────────┘  └──────────────────┘                     │
└─────────────────────────────────────────────────────────────────┘
```

## Directory Structure

```
switch-server/
├── pkg/device/                          # 设备交互核心层
│   ├── adapter.go                       # 协议适配器接口定义
│   ├── pool.go                          # 连接池 + 命令队列
│   ├── executor.go                      # 命令执行器
│   ├── logger.go                        # 执行日志记录
│   ├── error.go                         # 错误定义
│   ├── config.go                        # 配置结构定义
│   └── protocol/                        # 协议实现
│       ├── ssh.go                       # SSH 协议适配器
│       ├── telnet.go                    # Telnet 协议适配器
│       └── netconf.go                   # NETCONF 协议适配器（未来）
│
├── app/device/                          # 命令执行业务层
│   ├── apis/                            # API 接口
│   │   └── command.go                   # 命令执行 API
│   ├── service/                         # 业务逻辑
│   │   ├── command_service.go           # 命令执行服务
│   │   └── dto/                         # 数据传输对象
│   │       └── command.go               # 命令 DTO
│   └── router/                          # 路由配置
│       └── device.go                    # 设备路由
│
├── config/                              # 配置文件
│   └── settings.yml                     # 添加 device 配置段
│
└── logs/                                # 日志目录
    └── command.log                      # 命令执行日志
```

## Configuration

### config/settings.yml

```yaml
# 设备连接配置
device:
  # 连接信息
  connection:
    host: "127.0.0.1"        # 本地回连 SSH
    port: 22
    protocol: "ssh"          # ssh / telnet
    username: "admin"
    password: "your_password"
    timeout: 30s             # 连接超时

  # 连接池配置
  pool:
    max_connections: 3       # 最大并发连接数（可配置）
    min_connections: 1       # 最小保持连接数
    idle_timeout: 300s       # 空闲连接超时
    command_timeout: 30s     # 命令执行超时
    queue_timeout: 60s       # 队列等待超时
    max_queue_size: 100      # 最大队列长度

  # 日志配置
  log:
    enabled: true
    file: "logs/command.log"
    max_size: 100           # MB
    max_backups: 3          # 保留 3 个历史文件
    max_age: 7              # 保留 7 天
    compress: true          # 压缩旧文件
    include_output: true    # 是否包含命令输出（可能很大）
    max_output_size: 10240  # 最大输出记录 10KB
```

## API Endpoints

### 命令执行
```
POST /api/v1/device/command/execute
# 执行单个命令
Request: { "command": "show version" }
Response: { "output": "...", "duration": "123ms" }

POST /api/v1/device/command/batch
# 批量执行命令（顺序执行）
Request: { "commands": ["show version", "show interface"] }
Response: [{ "command": "show version", "output": "..." }, ...]

GET /api/v1/device/command/history?limit=50&offset=0
# 查询执行历史（从日志文件）
Response: [{ "timestamp": "...", "command": "...", "user": "..." }, ...]

GET /api/v1/device/status
# 获取连接状态
Response: { "connected": true, "active_connections": 2, "queue_size": 0 }
```

## Security Considerations

1. **密码加密存储**
   - 配置文件中的密码支持加密格式：`encrypted:base64(aes_encrypted_password)`
   - 密钥通过环境变量 `DEVICE_ENCRYPTION_KEY` 传入
   - 启动时解密，明文仅在内存中

2. **访问控制**
   - 集成现有 Casbin RBAC 权限系统
   - 命令执行需要特定权限

3. **命令执行安全**
   - 记录所有命令执行日志（用户、时间、命令）
   - 可选命令白名单（后续）

## Performance Considerations

1. **连接池配置**
   - 默认最大并发：3（可配置）
   - 空闲超时：5 分钟
   - 命令超时：30 秒

2. **内存优化**
   - 命令结果限制大小（最大 1MB）
   - 日志文件滚动，避免无限增长

3. **队列管理**
   - 最大队列长度：100
   - 队列超时：60 秒

## Implementation Flow

### 单命令执行流程

```
1. 用户 POST /api/v1/device/command/execute
   ↓
2. API 层验证请求，获取用户信息
   ↓
3. Service 层调用 ConnectionPool.GetConnection(timeout)
   ↓
4. 如果有空闲连接：
   - 获取连接
   - 执行命令
   - 释放连接
   - 返回结果
   ↓
5. 如果连接池满：
   - 加入队列等待
   - 等待超时 → 返回 ErrQueueTimeout
   - 获得连接 → 执行步骤 4
   ↓
6. 记录执行日志（异步）
   ↓
7. 返回结果给用户
```

## Migration Plan

### Phase 1: 核心层 (1-2 天)
1. 实现 ProtocolAdapter 接口
2. 实现 SSH 和 Telnet 适配器
3. 实现连接池 + 命令队列
4. 实现日志记录器

### Phase 2: 业务层 (1 天)
1. 实现命令执行 Service
2. 实现 API 层
3. 添加配置文件

### Phase 3: 测试和文档 (1 天)
1. 单元测试
2. 集成测试
3. API 文档

**总计：约 3-4 天**

## Risks / Trade-offs

| Risk | Impact | Mitigation |
|------|--------|------------|
| SSH 连接不稳定 | 高 | 实现重连机制、超时控制 |
| 队列满导致请求失败 | 中 | 可配置队列大小，前端显示排队状态 |
| 日志文件过大 | 中 | 滚动压缩，限制输出记录大小 |
| 并发限制影响性能 | 低 | 可配置并发数，根据设备性能调整 |
| 配置文件密码泄露 | 高 | 支持加密格式，密钥通过环境变量 |

## Open Questions

1. **本地 CLI 执行方式**：
   - 选项 A：SSH 回连 127.0.0.1（推荐，保持协议一致性）
   - 选项 B：直接执行本地 shell 命令（需要 root 权限）

2. **批量命令是否支持事务**：
   - 当前设计：顺序执行，遇到错误继续
   - 可选：遇到错误停止（可配置）

## Future Enhancements

1. **NETCONF 协议支持**：实现 NETCONF 适配器
2. **命令模板系统**：预定义常用命令模板
3. **命令调度任务**：定时执行命令
4. **命令白名单**：限制可执行的危险命令
5. **配置备份和回滚**：自动备份配置，支持回滚
