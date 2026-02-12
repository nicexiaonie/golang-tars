package interceptor

import (
	"context"
	"fmt"
	logger "go-base/pkg"
	"runtime/debug"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RecoveryUnaryInterceptor 错误恢复拦截器（一元调用）
func RecoveryUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				// 使用 context 中的 logger，自动包含 traceid
				log := logger.FromContext(ctx)
				log.Error(fmt.Sprintf("panic recovered: %v\n%s", r, debug.Stack()))

				// 返回 gRPC 错误
				err = status.Errorf(codes.Internal, "internal server error")
			}
		}()

		return handler(ctx, req)
	}
}

// RecoveryStreamInterceptor 错误恢复拦截器（流式调用）
func RecoveryStreamInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		defer func() {
			if r := recover(); r != nil {
				// 使用 context 中的 logger
				log := logger.FromContext(ss.Context())
				log.Error(fmt.Sprintf("panic recovered: %v\n%s", r, debug.Stack()))

				// 返回 gRPC 错误
				err = status.Errorf(codes.Internal, "internal server error")
			}
		}()

		return handler(srv, ss)
	}
}
