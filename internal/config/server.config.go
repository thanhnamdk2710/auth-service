package config

import "time"

type ServerConfig struct {
	Environment    string
	LogLevel       string
	Port           string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	IdleTimeout    time.Duration
	MaxHeaderBytes int
}

const (
	DefaultReadTimeoutSec  = 15
	DefaultWriteTimeoutSec = 15
	DefaultIdleTimeoutSec  = 60
	DefaultMaxHeaderBytes  = 1 << 20 // 1MB
)

func NewServerConfig() (*ServerConfig, error) {
	return &ServerConfig{
		Environment:    getEnv("APP_ENV", "development"),
		Port:           getEnv("APP_PORT", "8000"),
		LogLevel:       getEnv("LOG_LEVEL", "info"),
		ReadTimeout:    time.Duration(getEnvAsInt("HTTP_READ_TIMEOUT_SEC", DefaultReadTimeoutSec)) * time.Second,
		WriteTimeout:   time.Duration(getEnvAsInt("HTTP_WRITE_TIMEOUT_SEC", DefaultWriteTimeoutSec)) * time.Second,
		IdleTimeout:    time.Duration(getEnvAsInt("HTTP_IDLE_TIMEOUT_SEC", DefaultIdleTimeoutSec)) * time.Second,
		MaxHeaderBytes: getEnvAsInt("HTTP_MAX_HEADER_BYTES", DefaultMaxHeaderBytes),
	}, nil
}
