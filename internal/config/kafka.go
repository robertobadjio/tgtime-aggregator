package config

import (
	"fmt"
	"net"
	"os"
)

const (
	kafkaHostEnvName = "KAFKA_HOST"
	kafkaPortEnvName = "KAFKA_PORT"
)

// KafkaConfig Конфиг для подключения к Kafka
type KafkaConfig interface {
	Address() string
}

type kafkaConfig struct {
	host string
	port string
}

// NewKafkaConfig Конструктор конфига для подключения к Kafka
func NewKafkaConfig() (KafkaConfig, error) {
	host := os.Getenv(kafkaHostEnvName)
	if len(host) == 0 {
		return nil, fmt.Errorf("environment variable %s must be set", kafkaHostEnvName)
	}

	port := os.Getenv(kafkaPortEnvName)
	if len(port) == 0 {
		return nil, fmt.Errorf("environment variable %s must be set", kafkaPortEnvName)
	}

	return &kafkaConfig{
		host: host,
		port: port,
	}, nil
}

// Address Возвращает адрес для подключения
func (cfg *kafkaConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}
