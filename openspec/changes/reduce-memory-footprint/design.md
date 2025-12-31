# Design: Reduce Memory Footprint to 50MB

## Context

### 当前内存占用分析

根据现有配置 `settings.switch.yml` 的内存估算：

| 组件 | 当前占用 | 说明 |
|------|----------|------|
| Go 运行时 | 30-50MB | 基础运行时、GC、调度器 |
| SQLite + 数据 | 10-20MB | 数据库连接、缓存 |
| 连接池 | 5-10MB | maxOpenConns: 5, maxIdleConns: 2 |
| 内存队列 | 5-10MB | poolSize: 20 |
| Gin/HTTP | 5-10MB | HTTP 服务、中间件 |
| 前端静态文件 | 10-20MB | embed.FS 缓存 |
| **总计** | **65-110MB** | 实际运行约 80-100MB |

### 目标：~50MB

需要减少约 **30-50MB**，优化策略：
1. **Go 运行时调优**：-10MB
2. **连接池/队列**：-10MB
3. **前端文件**：-15MB（可选禁用）
4. **中间件优化**：-5MB

## Goals / Non-Goals

**Goals**:
- 将内存占用降低到 ~50MB（空闲状态）
- 保持核心功能完整（用户管理、配置管理、API）
- 提供可配置的内存级别（normal/minimal）
- 添加内存测量工具

**Non-Goals**:
- 不降低系统稳定性
- 不移除核心业务功能
- 不影响高并发场景（使用 normal 配置）

## Decisions

### 1. Go 运行时内存调优

**决策**: 在 `cmd/api/server.go` 启动时配置 Go 运行时参数

```go
// 在 main() 函数开始时添加
func initRuntime() {
    // 限制 GOMAXPROCS 到 CPU 核心数（或更少）
    // 单核设备设置为 1，减少调度开销
    runtime.GOMAXPROCS(1)

    // 设置 GOGC 降低 GC 频率
    // 默认 100，设置为 200 表示堆增长 200% 才触发 GC
    // 注意：设置过高会导致单次 GC 暂停时间增加
    debug.SetGCPercent(200)

    // Go 1.19+: 设置软内存限制
    // 设置为 60MB，超过此值时 GC 更积极
    const memoryLimit = 60 * 1024 * 1024
    debug.SetMemoryLimit(memoryLimit)

    // 设置最大线程数（减少栈内存）
    debug.SetMaxThreads(100)
}
```

**原因**:
- `GOMAXPROCS=1`：减少 goroutine 调度开销和内存
- `GOGC=200`：减少 GC 频率，降低 CPU 占用
- `SetMemoryLimit`：Go 1.19+ 特性，软限制内存使用

**权衡**:
- ✅ 优势：减少 5-15MB 内存
- ⚠️ 劣势：GC 频率降低可能导致内存峰值略高
- ⚠️ 劣势：单核 CPU 可能影响并发性能（但低并发场景可接受）

**内存节省**: ~10MB

### 2. 进一步降低连接池和队列

**决策**: 在 `settings.minimal.yml` 中设置

```yaml
settings:
  database:
    maxOpenConns: 2       # 从 5 降低到 2
    maxIdleConns: 1       # 从 2 降低到 1
    connMaxLifetime: 300
    connMaxIdleTime: 60

  queue:
    memory:
      poolSize: 5         # 从 20 降低到 5
```

**原因**:
- 低并发场景（1-3 个用户）不需要 5 个数据库连接
- 队列主要用于异步日志，低并发时 5 个足够

**权衡**:
- ✅ 优势：减少 ~10MB 内存
- ⚠️ 劣势：高并发时可能排队等待

**内存节省**: ~10MB

### 3. 可选禁用前端静态文件

**决策**: 添加配置选项控制是否加载前端

```yaml
settings:
  application:
    enableFrontend: false    # 禁用前端，仅 API 模式
```

代码修改 `cmd/api/server.go`:

```go
// 条件编译前端路由
if config.ApplicationConfig.EnableFrontend {
    // 注册前端静态文件路由
    registerFrontendRoutes(r)
}
```

