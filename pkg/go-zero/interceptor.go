package gozero

import (
	"context"
	"fmt"
	"math/rand/v2"
	"time"

	logger "golang-tars/pkg"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// LogInterceptor 日志拦截器 - 记录 RPC 调用信息
// 使用项目 logger 包，自动从 context 获取 TraceId/RequestId
func LogInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		start := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)
		duration := time.Since(start)

		log := logger.FromContext(ctx)
		if err != nil {
			st, _ := status.FromError(err)
			log.WithFields(logrus.Fields{
				"method":   method,
				"duration": duration.String(),
				"code":     st.Code().String(),
			}).Error(fmt.Sprintf("[RPC Client] %s failed: %v", method, err))
		} else {
			log.WithFields(logrus.Fields{
				"method":   method,
				"duration": duration.String(),
			}).Info(fmt.Sprintf("[RPC Client] %s OK", method))
		}

		return err
	}
}

// RetryInterceptor 重试拦截器 - 在失败时自动重试（指数退避 + 随机抖动）
func RetryInterceptor(maxRetries int) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		var lastErr error
		for i := 0; i <= maxRetries; i++ {
			lastErr = invoker(ctx, method, req, reply, cc, opts...)
			if lastErr == nil {
				return nil
			}

			// 检查上下文是否已取消
			if ctx.Err() != nil {
				return lastErr
			}

			// 检查是否可重试的错误
			if !isRetryable(lastErr) {
				return lastErr
			}

			// 指数退避 + 随机抖动
			if i < maxRetries {
				base := time.Duration(1<<uint(i)) * 100 * time.Millisecond // 100ms, 200ms, 400ms, 800ms...
				jitter := time.Duration(rand.Int64N(int64(base)/2 + 1))    // 0 ~ base/2
				backoff := base + jitter

				// 限制最大退避时间为 5 秒
				if backoff > 5*time.Second {
					backoff = 5 * time.Second
				}

				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(backoff):
				}
			}
		}
		return lastErr
	}
}

// MetadataInterceptor 元数据拦截器 - 自动传递指定的 metadata
func MetadataInterceptor(kvPairs map[string]string) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		md := metadata.New(kvPairs)

		// 合并已有的 metadata
		if existingMD, ok := metadata.FromOutgoingContext(ctx); ok {
			md = metadata.Join(existingMD, md)
		}

		ctx = metadata.NewOutgoingContext(ctx, md)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// isRetryable 判断错误是否可重试
func isRetryable(err error) bool {
	st, ok := status.FromError(err)
	if !ok {
		return false
	}

	switch st.Code() {
	case codes.Unavailable:
		return true
	case codes.Internal:
		return true
	case codes.DeadlineExceeded:
		return false // 超时不重试
	default:
		return false
	}
}
