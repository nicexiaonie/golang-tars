package user

import (
	"bufio"
	"context"
	"fmt"
	"go-base/internal/user/interceptor"
	logger "go-base/pkg"
	"go-base/pkg/config"
	"go-base/pkg/proto/demo/v1/demov1connect"
	"net"
	"net/http"
	"time"

	"connectrpc.com/connect"
	"connectrpc.com/vanguard"
	"github.com/TarsCloud/TarsGo/tars"
	"github.com/nicexiaonie/gconf"
	"github.com/nicexiaonie/ghelper"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// connectServer 全局变量，用于统一管理关闭
var connectServer *http.Server

// protoJSONCodec 自定义JSON编解码器，配置为输出零值字段
type protoJSONCodec struct {
	name string
}

// newProtoJSONCodec 创建支持输出零值字段的JSON编解码器
func newProtoJSONCodec() connect.Codec {
	return &protoJSONCodec{name: "json"}
}

func (c *protoJSONCodec) Name() string {
	return c.name
}

func (c *protoJSONCodec) Marshal(v any) ([]byte, error) {
	msg, ok := v.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("not a proto.Message")
	}
	// 使用 protojson，启用 EmitDefaultValues 输出零值字段
	return protojson.MarshalOptions{
		UseProtoNames:     false, // 使用 JSON 字段名（驼峰）
		EmitDefaultValues: true,  // 输出零值字段（关键！）
		EmitUnpopulated:   true,  // 输出未设置的字段
	}.Marshal(msg)
}

func (c *protoJSONCodec) Unmarshal(data []byte, v any) error {
	msg, ok := v.(proto.Message)
	if !ok {
		return fmt.Errorf("not a proto.Message")
	}
	return protojson.UnmarshalOptions{
		DiscardUnknown: true, // 忽略未知字段
	}.Unmarshal(data, msg)
}

// NewTarsPb 启动 Tars 集成的 Connect 服务（Connect 直接集成在 Tars HTTP Servant 下）
func NewTarsPb() {
	cfg := tars.GetServerConfig()
	if cfg == nil {
		panic("get server config failed")
	}
	// 初始化日志
	logPath := cfg.LogPath + "/" + cfg.App + "/" + cfg.Server + "/"
	if err := logger.LoggerInit(logPath); err != nil {
		panic(fmt.Sprintf("failed to initialize logger: %v", err))
	}
	logger.Logger.Info(fmt.Sprintf("NewTarsPb配置: %s", ghelper.ToString(cfg)))
	logger.Logger.Info(cfg.BasePath + "config/")

	// 初始化配置
	config.InitConfig()

	// 创建 Connect Handler 并集成到 Tars HTTP Servant
	connectMux := createConnectHandler()

	// 使用 h2c 包装以支持 HTTP/2（gRPC 协议需要）
	h2cHandler := h2c.NewHandler(connectMux, &http2.Server{})

	// 将 Connect Handler 包装到 TarsHttpMux
	tarsHttpMux := &tars.TarsHttpMux{}
	tarsHttpMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 包装 ResponseWriter 以支持 Flusher 接口（Vanguard 要求）
		h2cHandler.ServeHTTP(newFlushableResponseWriter(w), r)
	})

	// 注册到 Tars HTTP Servant
	tars.AddHttpServant(tarsHttpMux, cfg.App+"."+cfg.Server+".HttpServer")

	logger.Logger.Info("Connect service integrated with Tars HTTP Servant")
	logger.Logger.Info("Supported protocols: gRPC, gRPC-Web, Connect (HTTP/JSON)")

	// 启动独立的 Connect 服务器（非阻塞）
	// 等待启动完成确保服务就绪
	ready := startConnectServer()
	<-ready
	logger.Logger.Info("Connect server started successfully")

	// Run Tars application（阻塞调用，收到关闭信号后返回）
	tars.Run()

	// Tars 退出后关闭 Connect 服务器
	// 注意：不使用 defer，因为 tars.Run() 可能通过 os.Exit() 退出导致 defer 不执行
	shutdownConnectServer()
}

// createConnectHandler 创建 Connect HTTP Handler（使用 Vanguard 支持 RESTful API）
func createConnectHandler() *http.ServeMux {
	mux := http.NewServeMux()

	// 健康检查接口
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// 使用 Wire 自动注入依赖
	// Wire 会自动解析依赖关系：Handler -> Service -> DAO/Proxy
	demoHandler, err := InitializeHandler()
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("Failed to initialize handler: %v", err))
		panic(err)
	}

	// 注册 Connect 服务处理器（带拦截器和JSON序列化配置）
	path, connectHandler := demov1connect.NewDemoServiceHandler(
		demoHandler,
		interceptor.NewInterceptors(),
		connect.WithCodec(newProtoJSONCodec()), // 使用自定义JSON编解码器，输出零值字段
	)

	// 使用 Vanguard 转换器，支持 RESTful 风格的 API
	// Vanguard 会自动将 HTTP 请求转换为 gRPC/Connect 调用
	transcoder, err := vanguard.NewTranscoder([]*vanguard.Service{
		vanguard.NewService(demov1connect.DemoServiceName, connectHandler),
	})
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("Failed to create Vanguard transcoder: %v", err))
		panic(err)
	}

	// 注册 Vanguard 处理器到根路径
	// 这样可以同时支持：
	// 1. gRPC 路径: /demo.v1.DemoService/CreateDemo
	// 2. RESTful 路径: /api/v1/demo (需要在 proto 中配置 google.api.http 注解)
	mux.Handle("/", transcoder)

	logger.Logger.Info(fmt.Sprintf("Registered Connect service with Vanguard: %s", path))
	logger.Logger.Info("Vanguard transcoder enabled - supports both gRPC and RESTful APIs")

	return mux
}

