package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims JWT声明
type Claims struct {
	UserID            uint64 `json:"user_id"`
	Username          string `json:"username"`
	PhoneNumber       string `json:"phone_number"`
	DeviceFingerprint string `json:"device_fingerprint"` // 设备指纹
	jwt.RegisteredClaims
}

// Manager JWT管理器
type Manager struct {
	config *Config
}

// NewManager 创建JWT管理器
func NewManager(config *Config) *Manager {
	if config == nil {
		config = DefaultConfig()
	}
	return &Manager{
		config: config,
	}
}

// GenerateToken 生成访问token
func (m *Manager) GenerateToken(userID uint64, username, phoneNumber, deviceFingerprint string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:            userID,
		Username:          username,
		PhoneNumber:       phoneNumber,
		DeviceFingerprint: deviceFingerprint,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.config.Issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.config.ExpireDuration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.config.SecretKey))
}

// ParseToken 解析token
func (m *Manager) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(m.config.SecretKey), nil
	})

	if err != nil {
		// 检查是否是token过期错误
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, jwt.ErrTokenExpired
		}
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// ValidateToken 验证token是否有效
func (m *Manager) ValidateToken(tokenString string) bool {
	_, err := m.ParseToken(tokenString)
	return err == nil
}


// IsTokenExpiringSoon 检查token是否即将过期（默认5分钟内）
func (m *Manager) IsTokenExpiringSoon(tokenString string, threshold time.Duration) (bool, error) {
	claims, err := m.ParseToken(tokenString)
	if err != nil {
		return false, err
	}

	if claims.ExpiresAt == nil {
		return false, errors.New("token has no expiration time")
	}

	// 计算剩余时间
	remaining := time.Until(claims.ExpiresAt.Time)
	return remaining <= threshold, nil
}

// GetTokenRemainingTime 获取token剩余有效时间
func (m *Manager) GetTokenRemainingTime(tokenString string) (time.Duration, error) {
	claims, err := m.ParseToken(tokenString)
	if err != nil {
		return 0, err
	}

	if claims.ExpiresAt == nil {
		return 0, errors.New("token has no expiration time")
	}

	return time.Until(claims.ExpiresAt.Time), nil
}
