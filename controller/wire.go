//go:build wireinject

package controller

import (
	"chatroom/dao"
	"chatroom/service"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	wire.Struct(new(UserController), "*"),
	dao.NewUsers,
	service.ProviderSet,
	wire.Struct(new(Deps), "*"),
	wire.Struct(new(Controllers), "*"),
)
