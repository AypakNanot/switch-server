# opt-switch 光交换机管理系统架构设计文档

## 文档信息

| 项目名称 | opt-switch |
|---------|-----------|
| 文档版本 | 1.0.0 |
| 创建日期 | 2025-01-25 |
| 文档作者 | 架构设计团队 |
| 项目地址 | https://github.com/opt-switch-team/opt-switch |

---

## 目录

1. [项目概述](#1-项目概述)
2. [系统架构](#2-系统架构)
3. [技术栈](#3-技术栈)
4. [目录结构](#4-目录结构)
5. [核心模块设计](#5-核心模块设计)
6. [数据库设计](#6-数据库设计)
7. [API设计](#7-api设计)
8. [设备管理模块](#8-设备管理模块)
9. [权限管理系统](#9-权限管理系统)
10. [部署架构](#10-部署架构)
11. [性能优化](#11-性能优化)
12. [安全设计](#12-安全设计)
13. [可扩展性设计](#13-可扩展性设计)

---

## 1. 项目概述

### 1.1 项目简介

**opt-switch** 是一款专为光交换机、路由器等网络设备设计的超轻量级后台管理系统。与传统需要 200-500MB 内存的权限管理系统不同，opt-switch 通过深度优化，仅需 **50-100MB 内存**即可运行，可直接部署在交换机设备上。

### 1.2 设计目标

| 设计目标 | 说明 |
|---------|------|
| **超低内存占用** | 50-100MB 内存占用，适合资源受限的交换机设备 |
| **静态编译** | 无需 CGO，纯 Go 实现，支持多架构交叉编译 |
| **零外部依赖** | 使用 SQLite 内置数据库，无需独立数据库服务器 |
| **多架构支持** | 支持 ARM64、ARMv7、MIPS、MIPS64 等多种 CPU 架构 |
| **企业级权限** | 基于 Casbin 的 RBAC 权限控制系统 |
| **Web 可视化** | 内嵌 Web 管理界面，提供友好的操作体验 |

### 1.3 适用场景

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

---

## 2. 系统架构

### 2.1 整体架构图

```
┌─────────────────────────────────────────────────────────────────┐
│                        opt-switch 系统架构                        │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│                        客户端层 (Client Layer)                    │
├─────────────────────────────────────────────────────────────────┤
│  Web Browser (HTTP/HTTPS)  │  CLI Tools  │  API Clients        │
└─────────────────────────────────────────────────────────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────────┐
│                       Web 服务器层 (Gin Framework)               │
├─────────────────────────────────────────────────────────────────┤
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐        │
│  │ 静态文件  │  │ Swagger  │  │ API 路由 │  │ 中间件   │        │
│  │  服务    │  │  文档    │  │  服务    │  │  链      │        │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘        │
└─────────────────────────────────────────────────────────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────────┐
│                        应用服务层 (Application Layer)             │
├─────────────────────────────────────────────────────────────────┤
│  ┌──────────────────┐  ┌──────────────────┐  ┌───────────────┐  │
│  │   系统管理模块    │  │   设备管理模块    │  │  工具/任务    │  │
│  │  System Module   │  │  Device Module   │  │  Tools/Jobs   │  │
│  ├──────────────────┤  ├──────────────────┤  ├───────────────┤  │
│  │ • 用户管理        │  │ • SSH/Telnet     │  │ • 代码生成    │  │
│  │ • 角色管理        │  │ • 连接池         │  │ • 定时任务    │  │
│  │ • 菜单管理        │  │ • 命令执行       │  │ • 文件上传    │  │
│  │ • 权限管理        │  │ • 日志记录       │  │ • 系统监控    │  │
│  │ • 字典管理        │  │ • 配置管理       │  │               │  │
│  └──────────────────┘  └──────────────────┘  └───────────────┘  │
└─────────────────────────────────────────────────────────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────────┐
│                         业务逻辑层 (Service Layer)                │
├─────────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐              │
│  │  DTO 对象    │  │  业务服务    │  │  数据验证    │              │
│  │  (数据传输)  │  │  (Business) │  │ (Validator) │              │
│  └─────────────┘  └─────────────┘  └─────────────┘              │
└─────────────────────────────────────────────────────────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────────┐
│                        数据访问层 (Data Access Layer)             │
├─────────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐              │
│  │   GORM      │  │   Casbin    │  │  Redis/     │              │
│  │  (ORM)      │  │  (权限)     │  │  Memory     │              │
│  └─────────────┘  └─────────────┘  └─────────────┘              │
└─────────────────────────────────────────────────────────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────────┐
│                        数据存储层 (Storage Layer)                 │
├─────────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐              │
│  │  SQLite     │  │  File Store │  │  Queue      │              │
│  │ (主数据库)  │  │ (文件存储)   │  │ (内存队列)   │              │
│  └─────────────┘  └─────────────┘  └─────────────┘              │
└─────────────────────────────────────────────────────────────────┘
```

### 2.2 分层架构说明

| 层级 | 名称 | 职责 | 技术实现 |
|------|------|------|----------|
| **客户端层** | Client Layer | 用户交互界面 | Web Browser、CLI、API Client |
| **Web服务层** | Presentation | HTTP服务、路由、中间件 | Gin Framework |
| **应用服务层** | Application | 业务功能模块 | System、Device、Tools |
| **业务逻辑层** | Service | 具体业务逻辑、DTO处理 | Service Interface + DTO |
| **数据访问层** | Data Access | 数据库操作、权限控制 | GORM、Casbin |
| **数据存储层** | Storage | 数据持久化 | SQLite、File、Memory Queue |

---

## 3. 技术栈

### 3.1 后端技术栈

| 技术/框架 | 版本 | 用途 | 说明 |
|-----------|------|------|------|
| **Go** | 1.24+ | 编程语言 | 高性能、静态编译 |
| **Gin** | 1.10.0 | Web框架 | 高性能HTTP框架 |
| **GORM** | 1.25.12 | ORM框架 | 数据库操作 |
| **SQLite** | - | 数据库 | 纯Go驱动(glebarez/sqlite)，无需CGO |
| **Casbin** | 2.104.0 | 权限控制 | RBAC权限模型 |
| **JWT** | - | 身份认证 | golang-jwt/jwt |
| **Cobra** | 1.9.1 | CLI框架 | 命令行工具 |
| **Swagger** | 1.16.4 | API文档 | 自动生成API文档 |
| **Zap** | 1.27.0 | 日志框架 | 高性能结构化日志 |
| **Cron** | 3.0.1 | 定时任务 | robfig/cron |
| **Sentinel** | 1.0.4 | 限流熔断 | 阿里开源 |

### 3.2 支持的CPU架构

| 架构 | 字长 | 端序 | 典型设备 | 编译命令 |
|------|------|------|----------|----------|
| **ARM64** | 64位 | - | 新型光交换机 | `make build-arm64` |
| **ARMv7** | 32位 | - | 主流光交换机 | `make build-armv7` |
| **MIPSLE** | 32位 | 小端 | TP-Link、Netgear 路由器 | `make build-mipsle` |
| **MIPS** | 32位 | 大端 | 企业级交换机 | `make build-mips` |
| **MIPS64** | 64位 | 大端 | 电信级设备 | `make build-mips64` |
| **MIPS64LE** | 64位 | 小端 | 新型路由器 | `make build-mips64le` |

### 3.3 第三方依赖

```
github.com/gin-gonic/gin
github.com/glebarez/sqlite           # 纯Go SQLite驱动，无需CGO
github.com/casbin/casbin/v2          # 权限控制
github.com/go-admin-team/go-admin-core/sdk
github.com/spf13/cobra               # CLI框架
github.com/swaggo/gin-swagger        # Swagger文档
github.com/robfig/cron/v3            # 定时任务
github.com/alibaba/sentinel-golang   # 限流熔断
go.uber.org/zap                      # 日志框架
gopkg.in/natefinch/lumberjack.v2     # 日志轮转
```

---

## 4. 目录结构

```
switch-server/
├── app/                              # 应用模块目录
│   ├── admin/                        # 系统管理模块
│   │   ├── apis/                     # API 控制器
│   │   │   ├── sys_user.go          # 用户API
│   │   │   ├── sys_role.go          # 角色API
│   │   │   ├── sys_menu.go          # 菜单API
│   │   │   ├── sys_dept.go          # 部门API
│   │   │   ├── sys_dict_data.go     # 字典数据API
│   │   │   ├── sys_config.go        # 配置API
│   │   │   └── ...
│   │   ├── models/                   # 数据模型
│   │   │   ├── sys_user.go
│   │   │   ├── sys_role.go
│   │   │   ├── casbin_rule.go       # Casbin规则模型
│   │   │   └── ...
│   │   ├── service/                  # 业务服务层
│   │   │   ├── dto/                 # 数据传输对象
│   │   │   │   ├── sys_user.go
│   │   │   │   ├── sys_role.go
│   │   │   │   └── ...
│   │   │   ├── sys_user.go          # 用户服务
│   │   │   ├── sys_role.go          # 角色服务
│   │   │   └── ...
│   │   └── router/                   # 路由定义
│   │       ├── router.go            # 路由初始化
│   │       ├── sys_user.go
│   │       ├── sys_role.go
│   │       └── ...
│   │
│   ├── device/                       # 设备管理模块 (扩展功能)
│   │   ├── apis/
│   │   │   └── command.go           # 命令执行API
│   │   ├── service/
│   │   │   ├── dto/
│   │   │   └── command_service.go   # 命令服务
│   │   └── router/
│   │       └── device_router.go     # 设备路由
│   │
│   ├── jobs/                         # 定时任务模块
│   │   ├── models/
│   │   ├── service/
│   │   └── router/
│   │
│   └── other/                        # 其他功能模块
│       ├── apis/
│       │   ├── file.go              # 文件上传API
│       │   ├── tools/               # 代码生成工具API
│       │   └── sys_server_monitor.go # 服务器监控API
│       ├── models/
│       └── router/
│
├── cmd/                              # 命令行工具目录
│   ├── api/                          # API服务启动
│   ├── app/                          # 应用生成工具
│   ├── config/                       # 配置查看工具
│   ├── migrate/                      # 数据库迁移工具
│   ├── cobra.go                      # Cobra命令定义
│   └── ...
│
├── common/                           # 公共模块目录
│   ├── middleware/                   # 中间件
│   │   ├── auth.go                  # JWT认证中间件
│   │   ├── database.go              # 数据库中间件
│   │   ├── error.go                 # 错误处理中间件
│   │   └── ...
│   ├── database/                     # 数据库初始化
│   ├── models/                       # 公共模型
│   ├── actions/                      # 公共动作
│   ├── dto/                          # 公共DTO
│   ├── global/                       # 全局变量
│   └── response/                     # 响应封装
│
├── config/                           # 配置文件目录
│   ├── settings.yml                  # 主配置文件
│   ├── settings.minimal.yml          # 极简配置 (256MB设备)
│   ├── settings.switch.yml           # 标准配置 (512MB设备)
│   └── extend.go                     # 扩展配置定义
│
├── pkg/                              # 自定义包目录
│   └── device/                       # 设备管理包
│       ├── adapter.go                # 设备适配器接口
│       ├── pool.go                   # 连接池
│       ├── ssh.go                    # SSH连接
│       ├── telnet.go                 # Telnet连接
│       ├── config.go                 # 设备配置
│       ├── logger.go                 # 设备日志
│       └── init.go                   # 初始化
│
├── deploy/                           # 部署相关
│   └── ...
│
├── docs/                             # 文档目录
│   └── ...
│
├── static/                           # 静态资源 (内嵌Web UI)
│   └── ...
│
├── template/                         # 代码模板
│   └── ...
│
├── test/                             # 测试文件
│   └── ...
│
├── web/                              # Web前端源码
│   └── ...
│
├── .github/                          # GitHub配置
│   └── workflows/                    # CI/CD工作流
│
├── main.go                           # 程序入口
├── go.mod                            # Go模块定义
├── go.sum                            # Go依赖校验
├── Makefile                          # 编译脚本
├── Dockerfile                        # Docker镜像构建
└── README.md                         # 项目说明
```

---

## 5. 核心模块设计

### 5.1 系统管理模块 (app/admin)

系统管理模块是 opt-switch 的核心基础模块，提供完整的后台管理功能。

#### 5.1.1 模块组成

```
app/admin/
├── apis/          # API控制器层 - 处理HTTP请求
├── models/        # 数据模型层 - 定义数据库表结构
├── service/       # 业务逻辑层 - 核心业务处理
│   └── dto/      # 数据传输对象 - 请求/响应数据定义
└── router/        # 路由层 - URL路由定义
```

#### 5.1.2 功能列表

| 功能模块 | API路径 | 主要功能 |
|----------|---------|----------|
| **用户管理** | /api/v1/user | 用户增删改查、密码修改、用户状态管理 |
| **角色管理** | /api/v1/role | 角色增删改查、角色权限分配 |
| **菜单管理** | /api/v1/menu | 菜单树管理、按钮权限配置 |
| **部门管理** | /api/v1/dept | 部门树管理、部门人员管理 |
| **岗位管理** | /api/v1/post | 岗位增删改查 |
| **字典管理** | /api/v1/dict | 字典类型、字典数据管理 |
| **参数配置** | /api/v1/config | 系统参数配置管理 |
| **操作日志** | /api/v1/opera-log | 操作日志查询 |
| **登录日志** | /api/v1/login-log | 登录日志查询 |
| **API管理** | /api/v1/api | API接口管理 |

#### 5.1.3 数据模型关系

```
┌─────────────┐         ┌─────────────┐         ┌─────────────┐
│   sys_user  │────────▶│  sys_role   │────────▶│  sys_menu   │
│   用户表     │  多对多  │   角色表     │  多对多  │   菜单表     │
└─────────────┘         └─────────────┘         └─────────────┘
       │                       │
       ▼                       ▼
┌─────────────┐         ┌─────────────┐
│  sys_dept   │         │casbin_rule  │
│   部门表     │         │  权限规则表   │
└─────────────┘         └─────────────┘
       │
       ▼
┌─────────────┐
│  sys_post   │
│   岗位表     │
└─────────────┘
```

### 5.2 设备管理模块 (app/device + pkg/device)

设备管理模块是 opt-switch 的特色扩展功能，用于管理本地光交换机设备。

#### 5.2.1 模块架构

```
pkg/device/                         # 设备管理核心包
├── adapter.go                      # 设备适配器接口定义
├── pool.go                         # 连接池实现
├── ssh.go                          # SSH连接实现
├── telnet.go                       # Telnet连接实现
├── config.go                       # 设备配置
├── logger.go                       # 设备操作日志
├── init.go                         # 初始化逻辑
└── error.go                        # 错误定义

app/device/                         # 设备管理API模块
├── apis/
│   └── command.go                  # 命令执行API
├── service/
│   ├── dto/
│   │   └── command.go              # 命令请求/响应DTO
│   └── command_service.go          # 命令服务
└── router/
    └── device_router.go            # 设备路由
```

#### 5.2.2 核心接口

```go
// Adapter 设备适配器接口
type Adapter interface {
    // Connect 建立连接
    Connect() error

    // Execute 执行命令
    Execute(cmd string) (string, error)

    // ExecuteWithTimeout 带超时的命令执行
    ExecuteWithTimeout(cmd string, timeout time.Duration) (string, error)

    // Close 关闭连接
    Close() error

    // IsConnected 检查连接状态
    IsConnected() bool
}
```

#### 5.2.3 连接池设计

```
┌─────────────────────────────────────────────────────┐
│                  DeviceConnectionPool               │
├─────────────────────────────────────────────────────┤
│  max_connections: 3     # 最大连接数                 │
│  min_connections: 1     # 最小连接数                 │
│  idle_timeout: 300s     # 空闲超时                   │
│  command_timeout: 30s   # 命令超时                   │
│  max_queue_size: 100    # 队列大小                   │
└─────────────────────────────────────────────────────┘
          │
          ├─ conn1 (SSH) ─────▶ [光交换机设备]
          │
          ├─ conn2 (SSH) ─────▶ [光交换机设备]
          │
          └─ conn3 (Telnet) ───▶ [光交换机设备]
```

### 5.3 定时任务模块 (app/jobs)

基于 robfig/cron 的定时任务管理模块。

#### 5.3.1 功能特性

- Cron表达式任务调度
- 任务启停控制
- 任务执行日志
- 任务执行历史

### 5.4 工具模块 (app/other)

提供系统开发辅助工具。

#### 5.4.1 代码生成器

- 数据库表结构扫描
- 代码模板生成
- 支持前后端代码生成

#### 5.4.2 系统监控

- CPU、内存、磁盘监控
- 在线用户统计
- 服务器信息查询

---

## 6. 数据库设计

### 6.1 数据库选择

| 数据库 | 说明 | 适用场景 |
|--------|------|----------|
| **SQLite** | 默认，嵌入式数据库 | 交换机设备部署 |
| **MySQL** | 可选 | 高并发场景 |
| **PostgreSQL** | 可选 | 企业级应用 |
| **SQL Server** | 可选 | Windows环境 |

### 6.2 核心数据表

#### 6.2.1 用户权限相关表

| 表名 | 说明 | 主要字段 |
|------|------|----------|
| **sys_user** | 用户表 | user_id, username, password, nick_name, dept_id, role_ids |
| **sys_role** | 角色表 | role_id, role_name, role_key, data_scope |
| **sys_menu** | 菜单表 | menu_id, menu_name, path, component, perms |
| **sys_dept** | 部门表 | dept_id, parent_id, dept_name, ancestors |
| **sys_post** | 岗位表 | post_id, post_code, post_name |
| **sys_user_role** | 用户角色关联表 | user_id, role_id |
| **sys_role_menu** | 角色菜单关联表 | role_id, menu_id |
| **casbin_rule** | Casbin权限规则 | ptype, v0, v1, v2, v3 |

#### 6.2.2 系统配置相关表

| 表名 | 说明 | 主要字段 |
|------|------|----------|
| **sys_dict_data** | 字典数据表 | dict_code, dict_type, dict_label, dict_value |
| **sys_dict_type** | 字典类型表 | dict_id, dict_name, dict_type |
| **sys_config** | 参数配置表 | config_id, config_name, config_key, config_value |
| **sys_opera_log** | 操作日志表 | opera_id, title, business_type, method, request_method |
| **sys_login_log** | 登录日志表 | info_id, username, ipaddr, login_location, browser |

#### 6.2.3 API管理表

| 表名 | 说明 | 主要字段 |
|------|------|----------|
| **sys_api** | API接口表 | api_id, path, method, description, group |

### 6.3 数据库初始化

数据库初始化流程位于 `cmd/migrate/migration/init.go`：

```go
// 初始化流程
1. 连接数据库
2. 自动迁移数据表结构
3. 插入初始管理员账号 (admin/123456)
4. 插入默认角色 (管理员、普通用户)
5. 插入系统菜单
6. 初始化Casbin权限规则
```

---

## 7. API设计

### 7.1 API版本控制

```
基础路径: /api/v1
文档路径: /swagger/admin/index.html
```

### 7.2 认证方式

```
Header: Authorization: Bearer {token}
Token类型: JWT
Token有效期: 3600秒 (可配置)
```

### 7.3 主要API端点

#### 7.3.1 认证相关

| 方法 | 端点 | 描述 | 认证 |
|------|------|------|------|
| POST | /api/v1/base/login | 用户登录 | 否 |
| POST | /api/v1/base/logout | 退出登录 | 是 |
| GET | /api/v1/base/info | 获取当前用户信息 | 是 |
| GET | /api/v1/base/captcha | 获取验证码 | 否 |

#### 7.3.2 用户管理

| 方法 | 端点 | 描述 | 认证 |
|------|------|------|------|
| GET | /api/v1/user/list | 用户列表 | 是 |
| GET | /api/v1/user/:id | 获取用户详情 | 是 |
| POST | /api/v1/user | 创建用户 | 是 |
| PUT | /api/v1/user/:id | 更新用户 | 是 |
| DELETE | /api/v1/user/:id | 删除用户 | 是 |
| PUT | /api/v1/user/:id/password | 修改密码 | 是 |
| PUT | /api/v1/user/:id/status | 修改用户状态 | 是 |

#### 7.3.3 角色管理

| 方法 | 端点 | 描述 | 认证 |
|------|------|------|------|
| GET | /api/v1/role/list | 角色列表 | 是 |
| GET | /api/v1/role/:id | 获取角色详情 | 是 |
| POST | /api/v1/role | 创建角色 | 是 |
| PUT | /api/v1/role/:id | 更新角色 | 是 |
| DELETE | /api/v1/role/:id | 删除角色 | 是 |
| PUT | /api/v1/role/menu | 更新角色菜单 | 是 |

#### 7.3.4 设备管理

| 方法 | 端点 | 描述 | 认证 |
|------|------|------|------|
| POST | /api/v1/device/execute | 执行设备命令 | 是 |
| GET | /api/v1/device/status | 获取设备状态 | 是 |
| GET | /api/v1/device/logs | 获取命令日志 | 是 |

### 7.4 响应格式

```json
// 成功响应
{
    "code": 0,
    "msg": "操作成功",
    "data": { ... }
}

// 错误响应
{
    "code": 400,
    "msg": "请求参数错误",
    "data": null
}
```

---

## 8. 设备管理模块

### 8.1 设备连接架构

```
┌─────────────────────────────────────────────────────────┐
│                    Device Service                       │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  ┌─────────────────────────────────────────────────┐   │
│  │           Connection Pool (连接池)               │   │
│  │  • max_connections: 3                           │   │
│  │  • min_connections: 1                           │   │
│  │  • idle_timeout: 300s                           │   │
│  │  • command_timeout: 30s                         │   │
│  └─────────────────────────────────────────────────┘   │
│                          │                               │
│         ┌────────────────┼────────────────┐             │
│         ▼                ▼                ▼             │
│  ┌──────────┐     ┌──────────┐     ┌──────────┐        │
│  │   SSH    │     │   SSH    │     │  Telnet  │        │
│  │ Adapter  │     │ Adapter  │     │ Adapter  │        │
│  └──────────┘     └──────────┘     └──────────┘        │
│         │                │                │             │
└─────────┼────────────────┼────────────────┼─────────────┘
          │                │                │
          ▼                ▼                ▼
     ┌────────────────────────────────────────┐
     │         光交换机设备                     │
     │    (SSH/Telnet 服务)                    │
     └────────────────────────────────────────┘
```

### 8.2 设备配置

```yaml
# config/settings.yml
extend:
  device:
    connection:
      host: 127.0.0.1
      port: 22
      protocol: ssh          # ssh or telnet
      username: admin
      password: admin
      timeout: 30
    pool:
      max_connections: 3
      min_connections: 1
      idle_timeout: 300
      command_timeout: 30
      queue_timeout: 60
      max_queue_size: 100
    log:
      enabled: true
      file: logs/command.log
      max_size: 100
      max_backups: 3
      max_age: 7
      compress: true
```

### 8.3 设备命令执行流程

```
1. API请求
   │
   ▼
2. JWT认证
   │
   ▼
3. 参数验证 (DTO)
   │
   ▼
4. 从连接池获取连接
   │
   ▼
5. 执行命令
   │
   ▼
6. 记录日志
   │
   ▼
7. 返回结果
   │
   ▼
8. 连接归还池
```

---

## 9. 权限管理系统

### 9.1 权限模型

opt-switch 采用基于 **Casbin** 的 RBAC (Role-Based Access Control) 权限模型。

```
┌─────────────────────────────────────────────────────────────┐
│                      RBAC权限模型                           │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│   ┌─────────┐         ┌─────────┐         ┌─────────┐       │
│   │  User   │────────▶│  Role   │────────▶│  Menu   │       │
│   │  用户    │   N:M   │  角色    │   N:M   │  菜单    │       │
│   └─────────┘         └─────────┘         └─────────┘       │
│      │                                         │            │
│      ▼                                         ▼            │
│   ┌─────────┐                             ┌─────────┐       │
│   │  Dept   │                             │ Permission     │
│   │  部门    │                             │  权限           │
│   └─────────┘                             └─────────┘       │
│      │                                         │            │
│      └─────────────────────────────────────────┘            │
│                          │                                  │
│                          ▼                                  │
│                   ┌─────────┐                               │
│                   │ Casbin  │  (权限规则引擎)                │
│                   │  Rules  │                               │
│                   └─────────┘                               │
└─────────────────────────────────────────────────────────────┘
```

### 9.2 权限维度

| 权限维度 | 说明 | 示例 |
|----------|------|------|
| **菜单权限** | 控制用户可见菜单 | 用户管理、角色管理 |
| **操作权限** | 控制按钮可见性 | 新增、编辑、删除 |
| **数据权限** | 控制数据访问范围 | 全部、本部门、本人 |
| **API权限** | 控制API接口访问 | /api/v1/user:list |

### 9.3 Casbin策略模型

```
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
```

### 9.4 JWT认证流程

```
1. 用户登录
   │
   ├─ 验证用户名密码
   ├─ 查询用户角色和权限
   ├─ 生成JWT Token
   └─ 返回Token给客户端
   │
2. 后续请求
   │
   ├─ Header携带: Authorization: Bearer {token}
   ├─ 中间件解析Token
   ├─ 验证Token有效性
   ├─ 提取用户信息
   ├─ Casbin权限校验
   └─ 通过/拒绝请求
```

---

## 10. 部署架构

### 10.1 独立部署架构

```
┌────────────────────────────────────────────────────┐
│            光交换机 / 路由器设备                     │
├────────────────────────────────────────────────────┤
│                                                     │
│  ┌──────────────────────────────────────────────┐  │
│  │         opt-switch 应用                      │  │
│  │                                              │  │
│  │  ┌──────────────┐  ┌──────────────┐         │  │
│  │  │ Web UI       │  │ API Server   │         │  │
│  │  │ (内嵌静态)    │  │  (Gin)       │         │  │
│  │  └──────────────┘  └──────────────┘         │  │
│  │                                              │  │
│  │  ┌──────────────┐  ┌──────────────┐         │  │
│  │  │ SQLite DB    │  │ 设备管理     │         │  │
│  │  │ (本地文件)    │  │ (SSH/Telnet) │         │  │
│  │  └──────────────┘  └──────────────┘         │  │
│  │                                              │  │
│  └──────────────────────────────────────────────┘  │
│                                                     │
│  内存占用: 50-100MB                                 │
│  存储占用: ~100MB                                   │
└────────────────────────────────────────────────────┘
           ▲
           │ HTTP/HTTPS
           │
┌──────────┴──────────┐
│   管理员 PC           │
│   (Web 浏览器)        │
└──────────────────────┘
```

### 10.2 编译与部署

#### 10.2.1 交叉编译

```bash
# ARMv7 (主流光交换机)
make build-armv7

# ARM64 (新型交换机)
make build-arm64

# MIPSLE (路由器)
make build-mipsle

# 一次性编译所有常见架构
make build-switch
```

#### 10.2.2 部署步骤

```bash
# 1. 上传二进制文件
scp opt-switch-armv7 admin@<switch-ip>:/tmp/

# 2. SSH登录交换机
ssh admin@<switch-ip>

# 3. 创建工作目录
mkdir -p /opt/opt-switch/config
mkdir -p /opt/opt-switch/data

# 4. 移动文件
mv /tmp/opt-switch-armv7 /opt/opt-switch/
chmod +x /opt/opt-switch/opt-switch-armv7

# 5. 上传配置文件
scp config/settings.minimal.yml admin@<switch-ip>:/opt/opt-switch/config/settings.yml

# 6. 启动服务
/opt/opt-switch/opt-switch-armv7 server -c /opt/opt-switch/config/settings.yml
```

#### 10.2.3 系统服务配置

**systemd 方式 (推荐)**

```ini
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
```

### 10.3 Docker部署

```dockerfile
FROM alpine:latest

COPY opt-switch-arm64 /opt/opt-switch/
COPY config/settings.yml /opt/opt-switch/config/

EXPOSE 8000

CMD ["/opt/opt-switch/opt-switch-arm64", "server", "-c", "/opt/opt-switch/config/settings.yml"]
```

---

## 11. 性能优化

### 11.1 内存优化

| 配置项 | 极简配置 (256MB设备) | 标准配置 (512MB设备) |
|--------|---------------------|---------------------|
| **目标内存** | 50-70MB | 60-100MB |
| **GOMAXPROCS** | 1 | 2 |
| **GOGC** | 200 | 100 |
| **内存限制** | 60MB | 100MB |
| **DB最大连接** | 2 | 5 |
| **DB空闲连接** | 1 | 2 |
| **队列池大小** | 5 | 20 |
| **日志级别** | error | warn |

### 11.2 优化技术

#### 11.2.1 静态编译

```bash
# 纯 Go SQLite 驱动 (glebarez/sqlite)
# 无需 CGO，完全静态编译
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build
```

#### 11.2.2 内存控制

```go
// 运行时内存限制
runtime.Debug.SetGCPercent(200)  // 调整GC频率
debug.SetMemoryLimit(60 * 1024 * 1024)  // 软内存限制
```

#### 11.2.3 数据库连接池

```yaml
database:
  maxOpenConns: 2      # 降低最大连接数
  maxIdleConns: 1      # 降低空闲连接数
  connMaxLifetime: 300
  connMaxIdleTime: 60
```

#### 11.2.4 SQLite WAL模式

```go
// 启用WAL模式，提高并发性能
db.Exec("PRAGMA journal_mode=WAL")
db.Exec("PRAGMA synchronous=NORMAL")
db.Exec("PRAGMA cache_size=-2000")  # 2MB缓存
```

### 11.3 性能指标

| 指标 | 数值 | 说明 |
|------|------|------|
| **启动时间** | 2-5秒 | 含数据库初始化 |
| **内存占用** | 50-100MB | 根据配置而定 |
| **并发用户** | 1-10人 | 根据设备配置 |
| **API响应** | <100ms | 本地API调用 |

---

## 12. 安全设计

### 12.1 安全特性

| 安全特性 | 实现方式 |
|----------|----------|
| **身份认证** | JWT Token |
| **权限控制** | Casbin RBAC |
| **密码加密** | BCrypt |
| **SQL注入防护** | GORM参数化查询 |
| **XSS防护** | 输入过滤、输出转义 |
| **CSRF防护** | Token验证 |
| **限流保护** | Sentinel |
| **审计日志** | 操作日志、登录日志 |

### 12.2 JWT配置

```yaml
jwt:
  secret: opt-switch     # 生产环境务必修改
  timeout: 3600          # Token过期时间(秒)
```

### 12.3 密码策略

```go
// 默认密码策略
- 最小长度: 6位
- 加密方式: BCrypt
- 默认密码: 123456 (首次登录后必须修改)
```

### 12.4 安全建议

**生产环境部署前检查清单**

- [ ] 修改默认 JWT secret
- [ ] 修改默认管理员密码
- [ ] 启用防火墙，限制管理端口访问
- [ ] 定期备份数据库
- [ ] 配置日志轮转
- [ ] 限制登录失败次数
- [ ] 启用 HTTPS (如有条件)

---

## 13. 可扩展性设计

### 13.1 模块化设计

```
opt-switch 采用模块化设计，各模块独立可扩展：

app/
├── admin/      # 系统管理模块 (核心)
├── device/     # 设备管理模块 (扩展)
├── jobs/       # 定时任务模块 (扩展)
└── other/      # 其他功能模块 (扩展)
```

### 13.2 设备适配器扩展

```go
// Adapter 接口定义，支持扩展新的设备类型
type Adapter interface {
    Connect() error
    Execute(cmd string) (string, error)
    ExecuteWithTimeout(cmd string, timeout time.Duration) (string, error)
    Close() error
    IsConnected() bool
}

// 可扩展实现:
// - SSH Adapter (已实现)
// - Telnet Adapter (已实现)
// - SNMP Adapter (待扩展)
// - HTTP API Adapter (待扩展)
```

### 13.3 中间件扩展

```go
// 支持自定义中间件
func InitRouter() {
    // 现有中间件
    r.Use(middleware.Cors())
    r.Use(middleware.RequestID())
    r.Use(middleware.Auth())

    // 可扩展自定义中间件
    // r.Use(customMiddleware())
}
```

### 13.4 配置扩展

```go
// config/extend.go 定义了扩展配置结构
type Extend struct {
    AMap           AMap               // 现有扩展
    Runtime        RuntimeConfig      // 运行时配置
    ApplicationEx  ApplicationExConfig // 应用扩展配置
    Device         DeviceConfig       // 设备配置

    // 可继续扩展...
    // CustomModule   CustomConfig
}
```

---

## 附录

### A. 配置文件参考

#### A.1 极简配置 (settings.minimal.yml)

适用于 256MB 内存的光交换机设备。

#### A.2 标准配置 (settings.switch.yml)

适用于 512MB 内存的光交换机设备。

### B. API完整列表

详见 Swagger 文档：`http://<switch-ip>:8000/swagger/admin/index.html`

### C. 常见问题

详见项目 README.md 中的 FAQ 章节。

---

**文档结束**

如有疑问或建议，请访问项目主页：https://github.com/opt-switch-team/opt-switch
