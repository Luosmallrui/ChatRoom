package process

import (
	"chatroom/types"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"chatroom/pkg/core"
	"chatroom/socket/consume"

	"chatroom/utils"
	"github.com/redis/go-redis/v9"
	"github.com/sourcegraph/conc/pool"
)

type MessageSubscribe struct {
	redis          *redis.Client
	defaultConsume *consume.ChatSubscribe
	exampleConsume *consume.ExampleSubscribe
}

func NewMessageSubscribe(redis *redis.Client, defaultConsume *consume.ChatSubscribe, exampleConsume *consume.ExampleSubscribe) *MessageSubscribe {
	return &MessageSubscribe{redis: redis, defaultConsume: defaultConsume, exampleConsume: exampleConsume}
}

type IConsume interface {
	Call(event string, data []byte)
}

func (m *MessageSubscribe) Setup(ctx context.Context) error {

	log.Println("Start MessageSubscribe")

	go m.subscribe(ctx, []string{types.ImTopicChat, fmt.Sprintf(types.ImTopicChatPrivate, core.GetServerId())}, m.defaultConsume)

	//go m.subscribe(ctx, []string{entity.ImTopicExample, fmt.Sprintf(entity.ImTopicExamplePrivate, server.ID())}, m.exampleConsume)

	<-ctx.Done()

	return nil
}

func (m *MessageSubscribe) subscribe(ctx context.Context, topic []string, consume IConsume) {
	sub := m.redis.Subscribe(ctx, topic...) //Publish
	//sub := m.kafka.Subscribe(ctx, topic...)

	//sub := m.kafka.Subscribe(ctx, sub)

	defer func() {
		_ = sub.Close()
	}()

	worker := pool.New().WithMaxGoroutines(10)

	for data := range sub.Channel(redis.WithChannelHealthCheckInterval(10 * time.Second)) {
		data := data
		worker.Go(func() {
			m.handle(data, consume)
		})
	}

	worker.Wait()
}

func (m *MessageSubscribe) handle(data *redis.Message, consume IConsume) {
	var in types.SubscribeMessage
	if err := json.Unmarshal([]byte(data.Payload), &in); err != nil {
		log.Println("SubscribeContent Unmarshal Err: ", err.Error())
		return
	}

	defer func() {
		if err := recover(); err != nil {
			log.Println("MessageSubscribe Call Err: ", utils.PanicTrace(err))
		}
	}()

	consume.Call(in.Event, []byte(in.Payload))
}
