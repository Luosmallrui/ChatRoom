//go:build wireinject

package controller

import (
	"chatroom/dao"
	"chatroom/service"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	wire.Struct(new(UserController), "*"),
	wire.Struct(new(AuthController), "*"),
	dao.NewCaptchaStorage,
	dao.NewBase64Captcha,
	dao.NewTokenSessionStorage,
	dao.NewUsers,
	dao.NewAdmin,
	service.ProviderSet,
	wire.Struct(new(Deps), "*"),
	wire.Struct(new(Controllers), "*"),
)
