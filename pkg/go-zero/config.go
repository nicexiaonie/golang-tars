package gozero

import (
	"fmt"
	"time"

	"google.golang.org/grpc"
)

// Config go-zero zrpc 客户端配置
type Config struct {
	// 连接目标（由 Target 设置，Etcd 与 Endpoints 二选一）
	Etcd      *EtcdConfig // Etcd 服务发现配置
	Endpoints []string    // 直连地址列表，如 ["127.0.0.1:8080"]

	// 应用配置
	App      string // 应用名称
	Token    string // 认证 Token
	NonBlock bool   // 是否非阻塞模式

	// 超时配置
	Timeout time.Duration // 请求超时时间

	// 重试配置
	MaxRetries int // 最大重试次数（0 表示不重试）

	// 连接池配置
	PoolSize int // 连接池大小（gRPC 连接数），默认 1

	// Keepalive 配置
	KeepaliveTime       time.Duration // Keepalive ping 间隔，默认 30s
	KeepaliveTimeout    time.Duration // Keepalive 超时时间，默认 10s
	PermitWithoutStream bool          // 无活跃流时是否发送 Keepalive ping

	// 消息大小限制
	MaxRecvMsgSize int // 最大接收消息大小（字节），默认 4MB
	MaxSendMsgSize int // 最大发送消息大小（字节），默认 4MB

	// 中间件配置
	EnableLog bool // 是否启用日志拦截器（默认 true）

	// 元数据（自动注入到每个 RPC 请求的 metadata 中）
	Metadata map[string]string

	// 自定义拦截器（追加到内置拦截器之后）
	UnaryInterceptors  []grpc.UnaryClientInterceptor
	StreamInterceptors []grpc.StreamClientInterceptor
}

// EtcdConfig Etcd 配置
type EtcdConfig struct {
	Hosts []string // Etcd 地址列表
	Key   string   // 服务注册的 Key

	// 认证配置（可选）
	User string // 用户名
	Pass string // 密码

	// TLS 配置（可选）
	CertFile           string // 客户端证书文件
	CertKeyFile        string // 客户端私钥文件
	CACertFile         string // CA 证书文件
	InsecureSkipVerify bool   // 是否跳过证书验证
}

// Target 连接目标（Etcd 服务发现 或 直连地址，二选一）
type Target struct {
	Etcd      *EtcdConfig
	Endpoints []string
}

// WithEtcd 创建 Etcd 服务发现目标
func WithEtcd(hosts []string, key string) Target {
	return Target{Etcd: &EtcdConfig{Hosts: hosts, Key: key}}
}

// WithEtcdConfig 创建带完整配置的 Etcd 服务发现目标
func WithEtcdConfig(etcd *EtcdConfig) Target {
	return Target{Etcd: etcd}
}

// WithEndpoints 创建直连目标
func WithEndpoints(endpoints ...string) Target {
	return Target{Endpoints: endpoints}
}

// 默认消息大小限制（4MB）
const defaultMaxMsgSize = 4 * 1024 * 1024

// Option 函数式选项，用于覆盖默认配置
type Option func(*Config)

// defaultConfig 返回默认配置（内部使用）
func defaultConfig() *Config {
	return &Config{
		NonBlock:         true,
		Timeout:          3 * time.Second,
		MaxRetries:       0,
		PoolSize:         1,
		KeepaliveTime:    30 * time.Second,
		KeepaliveTimeout: 10 * time.Second,
		MaxRecvMsgSize:   defaultMaxMsgSize,
		MaxSendMsgSize:   defaultMaxMsgSize,
		EnableLog:        true,
	}
}

// ==================== Option 函数 ====================

// WithApp 设置应用名称
func WithApp(app string) Option {
	return func(c *Config) { c.App = app }
}

// WithToken 设置认证 Token
func WithToken(token string) Option {
	return func(c *Config) { c.Token = token }
}

// WithNonBlock 设置是否非阻塞模式
func WithNonBlock(nonBlock bool) Option {
	return func(c *Config) { c.NonBlock = nonBlock }
}

// WithTimeout 设置请求超时时间
func WithTimeout(d time.Duration) Option {
	return func(c *Config) { c.Timeout = d }
}

// WithMaxRetries 设置最大重试次数
func WithMaxRetries(n int) Option {
	return func(c *Config) { c.MaxRetries = n }
}

// WithPoolSize 设置连接池大小
func WithPoolSize(n int) Option {
	return func(c *Config) { c.PoolSize = n }
}

// WithKeepalive 设置 Keepalive 参数
func WithKeepalive(time, timeout time.Duration, permitWithoutStream bool) Option {
	return func(c *Config) {
		c.KeepaliveTime = time
		c.KeepaliveTimeout = timeout
		c.PermitWithoutStream = permitWithoutStream
	}
}

// WithMaxRecvMsgSize 设置最大接收消息大小
func WithMaxRecvMsgSize(size int) Option {
	return func(c *Config) { c.MaxRecvMsgSize = size }
}

// WithMaxSendMsgSize 设置最大发送消息大小
func WithMaxSendMsgSize(size int) Option {
	return func(c *Config) { c.MaxSendMsgSize = size }
}

// WithLog 启用/禁用日志拦截器
func WithLog(enable bool) Option {
	return func(c *Config) { c.EnableLog = enable }
}

// WithMetadata 设置自动注入的 metadata 键值对
func WithMetadata(kvPairs map[string]string) Option {
	return func(c *Config) { c.Metadata = kvPairs }
}

// WithUnaryInterceptors 添加自定义 Unary 拦截器
func WithUnaryInterceptors(interceptors ...grpc.UnaryClientInterceptor) Option {
	return func(c *Config) {
		c.UnaryInterceptors = append(c.UnaryInterceptors, interceptors...)
	}
}

// WithStreamInterceptors 添加自定义 Stream 拦截器
func WithStreamInterceptors(interceptors ...grpc.StreamClientInterceptor) Option {
	return func(c *Config) {
		c.StreamInterceptors = append(c.StreamInterceptors, interceptors...)
	}
}

// Validate 验证配置合法性
func (c *Config) Validate() error {
	if c.Etcd == nil && len(c.Endpoints) == 0 {
		return fmt.Errorf("either etcd or endpoints must be configured")
	}
	if c.Etcd != nil && len(c.Endpoints) > 0 {
		return fmt.Errorf("etcd and endpoints cannot be configured at the same time")
	}
	if c.Etcd != nil {
		if len(c.Etcd.Hosts) == 0 {
			return fmt.Errorf("etcd hosts cannot be empty")
		}
		if c.Etcd.Key == "" {
			return fmt.Errorf("etcd key cannot be empty")
		}
	}
	if c.PoolSize < 1 {
		return fmt.Errorf("pool_size must be at least 1, got %d", c.PoolSize)
	}
	if c.PoolSize > 64 {
		return fmt.Errorf("pool_size must not exceed 64, got %d", c.PoolSize)
	}
	return nil
}
