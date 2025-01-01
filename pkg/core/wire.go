//go:build wireinject
// +build wireinject

package core

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewGinServer,
	wire.Struct(new(AppProvider), "*"),
)
