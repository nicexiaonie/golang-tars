package interceptor

import (
	"context"
	logger "golang-tars/pkg"

	"connectrpc.com/connect"
	"github.com/nicexiaonie/ghelper"
)

// LogInterceptor 日志拦截器
func LogInterceptor() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {

			// 从 header 中获取或生成 request_id 和 trace_id
			requestID := req.Header().Get("X-Request-ID")
			traceID := req.Header().Get("X-Trace-ID")

			if requestID == "" {
				requestID = ghelper.UniqueId()
			}
			if traceID == "" {
				traceID = ghelper.UniqueId()
			}

			// 将追踪信息注入到 context 中（使用类型化 key 避免冲突）
			ctx = context.WithValue(ctx, logger.CtxKeyRequestID, requestID)
			ctx = context.WithValue(ctx, logger.CtxKeyTraceID, traceID)

			// 创建带有 traceid 的 logger 并注入 context
			log := logger.WithFields(map[string]any{
				"RequestID": requestID,
				"TraceID":   traceID,
			})
			ctx = context.WithValue(ctx, logger.CtxKeyLogger, log)

			// 执行处理器
			resp, err := next(ctx, req)

			return resp, err
		}
	}
}
