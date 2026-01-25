# Capability: Device Interaction (Simplified)

## ADDED Requirements

### Requirement: Protocol Adapter Interface

系统 MUST 提供统一的协议适配器接口，支持多种网络设备协议（SSH、Telnet、未来可能的 NETCONF）。

#### Scenario: SSH protocol support

- **WHEN** 系统需要通过 SSH 协议连接到本地设备（127.0.0.1）
- **THEN** 系统应使用 SSH 协议适配器建立连接
- **AND** 应支持密码认证
- **AND** 应支持终端模式（PTY）执行交互式命令

#### Scenario: Telnet protocol support

- **WHEN** 系统需要通过 Telnet 协议连接设备
- **THEN** 系统应使用 Telnet 协议适配器建立连接
- **AND** 应正确处理 Telnet 选项协商
- **AND** 应正确解析命令提示符和响应

#### Scenario: Future NETCONF protocol support

- **WHEN** 后续需要支持 NETCONF 协议
- **THEN** 应通过实现 ProtocolAdapter 接口添加 NETCONF 适配器
- **AND** 不应修改上层业务代码
- **AND** 应支持 NETCONF 的 XML 配置操作

---

### Requirement: Connection Pool and Command Queue

系统 MUST 实现连接池配合命令队列，支持多用户并发访问和排队机制。

#### Scenario: Connection pool management

- **WHEN** 系统启动时
- **THEN** 应根据配置创建连接池（默认最大 3 个连接）
- **AND** 应支持动态配置并发连接数
- **AND** 应监控连接状态

#### Scenario: Acquiring available connection

- **WHEN** 用户请求执行命令
- **AND** 连接池有空闲连接
- **THEN** 系统应立即分配连接
- **AND** 执行命令后应释放连接
- **AND** 应返回执行结果

#### Scenario: Queue when pool is full

- **WHEN** 用户请求执行命令
- **AND** 连接池已满（无空闲连接）
- **THEN** 系统应将请求加入队列
- **AND** 应按 FIFO 顺序处理队列
- **AND** 应等待连接释放后执行

#### Scenario: Queue timeout

- **WHEN** 请求在队列中等待时间超过配置的超时时间（默认 60 秒）
- **THEN** 系统应返回队列超时错误（ErrQueueTimeout）
- **AND** 应从队列中移除该请求
- **AND** 不应执行该命令

#### Scenario: Configurable concurrency

- **WHEN** 管理员需要调整并发连接数
- **THEN** 系统应支持通过配置文件修改 max_connections
- **AND** 应支持热加载配置（无需重启）
- **AND** 应在配置生效时优雅重建连接池

---

### Requirement: Device Configuration from File

系统 MUST 从配置文件读取设备连接信息，无需数据库存储。

#### Scenario: Load device configuration

- **WHEN** 系统启动时
- **THEN** 应从 config/settings.yml 读取 device.connection 配置
- **AND** 应验证配置的完整性
- **AND** 应初始化连接池

#### Scenario: Configuration validation

- **WHEN** 加载设备配置时
- **AND** 配置缺少必需字段（host, port, protocol, username）
- **THEN** 系统应返回配置错误（ErrInvalidConfig）
- **AND** 应记录错误日志
- **AND** 不应启动设备交互服务

#### Scenario: Hot reload configuration

- **WHEN** 管理员修改配置文件并发送 SIGHUP 信号
- **THEN** 系统应重新加载配置
- **AND** 应优雅关闭现有连接
- **AND** 应使用新配置重建连接池
- **AND** 不应影响正在执行的命令

#### Scenario: Optional password encryption

- **WHEN** 配置文件中的密码以 "encrypted:" 前缀开头
- **THEN** 系统应使用环境变量 DEVICE_ENCRYPTION_KEY 解密
- **AND** 应使用 AES-256-GCM 算法
- **AND** 解密失败时应返回错误

---

### Requirement: Command Execution

系统 MUST 支持向设备发送命令并获取执行结果。

#### Scenario: Single command execution

- **WHEN** 用户发送单个命令
- **THEN** 系统应从连接池获取连接（或排队）
- **AND** 应发送命令到设备
- **AND** 应等待命令执行完成
- **AND** 应返回命令输出和执行状态
- **AND** 应释放连接

#### Scenario: Batch command execution

- **WHEN** 用户发送多个命令
- **THEN** 系统应按顺序执行所有命令
- **AND** 应返回每个命令的执行结果
- **AND** 如果某个命令失败，应继续执行后续命令
- **AND** 应在所有命令完成后释放连接

#### Scenario: Command timeout

- **WHEN** 命令执行时间超过配置的超时时间（默认 30 秒）
- **THEN** 系统应终止命令执行
- **AND** 应返回超时错误（ErrCommandTimeout）
- **AND** 应释放连接
- **AND** 应记录超时事件

