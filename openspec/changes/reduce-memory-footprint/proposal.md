# Change: Reduce Memory Footprint to ~50-70MB

## Why

当前 opt-switch 在 ARM64 光交换机等资源受限设备上的内存占用约为 **80-100MB**。对于某些内存更受限的设备（如 256MB-512MB 的老旧交换机或嵌入式设备），这个内存占用仍然偏高。

用户反馈实际使用中**不需要太大的并发**，主要场景是单用户或少量用户管理光交换机配置。**前端 Web UI 是核心功能，必须保留**。因此可以进一步优化内存占用，目标控制在 **~50-70MB** 左右（包含前端）。

## What Changes

- **创建低内存配置文件**：`config/settings.minimal.yml`，针对 50-70MB 内存目标优化
- **添加 Go 运行时内存调优**：配置 GOMAXPROCS、GOGC、内存限制
- **进一步降低连接池和队列**：数据库连接池 5→2，队列池 20→5
- **选择性禁用非必要中间件**：禁用 Sentinel、Metrics（保留核心中间件）
- **保留前端静态文件**：Web UI 是核心功能，完整保留
- **优化 Gin 引擎配置**：减少并发缓冲区
- **添加内存监控和诊断工具**：帮助验证内存占用

## Impact

- **Affected specs**: deployment (new minimal-memory capability)
- **Affected code**:
  - `config/settings.minimal.yml` (new) - 低内存配置（包含前端）
  - `cmd/api/runtime.go` (new) - 运行时内存调优
  - `cmd/api/server.go` - 条件中间件加载
  - `config/extend.go` - 扩展配置结构
  - `scripts/measure-memory.sh` (new) - 内存监控脚本
  - `scripts/check-memory.sh` (new) - 内存检查脚本
  - `scripts/memory-test.sh` (new) - 内存测试脚本

- **User-visible changes**:
  - ✅ 内存占用降至 ~50-70MB（从 80-100MB）
  - ✅ 前端 Web UI 完整保留
  - ✅ 可在 256MB 内存设备上运行
  - ✅ 提供内存监控脚本

- **Trade-offs**:
  - ⚠️ 并发能力降低（适合低并发场景）
  - ⚠️ 某些高级功能（限流、监控）被禁用
  - ⚠️ 响应时间可能略有增加

- **Migration path**:
  - 现有部署不受影响
  - 新增 minimal 配置文件，用户可选择使用
  - 默认 settings.switch.yml 保持不变

## Why This Approach

1. **渐进式优化**：基于现有 settings.switch.yml 进一步优化，不破坏现有配置
2. **保留前端**：Web UI 是核心功能，必须保留，仅禁用非必要中间件
3. **Go 运行时调优**：利用 Go 1.20+ 的内存管理特性（GOMEMLIMIT）
4. **功能可选**：允许用户根据需要调整配置
5. **可验证性**：提供内存测量脚本，确保达到目标

## Success Criteria

- ✅ 使用 settings.minimal.yml 时，空闲内存占用 < 70MB
- ✅ 前端 Web UI 完整可用
- ✅ 核心功能正常（登录、菜单、数据操作）
- ✅ Docker 镜像测试验证
- ✅ 文档说明何时使用 minimal 配置
