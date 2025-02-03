package queue

import (
	"chatroom/pkg/core"
	"chatroom/pkg/core/socket"
	"chatroom/socket/consume"
	"context"
)

var _ consume.IConsumerHandle = (*GlobalMessage)(nil)

type GlobalMessage struct {
	Room *socket.RoomStorage
}

func (i *GlobalMessage) Topic() string {
	return "im.message.global"
}

func (i *GlobalMessage) Channel() string {
	return core.GetServerId()
}

func (i *GlobalMessage) Touch() bool {
	return false
}

func (i *GlobalMessage) Do(ctx context.Context, message []byte, attempts uint16) error {
	return nil
}