**原因**:
- 前端静态文件（CSS/JS/HTML）占用 10-20MB
- API-only 模式适合仅通过 API 调用的场景

**权衡**:
- ✅ 优势：减少 15-20MB 内存
- ⚠️ 劣势：无法通过浏览器访问 Web UI

**内存节省**: ~15MB（如果禁用）

### 4. 减少中间件和缓冲区

**决策**: 在 `settings.minimal.yml` 中禁用非必要中间件

```yaml
settings:
  application:
    enableMiddleware:
      requestID: false       # 禁用请求 ID 中间件
      sentinel: false        # 禁用限流中间件
      metrics: false         # 禁用指标收集
```

代码修改 `cmd/api/server.go`:

```go
// 条件启用中间件
if config.ApplicationConfig.EnableMiddleware.Sentinel {
    r.Use(common.Sentinel())
}
if config.ApplicationConfig.EnableMiddleware.RequestID {
    r.Use(common.RequestId(pkg.TrafficKey))
}
// ... 其他中间件
```

**原因**:
- Sentinel、Metrics 等中间件维护内部状态和缓冲区
- 低并发场景这些功能不是必需的

**权衡**:
- ✅ 优势：减少 5MB 内存
- ⚠️ 劣势：失去限流、监控等功能

**内存节省**: ~5MB

### 5. Gin 引擎优化

**决策**: 在 `cmd/api/server.go` 中配置 Gin 引擎

```go
h := gin.New()

// 减少并发缓冲区
h.Use(
    gin.Recovery(),                    // 仅保留恢复中间件
    // 其他中间件按需启用
)

// 设置读取缓冲区大小（默认 4KB，降低到 1KB）
// 注意：需要修改 Gin 源码或使用自定义 Reader
```

**原因**:
- Gin 默认为每个请求分配 4KB 缓冲区
- 低并发场景可以降低

**权衡**:
- ✅ 优势：减少 ~2MB 内存
- ⚠️ 劣势：大请求体可能需要分多次读取

**内存节省**: ~2MB

### 6. 配置文件对比

| 配置项 | settings.yml | settings.switch.yml | settings.minimal.yml |
|--------|--------------|---------------------|----------------------|
| maxOpenConns | 100 | 5 | 2 |
| maxIdleConns | 10 | 2 | 1 |
| poolSize | 100 | 20 | 5 |
| GOMAXPROCS | 自动 | 自动 | 1 |
| GOGC | 100 | 100 | 200 |
| MemoryLimit | 无 | 无 | 60MB |
| EnableFrontend | true | true | false |
| Sentinel | true | true | false |
| RequestID | true | true | false |

## Risks / Trade-offs

### Risk 1: 内存不足导致 OOM

**风险**: 即使优化后，某些设备可能仍然内存不足

**缓解措施**:
- 添加内存检测脚本，启动前检查可用内存
- 提供更激进的配置选项（禁用更多功能）
- 文档说明最低内存要求（128MB）

### Risk 2: 性能下降

**风险**: 降低连接池和 GOMAXPROCS 可能影响响应速度

**缓解措施**:
- 提供 performance 配置，用户可根据硬件选择
- 文档说明每个配置级别的适用场景
- 允许用户手动调整特定参数

### Risk 3: 功能缺失

**风险**: 禁用前端和中间件后，某些功能不可用

**缓解措施**:
- API-only 模式文档说明
- 提供独立的前端部署方案
- 清晰标注哪些功能在 minimal 模式下不可用

### Risk 4: GC 暂停时间增加

**风险**: GOGC=200 可能导致单次 GC 暂停时间变长

**缓解措施**:
- 监控 GC 暂停时间
- 如果影响用户体验，可以调整到 150
- 提供 GOGC 调优指南

## Migration Plan

### 1. 创建配置文件

- [ ] 创建 `config/settings.minimal.yml`
- [ ] 创建 `config/settings.performance.yml`（可选）

### 2. 修改代码

