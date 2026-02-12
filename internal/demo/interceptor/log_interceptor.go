package interceptor

import (
	"context"
	"fmt"
	logger "go-base/pkg"
	"time"

	"github.com/nicexiaonie/ghelper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// LogUnaryInterceptor 日志拦截器（一元调用）
func LogUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()

		// 从 metadata 中获取或生成 request_id 和 trace_id
		requestID, traceID := extractOrGenerateIDs(ctx)

		// 将追踪信息注入到 context 中
		ctx = context.WithValue(ctx, "RequestId", requestID)
		ctx = context.WithValue(ctx, "TraceId", traceID)

		// 创建带有 traceid 的 logger 并注入 context
		log := logger.WithFields(map[string]interface{}{
			"RequestId": requestID,
			"TraceId":   traceID,
		})
		ctx = context.WithValue(ctx, "Logger", log)

		// 记录请求开始
		log.Info(fmt.Sprintf("[gRPC] %s started", info.FullMethod))

		// 执行处理器
		resp, err := handler(ctx, req)

		// 记录请求结束
		duration := time.Since(start)
		if err != nil {
			log.Error(fmt.Sprintf("[gRPC] %s failed: %v (duration: %v)", info.FullMethod, err, duration))
		} else {
			log.Info(fmt.Sprintf("[gRPC] %s completed (duration: %v)", info.FullMethod, duration))
		}

		return resp, err
	}
}

// LogStreamInterceptor 日志拦截器（流式调用）
func LogStreamInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		start := time.Now()

		// 从 metadata 中获取或生成 request_id 和 trace_id
		ctx := ss.Context()
		requestID, traceID := extractOrGenerateIDs(ctx)

		// 将追踪信息注入到 context 中
		ctx = context.WithValue(ctx, "RequestId", requestID)
		ctx = context.WithValue(ctx, "TraceId", traceID)

		// 创建带有 traceid 的 logger 并注入 context
		log := logger.WithFields(map[string]interface{}{
			"RequestId": requestID,
			"TraceId":   traceID,
		})
		ctx = context.WithValue(ctx, "Logger", log)

		// 记录请求开始
		log.Info(fmt.Sprintf("[gRPC Stream] %s started", info.FullMethod))

		// 创建新的 ServerStream 包装器
		wrappedStream := &wrappedServerStream{
			ServerStream: ss,
			ctx:          ctx,
		}

		// 执行处理器
		err := handler(srv, wrappedStream)

		// 记录请求结束
		duration := time.Since(start)
		if err != nil {
			log.Error(fmt.Sprintf("[gRPC Stream] %s failed: %v (duration: %v)", info.FullMethod, err, duration))
		} else {
			log.Info(fmt.Sprintf("[gRPC Stream] %s completed (duration: %v)", info.FullMethod, duration))
		}

		return err
	}
}

// extractOrGenerateIDs 从 metadata 中提取或生成 request_id 和 trace_id
func extractOrGenerateIDs(ctx context.Context) (requestID, traceID string) {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		// 尝试从 metadata 中获取
		if ids := md.Get("x-request-id"); len(ids) > 0 {
			requestID = ids[0]
		}
		if ids := md.Get("x-trace-id"); len(ids) > 0 {
			traceID = ids[0]
		}
	}

	// 如果没有，则生成新的
	if requestID == "" {
		requestID = ghelper.UniqueId()
	}
	if traceID == "" {
		traceID = ghelper.UniqueId()
	}

	return requestID, traceID
}
