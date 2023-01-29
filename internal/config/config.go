package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type DBConfig struct {
	Host     string `env:"DB_HOST"`
	Port     string `env:"DB_PORT"`
	User     string `env:"DB_USER"`
	Password string `env:"DB_PASSWORD"`
	Name     string `env:"DB_NAME"`
}

type HTTPConfig struct {
	Port string `env:"HTTP_PORT"`
}

type GRPCConfig struct {
	Port string `env:"GRPC_PORT"`
}

type OTELConfig struct {
	Host string `env:"OTEL_HOST"`
	Port string `env:"OTEL_PORT"`
}

type KafkaConfig struct {
	Peers string `env:"KAFKA_PEERS"`
	Topic string `env:"KAFKA_TOPIC"`
}

type Config struct {
	HTTP  HTTPConfig
	GRPC  GRPCConfig
	DB    DBConfig
	OTEL  OTELConfig
	Kafka KafkaConfig
}

func New(envFiles ...string) (*Config, error) {
	err := godotenv.Load(envFiles...)
	if err != nil {
		return nil, errors.Wrap(err, "error while load from .env file")
	}

	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, errors.Wrap(err, "error while transfer env to config")
	}

	return &cfg, nil
}
