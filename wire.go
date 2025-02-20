//go:build wireinject
// +build wireinject

package main

import (
	"chatroom/config"
	"chatroom/controller"
	"chatroom/dao"
	"chatroom/dao/cache"
	"chatroom/pkg/client"
	"chatroom/pkg/core"
	"chatroom/pkg/kafka"
	"chatroom/service"
	"chatroom/socket"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	client.NewMySQLClient,
	client.NewEmailClient,
	client.NewRedisClient,
	config.NewFilesystem,
	kafka.NewKafkaClient,
	wire.Struct(new(client.Providers), "*"),
)

func NewHttpInjector(conf *config.Config) *core.AppProvider {
	panic(
		wire.Build(
			ProviderSet,
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
			ProviderSet,
			socket.ProviderSet,
			service.ProviderSet,
		),
	)
}
