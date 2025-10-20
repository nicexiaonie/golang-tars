package server

import (
	"fmt"
	logger "golang-tars/pkg"
	"golang-tars/pkg/consul"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/TarsCloud/TarsGo/tars"
	"github.com/nicexiaonie/gconf"
	"github.com/nicexiaonie/ghelper"
)

var consulClient *consul.Client

func NewTarsPb() {

	cfg := tars.GetServerConfig()
	if cfg == nil {
		panic("get server config failed")
	}

	// 初始化日志
	logPath := cfg.LogPath + "/" + cfg.App + "/" + cfg.Server + "/"
	logger.LoggerInit(logPath)
	logger.Logger.Info(fmt.Sprintf("NewTarsPb: %s", ghelper.ToString(cfg)))
	logger.Logger.Info(cfg.BasePath + "config/")

	// 初始化配置
	// 初始化全局配置
	gconf.Init(
		gconf.WithConfigPaths(cfg.BasePath + "config/"),
	)
	logger.Logger.Info(fmt.Sprintf("全局配置: %s", ghelper.ToString(gconf.GetStringMap("app"))))

	// 同步配置
	if cfg.Node != "" && len(cfg.Node) > 0 {

		// remoteServerConf := tars.NewRConf(cfg.App, cfg.Server, cfg.BasePath)
		// 1.拉取远程配置文件到本地 tars后台的配置
		// _, err = remoteServerConf.GetConfig(cfg.Server + ".yaml")
		// if err != nil {
		// 	logger.Logger.Error(fmt.Sprintf("get server config: %s", err.Error()))
		// }
	}

	// http服务  Demo.UserServer.HttpServerObj	tcp -h 172.17.239.228 -t 60000 -p 17191 -e 0	5	100000	50000	20000
	mux := &tars.TarsHttpMux{}
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello tars"))
	})
	tars.AddHttpServant(mux, cfg.App+"."+cfg.Server+".HttpServerObj") //Register http server

	// Run application
	tars.Run()

}

// setupGracefulShutdown 设置优雅关闭
func setupGracefulShutdown() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigCh
		logger.Logger.Info(fmt.Sprintf("Received signal: %v, shutting down gracefully...", sig))

		// 注销Consul服务
		if consulClient != nil {
			if err := consulClient.Close(); err != nil {
				logger.Logger.Error(fmt.Sprintf("deregister service from consul failed: %s", err.Error()))
			} else {
				logger.Logger.Info("Service deregistered from Consul")
			}
		}

		os.Exit(0)
	}()
}
