package handler

import (
	"chatroom/context"
	"chatroom/pkg/core/socket"
	"chatroom/pkg/core/socket/adapter"
	"chatroom/service"
	"chatroom/socket/handler/event"
	"log"
)

type ChatChannel struct {
	Storage service.IClientConnectService
	Event   *event.ChatEvent
}

// Conn 初始化连接
func (c *ChatChannel) Conn(ctx *context.Context) error {
	log.Printf("Attempting WebSocket connection with token: %s", "ade")
	conn, err := adapter.NewWsAdapter(ctx.Context.Writer, ctx.Context.Request)
	if err != nil {
		log.Printf("websocket connect error: %s", err.Error())
		return err
	}

	return c.NewClient(ctx.UserId(), conn)
}

func (c *ChatChannel) NewClient(uid int, conn socket.IConn) error {
	return socket.NewClient(conn, &socket.ClientOption{
		Uid:     uid,
		Channel: socket.Session.Chat,
		Storage: c.Storage,
		Buffer:  10,
	}, socket.NewEvent(
		// 连接成功回调
		socket.WithOpenEvent(c.Event.OnOpen), //推送自己已经上线
		// 接收消息回调
		socket.WithMessageEvent(c.Event.OnMessage), //发送消息的回调
		// 关闭连接回调
		socket.WithCloseEvent(c.Event.OnClose), //下线的回调
	))
}
