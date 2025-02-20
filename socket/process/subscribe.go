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
}

func NewMessageSubscribe(redis *redis.Client, defaultConsume *consume.ChatSubscribe, exampleConsume *consume.ExampleSubscribe,
	client *k.KafkaClient) *MessageSubscribe {
	return &MessageSubscribe{redis: redis, defaultConsume: defaultConsume, exampleConsume: exampleConsume, Kafka: client}
}

type IConsume interface {
	Call(event string, data []byte)
}

func (m *MessageSubscribe) Setup(ctx context.Context) error {

	log.Println("Start MessageSubscribe")

	go m.Subscribe(ctx, []string{types.ImTopicChat, fmt.Sprintf(types.ImTopicChatPrivate, core.GetServerId())}, m.defaultConsume)

	//go m.subscribe(ctx, []string{entity.ImTopicExample, fmt.Sprintf(entity.ImTopicExamplePrivate, server.ID())}, m.exampleConsume)

	<-ctx.Done()

	return nil
}

func (m *MessageSubscribe) Subscribe(ctx context.Context, topics []string, consume IConsume) {
	// 定义 Kafka 消费者配置
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{"localhost:9092"},
		GroupID:     "im_group",      // 消费者组 ID
		Topic:       "im_topic_chat", // 假设只消费一个主题
		MinBytes:    1,               // 最小消息字节数
		MaxBytes:    10e6,            // 最大消息字节数
		StartOffset: kafka.LastOffset,
		MaxWait:     10 * time.Millisecond, // 最大等待时间，默认是 500ms，可调低
	})
	defer reader.Close()

	fmt.Println("Kafka 消费者已启动，等待消息...")

	// 消息处理通道
	msgChan := make(chan kafka.Message, 100)

	// 使用 WaitGroup 来确保所有协程处理完成后退出
	var wg sync.WaitGroup
	workerCount := 10 // 最大并发数
	workerPool := make(chan struct{}, workerCount)

	// 启动消息处理协程
	go func() {
		for msg := range msgChan {
			workerPool <- struct{}{} // 控制并发数量
			wg.Add(1)
			go func(msg kafka.Message) {
				defer wg.Done()
				defer func() { <-workerPool }()

				// 调用处理逻辑
				if err := m.handleKafkaMessage(msg, consume); err != nil {
					log.Printf("处理消息时出错: %v", err)
				}
			}(msg)
		}
	}()

	// 主消费循环
	for {
		select {
		case <-ctx.Done():
			close(msgChan) // 关闭消息处理通道
			wg.Wait()      // 等待所有协程完成
			fmt.Println("消费已停止")
			return
		default:
			// 从 Kafka 读取消息
			msg, err := reader.ReadMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					// 上下文被取消，退出循环
					fmt.Println("消费已停止，收到上下文取消信号")
					close(msgChan)
					wg.Wait()
					return
				}
				log.Printf("读取 Kafka 消息时出错: %v", err)
				continue
			}

			// 将消息发送到处理通道
			msgChan <- msg

			// 提交偏移量
			if err := reader.CommitMessages(ctx, msg); err != nil {
				log.Printf("提交消息偏移量时出错: %v", err)
			}
		}
	}
}

func (m *MessageSubscribe) handleKafkaMessage(msg kafka.Message, consume IConsume) error {
	// 转换Redis消息格式为通用格式
	redisMsg := &redis.Message{
		Channel: msg.Topic,
		Pattern: msg.Topic,
		Payload: string(msg.Value),
	}
	go m.handle(redisMsg, consume)
	return nil
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
