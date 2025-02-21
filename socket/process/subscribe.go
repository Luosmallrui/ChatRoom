package process

import (
	"chatroom/pkg/core"
	k "chatroom/pkg/kafka"
	"chatroom/socket/consume"
	"chatroom/types"
	"chatroom/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
	"log"
	"sync"
	"time"
)

type MessageSubscribe struct {
	redis          *redis.Client
	Kafka          *k.KafkaClient
	defaultConsume *consume.ChatSubscribe
	exampleConsume *consume.ExampleSubscribe
	messageQueue   chan *redis.Message
}

func NewMessageSubscribe(redis *redis.Client, defaultConsume *consume.ChatSubscribe, exampleConsume *consume.ExampleSubscribe,
	client *k.KafkaClient) *MessageSubscribe {
	return &MessageSubscribe{redis: redis, defaultConsume: defaultConsume, exampleConsume: exampleConsume, Kafka: client}
}

type IConsume interface {
	Call(event string, data []byte)
}

func (m *MessageSubscribe) Setup(ctx context.Context) error {
	defer m.Kafka.Close()

	log.Println("Start MessageSubscribe")

	// 初始化消息队列
	m.messageQueue = make(chan *redis.Message, 1000) // 队列大小可根据需求调整
	// 启动工作池，设置并发数为 10（可根据需求调整）
	m.startWorkerPool(ctx, m.defaultConsume, 10)
	// 启动 Kafka 消费者
	go m.Subscribe(ctx, []string{
		types.ImTopicChat,
		fmt.Sprintf(types.ImTopicChatPrivate, core.GetServerId()),
	}, m.defaultConsume)
	<-ctx.Done() // 等待上下文取消信号
	return nil
}

// Subscribe 消费多个主题的消息
func (m *MessageSubscribe) Subscribe(ctx context.Context, topics []string, consume IConsume) {
	var wg sync.WaitGroup

	for _, topic := range topics {
		wg.Add(1)
		go func(topic string) {
			defer wg.Done()
			reader := kafka.NewReader(kafka.ReaderConfig{
				Brokers:     []string{"localhost:9092"},
				GroupID:     "im_group",
				Topic:       topic,
				MinBytes:    1,
				MaxBytes:    10e6,
				MaxWait:     50 * time.Millisecond,
				StartOffset: kafka.LastOffset,
			})
			defer reader.Close()

			for {
				select {
				case <-ctx.Done():
					return
				default:
					msg, err := reader.ReadMessage(ctx)
					if err != nil {
						continue
					}
					m.handleKafkaMessage(msg, consume)
				}
			}
		}(topic)
	}

	wg.Wait()
}

func (m *MessageSubscribe) handleKafkaMessage(msg kafka.Message, consume IConsume) error {
	m.messageQueue <- &redis.Message{
		Channel: msg.Topic,
		Pattern: msg.Topic,
		Payload: string(msg.Value),
	}
	return nil
}
func (m *MessageSubscribe) startWorkerPool(ctx context.Context, consume IConsume, workerCount int) {
	for i := 0; i < workerCount; i++ {
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				case msg := <-m.messageQueue:
					m.handle(msg, consume)
				}
			}
		}()
	}
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
