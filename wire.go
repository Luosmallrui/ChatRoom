//go:build wireinject

package main

import (
	"chatroom/config"
	"chatroom/controller"
	"chatroom/pkg/core"
	"github.com/google/wire"
)

func NewHttpInjector(conf *config.Config) *core.AppProvider {
	panic(
		wire.Build(
			core.ProviderSet,
			controller.ProviderSet,
		),
	)
}
