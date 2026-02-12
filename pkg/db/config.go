package db

import (
	"time"
)

// DBConfig 数据库配置
type DBConfig struct {
	// 基础配置
	Host     string `mapstructure:"host" yaml:"host"`
	Port     int    `mapstructure:"port" yaml:"port"`
	Username string `mapstructure:"username" yaml:"username"`
	Password string `mapstructure:"password" yaml:"password"`
	Database string `mapstructure:"database" yaml:"database"`
	Charset  string `mapstructure:"charset" yaml:"charset"`

	// 连接池配置
	MaxIdleConns    int           `mapstructure:"max_idle_conns" yaml:"max_idle_conns"`       // 最大空闲连接数
	MaxOpenConns    int           `mapstructure:"max_open_conns" yaml:"max_open_conns"`       // 最大打开连接数
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime" yaml:"conn_max_lifetime"` // 连接最大生命周期
	ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time" yaml:"conn_max_idle_time"` // 连接最大空闲时间

	// GORM 配置
	LogLevel                  int  `mapstructure:"log_level" yaml:"log_level"`                                       // 日志级别 1:Silent 2:Error 3:Warn 4:Info
	SlowThreshold             int  `mapstructure:"slow_threshold" yaml:"slow_threshold"`                             // 慢查询阈值(毫秒)
	IgnoreRecordNotFoundError bool `mapstructure:"ignore_record_not_found_error" yaml:"ignore_record_not_found_error"` // 忽略记录未找到错误
	PrepareStmt               bool `mapstructure:"prepare_stmt" yaml:"prepare_stmt"`                                 // 预编译语句
	DisableAutomaticPing      bool `mapstructure:"disable_automatic_ping" yaml:"disable_automatic_ping"`             // 禁用自动 ping

	// 主从配置（可选）
	Slaves []SlaveConfig `mapstructure:"slaves" yaml:"slaves"` // 从库配置
}

// SlaveConfig 从库配置
type SlaveConfig struct {
	Host     string `mapstructure:"host" yaml:"host"`
	Port     int    `mapstructure:"port" yaml:"port"`
	Username string `mapstructure:"username" yaml:"username"`
	Password string `mapstructure:"password" yaml:"password"`
	Database string `mapstructure:"database" yaml:"database"`
	Weight   int    `mapstructure:"weight" yaml:"weight"` // 权重（用于负载均衡）
}

// DefaultDBConfig 返回默认配置
func DefaultDBConfig() *DBConfig {
	return &DBConfig{
		Charset:                   "utf8mb4",
		MaxIdleConns:              10,
		MaxOpenConns:              100,
		ConnMaxLifetime:           time.Hour,
		ConnMaxIdleTime:           time.Minute * 10,
		LogLevel:                  4, // Info
		SlowThreshold:             200,
		IgnoreRecordNotFoundError: false,
		PrepareStmt:               true,
		DisableAutomaticPing:      false,
	}
}

// Merge 合并配置（用户配置覆盖默认配置）
func (c *DBConfig) Merge(userConfig *DBConfig) *DBConfig {
	if userConfig.Host != "" {
		c.Host = userConfig.Host
	}
	if userConfig.Port != 0 {
		c.Port = userConfig.Port
	}
	if userConfig.Username != "" {
		c.Username = userConfig.Username
	}
	if userConfig.Password != "" {
		c.Password = userConfig.Password
	}
	if userConfig.Database != "" {
		c.Database = userConfig.Database
	}
	if userConfig.Charset != "" {
		c.Charset = userConfig.Charset
	}
	if userConfig.MaxIdleConns != 0 {
		c.MaxIdleConns = userConfig.MaxIdleConns
	}
	if userConfig.MaxOpenConns != 0 {
		c.MaxOpenConns = userConfig.MaxOpenConns
	}
	if userConfig.ConnMaxLifetime != 0 {
		c.ConnMaxLifetime = userConfig.ConnMaxLifetime
	}
	if userConfig.ConnMaxIdleTime != 0 {
		c.ConnMaxIdleTime = userConfig.ConnMaxIdleTime
	}
	if userConfig.LogLevel != 0 {
		c.LogLevel = userConfig.LogLevel
	}
	if userConfig.SlowThreshold != 0 {
		c.SlowThreshold = userConfig.SlowThreshold
	}
	c.IgnoreRecordNotFoundError = userConfig.IgnoreRecordNotFoundError
	c.PrepareStmt = userConfig.PrepareStmt
	c.DisableAutomaticPing = userConfig.DisableAutomaticPing

	if len(userConfig.Slaves) > 0 {
		c.Slaves = userConfig.Slaves
	}

	return c
}
