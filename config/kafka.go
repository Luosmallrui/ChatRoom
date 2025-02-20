// Package config config/kafka.go
package config

type KafkaConfig struct {
	Brokers        []string `yaml:"brokers"`
	Topic          string   `yaml:"topic"`
	GroupID        string   `yaml:"group_id"`
	ClientID       string   `yaml:"client_id"`
	Username       string   `yaml:"username"`
	Password       string   `yaml:"password"`
	MinBytes       int      `yaml:"min_bytes"`
	MaxBytes       int      `yaml:"max_bytes"`
	RetryMax       int      `yaml:"retry_max"`
	BatchSize      int      `yaml:"batch_size"`
	BatchTimeout   int      `yaml:"batch_timeout"` // milliseconds
	ReadTimeout    int      `yaml:"read_timeout"`  // seconds
	WriteTimeout   int      `yaml:"write_timeout"` // seconds
	RequiredAcks   int      `yaml:"required_acks"`
	MaxAttempts    int      `yaml:"max_attempts"`
	CommitInterval int      `yaml:"commit_interval"` // seconds
}
