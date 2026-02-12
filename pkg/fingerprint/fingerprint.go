package fingerprint

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

// DeviceInfo 设备信息
type DeviceInfo struct {
	UserAgent  string // User-Agent
	IP         string // IP地址
	DeviceInfo string // 设备信息
	DeviceId   string // 设备ID
	Platform   string // 平台信息
}

// Generate 生成设备指纹
// 使用设备信息的组合生成唯一的设备指纹
func Generate(info *DeviceInfo) string {
	if info == nil {
		return ""
	}

	// 组合设备信息
	data := fmt.Sprintf("%s|%s|%s|%s|%s",
		strings.TrimSpace(info.UserAgent),
		strings.TrimSpace(info.IP),
		strings.TrimSpace(info.DeviceInfo),
		strings.TrimSpace(info.DeviceId),
		strings.TrimSpace(info.Platform),
	)

	// 使用SHA256生成指纹
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// GenerateSimple 简化版设备指纹生成（只使用UserAgent和IP）
func GenerateSimple(userAgent, ip string) string {
	return Generate(&DeviceInfo{
		UserAgent: userAgent,
		IP:        ip,
	})
}
