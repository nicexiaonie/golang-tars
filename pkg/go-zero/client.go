package gozero

import (
	"fmt"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

var (
	// 全局服务管理器
	globalManager *Manager
	once          sync.Once
)

// Manager go-zero RPC 服务管理器 - 管理多个 zrpc.Client 实例
type Manager struct {
	clients map[string]*Client // 客户端集合
	mu      sync.RWMutex       // 读写锁
}

// Client go-zero RPC 客户端封装
// 支持连接池模式（PoolSize > 1 时启用）
type Client struct {
	pool   *connPool // 连接池（管理一个或多个 zrpc.Client）
	config *Config   // 配置
	name   string    // 服务名称
}

// GetManager 获取全局服务管理器（单例）
func GetManager() *Manager {
	once.Do(func() {
		globalManager = &Manager{
			clients: make(map[string]*Client),
		}
	})
	return globalManager
}

// Init 初始化一个 RPC 客户端连接
//
// name: 服务名称（如 "user-rpc", "order-rpc"）
// target: 连接目标（Etcd 服务发现 或 直连地址），使用 WithEtcd / WithEndpoints 构造
// opts: 可选配置项，使用 WithXxx 函数式选项覆盖默认值
//
// 示例:
//
//	gozero.Init("user-rpc",
//	    gozero.WithEtcd([]string{"127.0.0.1:2379"}, "user.rpc"),
//	    gozero.WithTimeout(5*time.Second),
//	    gozero.WithPoolSize(3),
//	)
func (m *Manager) Init(name string, target Target, opts ...Option) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 检查是否已存在
	if _, exists := m.clients[name]; exists {
		return fmt.Errorf("rpc client '%s' already exists", name)
	}

	// 构建配置：默认值 + Target + 用户选项
	config := defaultConfig()
	config.Etcd = target.Etcd
	config.Endpoints = target.Endpoints
	for _, opt := range opts {
		opt(config)
	}

	// 验证配置
	if err := config.Validate(); err != nil {
		return fmt.Errorf("invalid config for '%s': %w", name, err)
	}

	// 构建 zrpc 客户端配置
	zrpcConf := buildZrpcConfig(config)

	// 构建客户端选项（拦截器、Keepalive、消息大小等）
	clientOpts := buildClientOptions(config)

	// 创建连接池（poolSize 个 zrpc.Client）
	pool, err := newConnPool(config.PoolSize, func() (zrpc.Client, error) {
		return zrpc.NewClient(zrpcConf, clientOpts...)
	})
	if err != nil {
		return fmt.Errorf("failed to create rpc client pool '%s': %w", name, err)
	}

	// 保存客户端
	m.clients[name] = &Client{
		pool:   pool,
		config: config,
		name:   name,
	}

	return nil
}

// GetClient 获取指定名称的 RPC 客户端
func (m *Manager) GetClient(name string) (*Client, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	client, exists := m.clients[name]
	if !exists {
		return nil, fmt.Errorf("rpc client '%s' not found", name)
	}

	return client, nil
}

// MustGetClient 获取指定名称的 RPC 客户端（不存在则 panic）
func (m *Manager) MustGetClient(name string) *Client {
	client, err := m.GetClient(name)
	if err != nil {
		panic(err)
	}
	return client
}

// Close 关闭指定 RPC 客户端，释放底层 gRPC 连接
func (m *Manager) Close(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	client, exists := m.clients[name]
	if !exists {
		return fmt.Errorf("rpc client '%s' not found", name)
	}

	client.pool.close()
	delete(m.clients, name)
	return nil
}

// CloseAll 关闭所有 RPC 客户端，释放全部底层 gRPC 连接
func (m *Manager) CloseAll() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, client := range m.clients {
		client.pool.close()
	}
	m.clients = make(map[string]*Client)
}

// ListClients 列出所有已注册的客户端名称
func (m *Manager) ListClients() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	names := make([]string, 0, len(m.clients))
	for name := range m.clients {
		names = append(names, name)
	}
	return names
}

// buildZrpcConfig 构建 zrpc 客户端配置
func buildZrpcConfig(config *Config) zrpc.RpcClientConf {
	conf := zrpc.RpcClientConf{
		NonBlock: config.NonBlock,
		Timeout:  int64(config.Timeout / time.Millisecond),
	}

	// 设置应用信息
	if config.App != "" {
		conf.App = config.App
	}
	if config.Token != "" {
		conf.Token = config.Token
	}

	// 配置服务发现或直连
	if config.Etcd != nil {
		conf.Etcd = discov.EtcdConf{
			Hosts: config.Etcd.Hosts,
			Key:   config.Etcd.Key,
			User:  config.Etcd.User,
			Pass:  config.Etcd.Pass,
		}
		if config.Etcd.CertFile != "" {
			conf.Etcd.CertFile = config.Etcd.CertFile
		}
		if config.Etcd.CertKeyFile != "" {
			conf.Etcd.CertKeyFile = config.Etcd.CertKeyFile
		}
		if config.Etcd.CACertFile != "" {
			conf.Etcd.CACertFile = config.Etcd.CACertFile
		}
		conf.Etcd.InsecureSkipVerify = config.Etcd.InsecureSkipVerify
	} else if len(config.Endpoints) > 0 {
		conf.Endpoints = config.Endpoints
	}

	return conf
}

