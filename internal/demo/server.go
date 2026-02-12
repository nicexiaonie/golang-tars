package server

// import (
// 	"fmt"
// 	logger "go-base/pkg"
// 	"net/http"

// 	"github.com/TarsCloud/TarsGo/tars"
// 	"github.com/nicexiaonie/gconf"
// 	"github.com/nicexiaonie/ghelper"
// )

// func NewTarsPb() {

// 	cfg := tars.GetServerConfig()
// 	if cfg == nil {
// 		panic("get server config failed")
// 	}

// 	// 初始化日志
// 	logPath := cfg.LogPath + "/" + cfg.App + "/" + cfg.Server + "/"
// 	logger.LoggerInit(logPath)
// 	logger.Logger.Info(fmt.Sprintf("NewTarsPb: %s", ghelper.ToString(cfg)))
// 	logger.Logger.Info(cfg.BasePath + "config/")

// 	// 初始化配置
// 	// 初始化全局配置
// 	gconf.Init(
// 		gconf.WithConfigPaths(cfg.BasePath + "config/"),
// 	)
// 	logger.Logger.Info(fmt.Sprintf("全局配置: %s", ghelper.ToString(gconf.GetStringMap("app"))))

// 	// 同步配置
// 	if cfg.Node != "" && len(cfg.Node) > 0 {

// 		// remoteServerConf := tars.NewRConf(cfg.App, cfg.Server, cfg.BasePath)
// 		// 1.拉取远程配置文件到本地 tars后台的配置
// 		// _, err = remoteServerConf.GetConfig(cfg.Server + ".yaml")
// 		// if err != nil {
// 		// 	logger.Logger.Error(fmt.Sprintf("get server config: %s", err.Error()))
// 		// }
// 	}

// 	// 启动 gRPC 服务器（在独立的 goroutine 中）
// 	go startGRPCServer()

// 	// 启动 HTTP 服务器（使用 gRPC-Gateway）
// 	ginEngine := setupHTTPServer()

// 	// 将 Gin 引擎包装为 TarsHttpMux
// 	tarsHttpMux := &tars.TarsHttpMux{}
// 	tarsHttpMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 		ginEngine.ServeHTTP(w, r)
// 	})
// 	tars.AddHttpServant(tarsHttpMux, cfg.App+"."+cfg.Server+".HttpServer")

// 	// Run application
// 	tars.Run()

// }
