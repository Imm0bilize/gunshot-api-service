package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"time"
)

type DBConfig struct {
	Host     string `env:"DB_HOST"`
	Port     string `env:"DB_PORT"`
	Password string `env:"DB_PASSWORD"`
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

type IdemKeyConfig struct {
	TTL time.Duration `env:"IDEMKEY_TTL"`
}

type Config struct {
	HTTP    HTTPConfig
	GRPC    GRPCConfig
	DB      DBConfig
	OTEL    OTELConfig
	IdemKey IdemKeyConfig
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
