# Change: Rename Project from go-admin to opt-switch

## Why
项目需要更明确地反映其定位 - 针对网络交换机设备的管理系统。"go-admin" 名称过于通用，而 "opt-switch" 能更好地体现项目的实际用途（光交换机管理系统）。

## What Changes
- **重命名项目名称**：将所有 "go-admin" 引用替换为 "opt-switch"
- **更新模块名称**：go.mod 中的 module 名称从 `go-admin` 改为 `opt-switch`
- **更新配置文件**：数据库名称、配置文件路径、系统名称等
- **更新文档**：README、注释中的项目名称
- **更新二进制输出**：编译后的可执行文件名称

## Impact
- **Affected specs**: 无新功能，仅重构
- **Affected code**:
  - `go.mod` - 模块名称
  - `config/db.sql` - 系统配置中的名称
  - `README.md`, `README.Zh-cn.md` - 项目文档
  - 所有包含 "go-admin" 的代码注释和字符串
  - Makefile 中的二进制输出名称

- **User-visible changes**:
  - ✅ 数据库名称更新
  - ✅ 二进制文件名更新
  - ✅ 系统名称更新
  - ✅ 文档更新

- **Migration path**:
  - 现有部署需要：
    1. 更新数据库名称（可选，向后兼容）
    2. 更新配置文件引用
    3. 重新编译部署
  - 数据库表结构不变

## Why This Approach

1. **全面替换**：使用全局搜索替换确保所有引用都被更新
2. **保持兼容**：数据库表名和结构保持不变，只更新项目级配置
3. **分步实施**：先更新核心文件，再更新文档和注释
4. **验证测试**：每次更新后进行编译和运行测试
