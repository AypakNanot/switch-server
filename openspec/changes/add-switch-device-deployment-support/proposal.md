# Change: Add Switch Device Deployment Support

## Why
网络交换机是常见的网络设备，通常具有以下特点：
- **特殊的 CPU 架构**：ARM（ARMv5, ARMv6, ARMv7, ARM64/aarch64）、MIPS（mips、mipsle、mips64、mips64le）、PowerPC 等
- **资源受限**：内存通常在 256MB-1GB 之间，存储空间有限
- **精简的操作系统**：运行嵌入式 Linux（如 OpenWrt、BusyBox），缺少标准 glibc
- **无显示界面**：只能通过网络（SSH、Web）访问和管理

当前项目虽然已经使用纯 Go SQLite（无需 CGO），但缺少针对交换机等嵌入式设备的交叉编译支持、低内存优化配置和部署文档。

## What Changes
- **添加多架构交叉编译支持**：在 Makefile 中添加 ARM、MIPS、PowerPC 等架构的编译目标
- **创建交换机专用配置文件**：`config/settings.switch.yml`，针对低内存环境优化
- **添加部署脚本**：自动化上传、启动、停止流程
- **创建部署文档**：说明不同交换机平台的部署方法

## Impact
- **Affected specs**: deployment (new capability)
- **Affected code**:
  - `Makefile` - 添加多架构编译目标
  - `config/settings.switch.yml` (new) - 交换机专用配置
  - `scripts/deploy-to-switch.sh` (new) - 自动部署脚本
  - `scripts/switch-deployment.md` (new) - 部署文档

- **User-visible changes**:
  - ✅ 可编译到 ARM、MIPS、PowerPC 等架构
  - ✅ 低内存环境优化配置（256MB 可运行）
  - ✅ 单文件部署，无需额外依赖
  - ✅ 提供自动化部署脚本
  - ✅ 详细的交换机部署文档

- **Supported Platforms**:
  - ARM: ARMv5, ARMv6, ARMv7, ARM64 (树莓派、大多数嵌入式设备)
  - MIPS: MIPS (大端)、MIPSLE (小端) - 常见于路由器和交换机
  - PowerPC: 常见于部分企业级交换机
  - x86_64: 标准 PC 和服务器

- **Migration path**:
  - 现有部署不受影响
  - 新增交叉编译目标，不改变现有构建流程
  - 配置文件独立，与现有配置共存

## Why This Approach

1. **纯 Go 编译**：项目已使用 `glebarez/sqlite`（纯 Go SQLite 驱动），无需 CGO，适合交叉编译
2. **静态链接**：使用 `-ldflags="-w -s"` 去除调试信息，减小体积
3. **配置分离**：交换机专用配置文件，避免修改主配置
4. **架构优先**：覆盖最常见的嵌入式架构（ARM、MIPS、PowerPC）
