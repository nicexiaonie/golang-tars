package redis

import (
	"time"
)

// Config Redis配置
type Config struct {
	// 基础配置
	Host     string `mapstructure:"host" yaml:"host"`
	Port     int    `mapstructure:"port" yaml:"port"`
	Password string `mapstructure:"password" yaml:"password"`
	DB       int    `mapstructure:"db" yaml:"db"` // 数据库索引

	// 连接池配置
	PoolSize     int           `mapstructure:"pool_size" yaml:"pool_size"`           // 连接池大小
	MinIdleConns int           `mapstructure:"min_idle_conns" yaml:"min_idle_conns"` // 最小空闲连接数
	MaxIdleConns int           `mapstructure:"max_idle_conns" yaml:"max_idle_conns"` // 最大空闲连接数
	MaxRetries   int           `mapstructure:"max_retries" yaml:"max_retries"`       // 最大重试次数
	DialTimeout  time.Duration `mapstructure:"dial_timeout" yaml:"dial_timeout"`     // 连接超时
	ReadTimeout  time.Duration `mapstructure:"read_timeout" yaml:"read_timeout"`     // 读超时
	WriteTimeout time.Duration `mapstructure:"write_timeout" yaml:"write_timeout"`   // 写超时
	PoolTimeout  time.Duration `mapstructure:"pool_timeout" yaml:"pool_timeout"`     // 连接池超时
	IdleTimeout  time.Duration `mapstructure:"idle_timeout" yaml:"idle_timeout"`     // 空闲连接超时

	// 集群配置（可选）
	ClusterMode bool     `mapstructure:"cluster_mode" yaml:"cluster_mode"` // 是否集群模式
	Addrs       []string `mapstructure:"addrs" yaml:"addrs"`               // 集群地址列表
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Host:         "127.0.0.1",
		Port:         6379,
		DB:           0,
		PoolSize:     10,
		MinIdleConns: 5,
		MaxIdleConns: 10,
		MaxRetries:   3,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolTimeout:  4 * time.Second,
		IdleTimeout:  5 * time.Minute,
		ClusterMode:  false,
	}
}

// Merge 合并配置（用户配置覆盖默认配置）
func (c *Config) Merge(userConfig *Config) *Config {
	if userConfig.Host != "" {
		c.Host = userConfig.Host
	}
	if userConfig.Port != 0 {
		c.Port = userConfig.Port
	}
	if userConfig.Password != "" {
		c.Password = userConfig.Password
	}
	if userConfig.DB != 0 {
		c.DB = userConfig.DB
	}
	if userConfig.PoolSize != 0 {
		c.PoolSize = userConfig.PoolSize
	}
	if userConfig.MinIdleConns != 0 {
		c.MinIdleConns = userConfig.MinIdleConns
	}
	if userConfig.MaxIdleConns != 0 {
		c.MaxIdleConns = userConfig.MaxIdleConns
	}
	if userConfig.MaxRetries != 0 {
		c.MaxRetries = userConfig.MaxRetries
	}
	if userConfig.DialTimeout != 0 {
		c.DialTimeout = userConfig.DialTimeout
	}
	if userConfig.ReadTimeout != 0 {
		c.ReadTimeout = userConfig.ReadTimeout
	}
	if userConfig.WriteTimeout != 0 {
		c.WriteTimeout = userConfig.WriteTimeout
	}
	if userConfig.PoolTimeout != 0 {
		c.PoolTimeout = userConfig.PoolTimeout
	}
	if userConfig.IdleTimeout != 0 {
		c.IdleTimeout = userConfig.IdleTimeout
	}
	if userConfig.ClusterMode {
		c.ClusterMode = userConfig.ClusterMode
	}
	if len(userConfig.Addrs) > 0 {
		c.Addrs = userConfig.Addrs
	}

	return c
}
