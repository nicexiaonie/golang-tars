package user

import (
	"golang-tars/internal/user/proxy"
	"golang-tars/pkg/db"
	"time"

	"github.com/nicexiaonie/gconf"
)

// ProvideDBManager 提供数据库管理器
// Wire 会自动调用这个函数来获取数据库管理器
// 数据库管理器可以管理多个数据库实例（default, analytics, read_replica 等）
func ProvideDBManager() *db.Manager {
	return db.GetManager()
}

// ProvideDemoProxyConfig 提供 DemoProxy 配置
// 从配置文件中读取第三方服务配置
func ProvideDemoProxyConfig() *proxy.DemoProxyConfig {
	// 从配置读取超时时间（秒），如果没有配置则使用默认值
	timeoutSeconds := gconf.GetInt("proxy.demo.timeout_seconds")
	timeout := time.Duration(timeoutSeconds) * time.Second
	if timeoutSeconds == 0 {
		timeout = 10 * time.Second // 默认 10 秒
	}

	return &proxy.DemoProxyConfig{
		BaseURL: gconf.GetString("proxy.demo.base_url"),
		Timeout: timeout,
		APIKey:  gconf.GetString("proxy.demo.api_key"),
	}
}
