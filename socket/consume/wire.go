//go:build wireinject
// +build wireinject

package consume

import (
	"chatroom/socket/consume/chat"
	"chatroom/socket/consume/example"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewChatSubscribe,
	wire.Struct(new(chat.Handler), "*"),

	NewExampleSubscribe,
	example.NewHandler,
)
