package interceptor

import (
	"context"
	"errors"
	"fmt"
	logger "golang-tars/pkg"
	"golang-tars/pkg/jwt"

	"connectrpc.com/connect"
	jwtLib "github.com/golang-jwt/jwt/v5"
)

// 注意：context key 类型定义在 logger.ContextKey 中，统一管理

// authWhitelist 认证白名单（不需要认证的方法）
var connectAuthWhitelist = map[string]bool{
	// 可以添加不需要认证的方法
	"/demo.v1.DemoService/CreateDemo": true, // 测试用，允许不认证
}

// AuthInterceptor JWT 认证拦截器
func AuthInterceptor() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			procedure := req.Spec().Procedure

			// 检查是否在白名单中
			if connectAuthWhitelist[procedure] {
				return next(ctx, req)
			}

			// 从 header 中获取 Authorization
			authHeader := req.Header().Get("Authorization")
			if authHeader == "" {
				logger.Logger.Warn("missing authorization header")
				return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("missing authorization token"))
			}

			// 解析 Bearer token 格式: "Bearer <token>"
			var tokenString string
			if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
				tokenString = authHeader[7:]
			} else {
				logger.Logger.Warn("invalid authorization header format")
				return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("invalid authorization header format"))
			}

			// 验证并解析 token
			claims, err := jwt.ParseToken(tokenString)
			if err != nil {
				logger.Logger.Warn("token validation failed: ", err)

				// 区分 token 过期和 token 无效
				if errors.Is(err, jwtLib.ErrTokenExpired) {
					return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("token expired"))
				}
				return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("token invalid"))
			}

			// 将用户信息注入到 context 中（使用类型化 key 避免冲突）
			ctx = context.WithValue(ctx, logger.CtxKeyUserID, claims.UserID)
			ctx = context.WithValue(ctx, logger.CtxKeyUsername, claims.Username)
			ctx = context.WithValue(ctx, logger.CtxKeyPhoneNumber, claims.PhoneNumber)
			ctx = context.WithValue(ctx, logger.CtxKeyDeviceFingerprint, claims.DeviceFingerprint)
			ctx = context.WithValue(ctx, logger.CtxKeyToken, tokenString)

			logger.Logger.Info("authenticated user: ", claims.UserID, " (", claims.Username, ")")

			// 继续处理请求
			return next(ctx, req)
		}
	}
}