#### Scenario: Output size limit

- **WHEN** 命令返回大量输出（超过 1MB）
- **THEN** 系统应截断输出
- **AND** 应返回截断提示
- **AND** 应记录完整输出到日志（可选）

---

### Requirement: Execution Logging

系统 MUST 使用日志文件记录命令执行历史，无需数据库。

#### Scenario: Log command execution

- **WHEN** 命令执行完成
- **THEN** 系统应记录执行日志到文件（logs/command.log）
- **AND** 日志应包含：时间戳、用户ID、用户名、命令、输出大小、执行结果、耗时、客户端IP
- **AND** 日志格式应为 JSON
- **AND** 记录日志失败不应影响命令执行结果

#### Scenario: Log file rotation

- **WHEN** 日志文件超过配置的大小限制（默认 100MB）
- **THEN** 系统应创建新的日志文件
- **AND** 应压缩旧日志文件
- **AND** 应保留指定数量的历史文件（默认 3 个）
- **AND** 应删除超过保留期的日志文件（默认 7 天）

#### Scenario: Query execution history

- **WHEN** 用户查询执行历史
- **THEN** 系统应从日志文件尾部读取记录
- **AND** 应支持分页查询（limit, offset）
- **AND** 应按时间倒序返回
- **AND** 不应影响正在写入的日志

---

### Requirement: Error Handling

系统 MUST 提供统一的错误处理和错误码体系。

#### Scenario: Connection failure error

- **WHEN** 设备连接失败
- **THEN** 系统应返回明确的错误码（ErrConnectionFailed）
- **AND** 应包含错误详细信息
- **AND** 应记录错误日志

#### Scenario: Authentication failure error

- **WHEN** 设备认证失败
- **THEN** 系统应返回明确的错误码（ErrAuthFailed）
- **AND** 不应记录敏感信息（明文密码）
- **AND** 应记录认证失败日志

#### Scenario: Queue full error

- **WHEN** 命令队列已满（超过 max_queue_size）
- **THEN** 系统应返回明确的错误码（ErrQueueFull）
- **AND** 应提示用户稍后重试
- **AND** 应记录警告日志

---

### Requirement: API Endpoints

系统 MUST 提供完整的 RESTful API 接口。

#### Scenario: Execute command API

- **WHEN** 客户端 POST /api/v1/device/command/execute
- **THEN** 系统应执行单个命令
- **AND** 请求体应包含：{"command": "show version"}
- **AND** 响应应包含：{"output": "...", "duration": "123ms", "success": true}
- **AND** 应需要认证和权限

#### Scenario: Batch execute API

- **WHEN** 客户端 POST /api/v1/device/command/batch
- **THEN** 系统应执行多个命令
- **AND** 请求体应包含：{"commands": ["cmd1", "cmd2"]}
- **AND** 响应应包含每个命令的结果数组
- **AND** 应需要认证和权限

#### Scenario: Query history API

- **WHEN** 客户端 GET /api/v1/device/command/history?limit=50&offset=0
- **THEN** 系统应返回执行历史
- **AND** 响应应包含历史记录数组
- **AND** 应支持分页参数
- **AND** 应需要认证和权限

#### Scenario: Device status API

- **WHEN** 客户端 GET /api/v1/device/status
- **THEN** 系统应返回连接状态
- **AND** 响应应包含：connected, active_connections, queue_size
- **AND** 应需要认证和权限

---

### Requirement: Access Control

系统 MUST 集成现有 RBAC 权限系统，控制命令执行权限。

#### Scenario: API permission control

- **WHEN** 用户访问命令执行 API
- **THEN** 系统应验证用户是否有相应权限
- **AND** 无权限时应返回 403 错误
- **AND** 应支持按钮级权限控制

#### Scenario: Log user information

- **WHEN** 执行命令时
- **THEN** 系统应从 gin.Context 提取用户信息
- **AND** 应记录用户ID和用户名到执行日志
- **AND** 应记录客户端IP地址

---

### Requirement: Memory Optimization

系统 MUST 在资源受限环境下（50-70MB）稳定运行。

#### Scenario: Connection pool limits

- **WHEN** 系统运行在资源受限环境
- **THEN** 应限制最大连接数（默认 3）
- **AND** 应及时清理空闲连接
- **AND** 应监控内存使用情况

#### Scenario: Log output limit

- **WHEN** 记录命令输出到日志
- **THEN** 应限制记录的输出大小（默认 10KB）
- **AND** 超过限制时应截断
- **AND** 应记录实际输出大小

#### Scenario: Queue size limit

- **WHEN** 命令队列接近上限
- **THEN** 应记录警告日志
- **AND** 达到上限时应拒绝新请求
- **AND** 应返回明确的错误信息