// buildClientOptions 构建客户端选项（拦截器、Keepalive、消息大小等）
func buildClientOptions(config *Config) []zrpc.ClientOption {
	var opts []zrpc.ClientOption

	// ==================== Keepalive 配置 ====================
	if config.KeepaliveTime > 0 || config.KeepaliveTimeout > 0 {
		kaParams := keepalive.ClientParameters{
			Time:                config.KeepaliveTime,
			Timeout:             config.KeepaliveTimeout,
			PermitWithoutStream: config.PermitWithoutStream,
		}
		opts = append(opts, zrpc.WithDialOption(grpc.WithKeepaliveParams(kaParams)))
	}

	// ==================== 消息大小限制 ====================
	if config.MaxRecvMsgSize > 0 {
		opts = append(opts, zrpc.WithDialOption(grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(config.MaxRecvMsgSize),
		)))
	}
	if config.MaxSendMsgSize > 0 {
		opts = append(opts, zrpc.WithDialOption(grpc.WithDefaultCallOptions(
			grpc.MaxCallSendMsgSize(config.MaxSendMsgSize),
		)))
	}

	// ==================== 拦截器 ====================
	var interceptors []grpc.UnaryClientInterceptor

	// 内置拦截器
	if config.EnableLog {
		interceptors = append(interceptors, LogInterceptor())
	}
	if config.MaxRetries > 0 {
		interceptors = append(interceptors, RetryInterceptor(config.MaxRetries))
	}
	if len(config.Metadata) > 0 {
		interceptors = append(interceptors, MetadataInterceptor(config.Metadata))
	}

	// 用户自定义拦截器（追加到内置拦截器之后）
	interceptors = append(interceptors, config.UnaryInterceptors...)

	// 添加 Unary 拦截器选项
	for _, interceptor := range interceptors {
		opts = append(opts, zrpc.WithUnaryClientInterceptor(interceptor))
	}

	// 添加 Stream 拦截器选项
	for _, interceptor := range config.StreamInterceptors {
		opts = append(opts, zrpc.WithDialOption(grpc.WithChainStreamInterceptor(interceptor)))
	}

	return opts
}

// ==================== Client 方法 ====================

// Conn 获取 gRPC 客户端连接
// 当 PoolSize > 1 时，返回连接池代理（自动 round-robin 分发每次 RPC 调用）
// 当 PoolSize == 1 时，直接返回底层连接（零开销）
func (c *Client) Conn() grpc.ClientConnInterface {
	return c.pool.getPooledConn()
}

// GetZrpcClient 获取一个原始 zrpc.Client
// 注意：当 PoolSize > 1 时，返回连接池中的第一个客户端
func (c *Client) GetZrpcClient() zrpc.Client {
	clients := c.pool.getClients()
	if len(clients) > 0 {
		return clients[0]
	}
	return nil
}

// GetZrpcClients 获取连接池中所有 zrpc.Client
func (c *Client) GetZrpcClients() []zrpc.Client {
	return c.pool.getClients()
}

// GetName 获取客户端名称
func (c *Client) GetName() string {
	return c.name
}

// GetConfig 获取配置
func (c *Client) GetConfig() *Config {
	return c.config
}

// GetPoolSize 获取连接池大小
func (c *Client) GetPoolSize() int {
	return c.pool.getSize()
}

// Stats 获取连接池统计信息
func (c *Client) Stats() PoolStats {
	return c.pool.stats(c.name)
}

// ==================== 包级别便捷函数 ====================

// Init 初始化一个 RPC 客户端（使用全局管理器）
func Init(name string, target Target, opts ...Option) error {
	return GetManager().Init(name, target, opts...)
}

// GetClient 获取指定名称的 RPC 客户端（使用全局管理器）
func GetClient(name string) (*Client, error) {
	return GetManager().GetClient(name)
}

// MustGetClient 获取指定名称的 RPC 客户端，不存在则 panic（使用全局管理器）
func MustGetClient(name string) *Client {
	return GetManager().MustGetClient(name)
}

// Conn 获取指定服务的 gRPC 连接（使用全局管理器）
func Conn(name string) (grpc.ClientConnInterface, error) {
	client, err := GetManager().GetClient(name)
	if err != nil {
		return nil, err
	}
	return client.Conn(), nil
}

// MustConn 获取指定服务的 gRPC 连接，不存在则 panic（使用全局管理器）
// 适用于初始化后确定存在的场景，是最常用的便捷入口
//
// 使用示例:
//
//	userSvc := pb.NewUserServiceClient(gozero.MustConn("user-rpc"))
//	resp, err := userSvc.GetUser(ctx, req)
func MustConn(name string) grpc.ClientConnInterface {
	return GetManager().MustGetClient(name).Conn()
}
