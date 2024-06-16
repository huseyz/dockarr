package config

import (
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"
)

type DeleteBehaviour string

const (
	Delete  DeleteBehaviour = "delete"
	Disable DeleteBehaviour = "disable"
	Ignore  DeleteBehaviour = "ignore"
)

type DockarrConfig struct {
	LogLevel        slog.Level
	DeleteBehaviour DeleteBehaviour
	SyncInterval    time.Duration
}

func InitConfig() *DockarrConfig {
	return &DockarrConfig{
		LogLevel:        getLogLevel(),
		DeleteBehaviour: getDeleteBehaviour(),
		SyncInterval:    getSyncInterval(),
	}
}

func getLogLevel() slog.Level {
	level := os.Getenv("LOG_LEVEL")
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	case "":
		return slog.LevelInfo
	default:
		slog.Warn("Invalid log level found, defaulting to Info.")
		return slog.LevelInfo
	}
}

func getDeleteBehaviour() DeleteBehaviour {
	deleteBehaviour := os.Getenv("DELETE_BEHAVIOUR")
	switch strings.ToLower(deleteBehaviour) {
	case "delete", "disable", "ignore":
		return DeleteBehaviour(deleteBehaviour)
	case "":
		return Ignore
	default:
		slog.Warn("Invalid delete behaviour found, defaulting to Ignore.")
		return Delete
	}
}

func getSyncInterval() time.Duration {
	intervalStr := os.Getenv("SYNC_INTERVAL")
	if intervalStr == "" {
		return 1 * time.Minute
	}
	seconds, err := strconv.Atoi(intervalStr)
	if err != nil {
		slog.Warn("Invalid input for sync interval, defaulting to 1 minute.")
		return 1 * time.Minute
	}
	return time.Duration(seconds) * time.Second
}