- [ ] `cmd/api/server.go` - 添加运行时调优函数
- [ ] `cmd/api/server.go` - 条件启用前端路由
- [ ] `cmd/api/server.go` - 条件启用中间件
- [ ] `common/middleware/` - 添加配置选项

### 3. 创建工具脚本

- [ ] `scripts/measure-memory.sh` - 测量内存占用
- [ ] `scripts/check-memory.sh` - 检查可用内存

### 4. 测试验证

- [ ] 使用 minimal 配置启动服务
- [ ] 测量空闲内存（RSS）
- [ ] 测试核心功能
- [ ] 压力测试（低并发）
- [ ] Docker 容器测试

### 5. 文档更新

- [ ] 更新 README.md 添加内存优化说明
- [ ] 更新 deploy/docs/DEPLOYMENT_GUIDE.md
- [ ] 添加配置选择指南

## Implementation Notes

### 内存测量脚本

```bash
#!/bin/bash
# scripts/measure-memory.sh

PID=$1
INTERVAL=5

echo "Measuring memory for PID: $PID"
echo "Time,RSS(MB),VSZ(MB)"

while true; do
    if [ ! -d "/proc/$PID" ]; then
        echo "Process $PID not found"
        break
    fi

    stats=$(cat /proc/$PID/status | grep -E "VmRSS|VmSize")
    rss=$(echo "$stats" | grep VmRSS | awk '{print $2}')
    vsz=$(echo "$stats" | grep VmSize | awk '{print $2}')

    rss_mb=$((rss / 1024))
    vsz_mb=$((vsz / 1024))

    echo "$(date +%H:%M:%S),$rss_mb,$vsz_mb"
    sleep $INTERVAL
done
```

### 内存检查脚本

```bash
#!/bin/bash
# scripts/check-memory.sh

REQUIRED_MB=128
available=$(free -m | grep Mem | awk '{print $7}')

echo "Available memory: ${available}MB"
echo "Required memory: ${REQUIRED_MB}MB"

if [ "$available" -lt "$REQUIRED_MB" ]; then
    echo "ERROR: Not enough memory!"
    exit 1
else
    echo "OK: Sufficient memory available"
    exit 0
fi
```

### 使用示例

```bash
# 1. 检查内存
./scripts/check-memory.sh

# 2. 使用 minimal 配置启动
./opt-switch-arm64 server -c config/settings.minimal.yml &

# 3. 测量内存
./scripts/measure-memory.sh $(pgrep opt-switch-arm64)
```

## Open Questions

### Q1: 是否需要更激进的配置（< 30MB）？

**问题**: 某些极低内存设备可能需要 < 30MB

**答案**: 暂不实现，作为可选增强
**原因**:
- 需要禁用更多功能
- 可能影响稳定性
- 用户可以根据 minimal.yml 进一步定制

### Q2: 是否自动检测硬件并选择配置？

**问题**: 能否根据内存大小自动选择合适的配置？

**答案**: 暂不实现，作为可选增强
**原因**:
- 增加复杂性
- 用户可能有特殊需求
- 当前手动选择更灵活

### Q3: 前端是否可以按需加载？

**问题**: 能否首次请求时才加载前端，而不是启动时加载？

**答案**: 暂不实现，技术复杂度高
**原因**:
- embed.FS 不支持延迟加载
- 需要重构为外部文件或 HTTP 下载
- 收益不大（禁用即可）

## Expected Memory Breakdown (After Optimization)

| 组件 | 当前 | minimal.yml | 节省 |
|------|------|-------------|------|
| Go 运行时 | 30-50MB | 20-30MB | -10MB |
| SQLite | 10-20MB | 10-15MB | -5MB |
| 连接池 | 5-10MB | 2-3MB | -7MB |
| 队列 | 5-10MB | 1-2MB | -8MB |
| Gin/HTTP | 5-10MB | 3-5MB | -2MB |
| 前端文件 | 10-20MB | 0MB (禁用) | -15MB |
| 中间件 | 5MB | 0MB (禁用) | -5MB |
| **总计** | **65-110MB** | **35-55MB** | **-50MB** |

**目标**: 空闲内存 < 50MB ✅
