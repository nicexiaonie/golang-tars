package interceptor

import (
	"context"
	"fmt"
	logger "golang-tars/pkg"
	"runtime/debug"

	"connectrpc.com/connect"
)

// RecoveryInterceptor 错误恢复拦截器
func RecoveryInterceptor() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (resp connect.AnyResponse, err error) {
			defer func() {
				if r := recover(); r != nil {
					log := logger.FromContext(ctx)
					log.Error(fmt.Sprintf("panic recovered: %v\n%s", r, debug.Stack()))
					err = connect.NewError(connect.CodeInternal, fmt.Errorf("internal server error"))
				}
			}()
			return next(ctx, req)
		}
	}
}
