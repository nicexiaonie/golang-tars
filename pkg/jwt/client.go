package jwt

import (
	"sync"
)

var (
	// globalManager 全局JWT管理器
	globalManager *Manager
	once          sync.Once
)

// InitGlobalManager 初始化全局JWT管理器
func InitGlobalManager(config *Config) {
	once.Do(func() {
		globalManager = NewManager(config)
	})
}

// GetGlobalManager 获取全局JWT管理器
func GetGlobalManager() *Manager {
	if globalManager == nil {
		InitGlobalManager(nil)
	}
	return globalManager
}

// GenerateToken 使用全局管理器生成token（便捷方法）
func GenerateToken(userID uint64, username, phoneNumber, deviceFingerprint string) (string, error) {
	return GetGlobalManager().GenerateToken(userID, username, phoneNumber, deviceFingerprint)
}

// ParseToken 使用全局管理器解析token（便捷方法）
func ParseToken(tokenString string) (*Claims, error) {
	return GetGlobalManager().ParseToken(tokenString)
}

// ValidateToken 使用全局管理器验证token（便捷方法）
func ValidateToken(tokenString string) bool {
	return GetGlobalManager().ValidateToken(tokenString)
}

