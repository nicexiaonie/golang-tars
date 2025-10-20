# Consul 封装使用文档

## 功能特性

- ✅ **服务注册/注销** - 自动注册服务到 Consul，支持优雅关闭时自动注销
- ✅ **服务发现** - 发现并获取健康的服务实例
- ✅ **健康检查** - 支持 HTTP/TCP/GRPC 健康检查
- ✅ **KV 存储** - 支持 Consul KV 存储的读写删除
- ✅ **监听服务变化** - 实时监听服务实例变化
- ✅ **监听 KV 变化** - 实时监听 KV 配置变化

## 快速开始

### 1. 环境变量配置

```bash
export CONSUL_ADDR="127.0.0.1:8500"  # Consul 地址
export SERVICE_ADDR="127.0.0.1"      # 服务地址
```

### 2. 在 server.go 中已自动集成

服务启动时会自动：
- 连接到 Consul
- 注册服务
- 配置健康检查（TCP 检查）
- 注册优雅关闭处理

### 3. 服务注册示例

```go
consulConfig := &consul.Config{
    Address:        "127.0.0.1:8500",
    Scheme:         "http",
    ServiceName:    "Demo.UserServer",
    ServiceID:      "Demo.UserServer-127.0.0.1-10015",
    ServiceAddress: "127.0.0.1",
    ServicePort:    10015,
    Tags:           []string{"Demo", "UserServer", "tars"},
    HealthCheck: &consul.HealthCheck{
        CheckID:                        "health-Demo.UserServer",
        Type:                           "tcp",
        TCP:                            "127.0.0.1:10015",
        Interval:                       10 * time.Second,
        Timeout:                        3 * time.Second,
        DeregisterCriticalServiceAfter: 30 * time.Second,
    },
}

client, err := consul.NewClient(consulConfig)
if err != nil {
    panic(err)
}

// 注册服务
err = client.RegisterService()
```

### 4. 服务发现示例

```go
// 发现健康的服务实例
services, err := client.DiscoverService("Demo.UserServer", true)
if err != nil {
    log.Fatal(err)
}

for _, service := range services {
    fmt.Printf("Service: %s:%d\n", service.Service.Address, service.Service.Port)
}

// 获取一个服务地址（负载均衡）
addr, err := client.GetServiceAddress("Demo.UserServer")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Service address: %s\n", addr)
```

### 5. 监听服务变化

```go
stopCh := make(chan struct{})
serviceCh, err := client.WatchService("Demo.UserServer", stopCh)
if err != nil {
    log.Fatal(err)
}

go func() {
    for services := range serviceCh {
        fmt.Printf("Service instances changed: %d instances\n", len(services))
        for _, service := range services {
            fmt.Printf("  - %s:%d\n", service.Service.Address, service.Service.Port)
        }
    }
}()
```

### 6. KV 存储操作

```go
// 设置值
err := client.PutKV("config/database/host", []byte("localhost"))

// 获取值
value, err := client.GetKV("config/database/host")
fmt.Printf("Value: %s\n", string(value))

// 删除值
err := client.DeleteKV("config/database/host")
```

### 7. 监听 KV 变化

```go
stopCh := make(chan struct{})
kvCh, err := client.WatchKV("config/database/host", stopCh)
if err != nil {
    log.Fatal(err)
}

go func() {
    for value := range kvCh {
        fmt.Printf("KV changed: %s\n", string(value))
    }
}()
```

## 健康检查类型

### TCP 检查
```go
HealthCheck: &consul.HealthCheck{
    Type:     "tcp",
    TCP:      "127.0.0.1:10015",
    Interval: 10 * time.Second,
    Timeout:  3 * time.Second,
}
```

### HTTP 检查
```go
HealthCheck: &consul.HealthCheck{
    Type:     "http",
    HTTP:     "http://127.0.0.1:8080/health",
    Interval: 10 * time.Second,
    Timeout:  3 * time.Second,
}
```

### GRPC 检查
```go
HealthCheck: &consul.HealthCheck{
    Type:     "grpc",
    GRPC:     "127.0.0.1:50051",
    Interval: 10 * time.Second,
    Timeout:  3 * time.Second,
}
```

## 优雅关闭

服务会在接收到 `SIGINT` 或 `SIGTERM` 信号时自动从 Consul 注销服务：

```go
// 已在 server.go 中自动实现
setupGracefulShutdown()
```

## 注意事项

1. 确保 Consul 服务已启动并可访问
2. 服务注册前需要确保端口已监听
3. 健康检查失败会导致服务自动注销
4. 建议在生产环境配置合适的健康检查间隔
