package kafka

import (
	"chatroom/config"
	"context"
	"github.com/segmentio/kafka-go"
)

type KafkaClient struct {
	config   *config.KafkaConfig
	Producer *kafka.Writer
}

// NewKafkaClient initializes a Kafka client using Confluent Kafka library
func NewKafkaClient(c *config.Config) *KafkaClient {
	cfg := c.Kafka
	writer := &kafka.Writer{
		Addr:  kafka.TCP("localhost:9092"), // Kafka broker 地址
		Topic: "im_topic_chat",             // 目标的 Kafka 主题
		Async: true,
	}
	return &KafkaClient{
		Producer: writer,
		config:   cfg,
	}
}

// ProduceMessage sends a message to the Kafka topic
func (k *KafkaClient) ProduceMessage(key, value string) error {
	k.Producer.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(key), // 可选，用于分区
		Value: []byte(value),
	})

	return nil
}
