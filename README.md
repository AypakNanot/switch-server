# opt-switch

<img align="right" width="200" src="https://doc-image.zhangwj.com/img/opt-switch.svg">

[![Build Status](https://github.com/wenjianzhang/opt-switch/workflows/build/badge.svg)](https://github.com/opt-switch-team/opt-switch)
[![Release](https://img.shields.io/github/release/opt-switch-team/opt-switch.svg?style=flat-square)](https://github.com/opt-switch-team/opt-switch/releases)
[![License](https://img.shields.io/github/license/mashape/apistatus.svg)](https://github.com/opt-switch-team/opt-switch/blob/master/LICENSE.md)

English | [简体中文](https://github.com/opt-switch-team/opt-switch/blob/master/README.Zh-cn.md)

> **轻量级光交换机管理系统** - 专为网络设备和嵌入式平台设计的权限管理解决方案

---

## 项目定位

opt-switch 是一款超轻量级的后台管理系统，专门针对**光交换机、路由器和嵌入式网络设备**进行优化。与传统需要 200-500MB 内存的系统不同，opt-switch 仅需 **50-100MB 内存**即可运行，可部署在资源受限的网络边缘设备上。

### 业务场景

| 场景 | 传统方案 | opt-switch 方案 |
|------|----------|-----------------|
| 光交换机管理 | 需要独立服务器 | 直接部署在交换机上 |
| 网络设备配置 | 串口/命令行 | Web 可视化界面 |
| 边缘计算管理 | 云端集中管理 | 本地就近管理 |
| 低功耗设备 | 无法运行完整系统 | 静态二进制，无需依赖 |

---

## 核心特性

### 针对嵌入式平台优化

- **超低内存占用**: 50-100MB（相比同类系统的 200-500MB）
- **多架构支持**: ARMv5/6/7、ARM64、MIPS、MIPS64、PowerPC
- **静态编译**: 无 CGO 依赖，纯 Go 实现
- **无数据库依赖**: 内置 SQLite，开箱即用

### 企业级权限管理

- 基于 **Casbin** 的 RBAC 权限控制
- 支持数据权限（按部门/组织）
- JWT 用户认证
- 完整的审计日志（操作日志、登录日志）

### 完整的管理功能

- **用户管理**: 用户配置、密码策略、状态管理
- **组织架构**: 部门管理、岗位管理、树形结构展示
- **权限控制**: 菜单管理、角色管理、按钮级权限
- **系统配置**: 字典管理、参数管理、动态配置
- **开发工具**: 代码生成器、表单构建器、Swagger 文档

---

## 在线体验

| 版本 | 地址 | 账号密码 |
|------|------|----------|
| Element UI (Vue2) | [vue2.opt-switch.dev](https://vue2.opt-switch.dev/#/login) | admin / 123456 |
| Arco Design (Vue3) | [vue3.opt-switch.dev](https://vue3.opt-switch.dev/#/login) | admin / 123456 |
| Ant Design | [antd.opt-switch.pro](https://antd.opt-switch.pro/) | admin / 123456 |

---

## 支持的平台

### ARM 系列（推荐用于光交换机）

| 架构 | 典型设备 | 编译命令 |
|------|----------|----------|
| ARM64 | 树莓派 4/5、新型光交换机 | `make build-arm64` |
| ARMv7 | 树莓派 2/3、主流交换机 | `make build-armv7` |
| ARMv6 | 树莓派 Zero | `make build-armv6` |

### MIPS 系列（路由器设备）

| 架构 | 典型设备 | 编译命令 |
|------|----------|----------|
| MIPSLE | TP-Link、Netgear 路由器 | `make build-mipsle` |
| MIPS | 企业级交换机 | `make build-mips` |

### 最低硬件要求

- **内存**: 256MB（推荐 512MB）
- **存储**: 100MB 可用空间
- **CPU**: 600MHz 单核（推荐 1GHz 双核）

---

## 快速开始

### 开发环境

```bash
# 克隆后端代码
git clone https://github.com/AypakNanot/switch-server.git
cd switch-server

# 克隆前端代码（同级目录）
git clone https://github.com/opt-switch-team/opt-switch-ui.git

# 后端启动
go mod tidy
go build
./opt-switch migrate -c config/settings.yml
./opt-switch server -c config/settings.yml

# 前端启动
cd ../opt-switch-ui
npm install
npm run dev
```

访问: http://localhost:8000

### Docker 部署（推荐）

```bash
# 构建镜像
docker build -t opt-switch:latest .

# 运行容器
docker run -d --name opt-switch -p 8000:8000 opt-switch:latest
```

### ARM64 设备部署

```bash
# 交叉编译 ARM64 二进制
env CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o opt-switch-arm64 .

# 上传到设备
scp opt-switch-arm64 root@<switch-ip>:/opt/opt-switch/

# SSH 登录并启动
ssh root@<switch-ip>
cd /opt/opt-switch
chmod +x opt-switch-arm64
./opt-switch-arm64 server -c config/settings.yml
```

详细部署指南请参考: [docs/switch-deployment.md](docs/switch-deployment.md)

---

## 配置文件

### 内存优化配置

项目提供了针对不同场景的配置文件：

| 配置文件 | 内存占用 | 适用场景 |
|----------|----------|----------|
| `config/minimal.yml` | 50-70MB | 256MB 内存设备 |
| `config/switch.yml` | 60-100MB | 光交换机/路由器 |
| `config/settings.yml` | 80-100MB | 标准部署 |

### 数据库支持

- **SQLite** (默认): 无需安装，开箱即用
- **MySQL**: 适用于高并发场景
- **PostgreSQL**: 企业级部署
- **SQL Server**: Windows 环境集成

---

## 文档资源

- [官方文档](https://www.opt-switch.dev)
- [部署指南](docs/switch-deployment.md)
- [API 文档](http://localhost:8000/swagger/index.html)
- [视频教程](https://space.bilibili.com/565616721/channel/detail?cid=125737)

---

## 技术栈

### 后端

- **语言**: Go 1.20+
- **框架**: Gin (Web 框架)
- **ORM**: GORM
- **权限**: Casbin
- **数据库**: SQLite (内置) / MySQL / PostgreSQL
- **认证**: JWT
- **文档**: Swagger

### 前端

- **框架**: Vue 2/3
- **UI 组件**: Element UI / Arco Design / Ant Design
- **构建**: Vite

---

## 项目结构

```
opt-switch/
├── cmd/                    # 应用入口
│   ├── api/               # API 服务
│   ├── migrate/           # 数据库迁移
│   └── config/            # 配置管理
├── app/                    # 业务逻辑
│   ├── admin/             # 管理模块
│   │   ├── apis/          # API 处理器
│   │   ├── models/        # 数据模型
│   │   └── router/        # 路由定义
│   └── jobs/              # 后台任务
├── common/                 # 公共组件
│   ├── middleware/        # 中间件
│   ├── database/          # 数据库
│   └── storage/           # 存储抽象
├── config/                 # 配置文件
├── docs/                   # 文档
└── scripts/                # 部署脚本
```

---

## 主要功能模块

### 1. 用户管理
系统用户配置，包括用户信息、角色分配、状态管理等。

### 2. 部门管理
配置系统组织架构（公司、部门、小组），支持树形结构展示和数据权限。

### 3. 岗位管理
配置系统用户的岗位信息。

### 4. 菜单管理
配置系统菜单、操作权限、按钮权限标识、接口权限等。

### 5. 角色管理
角色菜单权限分配，支持按组织划分数据范围权限。

### 6. 字典管理
维护系统中常用的固定数据。

### 7. 参数管理
动态配置系统常用参数。

### 8. 操作日志
系统正常运行日志记录和查询，异常信息日志记录。

### 9. 登录日志
系统登录日志记录查询，包含登录异常。

### 10. 代码生成
根据数据表结构，自动生成增删改查业务代码，可视化操作，零代码实现基础业务。

### 11. 表单构建
自定义页面样式，拖拽实现页面布局。

---

## 社区与支持

<table>
  <tr>
    <td><img src="https://raw.githubusercontent.com/wenjianzhang/image/master/img/wx.png" width="180px"></td>
    <td><img src="https://doc-image.zhangwj.com/img/qrcode_for_gh_b798dc7db30c_258.jpg" width="180px"></td>
    <td><img src="https://raw.githubusercontent.com/wenjianzhang/image/master/img/qq2.png" width="200px"></td>
    <td><a href="https://space.bilibili.com/565616721">Bilibili</a></td>
  </tr>
  <tr>
    <td>微信</td>
    <td>微信公众号</td>
    <td>QQ群</td>
    <td>视频教程</td>
  </tr>
</table>

---

## 贡献者

感谢所有为 opt-switch 做出贡献的开发者：

<!-- 贡献者头像列表保持原样 -->
<span style="margin: 0 5px;" ><a href="https://github.com/wenjianzhang" ><img src="https://images.weserv.nl/?url=avatars.githubusercontent.com/u/3890175?v=4&h=60&w=60&fit=cover&mask=circle&maxage=7d" /></a></span>
<span style="margin: 0 5px;" ><a href="https://github.com/G-Akiraka" ><img src="https://images.weserv.nl/?url=avatars.githubusercontent.com/u/45746659?s=64&v=4&w=60&fit=cover&mask=circle&maxage=7d" /></a></span>
<span style="margin: 0 5px;" ><a href="https://github.com/lwnmengjing" ><img src="https://images.weserv.nl/?url=avatars.githubusercontent.com/u/12806223?s=64&v=4&w=60&fit=cover&mask=circle&maxage=7d" /></a></span>
<span style="margin: 0 5px;" ><a href="https://github.com/bing127" ><img src="https://images.weserv.nl/?url=avatars.githubusercontent.com/u/31166183?s=60&v=4&w=60&fit=cover&mask=circle&maxage=7d" /></a></span>
<span style="margin: 0 5px;" ><a href="https://github.com/chengxiao" ><img src="https://images.weserv.nl/?url=avatars.githubusercontent.com/u/1379545?s=64&v=4&w=60&fit=cover&mask=circle&maxage=7d" /></a></span>

---

## JetBrains 开源支持

opt-switch 项目使用 JetBrains GoLand 进行开发，感谢 JetBrains 提供的开源许可支持。

<a href="https://www.jetbrains.com/?from=opt-switch" target="_blank"><img src="https://raw.githubusercontent.com/panjf2000/illustrations/master/jetbrains/jetbrains-variant-4.png" width="250" align="middle"/></a>

---

## 致谢

- [Gin](https://github.com/gin-gonic/gin) - 高性能 Go Web 框架
- [Casbin](https://github.com/casbin/casbin) - 权限控制库
- [GORM](https://github.com/jinzhu/gorm) - Go ORM 库
- [Vue.js](https://vuejs.org/) - 渐进式 JavaScript 框架
- [Element UI](https://element.eleme.io/) - Vue 组件库
- [Arco Design](https://arco.design/) - 字节跳动 UI 组件库

---

## 许可证

[MIT License](LICENSE.md)

Copyright (c) 2022 wenjianzhang
