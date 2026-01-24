# opt-switch 开发指南

## 目录

- [1. 环境搭建](#1-环境搭建)
- [2. 项目结构](#2-项目结构)
- [3. 开发流程](#3-开发流程)
- [4. 示例：添加设备管理模块](#4-示例添加设备管理模块)
- [5. 代码生成器使用](#5-代码生成器使用)
- [6. API 开发规范](#6-api-开发规范)
- [7. 数据库操作](#7-数据库操作)
- [8. 路由注册](#8-路由注册)
- [9. 调试与测试](#9-调试与测试)
- [10. 构建与部署](#10-构建与部署)

---

## 1. 环境搭建

### 1.1 系统要求

| 组件 | 版本要求 | 说明 |
|------|----------|------|
| Go | 1.24+ | 开发语言 |
| Git | 2.20+ | 版本控制 |
| SQLite | 3.x | 数据库（内置） |
| 编辑器 | VS Code / GoLand | 推荐 |

### 1.2 克隆项目

```bash
# 克隆后端代码
git clone https://github.com/AypakNanot/switch-server.git
cd switch-server

# 查看项目结构
tree -L 2
```

### 1.3 安装依赖

```bash
# 下载 Go 依赖
go mod download

# 验证依赖
go mod verify

# 整理依赖（可选）
go mod tidy
```

### 1.4 配置开发环境

```bash
# 复制开发配置文件
cp config/settings.yml config/settings.dev.yml

# 编辑配置文件，设置开发环境参数
# - mode: dev (开发模式，启用调试日志)
# - port: 8000
# - database: sqlite3 (使用本地数据库)
```

### 1.5 初始化数据库

```bash
# 编译项目
go build -o opt-switch

# 初始化数据库表结构
./opt-switch migrate -c config/settings.yml

# 启动开发服务器
./opt-switch server -c config/settings.yml
```

### 1.6 验证安装

```bash
# 访问健康检查接口
curl http://localhost:8000/info

# 访问 Swagger 文档
# 浏览器打开: http://localhost:8000/swagger/admin/index.html

# 访问前端界面
# 浏览器打开: http://localhost:8000
```

---

## 2. 项目结构

### 2.1 目录结构详解

```
opt-switch/
├── app/                        # 应用层（业务逻辑）
│   ├── admin/                  # 管理后台模块
│   │   ├── apis/               # API 处理器（控制器层）
│   │   │   ├── sys_user.go     # 用户 API
│   │   │   ├── sys_role.go     # 角色 API
│   │   │   └── ...
│   │   ├── models/             # 数据模型（ORM 层）
│   │   │   ├── sys_user.go     # 用户模型
│   │   │   └── ...
│   │   ├── router/             # 路由定义（路由层）
│   │   │   ├── init_router.go  # 路由初始化
│   │   │   ├── sys_user.go     # 用户路由
│   │   │   └── ...
│   │   └── service/            # 业务服务（服务层）
│   │       ├── sys_user.go     # 用户服务
│   │       ├── dto/            # 数据传输对象
│   │       │   └── sys_user.go # 用户 DTO
│   │       └── ...
│   ├── jobs/                   # 定时任务
│   └── other/                  # 其他业务模块
│
├── cmd/                        # 命令行工具
│   ├── api/                    # API 服务
│   │   ├── server.go           # 服务器启动
│   │   ├── runtime.go          # 运行时配置
│   │   └── jobs.go             # 任务调度
│   ├── migrate/                # 数据库迁移工具
│   ├── config/                 # 配置管理工具
│   └── version/                # 版本信息
│
├── common/                     # 公共组件
│   ├── middleware/             # 中间件
│   │   ├── jwt.go              # JWT 认证
│   │   ├── auth.go             # 权限检查
│   │   ├── logger.go           # 日志记录
│   │   └── ...
│   ├── database/               # 数据库连接
│   ├── storage/                # 存储抽象
│   ├── models/                 # 公共模型基类
│   ├── dto/                    # 公共 DTO
│   └── actions/                # 公共操作
│
├── config/                     # 配置文件
│   ├── settings.yml            # 默认配置
│   ├── settings.minimal.yml    # 极简配置（50-70MB）
│   └── settings.switch.yml     # 交换机配置（60-100MB）
│
├── web/                        # 前端资源（嵌入）
│   └── dist/                   # 构建产物
│
├── static/                     # 静态文件
├── template/                   # 代码生成模板
├── docs/                       # 文档
├── main.go                     # 程序入口
├── go.mod                      # Go 模块定义
├── go.sum                      # 依赖锁定
├── Makefile                    # 构建脚本
└── Dockerfile                  # Docker 镜像
```

### 2.2 分层架构

```
┌─────────────────────────────────────────────────────────┐
│                      Client (前端/API)                   │
└─────────────────────────────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────┐
│                    Router Layer (路由层)                  │
│  app/admin/router/ - 路由定义、中间件配置                 │
└─────────────────────────────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────┐
│                   API Layer (控制器层)                    │
│  app/admin/apis/ - 请求处理、参数绑定、响应封装           │
└─────────────────────────────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────┐
│                  Service Layer (服务层)                   │
│  app/admin/service/ - 业务逻辑、数据处理                 │
└─────────────────────────────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────┐
│                   Model Layer (模型层)                    │
│  app/admin/models/ - 数据模型、数据库映射                │
└─────────────────────────────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────┐
│                    Database (SQLite)                     │
└─────────────────────────────────────────────────────────┘
```

---

## 3. 开发流程

### 3.1 完整开发流程图

```
1. 需求分析
   │
   ▼
2. 数据库设计
   │
   ▼
3. 创建数据表
   │
   ▼
4. 使用代码生成器生成基础代码
   │
   ▼
5. 根据需求自定义修改代码
   │   ├── Model 层：定义数据结构
   │   ├── DTO 层：定义请求/响应结构
   │   ├── Service 层：实现业务逻辑
   │   ├── API 层：实现控制器
   │   └── Router 层：注册路由
   │
   ▼
6. 本地测试
   │
   ▼
7. 代码审查
   │
   ▼
8. 提交代码
   │
   ▼
9. 构建 & 部署
```

### 3.2 使用代码生成器的优势

| 手动开发 | 使用代码生成器 |
|----------|----------------|
| 需要手动编写所有层 | 自动生成完整 CRUD |
| 容易出错 | 代码规范统一 |
| 耗时长 | 几分钟完成基础功能 |
| 需要记住命名规范 | 自动遵守项目规范 |

---

## 4. 示例：添加设备管理模块

以下是一个完整的开发示例，演示如何添加一个新的设备管理功能。

### 4.1 需求分析

**功能描述**：管理光交换机设备信息

**核心功能**：
- 设备列表查询（分页、筛选）
- 添加新设备
- 编辑设备信息
- 删除设备
- 设备详情查看

### 4.2 数据库设计

```sql
-- 创建设备表
CREATE TABLE sys_device (
    device_id INTEGER PRIMARY KEY AUTOINCREMENT,
    device_name VARCHAR(100) NOT NULL,
    device_type VARCHAR(50),
    ip_address VARCHAR(15),
    port INTEGER,
    username VARCHAR(50),
    password VARCHAR(100),
    location VARCHAR(200),
    status VARCHAR(20) DEFAULT 'online',
    remark VARCHAR(500),
    create_by INTEGER,
    update_by INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建索引
CREATE INDEX idx_device_name ON sys_device(device_name);
CREATE INDEX idx_device_type ON sys_device(device_type);
CREATE INDEX idx_ip_address ON sys_device(ip_address);
```

### 4.3 使用代码生成器

#### 方法 1：通过 Web UI（推荐）

1. 启动服务并登录管理后台
2. 进入「系统工具」→「代码生成」
3. 选择 `sys_device` 表
4. 配置生成选项：
   - 包名：`admin`
   - 模块名：`device`
   - 功能名：`设备管理`
5. 点击「预览代码」查看生成的代码
6. 点击「生成代码」自动创建文件

#### 方法 2：手动生成（如果需要自定义）

```bash
# 运行代码生成命令（假设有命令行工具）
./opt-switch generate -t sys_device -p admin -m device
```

### 4.4 生成的文件结构

```
app/admin/
├── models/
│   └── sys_device.go          # 数据模型
├── service/dto/
│   └── sys_device.go          # 数据传输对象
├── service/
│   └── sys_device.go          # 业务服务
├── apis/
│   └── sys_device.go          # API 处理器
└── router/
    └── sys_device.go          # 路由定义
```

### 4.5 自定义代码

#### 4.5.1 Model 层（app/admin/models/sys_device.go）

```go
package models

import (
    "opt-switch/common/models"
)

type SysDevice struct {
    DeviceId   int    `gorm:"primaryKey;autoIncrement;comment:设备ID" json:"deviceId"`
    DeviceName string `json:"deviceName" gorm:"size:100;comment:设备名称"`
    DeviceType string `json:"deviceType" gorm:"size:50;comment:设备类型"`
    IpAddress  string `json:"ipAddress" gorm:"size:15;comment:IP地址"`
    Port       int    `json:"port" gorm:"comment:端口号"`
    Username   string `json:"username" gorm:"size:50;comment:用户名"`
    Password   string `json:"-" gorm:"size:100;comment:密码"`
    Location   string `json:"location" gorm:"size:200;comment:位置"`
    Status     string `json:"status" gorm:"size:20;comment:状态"`
    Remark     string `json:"remark" gorm:"size:500;comment:备注"`
    models.ControlBy
    models.ModelTime
}

func (*SysDevice) TableName() string {
    return "sys_device"
}

func (e *SysDevice) Generate() models.ActiveRecord {
    o := *e
    return &o
}

func (e *SysDevice) GetId() interface{} {
    return e.DeviceId
}
```

#### 4.5.2 DTO 层（app/admin/service/dto/sys_device.go）

```go
package dto

import (
    "opt-switch/app/admin/models"
    "opt-switch/common/dto"
    common "opt-switch/common/models"
)

// 分页查询请求
type SysDeviceGetPageReq struct {
    dto.Pagination `search:"-"`
    DeviceId       int    `form:"deviceId" search:"type:exact;column:device_id;table:sys_device" comment:"设备ID"`
    DeviceName     string `form:"deviceName" search:"type:contains;column:device_name;table:sys_device" comment:"设备名称"`
    DeviceType     string `form:"deviceType" search:"type:exact;column:device_type;table:sys_device" comment:"设备类型"`
    IpAddress      string `form:"ipAddress" search:"type:contains;column:ip_address;table:sys_device" comment:"IP地址"`
    Status         string `form:"status" search:"type:exact;column:status;table:sys_device" comment:"状态"`
    SysDeviceOrder
}

type SysDeviceOrder struct {
    DeviceIdOrder     string `search:"type:order;column:device_id;table:sys_device" form:"deviceIdOrder"`
    DeviceNameOrder   string `search:"type:order;column:device_name;table:sys_device" form:"deviceNameOrder"`
    CreatedAtOrder   string `search:"type:order;column:created_at;table:sys_device" form:"createdAtOrder"`
}

func (m *SysDeviceGetPageReq) GetNeedSearch() interface{} {
    return *m
}

// 插入请求
type SysDeviceInsertReq struct {
    DeviceId   int    `json:"deviceId" comment:"设备ID"`
    DeviceName string `json:"deviceName" comment:"设备名称" vd:"len($)>0"`
    DeviceType string `json:"deviceType" comment:"设备类型"`
    IpAddress  string `json:"ipAddress" comment:"IP地址" vd:"len($)>0"`
    Port       int    `json:"port" comment:"端口号"`
    Username   string `json:"username" comment:"用户名"`
    Password   string `json:"password" comment:"密码"`
    Location   string `json:"location" comment:"位置"`
    Status     string `json:"status" comment:"状态" default:"online"`
    Remark     string `json:"remark" comment:"备注"`
    common.ControlBy
}

func (s *SysDeviceInsertReq) Generate(model *models.SysDevice) {
    if s.DeviceId != 0 {
        model.DeviceId = s.DeviceId
    }
    model.DeviceName = s.DeviceName
    model.DeviceType = s.DeviceType
    model.IpAddress = s.IpAddress
    model.Port = s.Port
    model.Username = s.Username
    model.Password = s.Password
    model.Location = s.Location
    model.Status = s.Status
    model.Remark = s.Remark
    model.CreateBy = s.CreateBy
}

func (s *SysDeviceInsertReq) GetId() interface{} {
    return s.DeviceId
}

// 更新请求
type SysDeviceUpdateReq struct {
    DeviceId   int    `json:"deviceId" comment:"设备ID" vd:"$>0"`
    DeviceName string `json:"deviceName" comment:"设备名称" vd:"len($)>0"`
    DeviceType string `json:"deviceType" comment:"设备类型"`
    IpAddress  string `json:"ipAddress" comment:"IP地址" vd:"len($)>0"`
    Port       int    `json:"port" comment:"端口号"`
    Username   string `json:"username" comment:"用户名"`
    Password   string `json:"password" comment:"密码"`
    Location   string `json:"location" comment:"位置"`
    Status     string `json:"status" comment:"状态"`
    Remark     string `json:"remark" comment:"备注"`
    common.ControlBy
}

func (s *SysDeviceUpdateReq) Generate(model *models.SysDevice) {
    if s.DeviceId != 0 {
        model.DeviceId = s.DeviceId
    }
    model.DeviceName = s.DeviceName
    model.DeviceType = s.DeviceType
    model.IpAddress = s.IpAddress
    model.Port = s.Port
    model.Username = s.Username
    if s.Password != "" {
        model.Password = s.Password
    }
    model.Location = s.Location
    model.Status = s.Status
    model.Remark = s.Remark
}

func (s *SysDeviceUpdateReq) GetId() interface{} {
    return s.DeviceId
}

// 删除请求
type SysDeviceById struct {
    dto.ObjectById
    common.ControlBy
}

func (s *SysDeviceById) GetId() interface{} {
    if len(s.Ids) > 0 {
        s.Ids = append(s.Ids, s.Id)
        return s.Ids
    }
    return s.Id
}

func (s *SysDeviceById) GenerateM() (common.ActiveRecord, error) {
    return &models.SysDevice{}, nil
}
```

#### 4.5.3 Service 层（app/admin/service/sys_device.go）

```go
package service

import (
    "errors"
    "opt-switch/app/admin/models"
    "opt-switch/app/admin/service/dto"
    "opt-switch/common/actions"
    cDto "opt-switch/common/dto"

    "github.com/go-admin-team/go-admin-core/sdk/service"
    "gorm.io/gorm"
)

type SysDevice struct {
    service.Service
}

// GetPage 获取设备列表
func (e *SysDevice) GetPage(c *dto.SysDeviceGetPageReq, p *actions.DataPermission, list *[]models.SysDevice, count *int64) error {
    var err error
    var data models.SysDevice

    err = e.Orm.
        Scopes(
            cDto.MakeCondition(c.GetNeedSearch()),
            cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
            actions.Permission(data.TableName(), p),
        ).
        Find(list).Limit(-1).Offset(-1).
        Count(count).Error
    if err != nil {
        e.Log.Errorf("db error: %s", err)
        return err
    }
    return nil
}

// Get 获取设备详情
func (e *SysDevice) Get(d *dto.SysDeviceById, p *actions.DataPermission, model *models.SysDevice) error {
    var data models.SysDevice

    err := e.Orm.Model(&data).
        Scopes(
            actions.Permission(data.TableName(), p),
        ).
        First(model, d.GetId()).Error
    if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
        err = errors.New("查看对象不存在或无权查看")
        e.Log.Errorf("db error: %s", err)
        return err
    }
    if err != nil {
        e.Log.Errorf("db error: %s", err)
        return err
    }
    return nil
}

// Insert 添加设备
func (e *SysDevice) Insert(c *dto.SysDeviceInsertReq) error {
    var err error
    var data models.SysDevice

    // 检查设备名称是否重复
    var count int64
    err = e.Orm.Model(&data).Where("device_name = ?", c.DeviceName).Count(&count).Error
    if err != nil {
        e.Log.Errorf("db error: %s", err)
        return err
    }
    if count > 0 {
        err = errors.New("设备名称已存在")
        e.Log.Errorf("db error: %s", err)
        return err
    }

    c.Generate(&data)
    err = e.Orm.Create(&data).Error
    if err != nil {
        e.Log.Errorf("db error: %s", err)
        return err
    }
    return nil
}

// Update 更新设备
func (e *SysDevice) Update(c *dto.SysDeviceUpdateReq, p *actions.DataPermission) error {
    var err error
    var model models.SysDevice

    db := e.Orm.Scopes(
        actions.Permission(model.TableName(), p),
    ).First(&model, c.GetId())

    if err = db.Error; err != nil {
        e.Log.Errorf("Service Update error: %s", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权更新该数据")
    }

    c.Generate(&model)

    err = e.Orm.Model(&model).Omit("created_at").Updates(&model).Error
    if err != nil {
        e.Log.Errorf("db error: %s", err)
        return err
    }
    return nil
}

// Remove 删除设备
func (e *SysDevice) Remove(c *dto.SysDeviceById, p *actions.DataPermission) error {
    var err error
    var data models.SysDevice

    db := e.Orm.Model(&data).
        Scopes(
            actions.Permission(data.TableName(), p),
        ).Delete(&data, c.GetId())

    if err = db.Error; err != nil {
        e.Log.Errorf("Error found in Remove: %s", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权删除该数据")
    }
    return nil
}
```

#### 4.5.4 API 层（app/admin/apis/sys_device.go）

```go
package apis

import (
    "github.com/gin-gonic/gin/binding"
    "opt-switch/app/admin/models"
    "opt-switch/app/admin/service"
    "opt-switch/app/admin/service/dto"
    "opt-switch/common/actions"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/go-admin-team/go-admin-core/sdk/api"
    _ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"
)

type SysDevice struct {
    api.Api
}

// GetPage 获取设备列表
// @Summary 获取设备列表
// @Description 获取设备列表数据
// @Tags 设备管理
// @Param deviceName query string false "设备名称"
// @Param deviceType query string false "设备类型"
// @Success 200 {string} {object} response.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/sys-device [get]
// @Security Bearer
func (e SysDevice) GetPage(c *gin.Context) {
    s := service.SysDevice{}
    req := dto.SysDeviceGetPageReq{}
    err := e.MakeContext(c).
        MakeOrm().
        Bind(&req).
        MakeService(&s.Service).
        Errors
    if err != nil {
        e.Logger.Error(err)
        e.Error(500, err, err.Error())
        return
    }

    p := actions.GetPermissionFromContext(c)
    list := make([]models.SysDevice, 0)
    var count int64

    err = s.GetPage(&req, p, &list, &count)
    if err != nil {
        e.Error(500, err, "查询失败")
        return
    }

    e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取设备详情
// @Summary 获取设备详情
// @Description 获取设备详情
// @Tags 设备管理
// @Param deviceId path int true "设备ID"
// @Success 200 {object} response.Response "{"code": 200, "data": {...}}"
// @Router /api/v1/sys-device/{deviceId} [get]
// @Security Bearer
func (e SysDevice) Get(c *gin.Context) {
    s := service.SysDevice{}
    req := dto.SysDeviceById{}
    err := e.MakeContext(c).
        MakeOrm().
        Bind(&req, nil).
        MakeService(&s.Service).
        Errors
    if err != nil {
        e.Logger.Error(err)
        e.Error(500, err, err.Error())
        return
    }

    var object models.SysDevice
    p := actions.GetPermissionFromContext(c)
    err = s.Get(&req, p, &object)
    if err != nil {
        e.Error(http.StatusUnprocessableEntity, err, "查询失败")
        return
    }
    e.OK(object, "查询成功")
}

// Insert 添加设备
// @Summary 添加设备
// @Description 添加新设备
// @Tags 设备管理
// @Accept application/json
// @Product application/json
// @Param data body dto.SysDeviceInsertReq true "设备数据"
// @Success 200 {object} response.Response "{"code": 200, "data": {...}}"
// @Router /api/v1/sys-device [post]
// @Security Bearer
func (e SysDevice) Insert(c *gin.Context) {
    s := service.SysDevice{}
    req := dto.SysDeviceInsertReq{}
    err := e.MakeContext(c).
        MakeOrm().
        Bind(&req, binding.JSON).
        MakeService(&s.Service).
        Errors
    if err != nil {
        e.Logger.Error(err)
        e.Error(500, err, err.Error())
        return
    }

    err = s.Insert(&req)
    if err != nil {
        e.Logger.Error(err)
        e.Error(500, err, err.Error())
        return
    }

    e.OK(req.GetId(), "添加成功")
}

// Update 更新设备
// @Summary 更新设备
// @Description 更新设备信息
// @Tags 设备管理
// @Accept application/json
// @Product application/json
// @Param data body dto.SysDeviceUpdateReq true "设备数据"
// @Success 200 {object} response.Response "{"code": 200, "data": {...}}"
// @Router /api/v1/sys-device [put]
// @Security Bearer
func (e SysDevice) Update(c *gin.Context) {
    s := service.SysDevice{}
    req := dto.SysDeviceUpdateReq{}
    err := e.MakeContext(c).
        MakeOrm().
        Bind(&req).
        MakeService(&s.Service).
        Errors
    if err != nil {
        e.Logger.Error(err)
        e.Error(500, err, err.Error())
        return
    }

    p := actions.GetPermissionFromContext(c)
    err = s.Update(&req, p)
    if err != nil {
        e.Logger.Error(err)
        e.Error(500, err, "更新失败")
        return
    }
    e.OK(req.GetId(), "更新成功")
}

// Delete 删除设备
// @Summary 删除设备
// @Description 删除设备
// @Tags 设备管理
// @Param deviceId path int true "设备ID"
// @Success 200 {object} response.Response "{"code": 200, "data": {...}}"
// @Router /api/v1/sys-device/{deviceId} [delete]
// @Security Bearer
func (e SysDevice) Delete(c *gin.Context) {
    s := service.SysDevice{}
    req := dto.SysDeviceById{}
    err := e.MakeContext(c).
        MakeOrm().
        Bind(&req).
        MakeService(&s.Service).
        Errors
    if err != nil {
        e.Logger.Error(err)
        e.Error(500, err, err.Error())
        return
    }

    p := actions.GetPermissionFromContext(c)
    err = s.Remove(&req, p)
    if err != nil {
        e.Logger.Error(err)
        e.Error(500, err, "删除失败")
        return
    }
    e.OK(req.GetId(), "删除成功")
}
```

#### 4.5.5 Router 层（app/admin/router/sys_device.go）

```go
package router

import (
    "github.com/gin-gonic/gin"
    jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"
    "opt-switch/app/admin/apis"
    "opt-switch/common/actions"
    "opt-switch/common/middleware"
)

func init() {
    routerCheckRole = append(routerCheckRole, registerSysDeviceRouter)
}

// 需认证的路由
func registerSysDeviceRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
    api := apis.SysDevice{}
    r := v1.Group("/sys-device").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole()).Use(actions.PermissionAction())
    {
        r.GET("", api.GetPage)         // GET /api/v1/sys-device
        r.GET("/:id", api.Get)         // GET /api/v1/sys-device/:id
        r.POST("", api.Insert)         // POST /api/v1/sys-device
        r.PUT("", api.Update)          // PUT /api/v1/sys-device
        r.DELETE("/:id", api.Delete)   // DELETE /api/v1/sys-device/:id
    }
}
```

### 4.6 菜单权限配置

登录管理后台，配置菜单权限：

1. 进入「系统管理」→「菜单管理」
2. 添加新菜单：
   - 菜单名称：设备管理
   - 菜单类型：目录
   - 路由路径：device
   - 组件路径：Layout/ParentView

3. 添加子菜单：
   - 菜单名称：设备列表
   - 菜单类型：菜单
   - 路由路径：device/list
   - 组件路径：device/list/index
   - 权限标识：sys-device:list

4. 配置按钮权限：
   - 添加：sys-device:add
   - 编辑：sys-device:edit
   - 删除：sys-device:delete
   - 查询：sys-device:query

### 4.7 测试 API

```bash
# 1. 登录获取 Token
TOKEN=$(curl -X POST http://localhost:8000/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"123456"}' \
  | jq -r '.data.token')

# 2. 获取设备列表
curl -X GET "http://localhost:8000/api/v1/sys-device?page=1&pageSize=10" \
  -H "Authorization: Bearer $TOKEN"

# 3. 添加设备
curl -X POST http://localhost:8000/api/v1/sys-device \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "deviceName": "核心交换机01",
    "deviceType": "core",
    "ipAddress": "192.168.1.1",
    "port": 22,
    "username": "admin",
    "password": "password",
    "location": "机房A",
    "status": "online"
  }'

# 4. 更新设备
curl -X PUT http://localhost:8000/api/v1/sys-device \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "deviceId": 1,
    "deviceName": "核心交换机01-更新",
    "status": "offline"
  }'

# 5. 删除设备
curl -X DELETE http://localhost:8000/api/v1/sys-device/1 \
  -H "Authorization: Bearer $TOKEN"
```

---

## 5. 代码生成器使用

### 5.1 Web UI 方式

1. 启动服务并登录
2. 进入「系统工具」→「代码生成」
3. 选择数据库表
4. 配置生成选项：
   - 基础信息：包名、模块名、功能名
   - 字段配置：显示名称、查询方式、表单类型
   - 生成选项：是否生成前端代码
5. 预览代码
6. 生成代码

### 5.2 代码生成配置

| 配置项 | 说明 | 示例 |
|--------|------|------|
| 包名 | 代码所在包 | admin |
| 模块名 | 功能模块标识 | device |
| 功能名 | 中文名称 | 设备管理 |
| 业务名 | 数据库表名 | sys_device |
| 表前缀 | 自动去除前缀 | sys_ |

### 5.3 字段配置说明

| 配置项 | 说明 | 可选值 |
|--------|------|--------|
| 查询方式 | 列表查询条件 | exact(精确), contains(包含) |
| 表单类型 | 表单控件 | input, select, date, number |
| 是否显示 | 列表是否显示 | true/false |
| 是否必填 | 表单是否必填 | true/false |
| 验证规则 | 数据验证 | email, url, phone |

---

## 6. API 开发规范

### 6.1 命名规范

#### 文件命名
```
# Model 层
sys_{module}.go       // 示例：sys_device.go

# Service 层
sys_{module}.go       // 示例：sys_device.go

# DTO 层
sys_{module}.go       // 示例：sys_device.go

# API 层
sys_{module}.go       // 示例：sys_device.go

# Router 层
sys_{module}.go       // 示例：sys_device.go
```

#### 结构体命名
```
# Model
Sys{Module}           // 示例：SysDevice

# Service
Sys{Module}           // 示例：SysDevice

# DTO
Sys{Module}GetPageReq       // 分页查询请求
Sys{Module}ById             // 根据 ID 查询
Sys{Module}InsertReq        // 插入请求
Sys{Module}UpdateReq        // 更新请求

# API
Sys{Module}           // 示例：SysDevice
```

#### 方法命名
```
# Model 方法
TableName() string
Generate() ActiveRecord
GetId() interface{}

# Service 方法
GetPage()          // 分页查询
Get()              // 获取单个
Insert()           // 插入
Update()           // 更新
Remove()           // 删除

# API 方法
GetPage()          // GET /api/v1/sys-module
Get()              // GET /api/v1/sys-module/:id
Insert()           // POST /api/v1/sys-module
Update()           // PUT /api/v1/sys-module
Delete()           // DELETE /api/v1/sys-module/:id
```

### 6.2 路由规范

```
# RESTful 风格
GET    /api/v1/sys-module           # 列表
GET    /api/v1/sys-module/:id       # 详情
POST   /api/v1/sys-module           # 创建
PUT    /api/v1/sys-module           # 更新
DELETE /api/v1/sys-module/:id       # 删除

# 自定义操作
GET    /api/v1/module/profile       # 个人资料
POST   /api/v1/module/avatar        # 上传头像
PUT    /api/v1/module/password      # 修改密码
PUT    /api/v1/module/status        # 修改状态
```

### 6.3 响应格式

#### 成功响应
```json
{
  "code": 200,
  "msg": "操作成功",
  "data": {
    // 返回数据
  }
}
```

#### 分页响应
```json
{
  "code": 200,
  "msg": "查询成功",
  "data": [
    // 数据列表
  ],
  "page": {
    "currPage": 1,
    "pageSize": 10,
    "totalCount": 100,
    "totalPage": 10
  }
}
```

#### 错误响应
```json
{
  "code": 500,
  "msg": "操作失败",
  "data": "错误详情"
}
```

### 6.4 Swagger 注释规范

```go
// GetPage 获取{模块}列表
// @Summary {功能描述}
// @Description {详细描述}
// @Tags {模块名}
// @Param {参数名} query {类型} false "{参数说明}"
// @Success 200 {string} {object} response.Response
// @Router /api/v1/{path} [get]
// @Security Bearer
```

---

## 7. 数据库操作

### 7.1 Model 定义规范

```go
type SysDevice struct {
    // 主键
    DeviceId   int    `gorm:"primaryKey;autoIncrement;comment:设备ID" json:"deviceId"`

    // 字段定义格式：`gorm:"属性;配置;comment:注释" json:"字段名"`
    DeviceName string `json:"deviceName" gorm:"size:100;not null;comment:设备名称"`

    // 密码字段不返回给前端
    Password   string `json:"-" gorm:"size:100;comment:密码"`

    // 公共字段（所有表必备）
    models.ControlBy     // 创建人、更新人
    models.ModelTime     // 创建时间、更新时间
}

// 表名映射
func (*SysDevice) TableName() string {
    return "sys_device"
}

// 实现 ActiveRecord 接口
func (e *SysDevice) Generate() models.ActiveRecord {
    o := *e
    return &o
}

// 获取主键
func (e *SysDevice) GetId() interface{} {
    return e.DeviceId
}
```

### 7.2 GORM 常用操作

```go
// 查询
db.First(&device, 1)                    // 主键查询
db.Where("device_id = ?", 1).First(&device)
db.Where("device_name LIKE ?", "%核心%").Find(&devices)

// 分页查询
db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&list)

// 插入
db.Create(&device)

// 更新
db.Model(&device).Updates(map[string]interface{}{"status": "offline"})
db.Model(&device).Omit("password").Updates(&device)

// 删除
db.Delete(&device, deviceId)

// 事务
db.Transaction(func(tx *gorm.DB) error {
    // 业务逻辑
    return nil
})
```

### 7.3 数据库钩子

```go
// 创建前
func (e *SysDevice) BeforeCreate(tx *gorm.DB) error {
    // 数据验证、默认值设置
    return nil
}

// 更新前
func (e *SysDevice) BeforeUpdate(tx *gorm.DB) error {
    // 数据验证
    return nil
}

// 查询后
func (e *SysDevice) AfterFind(tx *gorm.DB) error {
    // 数据处理
    return nil
}
```

---

## 8. 路由注册

### 8.1 路由注册流程

```
1. 在 app/admin/router/ 目录创建路由文件
   ↓
2. 在 init() 函数中添加到 routerCheckRole
   ↓
3. 实现 register{Module}Router 函数
   ↓
4. 定义路由组和中间件
   ↓
5. 注册具体的路由和处理函数
```

### 8.2 路由注册示例

```go
package router

import (
    "github.com/gin-gonic/gin"
    jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"
    "opt-switch/app/admin/apis"
    "opt-switch/common/actions"
    "opt-switch/common/middleware"
)

// 1. 在 init 函数中注册
func init() {
    routerCheckRole = append(routerCheckRole, registerSysDeviceRouter)
}

// 2. 实现路由注册函数
func registerSysDeviceRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
    api := apis.SysDevice{}

    // 3. 定义路由组和中间件
    // authMiddleware.MiddlewareFunc() - JWT 认证
    // middleware.AuthCheckRole() - 角色权限检查
    // actions.PermissionAction() - 数据权限过滤
    r := v1.Group("/sys-device").
        Use(authMiddleware.MiddlewareFunc()).
        Use(middleware.AuthCheckRole()).
        Use(actions.PermissionAction())
    {
        // 4. 注册具体路由
        r.GET("", api.GetPage)         // 列表
        r.GET("/:id", api.Get)         // 详情
        r.POST("", api.Insert)         // 创建
        r.PUT("", api.Update)          // 更新
        r.DELETE("/:id", api.Delete)   // 删除
    }
}
```

### 8.3 中间件说明

| 中间件 | 说明 | 用途 |
|--------|------|------|
| authMiddleware.MiddlewareFunc() | JWT 认证 | 验证用户 Token |
| middleware.AuthCheckRole() | 角色检查 | 验证用户权限 |
| actions.PermissionAction() | 数据权限 | 过滤用户数据范围 |

---

## 9. 调试与测试

### 9.1 日志调试

```bash
# 查看实时日志
tail -f temp/logs/opt-switch.log

# 查看错误日志
tail -f temp/logs/error.log

# 搜索特定关键字
grep "ERROR" temp/logs/opt-switch.log
```

### 9.2 API 测试

#### 使用 curl

```bash
# 获取 Token
curl -X POST http://localhost:8000/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"123456"}'

# 使用 Token 访问 API
curl -X GET "http://localhost:8000/api/v1/sys-user?page=1&pageSize=10" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

#### 使用 Postman

1. 导入 Swagger 文档：`http://localhost:8000/swagger/admin/doc.json`
2. 设置环境变量：
   - `base_url`: http://localhost:8000
   - `token`: {{login_response.data.token}}
3. 创建请求并测试

### 9.3 数据库调试

```bash
# 使用 SQLite 命令行
sqlite3 opt-switch-db.db

# 查看所有表
.tables

# 查看表结构
.schema sys_device

# 查询数据
SELECT * FROM sys_device LIMIT 5;

# 退出
.quit
```

### 9.4 性能分析

```go
// 在代码中添加性能分析
import (
    "time"
    "github.com/gin-gonic/gin"
)

func (e SysDevice) GetPage(c *gin.Context) {
    start := time.Now()
    defer func() {
        e.Logger.Infof("GetPage 耗时: %v", time.Since(start))
    }()

    // 业务逻辑...
}
```

---

## 10. 构建与部署

### 10.1 本地构建

```bash
# 标准构建
go build -o opt-switch

# 优化构建（去除调试信息）
go build -ldflags="-w -s" -o opt-switch

# 交叉编译
env CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o opt-switch-arm64
```

### 10.2 使用 Makefile

```bash
# 查看所有命令
make help

# 常用命令
make build           # 当前平台
make build-arm64     # ARM64
make build-armv7     # ARMv7
make build-switch    # 所有交换机架构
```

### 10.3 测试构建结果

```bash
# 检查二进制文件信息
file opt-switch-arm64

# 检查文件大小
ls -lh opt-switch-arm64

# 测试运行
./opt-switch-arm64 server -c config/settings.minimal.yml
```

### 10.4 Docker 构建

```bash
# 构建镜像
docker build -t opt-switch:latest .

# 运行容器
docker run -d --name opt-switch -p 8000:8000 opt-switch:latest

# 查看日志
docker logs -f opt-switch
```

---

## 附录

### A. 常见问题

**Q: 代码生成器生成的代码在哪里？**

A: 生成的代码位于：
- Model: `app/admin/models/sys_{module}.go`
- Service: `app/admin/service/sys_{module}.go`
- DTO: `app/admin/service/dto/sys_{module}.go`
- API: `app/admin/apis/sys_{module}.go`
- Router: `app/admin/router/sys_{module}.go`

**Q: 如何添加自定义业务逻辑？**

A: 在 Service 层添加自定义方法，然后在 API 层调用。建议不要修改生成的代码，而是在生成的代码基础上扩展。

**Q: 如何实现复杂查询？**

A: 在 DTO 中定义查询条件，使用 `search` 标签配置查询方式，然后在 Service 层使用 GORM 的 Scope 功能。

**Q: 如何处理事务？**

A: 在 Service 层使用 `e.Orm.Transaction()` 方法包裹需要事务的操作。

### B. 参考资源

- [Gin 官方文档](https://gin-gonic.com/docs/)
- [GORM 官方文档](https://gorm.io/docs/)
- [Go-Admin-Core 文档](https://github.com/go-admin-team/go-admin-core)
- [Swagger 注解规范](https://swagger.io/docs/specification/2-0/)

### C. 最佳实践

1. **命名规范**：遵循 Go 语言的命名规范，使用驼峰命名法
2. **错误处理**：始终处理错误，使用 `e.Log.Errorf()` 记录错误日志
3. **数据验证**：在 DTO 中使用 `vd` 标签定义验证规则
4. **权限控制**：始终使用中间件保护需要认证的路由
5. **代码注释**：添加 Swagger 注释，方便生成 API 文档
6. **版本控制**：频繁提交，每次提交完成一个小功能
7. **测试覆盖**：为核心功能编写单元测试

---

<div align="center">

**Happy Coding!**

如有问题，请访问 [项目主页](https://www.opt-switch.dev) 或加入社区讨论

</div>
