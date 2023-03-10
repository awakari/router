package config

import (
	"github.com/kelseyhightower/envconfig"
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
	Uri string `envconfig:"API_CONSUMER_URI" default:"consumer:8080" required:"true"`
}

type QueueConfig struct {
	BatchSize          uint32 `envconfig:"QUEUE_BATCH_SIZE" default:"100" required:"true"`
	Limit              uint32 `envconfig:"QUEUE_LIMIT" default:"1000" required:"true"`
	Name               string `envconfig:"QUEUE_NAME" default:"router" required:"true"`
	SleepOnEmptyMillis uint32 `envconfig:"QUEUE_SLEEP_ON_EMPTY_MILLIS" default:"1000" required:"true"`
	SleepOnErrorMillis uint32 `envconfig:"QUEUE_SLEEP_ON_ERROR_MILLIS" default:"1000" required:"true"`
	Uri                string `envconfig:"QUEUE_URI" default:"queue:8080" required:"true"`
}

func NewConfigFromEnv() (cfg Config, err error) {
	err = envconfig.Process("", &cfg)
	return
}
