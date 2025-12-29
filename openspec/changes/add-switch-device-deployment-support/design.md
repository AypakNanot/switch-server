# Design: Switch Device Deployment Support

## Context
网络交换机和路由器通常运行嵌入式 Linux 系统（如 OpenWrt、DD-WRT、自定义固件），具有以下特点：

1. **CPU 架构多样**：
   - ARM: ARMv5-ARMv7 (32位)，ARM64/aarch64 (64位)
   - MIPS: 大端序、小端序（MIPSLE），32/64位
   - PowerPC: 主要用于企业级设备
   - x86_64: 部分高端交换机

2. **资源约束**：
   - 内存：256MB-1GB（典型 512MB）
   - 存储：8MB-128MB 闪存
   - CPU：单核或双核，低频率（600MHz-2GHz）

3. **系统环境**：
   - 使用 BusyBox 或精简 Linux
   - 可能缺少 glibc（使用 uClibc 或 musl）
   - 无显示界面，仅 SSH/串口/Web 管理界面

## Goals / Non-Goals

**Goals**:
- 支持最常见的交换机架构（ARMv7、ARM64、MIPSLE）
- 提供低内存优化配置
- 静态编译二进制，无外部依赖
- 自动化部署脚本
- 详细的部署文档

**Non-Goals**:
- 不支持所有可能的架构（专注最常见的）
- 不提供设备特定的二进制包（由用户自行编译）
- 不支持 Windows CE 或专有操作系统
- 不修改核心应用功能（仅部署相关）

## Decisions

### 1. 纯 Go SQLite 驱动（已完成）
**决策**: 项目已使用 `github.com/glebarez/sqlite` 纯 Go 驱动

**原因**:
- 无需 CGO，可静态编译
- 跨平台编译简单
- 性能约为 CGO 版本的 90-95%，可接受

**权衡**:
- ✅ 优势：无外部依赖，易部署
- ⚠️ 劣势：性能略低（但在交换机场景可接受）

### 2. 交叉编译目标架构
**决策**: 优先支持以下架构

| 架构 | GOARCH | GOARM | 典型设备 | 优先级 |
|------|--------|-------|----------|--------|
| ARMv7 | arm | 7 | 树莓派 2/3，大部分交换机 | P0 |
| ARM64 | arm64 | - | 树莓派 4/5，新型交换机 | P0 |
| MIPSLE | mipsle | - | 路由器、TP-Link、Netgear | P0 |
| MIPS | mips | - | 部分企业级设备 | P1 |
| ARMv6 | arm | 6 | 树莓派 Zero | P1 |
| ARMv5 | arm | 5 | 老旧设备 | P2 |
| MIPS64 | mips64 | - | 高端路由器 | P2 |
| MIPS64LE | mips64le | - | 高端路由器 | P2 |
| PPC64 | ppc64le | - | IBM 设备 | P2 |

**原因**:
- ARMv7/ARM64 和 MIPSLE 是最常见的交换机架构
- 覆盖 80% 以上的设备
- 可根据需求扩展

**权衡**:
- ✅ 覆盖主流设备
- ⚠️ 不支持罕见架构（可通过 Makefile 手动添加）

### 3. 低内存配置优化
**决策**: 创建 `config/settings.switch.yml`，优化参数如下

```yaml
database:
  maxOpenConns: 5       # 默认: 100
  maxIdleConns: 2       # 默认: 10
queue:
  memory:
    poolSize: 20        # 默认: 100
logger:
  level: warn           # 默认: info/trace
  enableddb: false      # 默认: false
  stdout: '1'           # 启用控制台日志，避免文件 I/O
application:
  mode: prod            # 禁用调试模式
  readtimeout: 30       # 增加超时，减少连接重建
  writertimeout: 30
```

**原因**:
- 连接池：5-2 个连接足以应对低并发场景
- 队列：减少内存队列大小
- 日志：WARN 级别足够，避免大量日志

**估算内存占用**:
- 基础运行时：~30-50MB
- SQLite + 数据：~10-20MB
- 连接池/队列：~10-20MB
- **总计：约 60-100MB**（在 256MB 系统中可接受）

### 4. 静态编译选项
**决策**: 使用以下编译选项

```bash
CGO_ENABLED=0 \
GOOS=linux \
GOARCH=arm \
GOARM=7 \
go build \
  -ldflags="-w -s" \
  -o go-admin-armv7 \
  main.go
```

**说明**:
- `CGO_ENABLED=0`: 禁用 CGO，纯静态编译
- `-ldflags="-w -s"`: 去除调试信息和符号表
- 结果：无外部依赖，可在任何 Linux 发行版上运行

**二进制大小**:
- 未压缩：约 40-50MB
- UPX 压缩（可选）：约 15-20MB（可能影响某些平台）

### 5. 部署脚本设计
**决策**: 提供 Bash 脚本 `scripts/deploy-to-switch.sh`

**功能**:
- SSH 连接到目标设备
- 备份现有二进制
- 上传新二进制
- 设置权限
- 重启服务
- 健康检查
- 失败回滚

**用法示例**:
```bash
# 部署
./scripts/deploy-to-switch.sh \
  --host=192.168.1.1 \
  --user=root \
  --port=22 \
  --arch=armv7 \
  --binary=./go-admin-armv7

# 仅重启
./scripts/deploy-to-switch.sh \
  --host=192.168.1.1 \
  --action=restart

# 停止服务
./scripts/deploy-to-switch.sh \
  --host=192.168.1.1 \
  --action=stop
```

## Risks / Trade-offs

### Risk 1: 内存不足
**风险**: 在 256MB 设备上，系统或其他进程可能占用大部分内存

