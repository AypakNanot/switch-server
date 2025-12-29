# Design: Rename Project to opt-switch

## Context
当前项目名为 "go-admin"，这是一个通用的后台管理系统名称。随着项目专注于光交换机设备管理，需要更具体的项目名称来反映其实际用途。

## Goals / Non-Goals

**Goals**:
- 将项目名称从 go-admin 改为 opt-switch
- 更新所有代码、文档、配置中的引用
- 保持向后兼容性
- 确保编译和运行正常

**Non-Goals**:
- 不改变数据库表结构
-不改变 API 接口路径（保持 /api/v1/...）
- 不改变核心功能逻辑

## Decisions

### 1. 命名规范
**决策**: 项目名称使用小写、连字符分隔的格式

**原因**:
- 符合 Go 项目命名惯例
- 避免特殊字符和空格问题
- 便于在 URL 和命令行中使用

### 2. 模块名称
**决策**: go.mod 中的 module 从 `go-admin` 改为 `opt-switch`

**影响**:
- 所有 import 路径需要更新
- 重新运行 `go mod tidy`

### 3. 二进制文件名
**决策**: 编译输出的二进制文件名为 `opt-switch`

**Makefile 更新**:
```makefile
BINARY_NAME=opt-switch
```

### 4. 数据库名称
**决策**: 将默认数据库名从 `go-admin-db.db` 改为 `opt-switch.db`

**兼容性**:
- 支持通过配置文件自定义数据库名
- 旧数据库文件可以通过重命名继续使用

### 5. 系统配置名称
**决策**: 更新系统默认配置中的名称

**配置文件更新**:
```yaml
# config/settings.yml
application:
  name: opt-switch管理系统
```

## Files to Change

### 核心文件
- `go.mod` - 模块名称
- `main.go` - 可能包含项目名称
- `Makefile` - 二进制输出名称

### 配置文件
- `config/settings.yml` - 系统名称
- `config/db.sql` - 系统配置数据
- `config/settings.switch.yml` - 交换机配置

### 文档文件
- `README.md` - 项目名称
- `README.Zh-cn.md` - 中文项目名称
- 所有包含 "go-admin" 的注释

## Migration Plan

### 开发环境
1. 更新 go.mod
2. 运行 `go mod tidy`
3. 更新本地导入路径
4. 重新编译测试

### 部署环境
1. 备份现有数据库
2. 重命名数据库文件（可选）
3. 更新配置文件
4. 部署新版本二进制

### 回滚计划
如需回滚：
1. 恢复旧版本二进制
2. 恢复配置文件
3. 数据库无需回滚（表结构未变）
