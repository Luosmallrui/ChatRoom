package kafka

import (
	"chatroom/config"
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
)

type KafkaClient struct {
	config    *config.KafkaConfig
	producers map[string]*kafka.Writer // 动态维护每个 topic 对应的 producer
}

// NewKafkaClient initializes a Kafka client
func NewKafkaClient(c *config.Config) *KafkaClient {
	cfg := c.Kafka

	return &KafkaClient{
		config:    cfg,
		producers: make(map[string]*kafka.Writer), // 初始化 producer map
	}
}

// getOrCreateProducer checks if a producer for the given topic exists, otherwise creates a new one
func (k *KafkaClient) getOrCreateProducer(topic string) *kafka.Writer {
	// 如果指定的 topic 已经存在对应的 producer，直接返回
	if producer, exists := k.producers[topic]; exists {
		return producer
	}

	// 创建一个新的 producer 并存储到 producers map
	producer := &kafka.Writer{
		Addr:  kafka.TCP("localhost:9092"), // Kafka broker 地址
		Topic: topic,                       // 动态指定 Kafka 主题
		Async: true,
	}
	k.producers[topic] = producer
	return producer
}

// ProduceMessage sends a message to the specified Kafka topic
func (k *KafkaClient) ProduceMessage(topic, key, value string) error {
	// 根据 topic 获取或创建 producer
	producer := k.getOrCreateProducer(topic)

	// 发送消息
	err := producer.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(key), // 可选，用于分区
		Value: []byte(value),
	})

	return err
}

// Close closes all Kafka producers
func (k *KafkaClient) Close() error {
	var err error
	for topic, producer := range k.producers {
		if closeErr := producer.Close(); closeErr != nil {
			err = closeErr
			// 打印错误日志，或者根据需要处理
			fmt.Printf("Failed to close producer for topic %s: %v\n", topic, closeErr)
		}
	}
	return err
}
