//go:build wireinject

package main

import (
	"chatroom/config"
	"chatroom/controller"
	"chatroom/dao"
	"chatroom/dao/cache"
	"chatroom/pkg/business"
	"chatroom/pkg/client"
	"chatroom/pkg/core"
	"chatroom/service"
	"chatroom/socket"
	"github.com/google/wire"
)

var providerSet = wire.NewSet(
	client.NewMySQLClient,
	client.NewEmailClient,
	client.NewRedisClient,
	config.NewFilesystem,
	// 基础服务
	//provider.NewMySQLClient,
	//provider.NewRedisClient,
	//provider.NewHttpClient,
	//provider.NewFilesystem,
	//provider.NewBase64Captcha,
	//provider.NewIpAddressClient,
	wire.Struct(new(client.Providers), "*"),
	//
	//repo.ProviderSet,     // 注入 Repo 依赖
	//business.ProviderSet, // 注入 Logic 依赖
	//service.ProviderSet,  // 注入 Service 依赖
)

func NewHttpInjector(conf *config.Config) *core.AppProvider {
	panic(
		wire.Build(
			providerSet,
			core.ProviderSet,
			controller.ProviderSet,
		),
	)
}

func NewSocketInjector(conf *config.Config) *socket.AppProvider {
	panic(
		wire.Build(
			dao.ProviderSet,
			cache.ProviderSet, // 注入 Cache 依赖
			providerSet,
			socket.ProviderSet,
			service.ProviderSet,
			business.ProviderSet,
		),
	)
}
