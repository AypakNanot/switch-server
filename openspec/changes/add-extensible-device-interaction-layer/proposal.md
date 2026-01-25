# Change: Add Extensible Device Interaction Layer

## Why

当前系统缺乏与网络设备交互的能力。系统运行在交换机设备上，需要管理自身这台交换机（通过本地回连或本地 CLI），提供配置管理、状态监控等功能。

当前存在的问题：
1. 没有统一的设备交互抽象
2. 缺乏协议扩展性（后续需支持 NETCONF）
3. 多用户同时操作时需要排队机制
4. 缺乏命令执行历史记录

## What Changes

- **ADDED** 设备交互层核心架构（`pkg/device/`）
  - 协议适配器接口（Protocol Adapter Interface）
  - **连接池 + 命令队列**（Connection Pool + Command Queue）
  - 命令执行器（Command Executor）
  - 执行日志记录器（Execution Logger）

- **ADDED** CLI 协议支持
  - SSH 协议适配器
  - Telnet 协议适配器
  - 命令解析和响应处理

- **ADDED** 命令执行 API（`app/device/`）
  - 单命令执行
  - 批量命令执行
  - 执行历史查询

- **ADDED** 配置文件支持
  - 设备连接配置（从配置文件读取，支持热加载）
  - 连接池配置（可配置并发数）

- **REMOVED** 无需数据库存储设备信息
  - 设备信息从配置文件读取
  - 执行历史使用日志文件

- **BREAKING** 无破坏性变更，纯新增功能

## Impact

- **Affected specs**: device-interaction (新增)
- **Affected code**:
  - 新增 `pkg/device/` - 设备交互核心层
  - 新增 `app/device/` - 命令执行 API
  - 修改 `config/settings.yml` - 添加设备配置
  - 修改 `app/admin/router/` - 添加命令路由

## Dependencies

- Go 1.24+
- golang.org/x/crypto/ssh (SSH 协议支持)
- 支持本地回连或本地 CLI 执行

## Migration Plan

无迁移需求，纯新增功能

## Open Questions

- 本地 CLI 执行是通过 SSH 回连 127.0.0.1，还是直接执行本地命令？（建议 SSH 回连保持协议一致性）
