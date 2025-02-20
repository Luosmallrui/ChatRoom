// Package kafka pkg/kafka/kafka.go
package kafka

import (
	"chatroom/config"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"strings"
	"time"
)

type KafkaClient struct {
	writer *kafka.Writer
	reader *kafka.Reader
	config *config.KafkaConfig
}

func NewKafkaClient(c *config.Config) *KafkaClient {
	cfg := c.Kafka
	writer := &kafka.Writer{
		Addr:         kafka.TCP(cfg.Brokers...),
		Topic:        cfg.Topic,
		Balancer:     &kafka.LeastBytes{},
		BatchSize:    cfg.BatchSize,
		BatchTimeout: time.Duration(cfg.BatchTimeout) * time.Millisecond,
		RequiredAcks: kafka.RequireAll,
		MaxAttempts:  cfg.MaxAttempts,
		Async:        true,
		Compression:  kafka.Snappy,
		ReadTimeout:  time.Duration(cfg.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Second,
		// 错误处理
		ErrorLogger: kafka.LoggerFunc(func(msg string, args ...interface{}) {
			// 过滤掉超时错误的日志
			if !strings.Contains(msg, "Request Timed Out") {
				log.Printf("Kafka Error: "+msg, args...)
			}
		}),

		// 普通日志
		Logger: kafka.LoggerFunc(func(msg string, args ...interface{}) {
			// 过滤掉超时相关的日志
			if !strings.Contains(msg, "no messages received") {
				log.Printf("Kafka Debug: "+msg, args...)
			}
		}),
	}

	// 配置 Reader
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        cfg.Brokers,
		Topic:          cfg.Topic,
		GroupID:        cfg.GroupID,
		MinBytes:       cfg.MinBytes,
		MaxBytes:       cfg.MaxBytes,
		CommitInterval: time.Duration(cfg.CommitInterval) * time.Second,
		StartOffset:    kafka.LastOffset,
		// 错误处理
		ErrorLogger: kafka.LoggerFunc(func(msg string, args ...interface{}) {
			// 过滤掉超时错误的日志
			if !strings.Contains(msg, "Request Timed Out") {
				log.Printf("Kafka Error: "+msg, args...)
			}
		}),

		// 普通日志
		Logger: kafka.LoggerFunc(func(msg string, args ...interface{}) {
			// 过滤掉超时相关的日志
			if !strings.Contains(msg, "no messages received") {
				log.Printf("Kafka Debug: "+msg, args...)
			}
		}),
	})
	fmt.Println("ok")

	return &KafkaClient{
		writer: writer,
		reader: reader,
		config: cfg,
	}
}

// 日志函数
func logKafkaDebug(msg string, args ...interface{}) {
	log.Printf("Kafka DEBUG: "+msg, args...)
}

func logKafkaError(msg string, args ...interface{}) {
	log.Printf("Kafka ERROR: "+msg, args...)
}

func (k *KafkaClient) Close() error {
	if err := k.writer.Close(); err != nil {
		return fmt.Errorf("error closing writer: %v", err)
	}
	if err := k.reader.Close(); err != nil {
		return fmt.Errorf("error closing reader: %v", err)
	}
	return nil
}
