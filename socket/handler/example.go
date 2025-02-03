package handler

import (
	"chatroom/dao/cache"
	"chatroom/pkg/core/socket"
	"chatroom/pkg/core/socket/adapter"
	"chatroom/socket/handler/event"
	"log"

	"chatroom/context"
)

// ExampleChannel 案例
type ExampleChannel struct {
	Storage *cache.ClientStorage
	Event   *event.ExampleEvent
}

func (c *ExampleChannel) Conn(ctx *context.Context) error {

	conn, err := adapter.NewWsAdapter(ctx.Context.Writer, ctx.Context.Request)
	if err != nil {
		log.Printf("websocket connect error: %s", err.Error())
		return err
	}

	return socket.NewClient(conn, &socket.ClientOption{
		Channel: socket.Session.Example,
		Uid:     0,
	}, socket.NewEvent(
		// 连接成功回调
		socket.WithOpenEvent(c.Event.OnOpen),
		// 接收消息回调
		socket.WithMessageEvent(c.Event.OnMessage),
		// 关闭连接回调
		socket.WithCloseEvent(c.Event.OnClose),
	))
}
