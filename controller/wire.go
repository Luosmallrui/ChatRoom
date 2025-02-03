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
	wire.Struct(new(User), "*"),
	wire.Struct(new(Auth), "*"),
	wire.Struct(new(Session), "*"),
	wire.Struct(new(Contact), "*"),
	wire.Struct(new(Group), "*"),
	wire.Struct(new(Emoticon), "*"),
	wire.Struct(new(Publish), "*"),
	dao.ProviderSet,
	cache.ProviderSet,
	service.ProviderSet,
	business.ProviderSet,
	wire.Struct(new(Deps), "*"),
	wire.Struct(new(Controllers), "*"),
)
