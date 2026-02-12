package interceptor

import (
	"context"
	"errors"
	logger "go-base/pkg"
	"go-base/pkg/jwt"

	jwtLib "github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// authWhitelist 认证白名单（不需要认证的方法）
var authWhitelist = map[string]bool{
	// 可以添加不需要认证的 gRPC 方法
	// "/demo.v1.DemoService/SomePublicMethod": true,
}

// AuthUnaryInterceptor JWT 认证拦截器（一元调用）
func AuthUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// 检查是否在白名单中
		if authWhitelist[info.FullMethod] {
			return handler(ctx, req)
		}

		// 从 metadata 中获取 token
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			logger.Logger.Warn("missing metadata")
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		// 获取 Authorization header
		authHeaders := md.Get("authorization")
		if len(authHeaders) == 0 {
			logger.Logger.Warn("missing authorization header")
			return nil, status.Error(codes.Unauthenticated, "missing authorization token")
		}

		authHeader := authHeaders[0]

		// 解析 Bearer token 格式: "Bearer <token>"
		var tokenString string
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenString = authHeader[7:]
		} else {
			logger.Logger.Warn("invalid authorization header format")
			return nil, status.Error(codes.Unauthenticated, "invalid authorization header format")
		}

		// 验证并解析 token
		claims, err := jwt.ParseToken(tokenString)
		if err != nil {
			logger.Logger.Warn("token validation failed: ", err)

			// 区分 token 过期和 token 无效
			if errors.Is(err, jwtLib.ErrTokenExpired) {
				return nil, status.Error(codes.Unauthenticated, "token expired")
			}
			return nil, status.Error(codes.Unauthenticated, "token invalid")
		}

		// 将用户信息注入到 context 中
		ctx = context.WithValue(ctx, "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "username", claims.Username)
		ctx = context.WithValue(ctx, "phone_number", claims.PhoneNumber)
		ctx = context.WithValue(ctx, "device_fingerprint", claims.DeviceFingerprint)
		ctx = context.WithValue(ctx, "token", tokenString)

		logger.Logger.Info("authenticated user: ", claims.UserID, " (", claims.Username, ")")

		// 继续处理请求
		return handler(ctx, req)
	}
}

// AuthStreamInterceptor JWT 认证拦截器（流式调用）
func AuthStreamInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// 检查是否在白名单中
		if authWhitelist[info.FullMethod] {
			return handler(srv, ss)
		}

		// 从 metadata 中获取 token
		ctx := ss.Context()
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			logger.Logger.Warn("missing metadata")
			return status.Error(codes.Unauthenticated, "missing metadata")
		}

		// 获取 Authorization header
		authHeaders := md.Get("authorization")
		if len(authHeaders) == 0 {
			logger.Logger.Warn("missing authorization header")
			return status.Error(codes.Unauthenticated, "missing authorization token")
		}

		authHeader := authHeaders[0]

		// 解析 Bearer token 格式
		var tokenString string
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenString = authHeader[7:]
		} else {
			logger.Logger.Warn("invalid authorization header format")
			return status.Error(codes.Unauthenticated, "invalid authorization header format")
		}

		// 验证并解析 token
		claims, err := jwt.ParseToken(tokenString)
		if err != nil {
			logger.Logger.Warn("token validation failed: ", err)

			if errors.Is(err, jwtLib.ErrTokenExpired) {
				return status.Error(codes.Unauthenticated, "token expired")
			}
			return status.Error(codes.Unauthenticated, "token invalid")
		}

		// 将用户信息注入到 context 中
		ctx = context.WithValue(ctx, "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "username", claims.Username)
		ctx = context.WithValue(ctx, "phone_number", claims.PhoneNumber)
		ctx = context.WithValue(ctx, "device_fingerprint", claims.DeviceFingerprint)
		ctx = context.WithValue(ctx, "token", tokenString)

		logger.Logger.Info("authenticated user: ", claims.UserID, " (", claims.Username, ")")

		// 创建新的 ServerStream 包装器
		wrappedStream := &wrappedServerStream{
			ServerStream: ss,
			ctx:          ctx,
		}

		// 继续处理请求
		return handler(srv, wrappedStream)
	}
}

// wrappedServerStream 包装 ServerStream 以注入新的 context
type wrappedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *wrappedServerStream) Context() context.Context {
	return w.ctx
}
