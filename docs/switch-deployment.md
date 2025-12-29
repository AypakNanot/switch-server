# Switch Device Deployment Guide
# 交换机设备部署指南

This guide explains how to deploy go-admin to network switches, routers, and other embedded devices.

本指南介绍如何将 go-admin 部署到网络交换机、路由器和其他嵌入式设备。

---

## Table of Contents / 目录

- [Supported Devices / 支持的设备](#supported-devices)
- [System Requirements / 系统要求](#system-requirements)
- [Quick Start / 快速开始](#quick-start)
- [Build Instructions / 编译说明](#build-instructions)
- [Deployment Methods / 部署方式](#deployment-methods)
- [Configuration / 配置](#configuration)
- [Service Management / 服务管理](#service-management)
- [Troubleshooting / 故障排除](#troubleshooting)
- [Performance Tuning / 性能优化](#performance-tuning)

---

## Supported Devices / 支持的设备

### CPU Architectures / CPU 架构

| Architecture | GOARCH | Typical Devices / 典型设备 | Priority / 优先级 |
|--------------|--------|----------------------------|-------------------|
| **ARMv7** | arm (GOARM=7) | Raspberry Pi 2/3, most switches | P0 |
| **ARM64** | arm64 | Raspberry Pi 4/5, newer switches | P0 |
| **MIPSLE** | mipsle | TP-Link, Netgear routers | P0 |
| **MIPS** | mips | Enterprise switches | P1 |
| **ARMv6** | arm (GOARM=6) | Raspberry Pi Zero, Pi 1 | P1 |
| **ARMv5** | arm (GOARM=5) | Old devices | P2 |
| **MIPS64** | mips64 | High-end routers | P2 |
| **MIPS64LE** | mips64le | High-end routers | P2 |
| **PPC64** | ppc64le | IBM equipment | P2 |

### Device Examples / 设备示例

**ARMv7 Devices:**
- Raspberry Pi 2, 3 (32-bit)
- Most network switches (2015-2020)
- OpenWrt routers with ARM chips

**ARM64 Devices:**
- Raspberry Pi 4, 5
- Newer enterprise switches
- High-end routers

**MIPS Devices:**
- TP-Link Archer series
- Netgear Nighthawk series
- Older OpenWrt routers

---

## System Requirements / 系统要求

### Minimum Requirements / 最低要求

| Resource | Minimum / 最低 | Recommended / 推荐 |
|----------|----------------|---------------------|
| **RAM** | 256 MB | 512 MB |
| **Storage** | 100 MB free | 200 MB free |
| **CPU** | 600 MHz single-core | 1 GHz dual-core |
| **OS** | Linux kernel 3.10+ | Linux kernel 4.0+ |

### Expected Memory Usage / 预期内存占用

With `config/settings.switch.yml`:
- Go runtime: ~30-50 MB
- SQLite + data: ~10-20 MB
- Connection pools/queues: ~10-20 MB
- HTTP server: ~5-10 MB
- **Total: ~60-100 MB**

### Software Requirements / 软件要求

- **SSH access** to the device (root or sudo privileges)
- **Linux** with standard POSIX utilities
- **Optional**: systemd or init.d for service management

---

## Quick Start / 快速开始

### Step 1: Identify Your Device Architecture / 第一步：识别设备架构

Login to your device and check the architecture:

登录到您的设备并检查架构：

```bash
ssh root@<switch-ip>
uname -m
```

Common outputs:
- `armv7l` → Use ARMv7 (`make build-armv7`)
- `aarch64` → Use ARM64 (`make build-arm64`)
- `mips` → Use MIPS (`make build-mips`)
- `mipsle` → Use MIPSLE (`make build-mipsle`)

### Step 2: Build the Binary / 第二步：编译二进制

```bash
# For ARMv7 (most common)
make build-armv7

# For ARM64 (newer devices)
make build-arm64

# For MIPSLE (routers)
make build-mipsle

# Build all common architectures
make build-switch
```

### Step 3: Deploy / 第三步：部署

```bash
# Deploy using the automated script
./scripts/deploy-to-switch.sh --host=192.168.1.1 --user=root --arch=armv7
```

Or deploy manually:

或手动部署：

```bash
# Copy the binary
scp go-admin-armv7 root@192.168.1.1:/usr/bin/go-admin

# SSH to the device
ssh root@192.168.1.1

# Set permissions
chmod +x /usr/bin/go-admin

# Create config directory
mkdir -p /etc/go-admin

# Start the service
/usr/bin/go-admin server -c /etc/go-admin/settings.yml
```

---

## Build Instructions / 编译说明

### Prerequisites / 前提条件

- Go 1.20 or later
- Make utility
- (Optional) UPX for binary compression

### Individual Architecture Builds / 单架构编译

```bash
# ARM variants
make build-armv5    # ARMv5 (old devices)
make build-armv6    # ARMv6 (Raspberry Pi Zero)
make build-armv7    # ARMv7 (Raspberry Pi 2/3, most switches)
make build-arm64    # ARM64 (Raspberry Pi 4/5)

# MIPS variants
make build-mips     # MIPS big-endian
make build-mipsle   # MIPS little-endian (most common)
make build-mips64   # MIPS64 big-endian
make build-mips64le # MIPS64 little-endian

# PowerPC variants
make build-ppc64    # PPC64 big-endian
make build-ppc64le  # PPC64 little-endian
```

### Batch Builds / 批量编译

```bash
# Build common switch architectures
make build-switch

# Build all supported architectures
make build-all

# List all available architectures
make list-arch
```

### Verifying the Binary / 验证二进制

```bash
# Check architecture
file go-admin-armv7
# Output: ELF 32-bit LSB executable, ARM, EABI5 version 1 (SYSV), statically linked

# Check for dynamic dependencies (should return "not a dynamic executable")
ldd go-admin-armv7 2>&1 | grep "not a dynamic"

# Check binary size
ls -lh go-admin-armv7
```

### Optional: Binary Compression / 可选：二进制压缩

```bash
# Install UPX (Linux)
sudo apt-get install upx

# Compress the binary (reduces size by ~60%)
upx --best --lzma go-admin-armv7

# Verify compressed binary
ls -lh go-admin-armv7
```

**Warning**: UPX compression may cause issues on some platforms. Test before deploying to production.

---

## Deployment Methods / 部署方式

### Method 1: Automated Script (Recommended) / 方法一：自动化脚本（推荐）

The `scripts/deploy-to-switch.sh` script provides:
- Automatic backup
- Safe deployment with rollback
- Health checks
- Service management

**Usage / 用法:**

```bash
# Deploy new version
./scripts/deploy-to-switch.sh \
  --host=192.168.1.1 \
  --user=root \
  --port=22 \
  --arch=armv7 \
  --binary=./go-admin-armv7

# Restart existing service
./scripts/deploy-to-switch.sh --host=192.168.1.1 --action=restart

# Stop service
./scripts/deploy-to-switch.sh --host=192.168.1.1 --action=stop

# Check status
./scripts/deploy-to-switch.sh --host=192.168.1.1 --action=status

# Rollback to previous version
./scripts/deploy-to-switch.sh --host=192.168.1.1 --action=rollback
```

**Features / 功能:**
- Backup existing binary before upload
- Automatic rollback on failure
- Health check after deployment
- Service status reporting

### Method 2: Manual Deployment / 方法二：手动部署

**Step 1: Prepare the device / 准备设备**

```bash
# SSH to the device
ssh root@<switch-ip>

# Create directories
mkdir -p /usr/bin
mkdir -p /etc/go-admin
mkdir -p /tmp/go-admin/logs
```

**Step 2: Copy files / 复制文件**

```bash
# From your local machine
scp go-admin-armv7 root@<switch-ip>:/usr/bin/go-admin
scp config/settings.switch.yml root@<switch-ip>:/etc/go-admin/settings.yml
```

**Step 3: Configure permissions / 配置权限**

```bash
# SSH to the device
ssh root@<switch-ip>

# Set executable permission
chmod +x /usr/bin/go-admin

# Verify
ls -l /usr/bin/go-admin
```

**Step 4: Test run / 测试运行**

```bash
# Run in foreground for testing
/usr/bin/go-admin server -c /etc/go-admin/settings.yml

# If successful, press Ctrl+C and run in background:
nohup /usr/bin/go-admin server -c /etc/go-admin/settings.yml > /tmp/go-admin.log 2>&1 &
```

### Method 3: OpenWrt IPK Package / 方法三：OpenWrt IPK 包

**Note**: This method is planned but not yet implemented. See [Open Questions](#open-questions) in the design document.

---

## Configuration / 配置

### Switch-Optimized Configuration / 交换机优化配置

Use `config/settings.switch.yml` for low-memory environments:

```yaml
settings:
  application:
    mode: prod                    # Production mode (less overhead)
    host: 0.0.0.0
    port: 8000
    readtimeout: 30               # Longer timeouts reduce reconnects
    writertimeout: 30
    enabledp: false               # Disable data permission (saves CPU)

  logger:
    path: /tmp/go-admin/logs      # Use /tmp to avoid flash writes
    stdout: '1'                   # Console logging (no file I/O)
    level: warn                   # Less verbose logging
    enableddb: false              # Disable DB query logging

  database:
    driver: sqlite3               # Pure Go SQLite (no CGO)
    source: /tmp/go-admin-db.db   # Use /tmp or external storage
    maxOpenConns: 5               # Reduced from 100
    maxIdleConns: 2               # Reduced from 10
    connMaxLifetime: 300
    connMaxIdleTime: 60

  queue:
    memory:
      poolSize: 20                # Reduced from 100
```

### Custom Configuration / 自定义配置

**1. Change the port:**
```yaml
settings:
  application:
    port: 8080  # Custom port
```

**2. Use external storage:**
```yaml
settings:
  database:
    source: /mnt/usb/go-admin-db.db  # USB storage
  logger:
    path: /mnt/usb/logs               # USB logs
```

**3. Adjust for more memory:**
```yaml
settings:
  database:
    maxOpenConns: 10     # Increase if you have 1GB+ RAM
  queue:
    memory:
      poolSize: 50
  logger:
    level: info          # More verbose logs
```

---

## Service Management / 服务管理

### Using systemd / 使用 systemd

**1. Create service file / 创建服务文件:**

```bash
# Copy the service template
scp scripts/go-admin.service root@<switch-ip>:/etc/systemd/system/

# Or create manually on the device
ssh root@<switch-ip>
cat > /etc/systemd/system/go-admin.service << 'EOF'
[Unit]
Description=Go-Admin Application Server
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/usr/bin
ExecStart=/usr/bin/go-admin server -c /etc/go-admin/settings.yml
Restart=on-failure
RestartSec=5
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
EOF
```

**2. Enable and start / 启用并启动:**

```bash
# Reload systemd
systemctl daemon-reload

# Enable on boot
systemctl enable go-admin

# Start service
systemctl start go-admin

# Check status
systemctl status go-admin

# View logs
journalctl -u go-admin -f
```

### Using init.d / 使用 init.d

**1. Create init script / 创建 init 脚本:**

```bash
scp scripts/go-admin.init root@<switch-ip>:/etc/init.d/go-admin

# Or create manually
ssh root@<switch-ip>
cat > /etc/init.d/go-admin << 'EOF'
#!/bin/sh
# Go-Admin init script

DAEMON="/usr/bin/go-admin"
CONFIG="/etc/go-admin/settings.yml"
PIDFILE="/var/run/go-admin.pid"
LOGFILE="/tmp/go-admin.log"

case "$1" in
  start)
    echo "Starting go-admin..."
    start-stop-daemon -S -x $DAEMON -- server -c $CONFIG -p $PIDFILE >> $LOGFILE 2>&1 &
    ;;
  stop)
    echo "Stopping go-admin..."
    start-stop-daemon -K -p $PIDFILE
    ;;
  restart)
    $0 stop
    sleep 2
    $0 start
    ;;
  status)
    if [ -f $PIDFILE ]; then
      echo "go-admin is running (PID: $(cat $PIDFILE))"
    else
      echo "go-admin is not running"
    fi
    ;;
  *)
    echo "Usage: $0 {start|stop|restart|status}"
    exit 1
    ;;
esac

exit 0
EOF

chmod +x /etc/init.d/go-admin
```

**2. Enable and start / 启用并启动:**

```bash
# Enable on boot (for OpenWrt)
/etc/init.d/go-admin enable

# Start service
/etc/init.d/go-admin start

# Check status
/etc/init.d/go-admin status

# View logs
tail -f /tmp/go-admin.log
```

### Manual Process Management / 手动进程管理

```bash
# Start in background
nohup /usr/bin/go-admin server -c /etc/go-admin/settings.yml > /tmp/go-admin.log 2>&1 &

# Find process ID
ps | grep go-admin

# Stop gracefully
kill -TERM $(cat /var/run/go-admin.pid)

# Force stop
kill -KILL $(ps | grep go-admin | grep -v grep | awk '{print $1}')

# Check if running
pgrep -x go-admin
```

---

## Troubleshooting / 故障排除

### Issue: "Cannot execute binary" / 无法执行二进制

**Cause**: Architecture mismatch

**Solution**:
```bash
# Check device architecture
uname -m

# Rebuild with correct architecture
make build-armv7    # for armv7l
make build-arm64    # for aarch64
make build-mipsle   # for mipsle
```

### Issue: Out of memory / 内存不足

**Symptoms**: Process killed, OOM errors

**Solutions**:
1. Use `config/settings.switch.yml`
2. Close unnecessary services:
   ```bash
   /etc/init.d/nginx stop
   /etc/init.d/cron stop
   ```
3. Add swap space:
   ```bash
   dd if=/dev/zero of=/tmp/swap bs=1M count=256
   chmod 600 /tmp/swap
   mkswap /tmp/swap
   swapon /tmp/swap
   ```

### Issue: SQLite database locked / 数据库锁定

**Cause**: Multiple processes accessing the database

**Solution**:
- Ensure only one instance is running: `pgrep -x go-admin`
- Reduce `maxOpenConns` in configuration
- Consider using external SQLite: `source: /mnt/usb/go-admin-db.db`

### Issue: Slow performance / 性能缓慢

**Solutions**:
1. Use production mode: `mode: prod`
2. Reduce logging: `level: warn`, `enableddb: false`
3. Increase timeouts: `readtimeout: 30`, `writertimeout: 30`
4. Check CPU usage: `top` or `htop`

### Issue: Cannot access web interface / 无法访问 Web 界面

**Checklist**:
```bash
# 1. Is the service running?
pgrep -x go-admin

# 2. Is the port listening?
netstat -tlnp | grep 8000

# 3. Check firewall
iptables -L -n | grep 8000

# 4. Check logs
tail -50 /tmp/go-admin.log

# 5. Test locally
curl http://localhost:8000/
```

### Issue: High flash wear / 闪存损耗高

**Solution**: Use `/tmp` or external storage

```yaml
settings:
  database:
    source: /tmp/go-admin-db.db        # RAM
    # or
    source: /mnt/usb/go-admin-db.db    # External storage
  logger:
    path: /tmp/go-admin/logs           # RAM
```

---

## Performance Tuning / 性能优化

### Memory Optimization / 内存优化

| Setting | Default | Switch | Impact |
|---------|---------|--------|--------|
| `maxOpenConns` | 100 | 5 | -95% connection memory |
| `maxIdleConns` | 10 | 2 | -80% idle overhead |
| `poolSize` | 100 | 20 | -80% queue memory |
| `log level` | info | warn | -60% log volume |

### Storage Optimization / 存储优化

1. **Binary compression** (saves ~60% size):
   ```bash
   upx --best --lzma go-admin-armv7
   ```

2. **Use external storage**:
   ```bash
   # Mount USB drive
   mount /dev/sda1 /mnt/usb

   # Update config
   database:
     source: /mnt/usb/go-admin-db.db
   logger:
     path: /mnt/usb/logs
   ```

3. **Periodic cleanup**:
   ```bash
   # Add to crontab
   echo "0 3 * * * rm -f /tmp/go-admin/logs/*.log" | crontab -
   ```

### CPU Optimization / CPU 优化

1. **Use production mode**: `mode: prod`
2. **Disable unnecessary features**:
   ```yaml
   enabledp: false     # Disable data permissions
   enableddb: false    # Disable query logging
   ```
3. **Increase timeouts**: Reduces connection overhead

---

## Advanced Topics / 高级主题

### Building Custom Architectures / 编译自定义架构

If your device uses an architecture not in the Makefile:

```bash
# Example: ARMv6 with hard float
env CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=arm \
    GOARM=6 \
    GOARM=softfloat \
    go build -ldflags="-w -s" -o go-admin-custom .
```

### Cross-Compilation Tools / 交叉编译工具

```bash
# Install qemu-user-static for testing
sudo apt-get install qemu-user-static

# Test ARM binary on x86_64
qemu-arm-static ./go-admin-armv7 version
```

### Docker-based Build / 基于 Docker 的编译

```dockerfile
FROM golang:1.21-alpine AS builder
RUN apk add --no-cache make git
WORKDIR /app
COPY . .
RUN make build-armv7
```

---

## Open Questions / 开放问题

See the [design document](../openspec/changes/add-switch-device-deployment-support/design.md) for:

- Q1: Should we provide OpenWrt IPK packages? (Planned)
- Q2: UPX compression support? (Optional)
- Q3: API-only mode to disable embedded frontend? (Future)

---

## References / 参考

- [Design Document](../openspec/changes/add-switch-device-deployment-support/design.md)
- [OpenSpec Proposal](../openspec/changes/add-switch-device-deployment-support/proposal.md)
- [Go Cross-Compilation](https://golang.org/doc/install/source#environment)
- [OpenWrt Documentation](https://openwrt.org/docs/start)

---

## Contributing / 贡献

If you have tested go-admin on a device not listed here, please contribute:

1. Device model and specifications
2. Architecture used (`uname -m`)
3. Any configuration changes required
4. Performance metrics

Submit a PR or issue to help improve this guide.

---

**Last Updated**: 2025-12-28
**Version**: 1.0.0
