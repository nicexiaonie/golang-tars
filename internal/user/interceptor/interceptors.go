package interceptor

import "connectrpc.com/connect"

// NewInterceptors 创建所有拦截器链
func NewInterceptors() connect.Option {
	interceptors := connect.WithInterceptors(
		RecoveryInterceptor(),           // 错误恢复（最外层）
		LogInterceptor(),                // 日志记录
		RequestResponseLogInterceptor(), // 请求响应参数日志
		AuthInterceptor(),               // JWT 认证
	)
	return interceptors
}
