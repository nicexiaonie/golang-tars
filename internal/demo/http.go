package server

// import (
// 	"context"
// 	"fmt"
// 	logger "golang-tars/pkg"
// 	demov1 "golang-tars/pkg/proto/demo/v1"
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// 	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
// 	"github.com/nicexiaonie/gconf"
// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/credentials/insecure"
// )

// // setupHTTPServer 设置 HTTP 服务器（使用 gRPC-Gateway）
// func setupHTTPServer() *gin.Engine {
// 	// 创建 Gin 引擎
// 	ginEngine := gin.New()

// 	// 添加 Gin 中间件
// 	ginEngine.Use(gin.Recovery())
// 	ginEngine.Use(gin.Logger())

// 	// 健康检查接口
// 	ginEngine.GET("/health", func(c *gin.Context) {
// 		c.JSON(200, gin.H{"status": "ok"})
// 	})

// 	// 创建 gRPC-Gateway mux
// 	grpcPort := gconf.GetInt("grpc.port")
// 	if grpcPort == 0 {
// 		grpcPort = 9090
// 	}

// 	mux := runtime.NewServeMux(
// 		// 配置 gRPC-Gateway 选项
// 		runtime.WithIncomingHeaderMatcher(customHeaderMatcher),
// 		runtime.WithErrorHandler(customErrorHandler),
// 	)

// 	// 连接到 gRPC 服务器
// 	opts := []grpc.DialOption{
// 		grpc.WithTransportCredentials(insecure.NewCredentials()),
// 	}

// 	grpcEndpoint := fmt.Sprintf("localhost:%d", grpcPort)

// 	// 注册 Demo 服务的 HTTP 处理器
// 	err := demov1.RegisterDemoServiceHandlerFromEndpoint(
// 		context.Background(),
// 		mux,
// 		grpcEndpoint,
// 		opts,
// 	)
// 	if err != nil {
// 		logger.Logger.Error(fmt.Sprintf("failed to register demo service handler: %v", err))
// 		panic(err)
// 	}

// 	// 将 gRPC-Gateway 挂载到 Gin
// 	ginEngine.Any("/api/*path", gin.WrapH(mux))

// 	logger.Logger.Info("HTTP server (gRPC-Gateway) configured")

// 	return ginEngine
// }

// // customHeaderMatcher 自定义 header 匹配器（用于传递 Authorization 等 header）
// func customHeaderMatcher(key string) (string, bool) {
// 	switch key {
// 	case "Authorization", "X-Request-Id", "X-Trace-Id":
// 		return key, true
// 	default:
// 		return runtime.DefaultHeaderMatcher(key)
// 	}
// }

// // customErrorHandler 自定义错误处理器
// func customErrorHandler(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
// 	// 使用默认的错误处理器
// 	runtime.DefaultHTTPErrorHandler(ctx, mux, marshaler, w, r, err)
// }
