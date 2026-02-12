package gozero

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

// connPool gRPC 连接池
// 持有多个 zrpc.Client，通过 round-robin 分发请求
// 注意：初始化后 clients/conns/size 不可变，无需加锁
type connPool struct {
	clients   []zrpc.Client              // 底层 zrpc 客户端列表
	conns     []grpc.ClientConnInterface // 对应的 gRPC 连接列表
	size      int                        // 连接池大小
	index     atomic.Uint64              // 原子计数器，用于 round-robin
	closeOnce sync.Once                  // 确保 close 只执行一次
}

// pooledConn 连接池代理
// 实现 grpc.ClientConnInterface 接口
// 每次 Invoke/NewStream 调用通过 round-robin 选择底层连接
// 对上层完全透明，即使缓存了 gRPC Service Client 也能自动分发
type pooledConn struct {
	pool *connPool
}

// newConnPool 创建连接池
// 初始化时创建 poolSize 个 zrpc.Client 实例
func newConnPool(poolSize int, builder func() (zrpc.Client, error)) (*connPool, error) {
	if poolSize < 1 {
		poolSize = 1
	}

	pool := &connPool{
		clients: make([]zrpc.Client, 0, poolSize),
		conns:   make([]grpc.ClientConnInterface, 0, poolSize),
		size:    poolSize,
	}

	for i := 0; i < poolSize; i++ {
		client, err := builder()
		if err != nil {
			// 创建失败时关闭已创建的连接
			pool.close()
			return nil, fmt.Errorf("failed to create pool connection [%d/%d]: %w", i+1, poolSize, err)
		}
		pool.clients = append(pool.clients, client)
		pool.conns = append(pool.conns, client.Conn())
	}

	return pool, nil
}

// next 通过 round-robin 获取下一个连接
func (p *connPool) next() grpc.ClientConnInterface {
	// Add(1) 返回递增后的值（从1开始），减1使索引从0开始，确保所有连接均匀分发
	idx := p.index.Add(1) - 1
	return p.conns[idx%uint64(p.size)]
}

// getPooledConn 获取连接池代理（实现 grpc.ClientConnInterface）
// 统一通过 pooledConn 代理，确保调用计数（PoolStats.Calls）在所有模式下准确
func (p *connPool) getPooledConn() grpc.ClientConnInterface {
	return &pooledConn{pool: p}
}

// getSize 获取连接池大小
func (p *connPool) getSize() int {
	return p.size
}

// getClients 获取所有 zrpc.Client（返回副本）
func (p *connPool) getClients() []zrpc.Client {
	result := make([]zrpc.Client, len(p.clients))
	copy(result, p.clients)
	return result
}

// close 关闭连接池中所有底层 gRPC 连接（幂等，可安全多次调用）
func (p *connPool) close() {
	p.closeOnce.Do(func() {
		for _, conn := range p.conns {
			if closer, ok := conn.(interface{ Close() error }); ok {
				_ = closer.Close()
			}
		}
		// 清空引用，防止关闭后误用
		p.clients = nil
		p.conns = nil
	})
}

// ==================== pooledConn 实现 grpc.ClientConnInterface ====================

// Invoke 实现 grpc.ClientConnInterface
// 通过 round-robin 选择一个底层连接执行 Unary RPC 调用
func (pc *pooledConn) Invoke(ctx context.Context, method string, args, reply interface{}, opts ...grpc.CallOption) error {
	conn := pc.pool.next()
	return conn.Invoke(ctx, method, args, reply, opts...)
}

// NewStream 实现 grpc.ClientConnInterface
// 通过 round-robin 选择一个底层连接创建 Stream RPC
func (pc *pooledConn) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	conn := pc.pool.next()
	return conn.NewStream(ctx, desc, method, opts...)
}

// ==================== 连接池统计 ====================

// PoolStats 连接池统计信息
type PoolStats struct {
	Name     string `json:"name"`      // 服务名称
	PoolSize int    `json:"pool_size"` // 连接池大小
	Calls    uint64 `json:"calls"`     // 累计调用次数（round-robin 轮转计数）
}

// stats 获取连接池统计信息
func (p *connPool) stats(name string) PoolStats {
	return PoolStats{
		Name:     name,
		PoolSize: p.size,
		Calls:    p.index.Load(),
	}
}
