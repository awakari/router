package config

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slog"
	"testing"
)

func TestConfig(t *testing.T) {
	cfg, err := NewConfigFromEnv()
	assert.Nil(t, err)
	assert.Equal(t, uint16(8080), cfg.Api.Port)
	assert.Equal(t, "matches:8080", cfg.Api.Matches.Uri)
	assert.Equal(t, slog.LevelDebug, cfg.Log.Level)
}
