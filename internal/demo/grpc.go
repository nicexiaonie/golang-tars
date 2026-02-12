package server

// import (
// 	"fmt"
// 	"go-base/internal/user/interceptor"
// 	"go-base/internal/user/service"
// 	logger "go-base/pkg"
// 	demov1 "go-base/pkg/proto/demo/v1"
// 	"net"
// 	"time"

// 	"github.com/nicexiaonie/gconf"
// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/keepalive"
// 	"google.golang.org/grpc/reflection"
// )

// // startGRPCServer 启动 gRPC 服务器
// func startGRPCServer() {
// 	// 监听端口
// 	grpcPort := gconf.GetInt("grpc.port")
// 	if grpcPort == 0 {
// 		grpcPort = 9090 // 默认端口
// 	}

// 	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
// 	if err != nil {
// 		logger.Logger.Error(fmt.Sprintf("failed to listen gRPC port %d: %v", grpcPort, err))
// 		panic(err)
// 	}

// 	// 从配置文件读取 gRPC 性能参数
// 	maxConcurrentStreams := gconf.GetInt("grpc.max_concurrent_streams")
// 	if maxConcurrentStreams == 0 {
// 		maxConcurrentStreams = 1000 // 默认值
// 	}

// 	maxRecvMsgSize := gconf.GetInt("grpc.max_recv_msg_size")
// 	if maxRecvMsgSize == 0 {
// 		maxRecvMsgSize = 4 * 1024 * 1024 // 默认 4MB
// 	}

// 	maxSendMsgSize := gconf.GetInt("grpc.max_send_msg_size")
// 	if maxSendMsgSize == 0 {
// 		maxSendMsgSize = 4 * 1024 * 1024 // 默认 4MB
// 	}

// 	// KeepAlive 配置
// 	maxConnectionIdle := gconf.GetInt("grpc.keepalive.max_connection_idle")
// 	maxConnectionAge := gconf.GetInt("grpc.keepalive.max_connection_age")
// 	maxConnectionAgeGrace := gconf.GetInt("grpc.keepalive.max_connection_age_grace")
// 	keepaliveTime := gconf.GetInt("grpc.keepalive.time")
// 	keepaliveTimeout := gconf.GetInt("grpc.keepalive.timeout")
// 	minTime := gconf.GetInt("grpc.keepalive.min_time")
// 	permitWithoutStream := gconf.GetBool("grpc.keepalive.permit_without_stream")

// 	connectionTimeout := gconf.GetInt("grpc.connection_timeout")
// 	if connectionTimeout == 0 {
// 		connectionTimeout = 120 // 默认 120 秒
// 	}

// 	// 创建 gRPC 服务器选项
// 	serverOptions := []grpc.ServerOption{
// 		// 1. 并发控制：限制每个连接的最大并发流数量
// 		grpc.MaxConcurrentStreams(uint32(maxConcurrentStreams)),

// 		// 2. 消息大小限制：防止恶意大消息攻击
// 		grpc.MaxRecvMsgSize(maxRecvMsgSize),
// 		grpc.MaxSendMsgSize(maxSendMsgSize),

// 		// 3. KeepAlive 配置：防止僵尸连接累积
// 		grpc.KeepaliveParams(keepalive.ServerParameters{
// 			MaxConnectionIdle:     time.Duration(maxConnectionIdle) * time.Second,
// 			MaxConnectionAge:      time.Duration(maxConnectionAge) * time.Second,
// 			MaxConnectionAgeGrace: time.Duration(maxConnectionAgeGrace) * time.Second,
// 			Time:                  time.Duration(keepaliveTime) * time.Second,
// 			Timeout:               time.Duration(keepaliveTimeout) * time.Second,
// 		}),

// 		// 4. KeepAlive 策略：限制客户端 ping 行为
// 		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
// 			MinTime:             time.Duration(minTime) * time.Second,
// 			PermitWithoutStream: permitWithoutStream,
// 		}),

// 		// 5. 连接超时：防止慢客户端长期占用连接
// 		grpc.ConnectionTimeout(time.Duration(connectionTimeout) * time.Second),

// 		// 6. 拦截器链
// 		grpc.ChainUnaryInterceptor(
// 			interceptor.RecoveryUnaryInterceptor(), // 错误恢复（最外层）
// 			interceptor.LogUnaryInterceptor(),      // 日志记录
// 			interceptor.AuthUnaryInterceptor(),     // JWT 认证
// 		),
// 		grpc.ChainStreamInterceptor(
// 			interceptor.RecoveryStreamInterceptor(), // 错误恢复（最外层）
// 			interceptor.LogStreamInterceptor(),      // 日志记录
// 			interceptor.AuthStreamInterceptor(),     // JWT 认证
// 		),
// 	}

// 	// 创建 gRPC 服务器
// 	grpcServer := grpc.NewServer(serverOptions...)

// 	// 注册 Demo 服务
// 	demoService := service.NewDemoService()
// 	demov1.RegisterDemoServiceServer(grpcServer, demoService)

// 	// 注册反射服务（只在非生产环境使用）
// 	env := gconf.GetString("app.env")
// 	if env != "production" {
// 		reflection.Register(grpcServer)
// 		logger.Logger.Info("gRPC reflection service enabled (non-production environment)")
// 	}

// 	logger.Logger.Info(fmt.Sprintf("gRPC server listening on :%d", grpcPort))
// 	logger.Logger.Info(fmt.Sprintf("gRPC server config: max_concurrent_streams=%d, max_recv_msg_size=%d, max_send_msg_size=%d",
// 		maxConcurrentStreams, maxRecvMsgSize, maxSendMsgSize))
// 	logger.Logger.Info(fmt.Sprintf("gRPC KeepAlive: max_connection_idle=%ds, max_connection_age=%ds, time=%ds, timeout=%ds",
// 		maxConnectionIdle, maxConnectionAge, keepaliveTime, keepaliveTimeout))

// 	// 启动服务器
// 	if err := grpcServer.Serve(lis); err != nil {
// 		logger.Logger.Error(fmt.Sprintf("failed to serve gRPC: %v", err))
// 		panic(err)
// 	}
// }
