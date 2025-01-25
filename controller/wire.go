//go:build wireinject

package controller

import (
	"chatroom/dao"
	"chatroom/dao/cache"
	"chatroom/pkg/business"
	"chatroom/service"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	wire.Struct(new(UserController), "*"),
	wire.Struct(new(AuthController), "*"),
	wire.Struct(new(SessionController), "*"),
	dao.ProviderSet,
	cache.ProviderSet,
	service.ProviderSet,
	business.ProviderSet,
	wire.Struct(new(Deps), "*"),
	wire.Struct(new(Controllers), "*"),
)
