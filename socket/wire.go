package socket

import (
	"chatroom/socket/consume"
	"chatroom/socket/handler"
	"chatroom/socket/handler/event"
	"chatroom/socket/process"
	"chatroom/socket/process/queue"
	"chatroom/socket/router"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	router.NewRouter,

	wire.Struct(new(handler.Handler), "*"),

	// process
	wire.Struct(new(process.SubServers), "*"),
	process.NewServer,
	process.NewHealthSubscribe,
	process.NewMessageSubscribe,
	wire.Struct(new(process.QueueSubscribe), "*"),
	wire.Struct(new(queue.GlobalMessage), "*"),
	wire.Struct(new(queue.LocalMessage), "*"),
	wire.Struct(new(queue.RoomControl), "*"),

	handler.ProviderSet,
	event.ProviderSet,
	consume.ProviderSet,

	// AppProvider
	wire.Struct(new(AppProvider), "*"),
)
