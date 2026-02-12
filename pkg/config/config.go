package config

import (
	"fmt"
	logger "go-base/pkg"

	"github.com/TarsCloud/TarsGo/tars"
	"github.com/nicexiaonie/gconf"
	"github.com/nicexiaonie/ghelper"
)

var CONFIG_LIST = map[string]*gconf.Gconf{}

func InitConfig() {
	cfg := tars.GetServerConfig()

	// 同步配置（如果在 Tars 平台上）
	if cfg.Node != "" && len(cfg.Node) > 0 {
		remoteConfigConf := tars.NewRConf(cfg.App, cfg.Server, cfg.BasePath)
		_, err := remoteConfigConf.GetConfig("config.yaml")
		if err != nil {
			logger.Logger.Error(fmt.Sprintf("get config config: %s", err.Error()))
		}

		// 如果拉取应用级别配置 则不要传入server值
		remoteGrpcConf := tars.NewRConf(cfg.App, "", cfg.BasePath)
		_, err = remoteGrpcConf.GetConfig("grpc.yaml")
		if err != nil {
			logger.Logger.Error(fmt.Sprintf("get grpc config: %s", err.Error()))
		}
	}

	// 初始化config配置
	err := gconf.Init(
		gconf.WithConfigPaths(cfg.BasePath),
		// 启用配置文件监听和热更新
		gconf.WithWatchConfig(true),
	)
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("init gconf: %s", err.Error()))
		panic(err)
	}
	logger.Logger.Info(fmt.Sprintf("全局配置: %s", ghelper.ToString(gconf.GetStringMap("app"))))
	CONFIG_LIST["config"] = gconf.GetInstance()

	// 初始化grpc配置
	gRpcConf, err := gconf.New(
		gconf.WithConfigName("grpc"),
		gconf.WithConfigType("yaml"),
		gconf.WithConfigPaths(cfg.BasePath),
		// 启用配置文件监听和热更新
		gconf.WithWatchConfig(true),
		// gconf.WithConfigPaths("./config"),
	)

	if err != nil {
		logger.Logger.Error(fmt.Sprintf("init gconf: %s", err.Error()))
		panic(err)
	}
	logger.Logger.Info(fmt.Sprintf("gRpc配置: %s", ghelper.ToString(gRpcConf.GetViper().AllKeys())))
	CONFIG_LIST["grpc"] = gRpcConf
}

func GetConfig(name string) *gconf.Gconf {
	return CONFIG_LIST[name]
}
