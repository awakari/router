package config

import (
	"github.com/kelseyhightower/envconfig"
	"time"
)

type Config struct {
	Api struct {
		Port     uint16 `envconfig:"API_PORT" default:"8080" required:"true"`
		Matches  MatchesConfig
		Consumer ConsumerConfig
	}
	Log struct {
		Level int `envconfig:"LOG_LEVEL" default:"-4" required:"true"`
	}
	Queue QueueConfig
}

type MatchesConfig struct {
	Uri       string `envconfig:"API_MATCHES_URI" default:"matches:8080" required:"true"`
	BatchSize uint32 `envconfig:"API_MATCHES_BATCH_SIZE" default:"100" required:"true"`
}

type ConsumerConfig struct {
	Backoff time.Duration `envconfig:"API_CONSUMER_BACKOFF" default:"1s" required:"true"`
	Uri     string        `envconfig:"API_CONSUMER_URI" default:"consumer:8080" required:"true"`
}

type QueueConfig struct {
	BatchSize uint32 `envconfig:"QUEUE_BATCH_SIZE" default:"100" required:"true"`
	FallBack  struct {
		Enabled bool   `envconfig:"QUEUE_FALLBACK_ENABLED" default:"true" required:"true"`
		Suffix  string `envconfig:"QUEUE_FALLBACK_SUFFIX" default:"fallback" required:"true""`
	}
	Limit        uint32        `envconfig:"QUEUE_LIMIT" default:"1000" required:"true"`
	Name         string        `envconfig:"QUEUE_NAME" default:"router" required:"true"`
	BackoffEmpty time.Duration `envconfig:"QUEUE_BACKOFF_EMPTY" default:"1s" required:"true"`
	BackoffError time.Duration `envconfig:"QUEUE_BACKOFF_ERROR" default:"1s" required:"true"`
	Uri          string        `envconfig:"QUEUE_URI" default:"queue:8080" required:"true"`
}

func NewConfigFromEnv() (cfg Config, err error) {
	err = envconfig.Process("", &cfg)
	return
}