// flushableResponseWriter 包装 ResponseWriter 以支持 Flusher 接口
// Vanguard 要求 ResponseWriter 必须实现 http.Flusher 接口
type flushableResponseWriter struct {
	http.ResponseWriter
}

// newFlushableResponseWriter 创建一个支持 Flush 的 ResponseWriter
func newFlushableResponseWriter(w http.ResponseWriter) http.ResponseWriter {
	return &flushableResponseWriter{ResponseWriter: w}
}

// Flush 实现 http.Flusher 接口
func (f *flushableResponseWriter) Flush() {
	// 如果底层 ResponseWriter 支持 Flush，则调用它
	if flusher, ok := f.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
	// 否则什么都不做（静默处理）
}

// Hijack 实现 http.Hijacker 接口（如果需要 WebSocket 支持）
func (f *flushableResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hijacker, ok := f.ResponseWriter.(http.Hijacker); ok {
		return hijacker.Hijack()
	}
	return nil, nil, fmt.Errorf("ResponseWriter does not support hijacking")
}

// Push 实现 http.Pusher 接口（HTTP/2 Server Push）
func (f *flushableResponseWriter) Push(target string, opts *http.PushOptions) error {
	if pusher, ok := f.ResponseWriter.(http.Pusher); ok {
		return pusher.Push(target, opts)
	}
	return fmt.Errorf("ResponseWriter does not support push")
}

// startConnectServer 启动独立的 Connect 服务器（支持 gRPC + HTTP/JSON 双协议）
// 返回启动完成的 channel，调用者可以等待启动完成
func startConnectServer() <-chan struct{} {
	ready := make(chan struct{})

	go func() {
		// 获取端口配置
		port := gconf.GetInt("grpc.port")
		if port == 0 {
			port = 9091 // 默认端口
		}

		// 创建 Connect Handler
		mux := createConnectHandler()

		// 创建支持 h2c（HTTP/2 without TLS）的服务器
		// 这样可以在单端口同时支持：
		// - HTTP/1.1 + JSON (Connect 协议)
		// - HTTP/2 + Protobuf (gRPC 协议)
		// - HTTP/1.1 + Protobuf (gRPC-Web 协议)
		h2cHandler := h2c.NewHandler(mux, &http2.Server{})

		connectServer = &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: h2cHandler,
			// 性能配置
			ReadTimeout:    30 * time.Second,
			WriteTimeout:   30 * time.Second,
			IdleTimeout:    120 * time.Second,
			MaxHeaderBytes: 1 << 20, // 1 MB
		}

		logger.Logger.Info(fmt.Sprintf("Connect server starting on :%d", port))
		logger.Logger.Info("Supported protocols: gRPC, gRPC-Web, Connect (HTTP/JSON)")

		// 标记启动完成（在实际监听之前）
		close(ready)

		// 启动服务器（阻塞调用）
		if err := connectServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Logger.Error(fmt.Sprintf("Connect server failed: %v", err))
			panic(err)
		}

		logger.Logger.Info("Connect server stopped")
	}()

	return ready
}

// shutdownConnectServer 优雅关闭 Connect 服务器
// 由 Tars 框架在关闭时调用
func shutdownConnectServer() {
	if connectServer == nil {
		return
	}

	logger.Logger.Info("Shutting down Connect server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := connectServer.Shutdown(ctx); err != nil {
		logger.Logger.Error(fmt.Sprintf("Connect server shutdown error: %v", err))
	} else {
		logger.Logger.Info("Connect server shutdown completed")
	}
}

// NewTarsPbWithSeparateConnect 启动 Tars + 独立 Connect 服务（双端口模式）
// 适用于需要完全独立的 gRPC 端口的场景
// func NewTarsPbWithSeparateConnect() {
// 	cfg := tars.GetServerConfig()
// 	if cfg == nil {
// 		panic("get server config failed")
// 	}

// 	// 初始化日志
// 	logPath := cfg.LogPath + "/" + cfg.App + "/" + cfg.Server + "/"
// 	logger.LoggerInit(logPath)
// 	logger.Logger.Info(fmt.Sprintf("NewTarsPbWithSeparateConnect: %s", ghelper.ToString(cfg)))

// 	// 初始化配置
// 	gconf.Init(
// 		gconf.WithConfigPaths(cfg.BasePath + "config/"),
// 	)

// 	// 启动独立的 Connect 服务器（在独立的 goroutine 中）
// 	go startConnectServer()

// 	// Tars HTTP 只提供健康检查
// 	tarsHttpMux := &tars.TarsHttpMux{}
// 	tarsHttpMux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("Content-Type", "application/json")
// 		w.WriteHeader(http.StatusOK)
// 		w.Write([]byte(`{"status":"ok"}`))
// 	})
// 	tars.AddHttpServant(tarsHttpMux, cfg.App+"."+cfg.Server+".HttpServer")

// 	// Run Tars application
// 	tars.Run()
// }

// NewConnectServer 启动独立的 Connect 服务（不依赖 Tars）
// func NewConnectServer() {
// 	// 初始化日志
// 	logger.LoggerInit("./logs/")

// 	// 初始化配置
// 	gconf.Init(
// 		gconf.WithConfigPaths("./config/"),
// 	)
// 	logger.Logger.Info(fmt.Sprintf("全局配置: %s", ghelper.ToString(gconf.GetStringMap("app"))))

// 	// 启动 Connect 服务器
// 	startConnectServer()
// }
