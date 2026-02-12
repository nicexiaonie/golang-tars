package proxy

import (
	"context"
	"fmt"
	"time"

	logger "golang-tars/pkg"
	demov1 "golang-tars/pkg/proto/demo/v1"

	"github.com/nicexiaonie/ghelper"
	"github.com/nicexiaonie/ghttp"
)

const (
	// DefaultTimeout 默认超时时间
	DefaultTimeout = 10 * time.Second
	// MaxRetryTimes 最大重试次数
	MaxRetryTimes = 3
)

// DemoProxy 第三方接口代理
type DemoProxy interface {
	// GetUserInfo 获取用户信息
	GetUserInfo(ctx context.Context, req *demov1.CreateDemoRequest) (*demov1.CreateDemoResponse, error)
}

// demoProxyImpl 第三方接口代理实现
type demoProxyImpl struct {
	client  *ghttp.Client
	baseURL string
	timeout time.Duration
}

// DemoProxyConfig 代理配置
type DemoProxyConfig struct {
	BaseURL string        // 基础URL
	Timeout time.Duration // 超时时间
	APIKey  string        // API密钥
}

// NewDemoProxy 创建第三方接口代理实例
func NewDemoProxy(config *DemoProxyConfig) DemoProxy {
	// 设置默认值
	if config.Timeout == 0 {
		config.Timeout = DefaultTimeout
	}

	// 创建 HTTP 客户端
	client := ghttp.NewClient(
		ghttp.WithBaseURL(config.BaseURL),
		ghttp.WithTimeout(config.Timeout),
		ghttp.WithHeader("Content-Type", "application/json"),
		ghttp.WithHeader("X-API-Key", config.APIKey),
	)

	return &demoProxyImpl{
		client:  client,
		baseURL: config.BaseURL,
		timeout: config.Timeout,
	}
}

// ========== 接口实现 ==========

// GetUserInfo 获取用户信息
func (p *demoProxyImpl) GetUserInfo(ctx context.Context, req *demov1.CreateDemoRequest) (*demov1.CreateDemoResponse, error) {
	// 需要使用ctx获取日志实体, 保证trace的连贯
	log := logger.FromContext(ctx)
	// 记录请求日志
	log.Info(fmt.Sprintf("[DemoProxy] GetUserInfo request: %s", ghelper.ToString(req)))
	// 发起HTTP GET请求
	resp, err := p.client.Get("/api/v1/user/info").
		Context(ctx).
		Query("name", req.Name).
		Do()
	if err != nil {
		log.Error(fmt.Sprintf("[DemoProxy] GetUserInfo request failed: %v", err))
		return nil, fmt.Errorf("获取用户信息失败: %w", err)
	}
	// 检查HTTP状态码
	if !resp.IsSuccess() {
		log.Error(fmt.Sprintf("[DemoProxy] GetUserInfo HTTP error: status=%d, body=%s",
			resp.StatusCode, resp.String()))
		return nil, fmt.Errorf("HTTP错误: status=%d", resp.StatusCode)
	}
	// 解析响应
	var result demov1.CreateDemoResponse
	if err := resp.JSON(&result); err != nil {
		log.Error(fmt.Sprintf("[DemoProxy] GetUserInfo decode response failed: %v", err))
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}
	return &result, nil
}
