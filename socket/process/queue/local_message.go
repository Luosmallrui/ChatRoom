package queue

import (
	"chatroom/pkg/core"
	"chatroom/pkg/core/socket"
	"chatroom/socket/consume"
	"context"
	"fmt"
)

var _ consume.IConsumerHandle = (*LocalMessage)(nil)

type LocalMessage struct {
	Room *socket.RoomStorage
}

func (i *LocalMessage) Topic() string {
	return fmt.Sprintf("im.message.local.%s", core.GetServerId())
}

func (i *LocalMessage) Channel() string {
	return "default"
}

func (i *LocalMessage) Touch() bool {
	return false
}

func (i *LocalMessage) Do(ctx context.Context, message []byte, attempts uint16) error {
	return nil
}
