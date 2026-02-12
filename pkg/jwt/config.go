package jwt

import "time"

// Config JWT配置
type Config struct {
	// 签名密钥
	SecretKey string `mapstructure:"secret_key" yaml:"secret_key"`

	// Token过期时间
	ExpireDuration time.Duration `mapstructure:"expire_duration" yaml:"expire_duration"`

	// 签发者
	Issuer string `mapstructure:"issuer" yaml:"issuer"`
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		SecretKey:      "your-secret-key-change-in-production",
		ExpireDuration: 24 * time.Hour, // 默认24小时
		Issuer:         "cue-words",
	}
}