**缓解措施**:
- 提供内存优化配置
- 文档说明最低内存要求（256MB）
- 建议在部署前关闭不必要的服务
- 提供内存监控脚本

### Risk 2: 存储空间不足
**风险**: 交换机闪存可能只有 16-32MB，二进制可能无法存放

**缓解措施**:
- 基础二进制约 40-50MB
- 可选 UPX 压缩（减小到 ~15-20MB）
- 文档说明存储空间要求
- 建议挂载外部存储（USB/SD 卡）

### Risk 3: 架构不兼容
**风险**: 某些交换机使用不常见的架构或特定厂商定制

**缓解措施**:
- 文档列出支持的架构
- 提供手动交叉编译指南
- 用户可根据 Makefile 自定义目标架构

### Risk 4: 性能问题
**风险**：嵌入式 CPU 性能较弱，响应可能较慢

**缓解措施**:
- 使用生产模式（mode: prod）
- 减少数据库查询日志
- 优化配置减少不必要的中间件
- 文档说明性能预期

## Migration Plan

### 开发环境准备

#### 统一版本要求
项目统一使用支持 Windows 7 的最后一个版本，确保最大兼容性：

- **Go 1.20.14**
  - 下载：https://go.dev/dl/go1.20.14.windows-386.zip 或 go1.20.14.windows-amd64.zip
  - 设置 GOPATH 和 GOROOT 环境变量

- **Node.js 16.20.2**
  - 下载：https://nodejs.org/dist/v16.20.2/
  - NPM 版本：8.19.4
  - 实测可在 Windows 7 上正常运行

#### 交叉编译要求
- 使用 `CGO_ENABLED=0` 确保纯静态编译
- 支持交叉编译到 Linux ARM/MIPS 等架构

#### 其他工具
- Linux/macOS 用户：安装 `sshpass`（用于脚本自动登录）

### 编译流程
1. 选择目标架构（参考设备文档或使用 `uname -m` 查看）
2. 运行对应的 make 命令
3. 验证二进制：`file go-admin-*`
4. （可选）压缩二进制：`upx go-admin-armv7`

### 部署流程
1. 确保目标设备有足够存储（至少 100MB）
2. 测试 SSH 连接：`ssh root@<switch-ip>`
3. 准备配置文件
4. 运行部署脚本或手动部署
5. 验证服务状态

### 回滚计划
**如果部署失败**:
1. 脚本自动回滚到备份的二进制
2. 手动回滚：恢复之前的二进制文件
3. 检查日志：`journalctl -u go-admin -n 50`

## Open Questions

### Q1: 是否需要 OpenWrt IPK 包？
**问题**: OpenWrt 用户更喜欢 IPK 包而不是手动部署

**答案**: 暂不实现，作为可选增强（tasks.md 8.1）
**原因**:
- IPK 包需要额外的构建基础设施
- 不同 OpenWrt 版本兼容性问题
- 手动部署已经足够简单

### Q2: 是否支持 UPX 压缩？
**问题**: UPX 可以显著减小二进制大小，但可能有兼容性问题

**答案**: 作为可选功能，不作为默认
**原因**:
- UPX 在某些架构上可能有问题
- 压缩后的首次启动较慢（解压开销）
- 用户可根据需要自行压缩

### Q3: 如何处理依赖的前端静态文件？
**问题**: 当前前端文件已嵌入二进制，会增加 20-30MB

**答案**: 保持现状，前端文件嵌入
**原因**:
- 简化部署（单文件）
- 交换机通常有足够存储（>64MB）
- 如需节省空间，可禁用前端（API-only 模式）

### Q4: 是否支持 musl libc？
**问题**: Alpine Linux 和某些嵌入式系统使用 musl 而非 glibc

**答案**: 由于使用静态编译（CGO_ENABLED=0），不依赖任何 libc
**优势**: 可在任何 Linux 发行版上运行，无需关心 libc 差异

## Implementation Notes

### Makefile 结构
```makefile
# Switch device builds
build-switch: build-armv7 build-arm64 build-mipsle
	@echo "Built for switch architectures"

build-armv7:
	env CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 \
		go build -ldflags="-w -s" -o go-admin-armv7 .
	@file go-admin-armv7

build-arm64:
	env CGO_ENABLED=0 GOOS=linux GOARCH=arm64 \
		go build -ldflags="-w -s" -o go-admin-arm64 .
	@file go-admin-arm64

build-mipsle:
	env CGO_ENABLED=0 GOOS=linux GOARCH=mipsle \
		go build -ldflags="-w -s" -o go-admin-mipsle .
	@file go-admin-mipsle

# ... 更多架构
```

### 部署脚本结构
```bash
#!/bin/bash
# deploy-to-switch.sh

HOST=""
USER="root"
PORT="22"
ARCH=""
ACTION="deploy"
BINARY=""

while [[ $# -gt 0 ]]; do
  case $1 in
    --host) HOST="$2"; shift ;;
    --user) USER="$2"; shift ;;
    --arch) ARCH="$2"; shift ;;
    --action) ACTION="$2"; shift ;;
    *) echo "Unknown option: $1"; exit 1 ;;
  esac
  shift
done

case $ACTION in
  deploy)
    # 部署逻辑
    ;;
  restart)
    # 重启逻辑
    ;;
  stop)
    # 停止逻辑
    ;;
esac
```

### 配置文件结构
```yaml
# config/settings.switch.yml
settings:
  application:
    mode: prod
    host: 0.0.0.0
    port: 8000
    readtimeout: 30
    writertimeout: 30
    enabledp: false
  logger:
    path: /tmp/go-admin/logs
    stdout: '1'
    level: warn
    enableddb: false
  database:
    driver: sqlite3
    source: /tmp/go-admin-db.db
  queue:
    memory:
      poolSize: 20
```
