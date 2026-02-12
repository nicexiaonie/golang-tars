package redis

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	defaultClient *Client
	mu            sync.RWMutex
)

// Client Redis客户端封装
type Client struct {
	client redis.UniversalClient
	config *Config
	ctx    context.Context
	cancel context.CancelFunc
}

// Init 初始化默认Redis客户端
// 使用互斥锁替代 sync.Once，允许初始化失败后重试
func Init(config *Config) error {
	mu.Lock()
	defer mu.Unlock()
	if defaultClient != nil {
		return nil // 已初始化成功，跳过
	}
	client, err := NewClient(config)
	if err != nil {
		return err
	}
	defaultClient = client
	return nil
}

// GetClient 获取默认Redis客户端
func GetClient() *Client {
	mu.RLock()
	defer mu.RUnlock()
	return defaultClient
}

// NewClient 创建新的Redis客户端
func NewClient(config *Config) (*Client, error) {
	if config == nil {
		config = DefaultConfig()
	}

	ctx, cancel := context.WithCancel(context.Background())

	var client redis.UniversalClient

	if config.ClusterMode && len(config.Addrs) > 0 {
		// 集群模式
		client = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:        config.Addrs,
			Password:     config.Password,
			PoolSize:     config.PoolSize,
			MinIdleConns: config.MinIdleConns,
			MaxRetries:   config.MaxRetries,
			DialTimeout:  config.DialTimeout,
			ReadTimeout:  config.ReadTimeout,
			WriteTimeout: config.WriteTimeout,
			PoolTimeout:  config.PoolTimeout,
		})
	} else {
		// 单机模式
		client = redis.NewClient(&redis.Options{
			Addr:         fmt.Sprintf("%s:%d", config.Host, config.Port),
			Password:     config.Password,
			DB:           config.DB,
			PoolSize:     config.PoolSize,
			MinIdleConns: config.MinIdleConns,
			MaxRetries:   config.MaxRetries,
			DialTimeout:  config.DialTimeout,
			ReadTimeout:  config.ReadTimeout,
			WriteTimeout: config.WriteTimeout,
			PoolTimeout:  config.PoolTimeout,
		})
	}

	c := &Client{
		client: client,
		config: config,
		ctx:    ctx,
		cancel: cancel,
	}

	// 测试连接
	if err := c.Ping(); err != nil {
		return nil, fmt.Errorf("redis ping failed: %w", err)
	}

	return c, nil
}

// Ping 测试连接
func (c *Client) Ping() error {
	return c.client.Ping(c.ctx).Err()
}

// Close 关闭连接
func (c *Client) Close() error {
	if c.cancel != nil {
		c.cancel()
	}
	return c.client.Close()
}

// GetRawClient 获取原始Redis客户端
func (c *Client) GetRawClient() redis.UniversalClient {
	return c.client
}

// GetContext 获取上下文
func (c *Client) GetContext() context.Context {
	return c.ctx
}

// ==================== 包级别便捷函数 ====================

// getClient 获取客户端和上下文（内部使用）
func getClient() (redis.UniversalClient, context.Context, error) {
	client := GetClient()
	if client == nil {
		return nil, nil, fmt.Errorf("redis client not initialized")
	}
	return client.client, client.ctx, nil
}

// ==================== 字符串操作 ====================

// Get 获取键值
func Get(key string) (string, error) {
	client, ctx, err := getClient()
	if err != nil {
		return "", err
	}
	return client.Get(ctx, key).Result()
}

// Set 设置键值
func Set(key string, value interface{}, expiration time.Duration) error {
	client, ctx, err := getClient()
	if err != nil {
		return err
	}
	return client.Set(ctx, key, value, expiration).Err()
}

// SetNX 设置键值（仅当键不存在时）
func SetNX(key string, value interface{}, expiration time.Duration) (bool, error) {
	client, ctx, err := getClient()
	if err != nil {
		return false, err
	}
	return client.SetNX(ctx, key, value, expiration).Result()
}

// ==================== 键操作 ====================

// Del 删除键
func Del(keys ...string) error {
	client, ctx, err := getClient()
	if err != nil {
		return err
	}
	return client.Del(ctx, keys...).Err()
}

// Exists 检查键是否存在
func Exists(keys ...string) (int64, error) {
	client, ctx, err := getClient()
	if err != nil {
		return 0, err
	}
	return client.Exists(ctx, keys...).Result()
}

// Expire 设置键过期时间
func Expire(key string, expiration time.Duration) error {
	client, ctx, err := getClient()
	if err != nil {
		return err
	}
	return client.Expire(ctx, key, expiration).Err()
}

// TTL 获取键的剩余过期时间
func TTL(key string) (time.Duration, error) {
	client, ctx, err := getClient()
	if err != nil {
		return 0, err
	}
	return client.TTL(ctx, key).Result()
}

// Incr 自增
func Incr(key string) (int64, error) {
	client, ctx, err := getClient()
	if err != nil {
		return 0, err
	}
	return client.Incr(ctx, key).Result()
}

// Decr 自减
func Decr(key string) (int64, error) {
	client, ctx, err := getClient()
	if err != nil {
		return 0, err
	}
	return client.Decr(ctx, key).Result()
}

// ==================== Hash操作 ====================

// HSet 设置Hash字段
func HSet(key string, values ...interface{}) error {
	client, ctx, err := getClient()
	if err != nil {
		return err
	}
	return client.HSet(ctx, key, values...).Err()
}

// HGet 获取Hash字段值
func HGet(key, field string) (string, error) {
	client, ctx, err := getClient()
	if err != nil {
		return "", err
	}
	return client.HGet(ctx, key, field).Result()
}

// HGetAll 获取Hash所有字段
func HGetAll(key string) (map[string]string, error) {
	client, ctx, err := getClient()
	if err != nil {
		return nil, err
	}
	return client.HGetAll(ctx, key).Result()
}

// HDel 删除Hash字段
func HDel(key string, fields ...string) error {
	client, ctx, err := getClient()
	if err != nil {
		return err
	}
	return client.HDel(ctx, key, fields...).Err()
}

// HExists 检查Hash字段是否存在
func HExists(key, field string) (bool, error) {
	client, ctx, err := getClient()
	if err != nil {
		return false, err
	}
	return client.HExists(ctx, key, field).Result()
}
