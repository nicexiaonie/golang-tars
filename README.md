# golang-tars

基于 TarsGo 的微服务基础代码库，作为新项目的起始模板。

## 特性

- **双协议支持** — gRPC + HTTP/JSON 共享单端口，基于 Connect-go + Vanguard
- **RESTful 映射** — Proto 注解自动生成 RESTful API，无需手动转换
- **依赖注入** — Google Wire 自动扫描生成，新增功能零配置接入
- **多数据库管理** — GORM/grds 封装，支持多实例、分页、事务
- **RPC 客户端** — go-zero 封装，连接池 + 重试 + 服务发现
- **统一中间件** — 认证(JWT)、日志、恢复、错误处理拦截器链
- **代码生成** — Proto、Wire、数据库模型一键生成

## 技术栈

| 类别 | 技术 |
|------|------|
| 语言 | Go 1.25 |
| 框架 | TarsGo |
| 协议 | gRPC / Connect-go / Vanguard |
| 数据库 | MySQL (GORM / grds) |
| 缓存 | Redis |
| 依赖注入 | Google Wire |
| RPC 客户端 | go-zero |
| 认证 | JWT |
| Proto 工具 | Buf / protoc-gen-validate |

## 项目结构

```
golang-tars/
├── cmd/                          # 应用入口
│   ├── user/                     #   User 服务（主服务示例）
│   └── demo/                     #   Demo 服务配置示例
├── internal/                     # 内部业务逻辑（不对外暴露）
│   ├── user/
│   │   ├── server.go             #   服务器启动（Tars + Connect）
│   │   ├── handler/              #   协议适配层（gRPC/HTTP → 业务）
│   │   ├── service/              #   业务逻辑层
│   │   ├── data/mysql/           #   数据访问层（DAO + Models）
│   │   ├── proxy/                #   外部服务代理
│   │   └── interceptor/          #   拦截器（认证/日志/恢复）
│   └── demo/                     #   Demo 服务（参考实现）
├── pkg/                          # 公共包（可跨项目复用）
│   ├── proto/                    #   生成的 Protobuf 代码
│   ├── db/                       #   数据库管理（多实例 + 分页）
│   ├── go-zero/                  #   gRPC 客户端管理
│   ├── jwt/                      #   JWT 认证
│   ├── redis/                    #   Redis 客户端
│   ├── response/                 #   统一响应与错误码
│   └── ...                       #   其他工具包
├── proto/                        # Protobuf 定义文件
│   ├── demo/v1/                  #   Demo 服务 API
│   └── common/v1/                #   公共消息
├── scripts/                      # 构建脚本
├── docs/                         # 生成的文档（OpenAPI）
└── Makefile                      # 构建入口
```

## 请求处理流程

```
客户端请求 (HTTP/JSON 或 gRPC)
    ↓
Tars HTTP Mux (h2c)
    ↓
Vanguard 转码器 (RESTful → gRPC)
    ↓
Interceptors (Recovery → Log → Auth)
    ↓
Handler (协议适配)
    ↓
Service (业务逻辑)
    ↓
DAO / Proxy (数据库 / 外部服务)
```

## 快速开始

```bash
# 安装依赖工具
make deps

# 编译 Proto 文件
make proto

# 生成 Wire 依赖注入代码
make wire-all

# 构建 User 服务
make build-user

# 部署到 Tars
make upload-user
```

## 新增功能流程

1. 在 `proto/` 下定义 API，运行 `make proto`
2. 在 `internal/user/` 下创建 Handler → Service → DAO
3. 运行 `make wire-all` 自动接入依赖注入
4. 构建部署 `make upload-user`

## 文档

- [工具与组件使用手册](README_TOOL.md) — 详细的安装配置和各组件使用说明
- [数据库模块文档](pkg/db/README.md) — 数据库管理模块详细说明

## 已知问题

详见项目分析报告，主要包括：

- 部分 `context.WithValue` 使用字符串 key（应使用类型化 key）
- DAO 层接收 context 但未传递给底层操作
- `sms_log_dao.go` 命名与实际类型不匹配
- 模型字段拼写错误（`ReigsterTime`、`LoginResule`）
- 部分错误未正确处理或包装
- 缺少单元测试
