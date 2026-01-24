# opt-switch

<div align="center">

![opt-switch Logo](https://doc-image.zhangwj.com/img/opt-switch.svg)

**轻量级光交换机管理系统**

[![Build Status](https://github.com/wenjianzhang/opt-switch/workflows/build/badge.svg)](https://github.com/opt-switch-team/opt-switch)
[![Release](https://img.shields.io/github/release/opt-switch-team/opt-switch.svg?style=flat-square)](https://github.com/opt-switch-team/opt-switch/releases)
[![License](https://img.shields.io/github/license/mashape/apistatus.svg)](https://github.com/opt-switch-team/opt-switch/blob/master/LICENSE.md)
[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://golang.org/)

[English](#english) | [简体中文](#%E7%AE%80%E4%BD%93%E4%B8%AD%E6%96%87)

</div>

---

## 简体中文

## 项目简介

**opt-switch** 是一款专为**光交换机、路由器等网络设备**设计的超轻量级后台管理系统。与传统需要 200-500MB 内存的权限管理系统不同，opt-switch 通过深度优化，仅需 **50-100MB 内存**即可运行，可直接部署在交换机设备上。

### 为什么选择 opt-switch

```
传统方案 vs opt-switch 方案：

传统方案：
┌─────────────┐         ┌──────────────┐
│  光交换机   │────────▶│  管理服务器  │
│  (华为/H3C) │  串口   │  (200-500MB) │
└─────────────┘         └──────────────┘
                                   │
                            ┌──────▼──────┐
                            │  管理员PC   │
                            └─────────────┘

opt-switch 方案：
┌─────────────────────────────────────┐
│         光交换机 / 路由器             │
│  ┌─────────────────────────────────┐ │
│  │    opt-switch (50-100MB)        │ │
│  │    ├─ Web UI (内嵌)             │ │
│  │    ├─ API 服务                  │ │
│  │    ├─ SQLite 数据库             │ │
│  │    └─ 权限管理系统              │ │
│  └─────────────────────────────────┘ │
│         ▲                            │
└─────────┼────────────────────────────┘
          │
     网络工程师
     (Web界面)
```

### 核心优势

| 特性 | 传统方案 | opt-switch |
|------|----------|------------|
| **内存占用** | 200-500MB | **50-100MB** |
| **部署位置** | 独立服务器 | **直接部署在交换机上** |
| **数据库依赖** | MySQL/PostgreSQL 服务器 | **SQLite 内置，零依赖** |
| **管理方式** | 串口/命令行 | **Web 可视化界面** |
| **启动时间** | 10-30秒 | **2-5秒** |

### 适用设备

| 设备类型 | 品牌 | 内存配置 | 推荐配置文件 |
|----------|------|----------|--------------|
| 光交换机 | 华为、H3C、锐捷 | 256MB | `config/settings.minimal.yml` |
| 光交换机 | 华为、H3C、锐捷 | 512MB | `config/settings.switch.yml` |
| 路由器 | TP-Link、Netgear | 256MB | `config/settings.minimal.yml` |
| 企业级交换机 | 各品牌 | 512MB+ | `config/settings.switch.yml` |

---

## 核心特性

### 1. 超低内存占用

- **50-70MB** 极简配置（256MB 内存设备）
- **60-100MB** 标准配置（512MB 内存设备）

### 2. 多架构支持

| 架构 | 字长 | 端序 | 典型设备 | 编译命令 |
|------|------|------|----------|----------|
| **ARM64** | 64位 | - | 新型光交换机 | `make build-arm64` |
| **ARMv7** | 32位 | - | 主流光交换机 | `make build-armv7` |
| **MIPSLE** | 32位 | 小端 | TP-Link、Netgear 路由器 | `make build-mipsle` |
| **MIPS** | 32位 | 大端 | 企业级交换机 | `make build-mips` |
| **MIPS64** | 64位 | 大端 | 电信级设备 | `make build-mips64` |
| **MIPS64LE** | 64位 | 小端 | 新型路由器 | `make build-mips64le` |

### 3. 静态编译

```bash
# 纯 Go SQLite 驱动 (glebarez/sqlite)
# 无需 CGO，完全静态编译
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build

# 生成的二进制文件可直接在目标设备运行
# 无需安装任何运行时依赖
./opt-switch-arm64 server -c config/settings.yml
```

### 4. 企业级权限管理

- **基于 Casbin 的 RBAC 权限控制**
- **多层权限体系**：菜单权限、操作权限、按钮权限、数据权限、API 接口权限
- **JWT 认证**：无状态 Token 认证机制
- **完整审计日志**：操作日志、登录日志

### 5. Web 管理界面

- 用户管理、角色管理、部门管理
- 菜单管理、字典管理、参数配置
- 操作日志、登录日志
- 代码生成器、表单构建器

---

## 快速开始

### 硬件要求

| 资源 | 最低要求 | 推荐配置 |
|------|----------|----------|
| 内存 | 256MB | 512MB |
| 存储 | 100MB | 500MB |
| CPU | 600MHz 单核 | 1GHz 双核 |

### 步骤 1：确定设备架构

在交换机上执行（如有 SSH 访问）：

```bash
uname -m
```

常见输出：
- `armv7l` → 使用 ARMv7 二进制
- `aarch64` → 使用 ARM64 二进制
- `mips` → 使用 MIPS 二进制
- `mipsle` → 使用 MIPSLE 二进制

### 步骤 2：交叉编译

在开发机上编译对应架构的二进制文件：

```bash
# ARMv7（最常见的交换机架构）
make build-armv7

# ARM64（新型交换机）
make build-arm64

# MIPSLE（部分路由器）
make build-mipsle

# MIPS（企业级交换机）
make build-mips

# 或一次编译所有常见架构
make build-switch
```

### 步骤 3：部署到交换机

```bash
# 上传到交换机
scp opt-switch-armv7 admin@<switch-ip>:/tmp/

# SSH 登录交换机
ssh admin@<switch-ip>

# 在交换机上操作
cd /tmp
chmod +x opt-switch-armv7

# 创建工作目录
mkdir -p /opt/opt-switch/config
mkdir -p /opt/opt-switch/data
mv opt-switch-armv7 /opt/opt-switch/

# 复制配置文件（使用低内存配置）
# 先在本地编辑好配置文件，然后上传
scp config/settings.minimal.yml admin@<switch-ip>:/opt/opt-switch/config/settings.yml

# 启动服务
/opt/opt-switch/opt-switch-armv7 server -c /opt/opt-switch/config/settings.yml
```

### 步骤 4：访问管理界面

打开浏览器访问：`http://<switch-ip>:8000`

默认账号密码：`admin / 123456`

**首次登录后请务必修改密码！**

---

## 配置说明

### 配置文件选择

| 配置文件 | 内存占用 | 适用设备 | 并发用户 |
|----------|----------|----------|----------|
| `config/settings.minimal.yml` | 50-70MB | 256MB 内存交换机 | 1-5 人 |
| `config/settings.switch.yml` | 60-100MB | 512MB 内存交换机 | 5-10 人 |

### settings.minimal.yml（极简配置）

适用于 256MB 内存的光交换机。

```yaml
# 目标内存占用: 50-70MB
# 适用场景: 256MB 内存设备，1-5 个并发用户

settings:
  application:
    mode: prod              # 生产模式
    host: 0.0.0.0          # 监听所有接口
    port: 8000             # 服务端口
    enabledp: false        # 禁用数据权限

  logger:
    path: /tmp/opt-switch/logs    # 日志路径
    stdout: '1'                   # 输出到控制台
    level: error                  # 仅记录错误日志
    enableddb: false              # 禁用数据库日志

  jwt:
    secret: opt-switch     # Token 密钥（生产环境务必修改）
    timeout: 3600          # Token 过期时间（秒）

  database:
    driver: sqlite3
    source: /opt/opt-switch/data.db
    maxOpenConns: 2       # 最大打开连接数
    maxIdleConns: 1       # 最大空闲连接数
    connMaxLifetime: 300
    connMaxIdleTime: 60

  queue:
    memory:
      poolSize: 5         # 内存队列池大小

extend:
  runtime:
    gomaxprocs: 1         # 单核模式
    gogc: 200             # 减少 GC 频率
    memoryLimit: 60       # 软内存限制 60MB

  applicationEx:
    enableFrontend: true  # 启用前端
    enableMiddleware:
      sentinel: false     # 禁用限流
      requestID: true
      metrics: false      # 禁用监控
```

### settings.switch.yml（标准配置）

适用于 512MB 内存的光交换机。

```yaml
# 目标内存占用: 60-100MB
# 适用场景: 512MB 内存设备，5-10 个并发用户

settings:
  application:
    mode: prod
    host: 0.0.0.0
    port: 8000
    enabledp: false

  logger:
    path: /tmp/opt-switch/logs
    stdout: '1'
    level: warn           # 警告级别
    enableddb: false

  jwt:
    secret: opt-switch
    timeout: 3600

  database:
    driver: sqlite3
    source: /opt/opt-switch/data.db
    maxOpenConns: 5       # 适中的连接数
    maxIdleConns: 2
    connMaxLifetime: 300
    connMaxIdleTime: 60

  queue:
    memory:
      poolSize: 20        # 适中的队列大小
```

---

## 设置开机自启

### 使用 systemd（推荐）

```bash
# 创建 systemd 服务文件
cat > /etc/systemd/system/opt-switch.service << 'EOF'
[Unit]
Description=opt-switch Management System
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/opt/opt-switch
ExecStart=/opt/opt-switch/opt-switch-armv7 server -c /opt/opt-switch/config/settings.yml
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF

# 启用并启动服务
systemctl daemon-reload
systemctl enable opt-switch
systemctl start opt-switch

# 查看服务状态
systemctl status opt-switch

# 查看日志
journalctl -u opt-switch -f
```

### 使用 init.d（传统方式）

```bash
# 创建 init 脚本
cat > /etc/init.d/opt-switch << 'EOF'
#!/bin/sh

case "$1" in
start)
    /opt/opt-switch/opt-switch-armv7 server -c /opt/opt-switch/config/settings.yml &
    echo "opt-switch started"
    ;;
stop)
    killall opt-switch-armv7
    echo "opt-switch stopped"
    ;;
restart)
    $0 stop
    sleep 2
    $0 start
    ;;
*)
    echo "Usage: $0 {start|stop|restart}"
    exit 1
    ;;
esac

exit 0
EOF

chmod +x /etc/init.d/opt-switch

# 添加开机自启（根据系统选择）
ln -s /etc/init.d/opt-switch /etc/rc.d/S99opt-switch

# 启动服务
/etc/init.d/opt-switch start
```

---

## Makefile 编译命令

### 交换机常用架构

```bash
# 编译最常见的三种交换机架构
make build-switch

# ARM 系列（推荐）
make build-arm64    # ARM64 (新型交换机)
make build-armv7    # ARMv7 (主流交换机)

# MIPS 系列
make build-mips     # MIPS 大端
make build-mipsle   # MIPS 小端（常见路由器）
make build-mips64   # MIPS64 大端
make build-mips64le # MIPS64 小端
```

### 查看所有支持的架构

```bash
make list-arch
```

---

## API 文档

启动服务后，访问 Swagger 文档：

```
http://<switch-ip>:8000/swagger/admin/index.html
```

### 主要 API 端点

#### 认证相关
| 方法 | 端点 | 描述 |
|------|------|------|
| POST | /api/v1/base/login | 用户登录 |
| POST | /api/v1/base/logout | 退出登录 |
| GET | /api/v1/base/info | 获取当前用户信息 |

#### 用户管理
| 方法 | 端点 | 描述 |
|------|------|------|
| GET | /api/v1/user/list | 用户列表 |
| POST | /api/v1/user | 创建用户 |
| PUT | /api/v1/user/:id | 更新用户 |
| DELETE | /api/v1/user/:id | 删除用户 |

#### 角色管理
| 方法 | 端点 | 描述 |
|------|------|------|
| GET | /api/v1/role/list | 角色列表 |
| POST | /api/v1/role | 创建角色 |
| PUT | /api/v1/role/:id | 更新角色 |
| DELETE | /api/v1/role/:id | 删除角色 |

---

## 常见问题

### 1. 内存占用过高

**问题**：运行后内存占用超过预期

**解决方案**：

1. 使用极简配置
```bash
/opt/opt-switch/opt-switch-armv7 server -c /opt/opt-switch/config/settings.minimal.yml
```

2. 检查数据库连接池配置
```yaml
database:
  maxOpenConns: 2  # 降低最大连接数
  maxIdleConns: 1  # 降低空闲连接数
```

### 2. 无法连接到交换机

**问题**：部署后无法访问 Web 界面

**解决方案**：

1. 检查防火墙设置
```bash
# 在交换机上开放端口
iptables -A INPUT -p tcp --dport 8000 -j ACCEPT
```

2. 检查服务是否正常运行
```bash
ps aux | grep opt-switch
netstat -tlnp | grep 8000
```

3. 检查日志
```bash
tail -f /tmp/opt-switch/logs/opt-switch.log
```

### 3. SQLite 数据库锁定

**问题**：多用户访问时数据库锁定

**解决方案**：

1. 启用 WAL 模式（在代码中已默认启用）

2. 考虑降低并发用户数

3. 如果并发需求高，建议使用 512MB 内存配置

### 4. 如何备份数据

**问题**：如何备份交换机上的配置和数据

**解决方案**：

```bash
# 备份数据库
cp /opt/opt-switch/data.db /backup/data-$(date +%Y%m%d).db

# 备份配置文件
cp /opt/opt-switch/config/settings.yml /backup/settings-$(date +%Y%m%d).yml

# 自动备份脚本
cat > /etc/cron.daily/opt-switch-backup << 'EOF'
#!/bin/bash
BACKUP_DIR="/backup/opt-switch"
mkdir -p $BACKUP_DIR
cp /opt/opt-switch/data.db $BACKUP_DIR/data-$(date +\%Y\%m\%d).db
# 保留最近 7 天的备份
find $BACKUP_DIR -name "data-*.db" -mtime +7 -delete
EOF

chmod +x /etc/cron.daily/opt-switch-backup
```

---

## 安全建议

### 生产环境安全检查清单

部署到生产环境前，请务必完成以下检查：

- [ ] 修改默认 JWT secret
- [ ] 修改默认管理员密码
- [ ] 启用防火墙，限制管理端口访问
- [ ] 定期备份数据库
- [ ] 配置日志轮转
- [ ] 限制登录失败次数

### 修改 JWT Secret

```yaml
# 生成随机密钥
openssl rand -hex 32

# 修改配置文件
jwt:
  secret: <your-random-secret-key>
  timeout: 3600
```

### 配置防火墙

```bash
# 仅允许特定 IP 访问管理界面
iptables -A INPUT -p tcp -s 192.168.1.0/24 --dport 8000 -j ACCEPT
iptables -A INPUT -p tcp --dport 8000 -j DROP

# 保存规则
iptables-save > /etc/iptables.rules
```

---

## 技术支持

### 文档资源
- [项目主页](https://www.opt-switch.dev)
- [API 文档](http://<switch-ip>:8000/swagger/admin/index.html)
- [视频教程](https://space.bilibili.com/565616721/channel/detail?cid=125737)

### 获取帮助
- **GitHub Issues**: [提交问题](https://github.com/opt-switch-team/opt-switch/issues)
- **GitHub Discussions**: [参与讨论](https://github.com/opt-switch-team/opt-switch/discussions)

### 社区

<table>
  <tr>
    <td><img src="https://raw.githubusercontent.com/wenjianzhang/image/master/img/wx.png" width="120px"></td>
    <td><img src="https://raw.githubusercontent.com/wenjianzhang/image/master/img/qq2.png" width="120px"></td>
  </tr>
  <tr>
    <td>微信</td>
    <td>QQ群</td>
  </tr>
</table>

---

## 许可证

本项目采用 [MIT License](LICENSE.md) 开源协议。

---

<div align="center">

**如果这个项目对您有帮助，请给我们一个 Star ⭐**

Made with ❤️ by opt-switch team

</div>

---

## English

## Introduction

**opt-switch** is an ultra-lightweight backend management system designed specifically for **optical switches, routers, and network devices**. Unlike traditional permission management systems that require 200-500MB of memory, opt-switch runs on just **50-100MB of memory** through deep optimization, allowing it to be deployed directly on switch devices.

### Why opt-switch

| Feature | Traditional Solution | opt-switch |
|---------|---------------------|------------|
| **Memory Usage** | 200-500MB | **50-100MB** |
| **Deployment** | Separate Server | **Directly on Switch** |
| **Database** | MySQL/PostgreSQL Server | **SQLite Built-in** |
| **Management** | Serial/CLI | **Web UI** |
| **Startup Time** | 10-30 seconds | **2-5 seconds** |

### Supported Devices

| Device Type | Brands | Memory | Recommended Config |
|-------------|--------|--------|-------------------|
| Optical Switch | Huawei, H3C, Ruijie | 256MB | `config/settings.minimal.yml` |
| Optical Switch | Huawei, H3C, Ruijie | 512MB | `config/settings.switch.yml` |
| Router | TP-Link, Netgear | 256MB | `config/settings.minimal.yml` |

---

## Core Features

### 1. Ultra-Low Memory Footprint
- **50-70MB** Minimal config (256MB devices)
- **60-100MB** Standard config (512MB devices)

### 2. Multi-Architecture Support

| Architecture | Typical Devices | Build Command |
|--------------|-----------------|---------------|
| **ARM64** | New optical switches | `make build-arm64` |
| **ARMv7** | Most optical switches | `make build-armv7` |
| **MIPSLE** | TP-Link, Netgear routers | `make build-mipsle` |
| **MIPS** | Enterprise switches | `make build-mips` |

### 3. Static Compilation
```bash
# Pure Go SQLite driver, no CGO required
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build
```

### 4. Enterprise Permission Management
- Casbin-based RBAC
- JWT authentication
- Complete audit logs

---

## Quick Start

### Hardware Requirements

| Resource | Minimum | Recommended |
|----------|---------|-------------|
| Memory | 256MB | 512MB |
| Storage | 100MB | 500MB |
| CPU | 600MHz single-core | 1GHz dual-core |

### Step 1: Determine Device Architecture

```bash
uname -m
```

Common outputs:
- `aarch64` → Use ARM64 binary
- `armv7l` → Use ARMv7 binary
- `mips` → Use MIPS binary
- `mipsle` → Use MIPSLE binary

### Step 2: Cross-Compile

```bash
# ARMv7 (most common)
make build-armv7

# ARM64 (new switches)
make build-arm64

# MIPSLE (routers)
make build-mipsle

# Or build all common architectures
make build-switch
```

### Step 3: Deploy to Switch

```bash
# Upload to switch
scp opt-switch-armv7 admin@<switch-ip>:/tmp/

# SSH login
ssh admin@<switch-ip>

# On switch
cd /tmp
chmod +x opt-switch-armv7
mkdir -p /opt/opt-switch/config
mkdir -p /opt/opt-switch/data
mv opt-switch-armv7 /opt/opt-switch/

# Upload config file
scp config/settings.minimal.yml admin@<switch-ip>:/opt/opt-switch/config/settings.yml

# Start service
/opt/opt-switch/opt-switch-armv7 server -c /opt/opt-switch/config/settings.yml
```

### Step 4: Access Web UI

Browser: `http://<switch-ip>:8000`

Default credentials: `admin / 123456`

**Change the password after first login!**

---

## Configuration

### Configuration Files

| Config File | Memory | Device | Users |
|-------------|--------|--------|-------|
| `config/settings.minimal.yml` | 50-70MB | 256MB switches | 1-5 |
| `config/settings.switch.yml` | 60-100MB | 512MB switches | 5-10 |

See the Chinese section for detailed configuration examples.

---

## Auto-Start on Boot

### Using systemd (Recommended)

```bash
cat > /etc/systemd/system/opt-switch.service << 'EOF'
[Unit]
Description=opt-switch Management System
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/opt/opt-switch
ExecStart=/opt/opt-switch/opt-switch-armv7 server -c /opt/opt-switch/config/settings.yml
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable opt-switch
systemctl start opt-switch
```

---

## Makefile Commands

```bash
# Build common switch architectures
make build-switch

# ARM builds
make build-arm64    # ARM64
make build-armv7    # ARMv7

# MIPS builds
make build-mips     # MIPS big-endian
make build-mipsle   # MIPS little-endian

# List all architectures
make list-arch
```

---

## API Documentation

After starting the service, access Swagger docs:

```
http://<switch-ip>:8000/swagger/admin/index.html
```

---

## FAQ

See the Chinese section for common issues and solutions.

---

## Security Checklist

Before deploying to production:

- [ ] Change default JWT secret
- [ ] Change default admin password
- [ ] Configure firewall rules
- [ ] Set up database backups
- [ ] Configure log rotation
- [ ] Limit login attempts

---

## Support

### Resources
- [Project Homepage](https://www.opt-switch.dev)
- [API Documentation](http://<switch-ip>:8000/swagger/admin/index.html)
- [Video Tutorials](https://space.bilibili.com/565616721/channel/detail?cid=125737)

### Get Help
- **GitHub Issues**: [Submit Issues](https://github.com/opt-switch-team/opt-switch/issues)
- **GitHub Discussions**: [Join Discussion](https://github.com/opt-switch-team/opt-switch/discussions)

---

## License

This project is licensed under the [MIT License](LICENSE.md).

---

<div align="center">

**If this project helps you, please give us a Star ⭐**

Made with ❤️ by opt-switch team

</div>
