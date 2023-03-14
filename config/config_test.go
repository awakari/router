package config

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slog"
	"testing"
	"time"
)

func TestConfig(t *testing.T) {
	cfg, err := NewConfigFromEnv()
	assert.Nil(t, err)
	assert.Equal(t, uint16(8080), cfg.Api.Port)
	assert.Equal(t, "matches:8080", cfg.Api.Matches.Uri)
	assert.Equal(t, uint32(100), cfg.Api.Matches.BatchSize)
	assert.Equal(t, slog.LevelDebug, slog.Level(cfg.Log.Level))
	assert.Equal(t, 1*time.Second, cfg.Api.Consumer.Backoff)
	assert.Equal(t, "consumer:8080", cfg.Api.Consumer.Uri)
	assert.Equal(t, "router", cfg.Queue.Name)
	assert.Equal(t, uint32(1000), cfg.Queue.Limit)
	assert.Equal(t, uint32(100), cfg.Queue.BatchSize)
	assert.Equal(t, 1*time.Second, cfg.Queue.BackoffEmpty)
	assert.Equal(t, 1*time.Second, cfg.Queue.BackoffError)
	assert.True(t, cfg.Queue.FallBack.Enabled)
	assert.Equal(t, "fallback", cfg.Queue.FallBack.Suffix)
	assert.Equal(t, "queue:8080", cfg.Queue.Uri)
}
