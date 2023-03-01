package config

import (
	"golang.org/x/exp/slog"
	"os"
	"strconv"
)

type (
	Config struct {
		//
		Api struct {
			//
			Port uint16
			//
			Matches Matches
			//
			Output Output
		}
		//
		Log struct {
			//
			Level slog.Level
		}
	}

	Matches struct {
		//
		Uri string
		//
		BatchSize uint32
	}

	Output struct {
		//
		Uri string
	}
)

const (
	envApiPort             = "API_PORT"
	defApiPort             = "8080"
	envApiMatchesUri       = "API_MATCHES_URI"
	defApiMatchesUri       = "matches:8080"
	envApiMatchesBatchSize = "API_MATCHES_BATCH_SIZE"
	defApiMatchesBatchSize = "100"
	envApiOutputUri        = "API_OUTPUT_URI"
	defApiOutputUri        = "output:8080"
	envLogLevel            = "LOG_LEVEL"
	defLogLevel            = "-4"
)

func NewConfigFromEnv() (cfg Config, err error) {
	apiPortStr := getEnvOrDefault(envApiPort, defApiPort)
	var apiPort uint64
	apiPort, err = strconv.ParseUint(apiPortStr, 10, 16)
	if err != nil {
		return
	}
	cfg.Api.Port = uint16(apiPort)
	cfg.Api.Matches.Uri = getEnvOrDefault(envApiMatchesUri, defApiMatchesUri)
	matchesBatchSizeStr := getEnvOrDefault(envApiMatchesBatchSize, defApiMatchesBatchSize)
	var matchesBatchSize uint64
	matchesBatchSize, err = strconv.ParseUint(matchesBatchSizeStr, 10, 32)
	if err != nil {
		return
	}
	cfg.Api.Matches.BatchSize = uint32(matchesBatchSize)
	cfg.Api.Output.Uri = getEnvOrDefault(envApiOutputUri, defApiOutputUri)
	logLevelStr := getEnvOrDefault(envLogLevel, defLogLevel)
	var logLevel int64
	logLevel, err = strconv.ParseInt(logLevelStr, 10, 16)
	if err != nil {
		return
	}
	cfg.Log.Level = slog.Level(logLevel)
	return
}

func getEnvOrDefault(envKey string, defVal string) (val string) {
	val = os.Getenv(envKey)
	if val == "" {
		val = defVal
	}
	return
}
