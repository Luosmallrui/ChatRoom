//go:build wireinject
// +build wireinject

// chatroom/socket/handler/event/wire/go
package event

import (
	"chatroom/socket/handler/event/chat"
	"chatroom/socket/handler/event/example"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	wire.Struct(new(ChatEvent), "*"),
	wire.Struct(new(chat.Handler), "*"),

	wire.Struct(new(ExampleEvent), "*"),
	example.NewHandler,
)
