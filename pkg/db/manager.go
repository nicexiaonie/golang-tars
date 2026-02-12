package db

import (
	"fmt"
	"sync"

	"github.com/nicexiaonie/grds"
)

var (
	// 全局数据库管理器
	globalManager *Manager
	once          sync.Once
)

// Manager 数据库管理器 - 只负责管理多个 grds.Client 实例
// 不重复封装 grds 已有的功能，直接使用 grds.Client 的方法
type Manager struct {
	clients map[string]*grds.Client // grds 客户端集合
	mu      sync.RWMutex            // 读写锁
}

// GetManager 获取全局数据库管理器（单例）
func GetManager() *Manager {
	once.Do(func() {
		globalManager = &Manager{
			clients: make(map[string]*grds.Client),
		}
	})
	return globalManager
}

// InitDB 初始化数据库连接
// name: 数据库实例名称（如 "default", "analytics"）
// config: 数据库配置
func (m *Manager) InitDB(name string, config *DBConfig) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 检查是否已存在
	if _, exists := m.clients[name]; exists {
		return fmt.Errorf("database instance '%s' already exists", name)
	}

	// 合并默认配置
	finalConfig := DefaultDBConfig().Merge(config)

	// 创建 grds 配置
	grdsConfig := grds.NewConfig(
		finalConfig.Host,
		finalConfig.Port,
		finalConfig.Username,
		finalConfig.Password,
		finalConfig.Database,
	)

	// 设置字符集
	if finalConfig.Charset != "" {
		grdsConfig.Charset = finalConfig.Charset
	}

	// 设置连接池参数
	grdsConfig.MaxIdleConns = finalConfig.MaxIdleConns
	grdsConfig.MaxOpenConns = finalConfig.MaxOpenConns
	grdsConfig.ConnMaxLifetime = finalConfig.ConnMaxLifetime
	grdsConfig.ConnMaxIdleTime = finalConfig.ConnMaxIdleTime

	// 创建 grds 客户端
	client, err := grds.NewClient(grdsConfig)
	if err != nil {
		return fmt.Errorf("failed to create database client '%s': %w", name, err)
	}

	// 保存客户端
	m.clients[name] = client

	return nil
}

// GetClient 获取指定名称的 grds 客户端
// 使用方式：client.Create(), client.Find(), client.Transaction() 等
func (m *Manager) GetClient(name string) (*grds.Client, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	client, exists := m.clients[name]
	if !exists {
		return nil, fmt.Errorf("database client '%s' not found", name)
	}

	return client, nil
}

// GetDefaultClient 获取默认 grds 客户端（名称为 "default"）
func (m *Manager) GetDefaultClient() (*grds.Client, error) {
	return m.GetClient("default")
}

// Close 关闭指定数据库连接
func (m *Manager) Close(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	client, exists := m.clients[name]
	if !exists {
		return fmt.Errorf("database client '%s' not found", name)
	}

	// 使用 grds.Client 的 Close 方法
	if err := client.Close(); err != nil {
		return fmt.Errorf("failed to close database '%s': %w", name, err)
	}

	delete(m.clients, name)
	return nil
}

// CloseAll 关闭所有数据库连接
func (m *Manager) CloseAll() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var errs []error
	for name, client := range m.clients {
		if err := client.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close database '%s': %w", name, err))
		}
	}

	m.clients = make(map[string]*grds.Client)

	if len(errs) > 0 {
		return fmt.Errorf("errors occurred while closing databases: %v", errs)
	}

	return nil
}

// ListInstances 列出所有数据库实例名称
func (m *Manager) ListInstances() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	names := make([]string, 0, len(m.clients))
	for name := range m.clients {
		names = append(names, name)
	}
	return names
}
