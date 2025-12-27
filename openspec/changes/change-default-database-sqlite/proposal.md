# Change: 将项目默认数据库配置改为 SQLite3

## Why
当前项目默认使用 MySQL 作为数据库，需要用户额外安装和配置 MySQL 服务才能运行项目。对于快速开发体验、演示环境、或者不需要高并发场景，SQLite3 是更简单的选择，因为它不需要独立的服务器进程，所有数据存储在单个文件中。

这将降低项目的上手难度，让开发者能够更快地运行和测试系统。

## What Changes
- 修改 `config/settings.yml` 中的默认数据库驱动从 `mysql` 改为 `sqlite3`
- 更新默认数据库连接字符串为 SQLite 文件路径
- 更新 Dockerfile 以支持 SQLite 编译标签
- 更新 Makefile 添加 SQLite 构建选项
- 更新 README 文档说明 SQLite 配置方式

## Impact
- **Affected specs**: database-config
- **Affected code**:
  - `config/settings.yml` - 主配置文件
  - `Dockerfile` - Docker 构建配置
  - `Makefile` - 构建脚本
  - `README.md`, `README.Zh-cn.md` - 文档

- **User-visible changes**:
  - 新用户可以直接运行项目而无需安装 MySQL
  - 数据存储在 `go-admin-db.db` 文件中
  - 生产环境仍可轻松切换回 MySQL/PostgreSQL

- **Migration path**:
  - 现有 MySQL 用户需要手动配置 `settings.yml` 中的数据库连接
  - 提供配置示例文件 `settings.mysql.yml`
