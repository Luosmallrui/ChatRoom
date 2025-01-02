package process

import (
	"chatroom/config"
	"chatroom/socket/consume"
	"chatroom/socket/process/queue"
	"context"
	"errors"
	"github.com/nsqio/go-nsq"
	"log"
)

type QueueSubscribe struct {
	Config        *config.Config
	GlobalMessage *queue.GlobalMessage
	LocalMessage  *queue.LocalMessage
	RoomControl   *queue.RoomControl
}

func (m *QueueSubscribe) Setup(ctx context.Context) error {

	c := consume.NewConsumer(m.Config.Nsq.Addr, nsq.NewConfig())

	c.Register("default", m.GlobalMessage)
	c.Register("default", m.RoomControl)
	c.Register("default", m.LocalMessage)

	if err := c.Start(ctx, "default"); err != nil {
		log.Fatal(err)
	}

	return errors.New("not implement")
}
