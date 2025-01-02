//go:build wireinject

package dao

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	wire.Struct(new(Users), "*"),
)
