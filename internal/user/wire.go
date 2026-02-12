//go:build wireinject
// +build wireinject

package user

import (
	"golang-tars/internal/user/data/mysql/dao"
	"golang-tars/internal/user/handler"
	"golang-tars/internal/user/proxy"
	"golang-tars/internal/user/service"

	"github.com/google/wire"
)

// ProviderSet 定义所有的依赖注入 Provider
// Wire 会按照依赖关系自动排序并生成初始化代码
var ProviderSet = wire.NewSet(
	// 基础设施层 Provider
	ProvideDBManager,
	ProvideDemoProxyConfig,

	// DAO 层（数据访问）
	dao.NewDemoDAO,
	dao.NewUserDAO,
	dao.NewUserLoginLogDAO,

	// Proxy 层（第三方服务调用）
	proxy.NewDemoProxy,

	// Service 层（业务逻辑）
	service.NewDemoService,

	// Handler 层（API 处理）
	handler.NewDemoHandler,
)

// InitializeHandler 初始化 Handler
// Wire 会自动生成这个函数的实现，解析所有依赖关系
// 生成的代码在 wire_gen.go 文件中
func InitializeHandler() (*handler.DemoHandler, error) {
	wire.Build(ProviderSet)
	return nil, nil
}
