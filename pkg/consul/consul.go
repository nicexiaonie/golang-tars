package consul

import (
	"fmt"
	"time"
	"github.com/hashicorp/consul/api"
)

// Client Consul客户端封装
type Client struct {
	client *api.Client
	config *Config
}

// Config Consul配置
type Config struct {
	Address        string        // Consul地址，例如: "127.0.0.1:8500"
	Scheme         string        // http 或 https
	Datacenter     string        // 数据中心
	Token          string        // ACL Token
	Timeout        time.Duration // 超时时间
	ServiceName    string        // 服务名称
	ServiceID      string        // 服务ID
	ServiceAddress string        // 服务地址
	ServicePort    int           // 服务端口
	Tags           []string      // 服务标签
	HealthCheck    *HealthCheck  // 健康检查配置
}

// HealthCheck 健康检查配置
type HealthCheck struct {
	CheckID                        string        // 健康检查ID
	Type                           string        // 检查类型: http, tcp, grpc
	Interval                       time.Duration // 检查间隔
	Timeout                        time.Duration // 检查超时
	DeregisterCriticalServiceAfter time.Duration // 注销临界服务的时间
	HTTP                           string        // HTTP检查地址
	TCP                            string        // TCP检查地址
	GRPC                           string        // GRPC检查地址
}

// NewClient 创建Consul客户端
func NewClient(cfg *Config) (*Client, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is nil")
	}

	// 设置默认值
	if cfg.Scheme == "" {
		cfg.Scheme = "http"
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = 10 * time.Second
	}

	// 创建Consul配置
	config := api.DefaultConfig()
	config.Address = cfg.Address
	config.Scheme = cfg.Scheme
	config.Datacenter = cfg.Datacenter
	config.Token = cfg.Token

	// 创建Consul客户端
	client, err := api.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("create consul client failed: %w", err)
	}

	return &Client{
		client: client,
		config: cfg,
	}, nil
}

// RegisterService 注册服务
func (c *Client) RegisterService() error {
	registration := &api.AgentServiceRegistration{
		ID:      c.config.ServiceID,
		Name:    c.config.ServiceName,
		Address: c.config.ServiceAddress,
		Port:    c.config.ServicePort,
		Tags:    c.config.Tags,
	}

	// 配置健康检查
	if c.config.HealthCheck != nil {
		check := &api.AgentServiceCheck{
			CheckID:                        c.config.HealthCheck.CheckID,
			Interval:                       c.config.HealthCheck.Interval.String(),
			Timeout:                        c.config.HealthCheck.Timeout.String(),
			DeregisterCriticalServiceAfter: c.config.HealthCheck.DeregisterCriticalServiceAfter.String(),
		}

		switch c.config.HealthCheck.Type {
		case "http":
			check.HTTP = c.config.HealthCheck.HTTP
		case "tcp":
			check.TCP = c.config.HealthCheck.TCP
		case "grpc":
			check.GRPC = c.config.HealthCheck.GRPC
		}

		registration.Check = check
	}

	// 注册服务
	if err := c.client.Agent().ServiceRegister(registration); err != nil {
		return fmt.Errorf("register service failed: %w", err)
	}

	return nil
}

// DeregisterService 注销服务
func (c *Client) DeregisterService() error {
	if err := c.client.Agent().ServiceDeregister(c.config.ServiceID); err != nil {
		return fmt.Errorf("deregister service failed: %w", err)
	}
	return nil
}

// DiscoverService 服务发现
func (c *Client) DiscoverService(serviceName string, healthy bool) ([]*api.ServiceEntry, error) {
	services, _, err := c.client.Health().Service(serviceName, "", healthy, nil)
	if err != nil {
		return nil, fmt.Errorf("discover service failed: %w", err)
	}
	return services, nil
}

// GetServiceAddress 获取服务地址（随机选择一个健康的服务实例）
func (c *Client) GetServiceAddress(serviceName string) (string, error) {
	services, err := c.DiscoverService(serviceName, true)
	if err != nil {
		return "", err
	}

	if len(services) == 0 {
		return "", fmt.Errorf("no healthy service found: %s", serviceName)
	}

	// 简单轮询，返回第一个
	service := services[0]
	return fmt.Sprintf("%s:%d", service.Service.Address, service.Service.Port), nil
}

// WatchService 监听服务变化
func (c *Client) WatchService(serviceName string, stopCh <-chan struct{}) (<-chan []*api.ServiceEntry, error) {
	resultCh := make(chan []*api.ServiceEntry)

	go func() {
		defer close(resultCh)

		var lastIndex uint64
		for {
			select {
			case <-stopCh:
				return
			default:
				queryOpts := &api.QueryOptions{
					WaitIndex: lastIndex,
					WaitTime:  30 * time.Second,
				}

				services, meta, err := c.client.Health().Service(serviceName, "", true, queryOpts)
				if err != nil {
					time.Sleep(1 * time.Second)
					continue
				}

				if lastIndex != meta.LastIndex {
					lastIndex = meta.LastIndex
					resultCh <- services
				}
			}
		}
	}()

	return resultCh, nil
}

// GetKV 获取KV存储的值
func (c *Client) GetKV(key string) ([]byte, error) {
	pair, _, err := c.client.KV().Get(key, nil)
	if err != nil {
		return nil, fmt.Errorf("get kv failed: %w", err)
	}
	if pair == nil {
		return nil, fmt.Errorf("key not found: %s", key)
	}
	return pair.Value, nil
}

// PutKV 设置KV存储的值
func (c *Client) PutKV(key string, value []byte) error {
	p := &api.KVPair{Key: key, Value: value}
	_, err := c.client.KV().Put(p, nil)
	if err != nil {
		return fmt.Errorf("put kv failed: %w", err)
	}
	return nil
}

// DeleteKV 删除KV存储的值
func (c *Client) DeleteKV(key string) error {
	_, err := c.client.KV().Delete(key, nil)
	if err != nil {
		return fmt.Errorf("delete kv failed: %w", err)
	}
	return nil
}

// WatchKV 监听KV变化
func (c *Client) WatchKV(key string, stopCh <-chan struct{}) (<-chan []byte, error) {
	resultCh := make(chan []byte)

	go func() {
		defer close(resultCh)

		var lastIndex uint64
		for {
			select {
			case <-stopCh:
				return
			default:
				queryOpts := &api.QueryOptions{
					WaitIndex: lastIndex,
					WaitTime:  30 * time.Second,
				}

				pair, meta, err := c.client.KV().Get(key, queryOpts)
				if err != nil {
					time.Sleep(1 * time.Second)
					continue
				}

				if lastIndex != meta.LastIndex {
					lastIndex = meta.LastIndex
					if pair != nil {
						resultCh <- pair.Value
					}
				}
			}
		}
	}()

	return resultCh, nil
}

// Close 关闭客户端并注销服务
func (c *Client) Close() error {
	return c.DeregisterService()
}
