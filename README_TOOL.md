[toc]

# 工具与组件使用手册

本文档详细介绍 go-base 项目中各工具、组件的安装配置和使用方法。

---

## 环境要求

| 工具 | 版本 | 说明 |
|------|------|------|
| Go | 1.25.0+ | 通过 asdf 管理版本 |
| Protoc | 3.x+ | Protocol Buffers 编译器 |
| Buf | latest | Proto 构建工具 |
| Wire | 0.7.0+ | 依赖注入代码生成 |
| grds-gen | latest | 数据库模型生成 |

---

## 环境版本管理 (asdf)

项目使用 [asdf](https://asdf-vm.com/) 统一管理开发工具版本，版本定义在 `.tool-versions` 文件中。

```bash
# 安装 asdf（macOS）
brew install asdf

# 添加 Go 插件
asdf plugin add golang

# 查看可用版本
asdf list all golang

# 安装指定版本
asdf install golang 1.25.0

# 设置项目版本（写入 .tool-versions）
asdf set golang 1.25.0

# 设置全局默认版本
asdf set -u golang 1.25.0
```

> asdf 会从当前目录逐层向上查找 `.tool-versions` 文件。若需全局默认版本，在 `$HOME/.tool-versions` 中定义。

---

## 双协议支持：gRPC + HTTP/JSON

### 核心特性

- **一套代码，两种协议**：基于 Connect-go + Vanguard 实现
- **单端口服务**：gRPC 和 HTTP/JSON 共享同一端口
- **零转换开销**：无需手动 JSON ↔ Protobuf 双向转换
- **RESTful 映射**：通过 proto 注解自动生成 RESTful API

### 架构对比

#### grpc-gateway 方案（未采用）

```
┌─────────────────────────────────────────┐
│              HTTP/JSON Client           │
└───────────────┬─────────────────────────┘
                │ HTTP/1.1 + JSON
                ▼
┌─────────────────────────────────────────┐
│           gRPC-Gateway                  │
│  ┌─────────────────────────────────┐    │
│  │  HTTP/JSON ↔ Protobuf 转换      │    │
│  └─────────────────────────────────┘    │
│                 │                        │
│                 ▼                        │
│  ┌─────────────────────────────────┐    │
│  │  代理转发到 gRPC 服务           │    │
│  └─────────────────────────────────┘    │
└───────────────┬─────────────────────────┘
                │ gRPC/HTTP2
                ▼
┌─────────────────────────────────────────┐
│             gRPC 服务                    │
└─────────────────────────────────────────┘
```

#### Connect-go 方案（当前采用）

```
┌─────────────────────────────────────────┐
│          客户端（多种协议）              │
│  ┌────────┬────────┬─────────────────┐  │
│  │ gRPC   │Connect │  HTTP/JSON      │  │
│  │ Client │Client  │  Client         │  │
│  └────────┴────────┴─────────────────┘  │
└───────────────┬─────────────────────────┘
                │
                ▼
┌─────────────────────────────────────────┐
│          Connect-Go 服务                │
│  ┌─────────────────────────────────┐    │
│  │    多协议统一处理层              │    │
│  │  ┌─────┬─────────┬──────────┐   │    │
│  │  │gRPC │ Connect │ HTTP/JSON│   │    │
│  │  │协议 │ 协议    │ 协议     │   │    │
│  │  └─────┴─────────┴──────────┘   │    │
│  └─────────────────────────────────┘    │
│                 │                        │
│                 ▼                        │
│  ┌─────────────────────────────────┐    │
│  │     业务逻辑实现                 │    │
│  └─────────────────────────────────┘    │
└─────────────────────────────────────────┘
```

**Connect-go 优势：**
- 零拷贝处理：直接处理原始协议
- 单次序列化：根据协议只进行一次序列化
- 减少网络跳数：客户端直接连接服务
- 统一连接池：复用相同连接

---

## Protobuf 工具链

### 安装依赖

```bash
# protoc 插件
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/envoyproxy/protoc-gen-validate@latest

# Connect-go 插件（如不使用 buf）
go install connectrpc.com/connect/cmd/protoc-gen-connect-go@latest
```

### 使用 buf 管理 Proto

```bash
# 安装 buf
brew install bufbuild/buf/buf

# 或使用 Docker
docker pull bufbuild/buf
```

### Proto 编译

```bash
# 编译所有 proto 文件
make proto

# 清理生成的 proto 代码
make proto-clean
```

生成文件输出到 `pkg/proto/` 目录，包含：
- `*.pb.go` — Protobuf 消息定义
- `*_grpc.pb.go` — gRPC 服务接口
- `*.pb.validate.go` — 字段验证代码
- `*connect/*.connect.go` — Connect-go 处理器接口

### Proto 定义示例

```protobuf
// proto/demo/v1/demo.proto
service DemoService {
  rpc CreateDemo(CreateDemoRequest) returns (CreateDemoResponse) {
    option (google.api.http) = {
      post: "/api/v1/demo"
      body: "*"
    };
  }
  
  rpc GetDemo(GetDemoRequest) returns (GetDemoResponse) {
    option (google.api.http) = {
      get: "/api/v1/demo/{id}"
    };
  }
}
```

### 生成 OpenAPI 文档

```bash
# 安装 buf 后执行
make openapi
```

生成文件位于 `docs/openapi/openapi.yaml`，可导入 Swagger UI 或 Apifox 查看。

---

## Wire 依赖注入

### 安装

```bash
go install github.com/google/wire/cmd/wire@latest
```

### 自动生成配置

本项目提供自动扫描和生成 Wire 配置的脚本，无需手动维护 `wire.go` 文件。

```bash
# 方式一：一键生成（推荐）
make wire-all

# 方式二：分步执行
make wire-gen   # 自动扫描并生成 wire.go
make wire       # 运行 wire 生成 wire_gen.go
```

### 扫描规则

脚本会自动扫描以下内容并生成 `wire.go`：

| 扫描目录 | 匹配规则 | 说明 |
|----------|---------|------|
| `internal/*/data/mysql/dao/` | `New*` 函数 | DAO 层构造函数 |
| `internal/*/service/` | `New*` 函数 | Service 层构造函数 |
| `internal/*/handler/` | `New*` 函数 | Handler 层构造函数 |
| `internal/*/proxy/` | `New*` 函数 | Proxy 层构造函数 |
| `internal/*/provider.go` | `Provide*` 函数 | Provider 函数 |

### 使用场景

**添加新功能时：**

1. 创建 DAO/Service/Handler，函数名以 `New` 开头
2. 运行 `make wire-all`
3. Wire 会自动发现并生成依赖注入代码

**示例：**

```go
// internal/user/service/article_service.go
func NewArticleService(articleDAO *dao.ArticleDAO) *ArticleService {
    return &ArticleService{articleDAO: articleDAO}
}
```

运行 `make wire-all` 后，`NewArticleService` 会自动添加到 `wire.go` 的 ProviderSet 中。

---

## 数据库 (MySQL + GORM/grds)

### 安装 grds

```bash
go get -u github.com/nicexiaonie/grds
```

### 模型文件生成

```bash
# 安装模型生成工具
go install github.com/nicexiaonie/grds/cmd/grds-gen@latest

# 在项目根目录或包含 .grds.yaml 的目录执行
grds-gen
```

配置文件 `.grds.yaml` 定义数据库连接和生成规则，项目根目录和 `internal/*/data/mysql/` 下各有一份。

### 数据库初始化

```go
import "go-base/pkg/db"

dbManager := db.GetManager()

dbConfig := &db.DBConfig{
    Host:         "localhost",
    Port:         3306,
    Username:     "root",
    Password:     "password",
    Database:     "mydb",
    MaxIdleConns: 100,
    MaxOpenConns: 200,
}

err := dbManager.InitDB("default", dbConfig)
```

### 多数据库管理

```go
dbManager := db.GetManager()

// 初始化多个数据库
dbManager.InitDB("main", &db.DBConfig{...})
dbManager.InitDB("analytics", &db.DBConfig{...})

// 使用不同的数据库
mainClient, _ := dbManager.GetClient("main")
analyticsClient, _ := dbManager.GetClient("analytics")
```

### DAO 使用

```go
// 获取 grds.Client
client, _ := db.GetManager().GetClient("default")

// 创建 DAO（通过 Wire 自动注入）
userDAO := dao.NewUserDAO(client)

// 业务方法
user, _ := userDAO.FindByUsername("test")

// 或直接使用 grds.Client
client.Create(&user)
client.Find(&users)
```

### 分页查询

```go
import "go-base/pkg/db"

client, _ := dbManager.GetClient("default")

var users []models.Users
result, err := db.Paginate(client.DB(), &users, 1, 10, func(db *gorm.DB) *gorm.DB {
    return db.Where("status = ?", 1)
})

fmt.Println(result.Total)      // 总记录数
fmt.Println(result.TotalPages) // 总页数
```

### 事务

```go
client, _ := dbManager.GetClient("default")

err := client.Transaction(func(tx *gorm.DB) error {
    if err := tx.Create(&user1).Error; err != nil {
        return err
    }
    if err := tx.Create(&user2).Error; err != nil {
        return err
    }
    return nil // 返回 nil 自动提交
})
```

### Manager API 参考

```go
InitDB(name string, config *DBConfig) error       // 初始化数据库
GetClient(name string) (*grds.Client, error)       // 获取客户端
GetDefaultClient() (*grds.Client, error)           // 获取默认客户端
Close(name string) error                           // 关闭连接
CloseAll() error                                   // 关闭所有连接
ListInstances() []string                           // 列出所有实例
```

---

## gRPC 客户端 (go-zero)

### 初始化客户端

#### 直连模式

```go
err := gozero.Init("order-rpc",
    gozero.WithEndpoints("127.0.0.1:8080", "127.0.0.1:8081"),
    gozero.WithTimeout(3*time.Second),
    gozero.WithPoolSize(3),           // 连接池大小，默认为1；>1时 round-robin 分发
    gozero.WithKeepalive(
        30*time.Second,  // ping 间隔
        10*time.Second,  // 超时时间
        false,           // 无活跃流时是否发送 ping
    ),
    gozero.WithMaxRecvMsgSize(10*1024*1024),  // 10MB
    gozero.WithMaxSendMsgSize(10*1024*1024),
    gozero.WithMaxRetries(3),  // 最大重试 3 次（指数退避 + 随机抖动）
    gozero.WithMetadata(map[string]string{
        "app-id":      "my-app",
        "app-version": "1.0.0",
    }),
    gozero.WithUnaryInterceptors(myInterceptor()),
)
```

**重试策略：**
- 仅重试 `Unavailable` 和 `Internal` 错误
- 指数退避：100ms → 200ms → 400ms → 800ms（最大 5s）
- 添加随机抖动避免惊群效应

#### 自定义拦截器

```go
func myInterceptor() grpc.UnaryClientInterceptor {
    return func(ctx context.Context, method string, req, reply interface{},
        cc *grpc.ClientConn, invoker grpc.UnaryInvoker,
        opts ...grpc.CallOption) error {
        // 自定义逻辑
        return invoker(ctx, method, req, reply, cc, opts...)
    }
}
```

### 创建 gRPC Service Client

```go
// 方式一：使用 MustConn（推荐，适合初始化后的场景）
userSvc := pb.NewUserServiceClient(gozero.MustConn("user-rpc"))
resp, err := userSvc.GetUser(ctx, req)

// 方式二：使用 GetClient（适合需要错误处理的场景）
client, err := gozero.GetClient("user-rpc")
if err != nil {
    return err
}
userSvc := pb.NewUserServiceClient(client.Conn())
```

### 客户端管理

```go
// 获取客户端（带错误处理）
client, err := gozero.GetClient("user-rpc")

// 获取客户端（panic if not found）
client := gozero.MustGetClient("user-rpc")

// 列出所有客户端
manager := gozero.GetManager()
names := manager.ListClients()

// 查看连接池统计
client := gozero.MustGetClient("user-rpc")
stats := client.Stats()
fmt.Printf("连接池大小: %d, 累计调用: %d\n", stats.PoolSize, stats.Calls)
```

---

## JWT 认证

### 配置

JWT 通过 Tars 配置文件加载：

```yaml
jwt:
  secret_key: "your-secret-key"
  issuer: "go-base"
  expire_hours: 24
```

### 使用

```go
import "go-base/pkg/jwt"

// 初始化
jwtClient := jwt.NewClient(config)

// 生成 Token
token, err := jwtClient.GenerateToken(userID, username)

// 验证 Token
claims, err := jwtClient.ValidateToken(tokenString)
```

### 拦截器集成

认证拦截器自动从请求头提取并验证 JWT Token：

```go
// 在 interceptor 中配置白名单（无需认证的接口）
whiteList := []string{
    "CreateDemo",  // 示例：跳过认证
}
```

---

## Redis

### 初始化

```go
import "go-base/pkg/redis"

err := redis.Init(&redis.Config{
    Addr:     "localhost:6379",
    Password: "",
    DB:       0,
})

// 获取客户端
client := redis.GetClient()
```

---

## 日志

项目使用 [glog](https://github.com/nicexiaonie/glog)（基于 Logrus 封装），支持：
- 结构化日志
- 请求 ID / 追踪 ID 传递
- 上下文感知日志

```go
import "go-base/pkg"

// 初始化
pkg.LoggerInit()

// 使用
pkg.Logger.Info("message")
pkg.Logger.WithField("key", "value").Error("error message")

// 带上下文的日志
pkg.LogEntry(ctx).Info("request processed")
```

---

## 构建与部署

### Makefile 命令

```bash
# Proto 相关
make proto          # 编译所有 proto 文件
make proto-clean    # 清理生成的 proto 代码
make openapi        # 生成 OpenAPI 文档

# Wire 相关
make wire-gen       # 自动扫描生成 wire.go
make wire           # 运行 wire 生成 wire_gen.go
make wire-all       # wire-gen + wire 一键执行

# 构建部署
make build-user     # 构建 user 服务
make upload-user    # 构建并上传 user 服务到 Tars

# 依赖
make deps           # 安装项目依赖
```

### Tars 部署

```bash
# 进入服务目录
cd cmd/user

# 构建并上传到 Tars 平台
make upload
```

部署流程：
1. `make upload` 编译 Go 二进制
2. 打包为 `.tgz` 文件
3. 通过 Tars API 上传到 Tars 管理平台
4. Tars 平台自动发布和管理服务实例

---

## IDE 配置

### VSCode / Cursor

项目支持 `.vscode/tasks.json` 配置任务快捷操作，可通过 Task Manager 扩展使用。

常用任务：
- 编译 Proto 文件
- 生成 Wire 代码
- 构建服务
- 部署服务

---

## 参考链接

| 组件 | 文档 |
|------|------|
| Connect-go | https://connectrpc.com/docs/go/getting-started |
| Vanguard | https://github.com/connectrpc/vanguard-go |
| TarsGo | https://github.com/TarsCloud/TarsGo |
| Google Wire | https://github.com/google/wire |
| GORM | https://gorm.io/docs/ |
| grds | https://github.com/nicexiaonie/grds |
| go-zero | https://go-zero.dev/ |
| Buf | https://buf.build/docs/ |
| asdf | https://asdf-vm.com/ |
| protoc-gen-validate | https://github.com/envoyproxy/protoc-gen-validate |
